"""Benchmarks machine translation with chrF (P3-05, DQ-9): our translation service
against Google Translate and NLLB-200, on FLORES-200 and on our own test split.

    uv run --group model vernacular-benchmark results/ --corpus corpus/ --flores flores200/devtest \\
        --language bem --language nya --service http://gpu-host:8090 --google-key $KEY --nllb

Each system translates every test sentence in both directions, English to the language
and back. It writes ``results/results.json`` and ``results/results.md``, a table for the
docs. A system that does not offer a language is left out of that row. The test sets are
never trained on (DQ-10): pass the same FLORES files to vernacular-corpus --exclude.
"""

from __future__ import annotations

import argparse
import html
import json
import urllib.error
import urllib.parse
import urllib.request
from collections.abc import Callable, Sequence
from concurrent.futures import ThreadPoolExecutor
from dataclasses import dataclass
from pathlib import Path

from sacrebleu.metrics import CHRF

from .corpus import read_split
from .prompt import ENGLISH

#: FLORES-200 and NLLB-200 codes.
FLORES = {ENGLISH: "eng_Latn", "bem": "bem_Latn", "nya": "nya_Latn"}
#: Google Translate codes; keep in step with apps/api/internal/candidates/google.go.
GOOGLE = {ENGLISH: "en", "bem": "bem", "nya": "ny"}


@dataclass(frozen=True)
class EvalSet:
    name: str  # flores or vernacular
    language: str
    english: list[str]
    text: list[str]  # aligned with english


#: Translates sentences from one API code to another; None when the system lacks the pair.
System = Callable[[Sequence[str], str, str], list[str] | None]


def flores(directory: Path, language: str) -> EvalSet | None:
    """FLORES-200 devtest for a language, aligned with English; None when FLORES lacks it."""
    if language not in FLORES:
        return None

    def lines(code: str) -> list[str]:
        return (directory / f"{FLORES[code]}.devtest").read_text(encoding="utf-8").splitlines()

    english, text = lines(ENGLISH), lines(language)
    if len(english) != len(text):
        raise ValueError(f"FLORES {language} has {len(text)} lines, English {len(english)}")
    return EvalSet("flores", language, english, text)


def own(corpus: Path, language: str) -> EvalSet:
    """Our held-out test split for a language."""
    pairs = list(read_split(corpus, language, "test"))
    return EvalSet("vernacular", language, [p["en"] for p in pairs], [p["text"] for p in pairs])


def score(systems: dict[str, System], sets: Sequence[EvalSet]) -> list[dict]:
    """chrF of every system on every test set and direction; None where a system lacks the pair."""
    chrf = CHRF()
    rows = []
    for ts in sets:
        for source, target, inputs, refs in [
            (ENGLISH, ts.language, ts.english, ts.text),
            (ts.language, ENGLISH, ts.text, ts.english),
        ]:
            row = {"test_set": ts.name, "direction": f"{source}-{target}", "sentences": len(inputs)}
            for name, translate in systems.items():
                out = translate(inputs, source, target) if inputs else None
                row[name] = None if out is None else round(chrf.corpus_score(out, [refs]).score, 2)
            rows.append(row)
    return rows


def markdown(rows: Sequence[dict], systems: Sequence[str]) -> str:
    """The results as a Markdown table, one row per test set and direction."""
    head = ["Test set", "Direction", "Sentences", *systems]
    lines = ["| " + " | ".join(head) + " |", "|" + "---|" * len(head)]
    for r in rows:
        cells = [r["test_set"], r["direction"], str(r["sentences"])]
        cells += ["–" if r[s] is None else f"{r[s]:.1f}" for s in systems]
        lines.append("| " + " | ".join(cells) + " |")
    return "\n".join(lines) + "\n\nchrF (sacrebleu defaults), higher is better. – means not offered.\n"


def _post_json(url: str, body: dict, timeout: float = 30) -> dict:
    """POSTs JSON and returns the JSON answer; urllib.error.HTTPError on any non-2xx."""
    req = urllib.request.Request(
        url, data=json.dumps(body).encode(), headers={"Content-Type": "application/json"}
    )
    with urllib.request.urlopen(req, timeout=timeout) as res:
        return json.load(res)


def service(url: str, workers: int = 8) -> System:
    """Our translation service, by its contract (docs/translation-service.md): one
    sentence a call, a few calls at a time."""
    endpoint = url.rstrip("/") + "/translate"

    def translate(texts: Sequence[str], source: str, target: str) -> list[str] | None:
        def one(text: str) -> str:
            return _post_json(endpoint, {"text": text, "from": source, "to": target})["translation"]

        try:
            with ThreadPoolExecutor(workers) as pool:
                return list(pool.map(one, texts))
        except urllib.error.HTTPError as e:
            if e.code == 422:  # the service does not offer the pair
                return None
            raise

    return translate


def google(key: str) -> System:
    """Google Translate (v2 REST), 100 sentences a call."""
    endpoint = "https://translation.googleapis.com/language/translate/v2?" + urllib.parse.urlencode(
        {"key": key}
    )

    def translate(texts: Sequence[str], source: str, target: str) -> list[str] | None:
        if source not in GOOGLE or target not in GOOGLE:
            return None
        out: list[str] = []
        for i in range(0, len(texts), 100):
            body = {
                "q": list(texts[i : i + 100]),
                "source": GOOGLE[source],
                "target": GOOGLE[target],
                "format": "text",
            }
            res = _post_json(endpoint, body)
            out += [html.unescape(t["translatedText"]) for t in res["data"]["translations"]]
        return out

    return translate


def nllb(model_name: str = "facebook/nllb-200-distilled-600M") -> System:  # pragma: no cover - needs weights
    """NLLB-200 run locally. Its licence is non-commercial: for comparison only."""
    import torch
    from transformers import AutoModelForSeq2SeqLM, AutoTokenizer

    tok = AutoTokenizer.from_pretrained(model_name)
    model = AutoModelForSeq2SeqLM.from_pretrained(
        model_name, torch_dtype=torch.bfloat16, device_map="auto"
    ).eval()

    def translate(texts: Sequence[str], source: str, target: str) -> list[str] | None:
        if source not in FLORES or target not in FLORES:
            return None
        tok.src_lang = FLORES[source]
        out: list[str] = []
        with torch.inference_mode():
            for i in range(0, len(texts), 16):
                inputs = tok(list(texts[i : i + 16]), return_tensors="pt", padding=True).to(model.device)
                gen = model.generate(
                    **inputs,
                    forced_bos_token_id=tok.convert_tokens_to_ids(FLORES[target]),
                    max_new_tokens=256,
                )
                out += tok.batch_decode(gen, skip_special_tokens=True)
        return out

    return translate


def main(argv: list[str] | None = None) -> None:  # pragma: no cover - calls real systems
    ap = argparse.ArgumentParser(description=__doc__.split("\n\n")[0])
    ap.add_argument("out", type=Path, help="where to write results.json and results.md")
    ap.add_argument("--language", action="append", required=True, help="API code, e.g. bem; repeat")
    ap.add_argument("--corpus", type=Path, help="a vernacular-corpus directory, for our own test split")
    ap.add_argument("--flores", type=Path, help="the FLORES-200 devtest directory")
    ap.add_argument("--service", help="our translation service's base URL")
    ap.add_argument("--google-key", help="a Google Cloud Translation API key")
    ap.add_argument("--nllb", action="store_true", help="also run NLLB-200 locally (needs the model group)")
    args = ap.parse_args(argv)

    systems: dict[str, System] = {}
    if args.service:
        systems["vernacular"] = service(args.service)
    if args.google_key:
        systems["google"] = google(args.google_key)
    if args.nllb:
        systems["nllb-200"] = nllb()
    if not systems:
        ap.error("give at least one of --service, --google-key, --nllb")
    sets: list[EvalSet] = []
    for lang in args.language:
        if args.flores and (fl := flores(args.flores, lang)):
            sets.append(fl)
        if args.corpus:
            sets.append(own(args.corpus, lang))
    rows = score(systems, sets)
    args.out.mkdir(parents=True, exist_ok=True)
    (args.out / "results.json").write_text(json.dumps(rows, indent=2) + "\n", encoding="utf-8")
    table = markdown(rows, list(systems))
    (args.out / "results.md").write_text(table, encoding="utf-8")
    print(table)


if __name__ == "__main__":
    main()
