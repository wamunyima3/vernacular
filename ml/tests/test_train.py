import pytest

from vernacular_ml.prompt import END, answer, prompt
from vernacular_ml.train import IGNORE, Example, chrf_report, encode, examples


def test_prompt_names_languages_with_our_codes():
    p = prompt("Good night", "eng", "toi")
    assert p.startswith("<bos><start_of_turn>user\n")
    assert "You are a professional English (en) to Tonga (toi) translator." in p
    assert "into Tonga:\n\n\nGood night<end_of_turn>" in p  # never the two-letter "to", which is Tongan
    assert p.endswith("<start_of_turn>model\n")


def test_prompt_refuses_unnamed_languages():
    with pytest.raises(ValueError, match="kqn"):
        prompt("water", "eng", "kqn")


def test_each_pair_trains_both_directions():
    got = list(examples([{"language": "bem", "en": "water", "text": "amenshi"}]))
    assert [(e.source, e.target, e.text, e.reference) for e in got] == [
        ("eng", "bem", "water", "amenshi"),
        ("bem", "eng", "amenshi", "water"),
    ]


def test_encode_masks_the_prompt():
    tokenize = lambda s: [ord(c) for c in s]  # noqa: E731 - one token per character
    ex = Example("bem", "eng", "bem", "water", "amenshi")
    enc = encode(ex, tokenize, max_tokens=10_000)
    head = len(prompt("water", "eng", "bem"))
    assert enc["labels"][:head] == [IGNORE] * head
    assert enc["labels"][head:] == tokenize("amenshi" + END) == enc["input_ids"][head:]
    assert len(encode(ex, tokenize, max_tokens=5)["labels"]) == 5


def test_chrf_report_per_language_and_direction():
    exs = [Example("bem", "eng", "bem", "water", "amenshi"), Example("bem", "bem", "eng", "amenshi", "water")]
    report = chrf_report(exs, ["amenshi", "fire"])
    assert report["bem"]["eng-bem"] == {"chrf": 100.0, "sentences": 1}
    assert report["bem"]["bem-eng"]["chrf"] < 20


def test_answer_stops_at_end_of_turn():
    assert answer(" amenshi<end_of_turn>\n<eos>") == "amenshi"
