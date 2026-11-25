import json

from vernacular_ml.benchmark import EvalSet, flores, markdown, own, score


def echo_reference(refs):
    """A system that knows the answers for one direction, and lacks Tonga."""

    def translate(texts, source, target):
        if "toi" in (source, target):
            return None
        return [refs.get(t, "") for t in texts]

    return translate


def test_scores_both_directions_and_marks_missing_pairs():
    sets = [
        EvalSet("vernacular", "bem", ["water"], ["amenshi"]),
        EvalSet("vernacular", "toi", ["water"], ["meenda"]),
    ]
    perfect = echo_reference({"water": "amenshi", "amenshi": "water"})
    rows = score({"ours": perfect, "blank": lambda texts, s, t: ["" for _ in texts]}, sets)
    by = {(r["test_set"], r["direction"]): r for r in rows}
    assert by[("vernacular", "eng-bem")]["ours"] == 100.0
    assert by[("vernacular", "bem-eng")]["ours"] == 100.0
    assert by[("vernacular", "eng-bem")]["blank"] == 0.0
    assert by[("vernacular", "eng-toi")]["ours"] is None
    table = markdown(rows, ["ours", "blank"])
    assert "| vernacular | eng-bem | 1 | 100.0 | 0.0 |" in table and "| – |" in table


def test_reads_flores_aligned_with_english(tmp_path):
    (tmp_path / "eng_Latn.devtest").write_text("Water.\nFire.\n")
    (tmp_path / "bem_Latn.devtest").write_text("Amenshi.\nUmulilo.\n")
    ts = flores(tmp_path, "bem")
    assert ts.english == ["Water.", "Fire."] and ts.text == ["Amenshi.", "Umulilo."]
    assert flores(tmp_path, "toi") is None  # FLORES has no Tonga


def test_reads_our_test_split(tmp_path):
    (tmp_path / "bem").mkdir()
    (tmp_path / "bem" / "test.jsonl").write_text(json.dumps({"en": "water", "text": "amenshi"}) + "\n")
    ts = own(tmp_path, "bem")
    assert (ts.name, ts.english, ts.text) == ("vernacular", ["water"], ["amenshi"])
