"""The translation service (P3-03): the fine-tuned model behind the contract in
docs/translation-service.md, which the API calls for labelled machine output (TRN-3).

    MODEL_DIR=runs/2026.1 uv run --group serve uvicorn --factory vernacular_ml.service:from_env --port 8090

MODEL_DIR is a directory written by vernacular-train: the LoRA adapter and report.json,
which names the base model and the languages it was trained on. Pairs between English
and those languages are served; anything else is 422. MODEL_NAME, the name API users see,
defaults to ``vernacular-mt-<directory name>``.
"""

from __future__ import annotations

import json
import math
import os
from collections.abc import Sequence
from dataclasses import dataclass
from pathlib import Path
from typing import Protocol

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

from .model import load_base
from .prompt import ENGLISH, answer, prompt


class Request(BaseModel):
    text: str = Field(min_length=1, max_length=500)
    from_: str = Field(alias="from")
    to: str


class Response(BaseModel):
    translation: str
    confidence: float = Field(ge=0, le=1)
    model: str


@dataclass(frozen=True)
class Output:
    text: str
    token_logprobs: Sequence[float]  # of each generated token


class Model(Protocol):
    def generate(self, prompt: str) -> Output: ...


def confidence(token_logprobs: Sequence[float]) -> float:
    """The geometric mean of the tokens' probabilities: 1 when the model was sure of every token."""
    if not token_logprobs:
        return 0.0
    return round(min(1.0, math.exp(sum(token_logprobs) / len(token_logprobs))), 4)


def create(model: Model, languages: set[str], name: str) -> FastAPI:
    """The service for model, translating between English and languages."""
    app = FastAPI(title="Vernacular translation service", docs_url=None, redoc_url=None)

    @app.post("/translate")
    def translate(req: Request) -> Response:
        pair = {req.from_, req.to}
        if ENGLISH not in pair or len(pair) != 2 or not (pair - {ENGLISH}) <= languages:
            raise HTTPException(422, f"Translates between {ENGLISH} and {', '.join(sorted(languages))}")
        out = model.generate(prompt(req.text, req.from_, req.to))
        text = answer(out.text)
        if not text:
            raise HTTPException(502, "The model gave no translation")
        return Response(translation=text, confidence=confidence(out.token_logprobs), model=name)

    @app.get("/healthz")
    def healthz() -> dict[str, str]:
        return {"status": "ok", "model": name}

    return app


class Transformers:  # pragma: no cover - needs the model weights
    """The base model with the LoRA adapter, greedy decoding."""

    def __init__(self, base: str, adapter: Path, max_new_tokens: int = 256) -> None:
        from peft import PeftModel

        self.tok, model = load_base(base)
        self.model = PeftModel.from_pretrained(model, adapter).eval()
        self.max_new_tokens = max_new_tokens

    def generate(self, prompt: str) -> Output:
        import torch

        inputs = self.tok(prompt, return_tensors="pt", add_special_tokens=False).to(self.model.device)
        with torch.inference_mode():
            out = self.model.generate(
                **inputs,
                max_new_tokens=self.max_new_tokens,
                do_sample=False,
                output_scores=True,
                return_dict_in_generate=True,
            )
        new = out.sequences[0][inputs["input_ids"].shape[1] :]
        scores = self.model.compute_transition_scores(out.sequences, out.scores, normalize_logits=True)[0]
        return Output(self.tok.decode(new), [float(s) for s in scores])


def from_env() -> FastAPI:  # pragma: no cover - needs the model weights
    """The app for MODEL_DIR, for `uvicorn --factory`."""
    model_dir = Path(os.environ["MODEL_DIR"])
    report = json.loads((model_dir / "report.json").read_text(encoding="utf-8"))
    name = os.environ.get("MODEL_NAME", f"vernacular-mt-{model_dir.name}")
    return create(Transformers(report["base"], model_dir), set(report["languages"]), name)
