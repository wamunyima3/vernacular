# User stories

Each story names its persona, the requirement IDs from [requirements.md](requirements.md) it covers, and acceptance criteria. Tasks in [tasks.md](../tasks.md) reference stories by ID.

Every response item carries the same trust fields: `status`, `sources` (with licence), `reviews` and `release`. Acceptance criteria below assume them and do not repeat them.

## Personas

| Persona | Who they are |
|---|---|
| **Developer** | App, web, SMS or USSD developer adding local language support. This includes government, health, agriculture and finance teams. |
| **Learner** | School pupil, teacher, researcher or member of the public who needs a reliable dictionary. |
| **Contributor** | Registered native speaker who proposes entries and corrections. |
| **Reviewer** | Native speaker approved for one language, who works through the review queue. |
| **Senior reviewer** | Experienced reviewer who resolves disagreements and reports, and qualifies new reviewers. |
| **Admin** | Project owner, who runs sources, AI jobs, releases and reviewer approval. |

## Coverage

| Requirement group | Stories |
|---|---|
| LEX-1…8 | US-03, US-04, US-05, US-16, US-33 |
| PHR-1…3 | US-08, US-12, US-32 |
| TRN-1…4 | US-06, US-09, US-10, US-31 |
| AI-1…7 | US-07, US-20, US-21, US-26, US-27 |
| PRV-1…4 | US-06, US-11, US-15, US-22, US-29 |
| CON-1…5 | US-18, US-19, US-20, US-23 |
| DEV-1…5 | US-01, US-02, US-13, US-14 |
| DQ-1…10 | US-20, US-22, US-23, US-24, US-25, US-28, US-30, US-31 |
| NFR | US-01 (NFR-5), US-05 (NFR-7), US-11 (NFR-3, NFR-4), US-17 (NFR-9), US-33 (NFR-6) |

---

## Developer

### US-01 Get an API key (Phase 1)
As a developer, I want to sign up and create API keys myself, so that I can start calling the API without waiting for anyone.
→ DEV-1, DEV-3, NFR-5
- Given I have signed up, when I create a key, then the key is shown once and stored only as a hash.
- Given I have a key, when I revoke it, then calls made with it return `401` from then on.
- Calls without a key return `401` with the standard error body and a link to the docs.

### US-02 Learn and try the API (Phase 1)
As a developer, I want interactive docs, a live "Try it" search and copy-ready examples, so that I can make my first call in minutes.
→ DEV-2
- The API reference is generated from `contracts/openapi.yaml` and lists every endpoint, parameter and response.
- "Try it" calls the live API and shows the raw JSON beside the formatted result.
- Examples exist for JavaScript, Go, Python and curl, and they run unchanged.

### US-03 Look up a Zambian word (Phase 1)
As a developer, I want to look up a Bemba or Nyanja word and get its meanings, part of speech and English equivalents, so that I can show correct local text.
→ LEX-1, LEX-3, LEX-6, PRV-1
- Given `amenshi` in Bemba, the result includes headword, part of speech, noun class, plural, root, variety, register and English equivalents.
- Every result carries its trust fields.

### US-04 Find equivalents for an English word (Phase 1)
As a developer, I want to send an English word and get equivalents in one or more Zambian languages, so that I can localise my labels.
→ LEX-2, LEX-7
- `GET /v1/lookup?q=water&from=eng&to=bem` returns matching Bemba entries.
- Results link to a shared `concept_id`, and `GET /v1/concepts/{id}` returns the concept across all languages.

### US-05 Search despite spelling differences (Phase 1)
As a developer, I want search to tolerate spelling variants, missing prefixes and typos, so that real user input still finds the word.
→ LEX-4, NFR-7
- A query that differs only in case or accent marks matches.
- A known spelling variant matches its entry.
- An inflected form matches through the stemmer, and the response says whether the headword, a variant or a root matched.

### US-06 Get verified results only (Phase 1)
As a developer, I want results to be verified by default, and to be able to forbid anything else, so that I never show a guess to my users.
→ TRN-4, PRV-1
- With no `min_status`, only `verified` items are returned.
- With `verified_only=true` and no match, the response is an empty result with a not-found code, never a guess.
- `disputed` items are never returned.

### US-07 Opt in to lower trust levels (Phase 1)
As a developer, I want to ask for `reviewed`, `sourced` or `machine` items explicitly, so that I get coverage when I accept the risk.
→ AI-6, `min_status`
- `machine` items appear only when requested, and they are always labelled `machine`.

### US-08 Download a phrase pack (Phase 1)
As a developer, I want curated packs for greetings, numbers, dates, interface labels, health and money, so that I can localise common text in one call.
→ PHR-1, NFR-3
- `GET /v1/phrases?pack=greetings&lang=bem` returns the pack with trust fields on every phrase.
- Packs are cacheable and usable offline.

### US-09 Detect the language of a text (Phase 1)
As a developer, I want to detect which supported language a text is written in, with a confidence score, so that I can route messages.
→ TRN-1
- `POST /v1/detect` returns a ranked list of languages with confidence scores between 0 and 1.

### US-10 Translate a known phrase (Phase 1)
As a developer, I want short texts translated by exact match against verified phrases, so that common sentences come back correct.
→ TRN-2
- An exact match returns `method: verified_phrase`.
- No match returns nothing, not machine output (that comes in US-31).

### US-11 Pin a dataset release (Phase 2)
As a developer, I want to pin a release and rely on caching, so that my app's text does not change under me.
→ PRV-4, NFR-3, NFR-4
- `release=2027.1` returns data as of that release.
- Responses carry an `ETag` tied to the release; a matching `If-None-Match` returns `304`.
- Breaking changes ship under a new path version, with six months of overlap.

### US-12 Export a phrase pack as a localisation file (Phase 2)
As a developer, I want a phrase pack as JSON, PO or XLIFF, so that I can drop it into my i18n pipeline.
→ PHR-3

### US-13 See my usage and limits (Phase 1 headers, Phase 2 dashboard)
As a developer, I want rate-limit headers and a usage dashboard per key, so that I stay within limits.
→ DEV-3, DEV-5
- Every response carries the limit, the remaining count and the reset time.
- Going over the limit returns `429` with the standard error body.

### US-14 Use a client library (Phase 2)
As a developer, I want official JavaScript and Python clients, so that I avoid writing HTTP code.
→ DEV-4

### US-31 Machine translation with a label (Phase 3)
As a developer, I want free-text translation that tells me how it got the answer, so that I can decide whether to show it.
→ TRN-3, DQ-9
- The response includes `method`, `confidence` and `model`. Verified phrase matches take priority.
- Published quality scores per language are linked from the docs.

### US-32 Phrases with placeholders (Phase 2)
As a developer, I want phrases with placeholders such as a name or an amount, with notes on how the surrounding words change, so that I can build dynamic messages correctly.
→ PHR-2

## Learner

### US-15 Report a wrong entry (Phase 1)
As anyone using the data, I want to report an entry as wrong, so that mistakes get fixed.
→ PRV-3, DQ-8
- `POST /v1/reports` and a portal button both create a report that appears in the reviewers' inbox.
- After a reviewer confirms that the report has merit, the entry becomes `disputed` and is due for resolution within 14 days.

### US-16 Read a full entry (Phase 1)
As a learner, I want an entry's senses, example sentences, variants and sources, so that I understand how the word is used.
→ LEX-1, LEX-3, LEX-5

### US-17 Download the open dataset (Phase 2)
As a researcher, I want to download the verified dataset in open formats with a change log, so that the data outlives the service.
→ NFR-9, PRV-4

### US-33 Hear a pronunciation (Phase 2)
As a learner, I want audio recorded by a native speaker, so that I can say the word correctly.
→ LEX-8, NFR-6
- Every clip has a consent record, and a speaker can withdraw their recordings.

## Contributor

### US-18 Propose content (Phase 1)
As a contributor, I want to propose new entries, senses, examples and corrections, so that the dictionary grows.
→ CON-1
- Proposals enter the review queue with status `sourced` or as a correction linked to the existing item.

### US-19 Get public credit (Phase 1)
As a contributor or reviewer, I want public credit for accepted items, so that volunteering is recognised.
→ CON-5
- A per-language contributors page and leaderboard, plus a named list in every release.

## Reviewer

### US-20 Review items blind on my phone (Phase 1)
As a reviewer, I want one item at a time, without seeing the other reviewer's decision or which AI wrote a candidate, so that my judgement is independent.
→ CON-2, CON-3, DQ-4, AI-4
- I only see items in languages I am approved for.
- I never see items I authored.
- AI candidates are shown side by side, with no provider names.
- The screen works on a phone in short sessions.

### US-21 Quick confirms first (Phase 1)
As a reviewer, I want items where all sources agree queued as quick confirms, so that I can move fast.
→ AI-3
- Full agreement means quick confirm; two agreeing means standard review; otherwise full review.

### US-22 Correct an item and keep its history (Phase 1)
As a reviewer, I want my correction to replace the text while keeping what was there before, so that nothing is lost.
→ PRV-2, NFR-8, DQ-2
- Two independent accepts make an item `verified`.
- Every change records who, when and what changed.

## Senior reviewer

### US-23 Resolve disputes (Phase 1)
As a senior reviewer, I want a view of reviewer disagreements and open user reports, so that I can settle them.
→ CON-4, DQ-8

### US-24 Qualify and monitor reviewers (Phase 1)
As a senior reviewer, I want new reviewers to pass a short test set, and I want to track every reviewer's agreement rate, so that review quality stays high.
→ DQ-3, DQ-7
- Reviewers below the agreement threshold lose approval until they are retrained.

## Admin

### US-25 Keep the source register (Phase 1)
As the admin, I want every source recorded with its rights holder, licence, permission document and allowed uses, so that nothing unlicensed gets in.
→ DQ-1
- Ingestion refuses a source that has no register row.

### US-26 Approve reviewers and see provider quality (Phase 1)
As the admin, I want to approve reviewers per language and see acceptance and correction rates per AI provider and language.
→ AI-5, DQ-7

### US-27 Generate AI candidates (Phase 1)
As the admin, I want to run candidate generation from Claude and Google for the core list, so that reviewers start from a draft.
→ AI-1, AI-2
- Each candidate stores provider, model version, date and status `machine`.

### US-28 Re-check a sample every quarter (Phase 2)
As the admin, I want 5% of verified entries re-checked every quarter, with a failing batch sent back to review, so that errors are caught.
→ DQ-6

### US-29 Build a release (Phase 2)
As the admin, I want to freeze the verified set into a numbered release with counts and a change log.
→ PRV-4

### US-30 Ingest a licensed source (Phase 1)
As the admin, I want to parse a licensed source into candidate entries that go through normalisation and review.
→ DQ-1, LEX-3, LEX-7
