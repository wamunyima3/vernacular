from vernacular_ml import licence


def test_writes_the_notice_and_a_card_marking_the_adapter_modified(tmp_path):
    licence.write(tmp_path, "google/translategemma-4b-it", ["bem", "nya"])
    assert (tmp_path / "NOTICE").read_text() == (
        "Gemma is provided under and subject to the Gemma Terms of Use found at ai.google.dev/gemma/terms\n"
    )
    card = (tmp_path / "README.md").read_text()
    assert card.startswith("---\nbase_model: google/translategemma-4b-it\n")
    assert "license: gemma" in card and "- bem\n- nya\n" in card
    assert "**Modified from [google/translategemma-4b-it]" in card
    assert "English pairs in\nBemba, Nyanja." in card
    assert licence.TERMS in card and licence.PROHIBITED_USE in card
