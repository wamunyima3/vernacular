# Phrase pack source list

What each phrase pack (PHR-1) should say, in English. `make seed` loads it into `phrase_list`, linking each line to the concept it expresses, so AI candidates (`make phrase-candidates LANG=bem`, P1-AI-04) and reviewers can supply each language's phrase.

- **Format:** CSV with the header `pack,source_en,notes`. Packs: `greetings`, `numbers`, `dates`, `interface`, `health`, `money`.
- **notes** tell translators the context, such as "button" or "said to someone leaving". Placeholders (`…`) are a Phase 2 feature (PHR-2).
- **Licence:** same as the dataset (CC BY 4.0 proposed, P0-OUT-08).
- Importing is idempotent; lines are matched by pack and English text.
