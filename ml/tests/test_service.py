import math

from fastapi.testclient import TestClient

from vernacular_ml.service import Output, confidence, create


class Fake:
    def __init__(self, text="Mulale bwino<end_of_turn>"):
        self.text, self.prompts = text, []

    def generate(self, prompt):
        self.prompts.append(prompt)
        return Output(self.text, [math.log(0.5), math.log(0.5)])


def client(model):
    return TestClient(create(model, {"bem"}, "vernacular-mt-test"))


def test_translates_to_the_contract():
    model = Fake()
    res = client(model).post("/translate", json={"text": "Good night", "from": "eng", "to": "bem"})
    assert res.status_code == 200
    assert res.json() == {"translation": "Mulale bwino", "confidence": 0.5, "model": "vernacular-mt-test"}
    assert "source_lang=en target_lang=bem" in model.prompts[0]


def test_refuses_pairs_it_was_not_trained_on():
    for pair in [("eng", "nya"), ("bem", "bem"), ("bem", "nya")]:
        res = client(Fake()).post("/translate", json={"text": "x", "from": pair[0], "to": pair[1]})
        assert res.status_code == 422, pair


def test_an_empty_answer_is_an_error():
    res = client(Fake("<end_of_turn>")).post("/translate", json={"text": "x", "from": "bem", "to": "eng"})
    assert res.status_code == 502


def test_confidence():
    assert confidence([]) == 0.0
    assert confidence([0.0, 0.0]) == 1.0
    assert confidence([math.log(0.25), math.log(1.0)]) == 0.5
