import csv
import io
import json
import zipfile
from pathlib import Path

from vernacular_ml.corpus import SPLITS, build, main, normalise, split_of


def _csv(rows: list[dict[str, str]]) -> str:
    out = io.StringIO()
    w = csv.DictWriter(out, fieldnames=list(rows[0]))
    w.writeheader()
    w.writerows(rows)
    return out.getvalue()


def release(path: Path) -> Path:
    """A release with one open source, one that forbids training, and items citing each."""
    files = {
        "release.json": json.dumps({"name": "2026.1"}),
        "sources.csv": _csv(
            [
                {"name": "Open", "allowed_uses": "redistribute derive train"},
                {"name": "Closed", "allowed_uses": "redistribute"},
            ]
        ),
        "entries.csv": _csv(
            [
                {"id": "e1", "language": "bem", "headword": "amenshi"},
                {"id": "e2", "language": "bem", "headword": "umulilo"},
            ]
        ),
        "senses.csv": _csv(
            [
                {"id": "s1", "entry_id": "e1", "gloss": "water"},
                {"id": "s2", "entry_id": "e2", "gloss": "fire"},
            ]
        ),
        "examples.csv": _csv(
            [{"id": "x1", "sense_id": "s1", "text": "Nalanwa amenshi.", "translation_en": "I drank water."}]
        ),
        "phrases.csv": _csv(
            [
                {"id": "p1", "language": "nya", "source_en": "Good morning", "text": "Mwadzuka bwanji"},
                {"id": "p2", "language": "nya", "source_en": "good  MORNING", "text": "mwadzuka bwanji"},
                {"id": "p3", "language": "nya", "source_en": "Held out", "text": "Sindikudziwa"},
            ]
        ),
        "item_sources.csv": _csv(
            [
                {"item_id": "s1", "source": "Open"},
                {"item_id": "s2", "source": "Open"},
                {"item_id": "s2", "source": "Closed"},  # one source forbids training: left out
                {"item_id": "x1", "source": "Open"},
                {"item_id": "p1", "source": "Open"},
                {"item_id": "p2", "source": "Open"},
                {"item_id": "p3", "source": "Open"},
            ]
        ),
    }
    with zipfile.ZipFile(path, "w") as z:
        for name, body in files.items():
            z.writestr(name, body)
    return path


def ids(corpus, language):
    return sorted(p.id for split in SPLITS for p in corpus[language][split])


def test_keeps_trainable_unique_pairs_not_held_out(tmp_path):
    with zipfile.ZipFile(release(tmp_path / "r.zip")) as z:
        corpus = build(z, held_out=["held out"])
    assert ids(corpus, "bem") == ["s1", "x1"]  # s2 also cites a source that forbids training
    assert ids(corpus, "nya") == ["p1"]  # p2 repeats p1, p3 is held out


def test_split_is_by_english_and_stable():
    assert split_of("Water") == split_of("  water ")
    counts = {s: 0 for s in SPLITS}
    for i in range(10_000):
        counts[split_of(f"sentence {i}")] += 1
    assert 400 < counts["test"] < 600 and 400 < counts["dev"] < 600


def test_normalise():
    assert normalise("  Ámenshi\tNO ") == normalise("ámenshi no")


def test_cli_writes_splits_and_manifest(tmp_path, capsys):
    main([str(release(tmp_path / "r.zip")), str(tmp_path / "out")])
    manifest = json.loads((tmp_path / "out" / "manifest.json").read_text())
    assert manifest["release"] == "2026.1"
    assert sum(manifest["languages"]["nya"][s]["total"] for s in SPLITS) == 2
    lines = [
        json.loads(line)
        for s in SPLITS
        for line in (tmp_path / "out" / "bem" / f"{s}.jsonl").read_text().splitlines()
    ]
    assert {line["kind"] for line in lines} == {"sense", "example"}
    assert "bem" in capsys.readouterr().out
