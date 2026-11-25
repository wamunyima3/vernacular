"""Loading the base model, shared by training and the service. Needs the `model`
dependency group; torch and transformers are imported on call."""

from __future__ import annotations

from typing import Any


def load_base(base: str) -> tuple[Any, Any]:  # pragma: no cover - needs the model weights
    """The tokenizer and the base model in bfloat16, placed on the available devices."""
    import torch
    from transformers import AutoModelForCausalLM, AutoTokenizer

    tok = AutoTokenizer.from_pretrained(base)
    model = AutoModelForCausalLM.from_pretrained(base, torch_dtype=torch.bfloat16, device_map="auto")
    return tok, model
