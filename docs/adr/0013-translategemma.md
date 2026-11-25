# 0013. TranslateGemma as the base translation model

- Status: Accepted
- Date: 2026-11-25

## Context

Phase 3 serves labelled machine translation when no verified phrase matches (TRN-3). We fine-tune an open model on our verified pairs (P3-02a corpus) and serve it as a separate service (P3-03, [translation-service.md](../translation-service.md)). The model's licence must allow commercial serving, because API users may be businesses.

Candidates considered:

| Model                                                | Licence                                                                         | Notes                                                                                           |
| ---------------------------------------------------- | ------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| TranslateGemma 4B / 12B / 27B (Google, January 2026) | Gemma Terms of Use: commercial use allowed, with Google's prohibited-use policy | Built for translation on Gemma 3; 55 languages fully evaluated, about 500 more pairs trained on |
| MADLAD-400 3B / 7B / 10B MT (Google, 2023)           | Apache 2.0 or ODC-BY (to confirm on the model card)                             | T5-style, 400+ languages, older                                                                 |
| NLLB-200, SeamlessM4T, Aya                           | CC BY-NC                                                                        | Non-commercial: ruled out, though NLLB covers Bemba and Nyanja                                  |
| Opus-MT (Helsinki)                                   | CC BY 4.0                                                                       | Zambian pairs mostly trained on JW300, withdrawn over copyright: fails the source register      |

## Decision

Use **TranslateGemma**, starting with **`google/translategemma-4b-it`**. It is the newest translation-specific open model, and 4B can be fine-tuned with LoRA on one 24 GB GPU and served with vLLM on a modest one. Move to 12B only if the P3-05 benchmark shows 4B falls short.

- **Language codes.** Its chat template expects two-letter ISO 639-1 codes. Bemba (`bem`), Tonga of Zambia (`toi`) and Lozi (`loz`) have none, and the two-letter `to` is Tongan, a different language. We fine-tune with the raw prompt, `<translate source_lang=… target_lang=…>`, using our own ISO 639-3 codes, and the service builds the same prompt. Nyanja uses `ny` if the base model already knows it, which P3-05 checks.
- **Training data** comes only from the P3-02a corpus: verified items whose sources allow `train`, with held-out test sets excluded (DQ-10). Google Translate and other machine output is never trained on unless its source allows it.
- **Licence obligations.** The Gemma Terms of Use pass on to anyone given the weights. We publish the fine-tuned weights with those terms and the prohibited-use policy, and say so in the API docs for machine output. The code stays MIT and the data CC BY; only the model weights carry Gemma's terms.

## Consequences

- Machine output depends on a non-OSI licence. If Google changes the Gemma terms, or a later Gemma moves to Apache 2.0, revisit this ADR. MADLAD-400 is the fallback, since the corpus and service contract do not depend on the model.
- The 55 evaluated languages do not include our languages, so quality rests on fine-tuning. P3-05 must show it beats Google Translate and NLLB-200 on FLORES-200 (Bemba, Nyanja) and on our own test split (all four languages) before `TRANSLATION_SERVICE_URL` is set in production.
- A GPU is needed for training and, for acceptable latency, for serving. That is a hosting cost for P3-03, recorded in [deploy.md](../deploy.md) when it is chosen.
