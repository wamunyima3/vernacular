# Google SMOL and GATITOS: Zambian language coverage

Checked on 2026-11-20 against the `google/smol` dataset on Hugging Face (task P0-CNT-01). GATITOS ships as one subset of SMOL.

- **Licence:** CC BY 4.0. Redistribution is allowed with attribution, so the data can enter the source register (DQ-1).
- **Format:** JSONL, one pair per line: `{"src", "trgs": [...], "sl", "tl", "is_source_orig"}`. `trgs` can hold several acceptable translations.

## Coverage

| Language  | Code  | GATITOS (words, short phrases) | SmolSent (sentences) | SmolDoc (documents) |
| --------- | ----- | ------------------------------ | -------------------- | ------------------- |
| Bemba     | `bem` | en→bem 3,977 · bem→en 2,864    | en→bem 863           | en→bem 130          |
| Nyanja    | `ny`  | **none**                       | en→ny 863            | en→ny 354           |
| Tonga     | `toi` | none                           | none                 | none                |
| Lozi      | `loz` | none                           | none                 | none                |
| Kaonde    | `kqn` | none                           | none                 | none                |
| Lunda     | `lun` | none                           | none                 | none                |
| Luvale    | `lue` | none                           | none                 | none                |
| _Tumbuka_ | `tum` | en→tum 3,983                   | en→tum               | en→tum              |

Traps when matching codes:

- `to` in SMOL is **Tongan** (Polynesian), not Zambian Tonga (`toi`).
- `lu` is **Luba-Katanga**, not Lunda (`lun`).
- Tumbuka is spoken in eastern Zambia but is not one of our seven languages. It is noted for later.

## What this means for the plan

1. **Bemba** has a full GATITOS list. It seeds the core concept list (P0-CNT-03), and each pair becomes a `sourced` candidate.
2. **Nyanja** has no GATITOS list. Its core list has to come from Claude/Google candidates (P1-AI-01, P1-AI-02) and licensed dictionaries. SmolSent gives 863 sentences for examples and phrases.
3. **Tonga, Lozi, Kaonde, Lunda and Luvale** have nothing in SMOL. They depend on outreach (P0-OUT-01 to P0-OUT-03) and volunteers.
4. SmolSent sentences are the same English source set across languages, so `en→bem` and `en→ny` rows line up through the English sentence. That gives a free bem↔ny alignment for 863 sentences.

## Attribution

> Caswell et al., "SMOL: Professionally translated parallel data for 115 under-represented languages", 2025. Google, CC BY 4.0.
