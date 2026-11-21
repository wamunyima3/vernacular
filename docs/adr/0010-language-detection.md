# 0010. Our own n-gram language detector, trained on our data

- Status: Accepted
- Date: 2026-11-21

## Context

`POST /v1/detect` must name which supported language a text is in, with a confidence score (TRN-1, US-09). The requirements ask us to evaluate existing language ID models first.

- **GlotLID v3** (fastText, Apache 2.0) has over 2,000 labels and is the strongest open option for small languages. The model is about 1.7 GB, and fastText has no maintained pure-Go inference, so serving it means a separate Python or C++ service next to the Go monolith (ADR-0002). Whether every Zambian language we plan (Tonga, Lozi, Kaonde, Lunda, Luvale) has a label still needs checking against its label list.
- **Pure-Go detectors** (lingua-go, whatlanggo) cover 75–84 major languages. We found no Zambian language among them.
- The API only needs to tell apart the languages it serves, plus English, which is mixed into most Zambian text. That is a small, closed set.

## Decision

- A small character n-gram classifier in Go (`internal/detect`): multinomial naive Bayes over 1–4-character n-grams of the folded text (NFR-7), with add-one smoothing. Confidence is the softmax of the per-language log-likelihoods, averaged per n-gram and scaled by a fixed temperature, so it reflects how distinct a text is rather than how long.
- It is trained when the API starts, from our own data: each active language's servable headwords, examples and phrases, and English from concept glosses and example translations. Languages with too little text are left out of the model, not guessed.
- GlotLID stays an offline tool: a baseline to measure our detector against (DQ), and a filter if we crawl text later.

## Consequences

- No new service, no large model file, no licence beyond our own data's. Detection improves as reviewed data grows, and a new language works once it has data.
- Accuracy is weak on very short text and between closely related languages (Bemba and Lunda, Nyanja and Chewa varieties). Responses rank languages rather than give one answer, so callers can see the margin.
- Training at start-up reads a few thousand rows; at larger sizes the model should be built per release and loaded from the release artefacts.
