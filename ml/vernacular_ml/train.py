"""Fine-tunes TranslateGemma with LoRA on the parallel corpus and reports chrF (P3-02b).

    uv run --group train vernacular-train corpus/ runs/2026.1 --language bem --language nya

Each pair trains both directions, English to the language and back. The loss covers the
answer only, not the prompt. After training, the dev split is translated and scored with
chrF per language and direction, written to ``<out>/report.json`` next to the LoRA
adapter. The test split is left for the benchmark (P3-05).

The base model is gated: accept the Gemma Terms of Use on Hugging Face and set HF_TOKEN.
Training needs a GPU; the ``train`` dependency group's packages are imported inside main().
"""

from __future__ import annotations

import argparse
import json
from collections import defaultdict
from collections.abc import Callable, Iterable, Iterator, Sequence
from dataclasses import dataclass
from pathlib import Path

from sacrebleu.metrics import CHRF

from . import licence
from .corpus import read_split
from .model import load_base
from .prompt import END, ENGLISH, answer, prompt

BASE_MODEL = "google/translategemma-4b-it"
IGNORE = -100  # label index the loss skips


@dataclass(frozen=True)
class Example:
    language: str  # the Zambian language of the pair
    source: str  # API code of the input's language
    target: str
    text: str
    reference: str

    @property
    def direction(self) -> str:
        return f"{self.source}-{self.target}"


def examples(pairs: Iterable[dict]) -> Iterator[Example]:
    """Both directions of each corpus pair."""
    for p in pairs:
        yield Example(p["language"], ENGLISH, p["language"], p["en"], p["text"])
        yield Example(p["language"], p["language"], ENGLISH, p["text"], p["en"])


def encode(ex: Example, tokenize: Callable[[str], list[int]], max_tokens: int) -> dict[str, list[int]]:
    """Token IDs and labels for one example, the prompt masked out of the loss."""
    head = tokenize(prompt(ex.text, ex.source, ex.target))
    tail = tokenize(ex.reference + END)
    ids = (head + tail)[:max_tokens]
    labels = ([IGNORE] * len(head) + tail)[:max_tokens]
    return {"input_ids": ids, "attention_mask": [1] * len(ids), "labels": labels}


Score = dict[str, float | int]  # {"chrf": ..., "sentences": ...}


def chrf_report(exs: Sequence[Example], hypotheses: Sequence[str]) -> dict[str, dict[str, Score]]:
    """chrF per language and direction, with how many sentences each covers."""
    chrf = CHRF()
    groups: dict[tuple[str, str], tuple[list[str], list[str]]] = defaultdict(lambda: ([], []))
    for ex, hyp in zip(exs, hypotheses, strict=True):
        hyps, refs = groups[(ex.language, ex.direction)]
        hyps.append(hyp)
        refs.append(ex.reference)
    report: dict[str, dict[str, Score]] = defaultdict(dict)
    for (language, direction), (hyps, refs) in sorted(groups.items()):
        score = chrf.corpus_score(hyps, [refs]).score
        report[language][direction] = {"chrf": round(score, 2), "sentences": len(hyps)}
    return dict(report)


def main(argv: list[str] | None = None) -> None:  # pragma: no cover - needs a GPU and the model
    ap = argparse.ArgumentParser(description=__doc__.split("\n\n")[0])
    ap.add_argument("corpus", type=Path, help="a directory written by vernacular-corpus")
    ap.add_argument("out", type=Path, help="where to write the adapter and report.json")
    ap.add_argument("--language", action="append", required=True, help="API code, e.g. bem; repeat")
    ap.add_argument("--base", default=BASE_MODEL)
    ap.add_argument("--epochs", type=float, default=3)
    ap.add_argument("--lr", type=float, default=2e-4)
    ap.add_argument("--lora-rank", type=int, default=16)
    ap.add_argument("--max-tokens", type=int, default=256)
    args = ap.parse_args(argv)

    import torch
    from datasets import Dataset
    from peft import LoraConfig, get_peft_model
    from transformers import DataCollatorForSeq2Seq, Trainer, TrainingArguments

    tok, model = load_base(args.base)
    model = get_peft_model(
        model,
        LoraConfig(
            r=args.lora_rank,
            lora_alpha=2 * args.lora_rank,
            lora_dropout=0.05,
            target_modules="all-linear",
            task_type="CAUSAL_LM",
        ),
    )

    def tokenize(s: str) -> list[int]:
        return tok(s, add_special_tokens=False)["input_ids"]

    def load(split: str) -> list[Example]:
        return [ex for lang in args.language for ex in examples(read_split(args.corpus, lang, split))]

    train = Dataset.from_list([encode(ex, tokenize, args.max_tokens) for ex in load("train")])
    Trainer(
        model=model,
        args=TrainingArguments(
            output_dir=str(args.out),
            num_train_epochs=args.epochs,
            learning_rate=args.lr,
            per_device_train_batch_size=8,
            gradient_accumulation_steps=2,
            bf16=True,
            logging_steps=20,
            save_strategy="no",
            report_to=[],
        ),
        train_dataset=train,
        data_collator=DataCollatorForSeq2Seq(tok, label_pad_token_id=IGNORE),
    ).train()
    model.save_pretrained(args.out)
    licence.write(args.out, args.base, args.language)

    dev = load("dev")
    model.eval()
    tok.padding_side = "left"  # so every prompt ends where generation starts
    hypotheses = []
    with torch.inference_mode():
        for i in range(0, len(dev), 32):
            batch = [prompt(ex.text, ex.source, ex.target) for ex in dev[i : i + 32]]
            inputs = tok(batch, return_tensors="pt", padding=True, add_special_tokens=False).to(model.device)
            out = model.generate(**inputs, max_new_tokens=args.max_tokens, do_sample=False)
            start = inputs["input_ids"].shape[1]
            hypotheses += [answer(tok.decode(row[start:])) for row in out]
    report = {"base": args.base, "languages": args.language, "dev_chrf": chrf_report(dev, hypotheses)}
    (args.out / "report.json").write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(report["dev_chrf"], indent=2))


if __name__ == "__main__":
    main()
