"""Builds an English-to-language parallel corpus from a Vernacular release zip.

Every pair is a verified item: a sense (its English gloss and headword), an
example sentence or a phrase. Only items whose every cited source allows
``train`` are used (the source register, DQ-1). Pairs are split into train,
dev and test by their English text, so one meaning never crosses splits and
the split is the same in every release. Anything in a held-out test set, such
as FLORES, is left out entirely (DQ-10).

    vernacular-corpus vernacular-2026.1.zip corpus/ --exclude flores200.bem_Latn.devtest

writes ``corpus/<language>/{train,dev,test}.jsonl`` and ``corpus/manifest.json``.
"""

from __future__ import annotations

import argparse
import csv
import hashlib
import io
import json
import unicodedata
import zipfile
from collections import Counter, defaultdict
from collections.abc import Iterable, Iterator
from dataclasses import asdict, dataclass
from pathlib import Path

SPLITS = ("train", "dev", "test")
#: Share of English texts, per thousand, in the dev and test splits.
DEV_PER_MILLE, TEST_PER_MILLE = 50, 50


@dataclass(frozen=True)
class Pair:
    id: str
    language: str
    kind: str  # sense, example or phrase
    en: str
    text: str


def normalise(s: str) -> str:
    """Text as compared for duplicates and exclusion: NFC, case-folded, spaces collapsed."""
    return " ".join(unicodedata.normalize("NFC", s).casefold().split())


def split_of(en: str) -> str:
    """The split an English text belongs to, fixed by its hash."""
    return _split(normalise(en))


def _split(normalised: str) -> str:
    bucket = int.from_bytes(hashlib.sha256(normalised.encode()).digest()[:4], "big") % 1000
    if bucket < TEST_PER_MILLE:
        return "test"
    return "dev" if bucket < TEST_PER_MILLE + DEV_PER_MILLE else "train"


def _rows(release: zipfile.ZipFile, name: str) -> Iterator[dict[str, str]]:
    with release.open(name) as f:
        yield from csv.DictReader(io.TextIOWrapper(f, encoding="utf-8", newline=""))


def trainable(release: zipfile.ZipFile) -> set[str]:
    """IDs of items whose every cited source allows training; an item citing none is left out."""
    allowed = {s["name"] for s in _rows(release, "sources.csv") if "train" in s["allowed_uses"].split()}
    cited, refused = set(), set()
    for r in _rows(release, "item_sources.csv"):
        cited.add(r["item_id"])
        if r["source"] not in allowed:
            refused.add(r["item_id"])
    return cited - refused


def pairs(release: zipfile.ZipFile) -> Iterator[Pair]:
    """Every pair in the release, before any filtering."""
    entries = {e["id"]: e for e in _rows(release, "entries.csv")}
    senses = {}  # a release is a snapshot, so every sense's entry is in it
    for s in _rows(release, "senses.csv"):
        entry = entries[s["entry_id"]]
        senses[s["id"]] = entry["language"]
        yield Pair(s["id"], entry["language"], "sense", s["gloss"], entry["headword"])
    for x in _rows(release, "examples.csv"):
        yield Pair(x["id"], senses[x["sense_id"]], "example", x["translation_en"], x["text"])
    for p in _rows(release, "phrases.csv"):
        yield Pair(p["id"], p["language"], "phrase", p["source_en"], p["text"])


def build(release: zipfile.ZipFile, held_out: Iterable[str] = ()) -> dict[str, dict[str, list[Pair]]]:
    """Pairs by language and split: trainable, not held out, without duplicates."""
    excluded = {normalise(s) for s in held_out if s.strip()}
    allowed = trainable(release)
    seen: set[tuple[str, str, str]] = set()
    out: dict[str, dict[str, list[Pair]]] = defaultdict(lambda: {s: [] for s in SPLITS})
    for p in pairs(release):
        en, text = normalise(p.en), normalise(p.text)
        if p.id not in allowed or not en or not text:
            continue
        key = (p.language, en, text)
        if en in excluded or text in excluded or key in seen:
            continue
        seen.add(key)
        out[p.language][_split(en)].append(p)
    return dict(out)


def write(corpus: dict[str, dict[str, list[Pair]]], out: Path, release_name: str) -> dict:
    """Writes one JSON line per pair, and a manifest of the counts."""
    manifest: dict = {"release": release_name, "languages": {}}
    for language, splits in sorted(corpus.items()):
        (out / language).mkdir(parents=True, exist_ok=True)
        counts = {}
        for split, items in splits.items():
            with (out / language / f"{split}.jsonl").open("w", encoding="utf-8") as f:
                for p in items:
                    f.write(json.dumps(asdict(p), ensure_ascii=False) + "\n")
            counts[split] = dict(Counter(p.kind for p in items), total=len(items))
        manifest["languages"][language] = counts
    (out / "manifest.json").write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
    return manifest


def read_split(corpus: Path, language: str, split: str) -> Iterator[dict]:
    """The pairs of one language's split, as write() stored them."""
    if split not in SPLITS:
        raise ValueError(f"no split {split!r}; use one of {SPLITS}")
    with (corpus / language / f"{split}.jsonl").open(encoding="utf-8") as f:
        for line in f:
            yield json.loads(line)


def main(argv: list[str] | None = None) -> None:
    ap = argparse.ArgumentParser(description=__doc__.split("\n\n")[0])
    ap.add_argument("release", type=Path, help="a release zip from GET /v1/releases/{name}/download")
    ap.add_argument("out", type=Path, help="the directory to write the corpus to")
    ap.add_argument(
        "--exclude", type=Path, action="append", default=[], help="a held-out test set, one sentence a line"
    )
    args = ap.parse_args(argv)
    held_out = [line for path in args.exclude for line in path.read_text(encoding="utf-8").splitlines()]
    with zipfile.ZipFile(args.release) as release:
        name = json.loads(release.read("release.json")).get("name", args.release.stem)
        manifest = write(build(release, held_out), args.out, name)
    for language, counts in manifest["languages"].items():
        print(language, " ".join(f"{s}={counts[s]['total']}" for s in SPLITS))


if __name__ == "__main__":
    main()
