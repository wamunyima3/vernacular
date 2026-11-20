# Core concept lists

Our own English concept lists, added to the core list (LEX-7) by `make seed`. GATITOS supplies about 3,800 general concepts; these lists add everyday Zambian terms that GATITOS lacks or has no domain for.

- **Format:** CSV with the header `domain,gloss_en`. One concept per row. Lower case, except "I", acronyms and proper nouns. Add a hint in brackets when the English word is ambiguous: `change (money)`.
- **Licence:** same as the dataset (CC BY 4.0 proposed, P0-OUT-08).
- Importing is idempotent. A row that matches an existing GATITOS concept gives it a domain and does not create a duplicate.
