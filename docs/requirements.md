> Snapshot of the Claude Doc **[Zambian Languages API: Design and Requirements](https://claude.ai/artifact/WmFptXsPw5GrF5ntVJFTFk)** (rev 26, taken 2026-11-20).
> The live doc is the source of truth. Re-export this file when the doc changes; do not edit requirement text here.

# Zambian Languages API: Design and Requirements

2026-11-15 · Wamunyima Mukelabai

## 1. Summary

We are building a public API that serves verified words, phrases and translations for Zambian languages, starting with Bemba and Nyanja and growing to all seven regional languages. Developers call it to add local language support to their own systems.

The rule that shapes everything else: every item the API returns says where it came from and who checked it. AI translations from Claude and Google are used as a source, but they stay labelled as machine output until two native speakers have accepted them.

Three findings from the research drive the design:

- Google did not get to Bemba and Nyanja by collecting dictionaries. It trained one large model on web text from over 1,000 languages. That gives fluent output but no guarantee that a given word is right, and it only covers two of the seven regional languages in its developer API.
- A lot of Zambian language data already exists, but most of it is speech rather than text, most of it is Bemba, and most of it is licensed for research only. We cannot simply scrape it into a commercial service.
- The gap nobody has filled is a clean, licensed, human-verified lexicon and phrase bank with an API in front of it. That is what we should build first. Machine translation comes later and sits on top of it.

## 2. How Google documents and translates a language

Google does not write down a language word by word. It collects very large amounts of text, trains one model on all languages at once, and lets related languages help each other. The steps below come from Google's own research posts.

1. **Find text in the language.** Google crawls the web and runs a language identification model over it to pick out sentences in each language. For small languages this is hard, because crawled data often holds more noise than usable text. Google trained its own identifier for over 1,000 languages and added filters to clean the result. ([Google Research, 2022](https://research.google/blog/unlocking-zero-resource-machine-translation-to-support-new-languages-in-google-translate/))
2. **Use translated text where it exists.** Big languages have millions of sentence pairs. Small languages often have none. In the 2022 work, the model saw plain text for 1,138 languages and translated pairs for only 112 of them.
3. **Train one model for everything.** The model learns the small language from plain text by filling in hidden words, and learns the task of translating from the big languages. Because one model holds all of them, skills transfer across. This is why being a Bantu language next to Swahili, Zulu and Shona helps.
4. **Make synthetic pairs and filter them.** The model translates plain text to create extra training pairs. Google keeps only pairs that survive a round trip, meaning the sentence still matches after being translated out and back.
5. **Check with native speakers.** Google says it worked with over 100 speakers for the 2022 launch. For the 2024 launch of 110 languages it used its PaLM 2 language model, chose the most widely used variety of each language, and relied on volunteers and linguists for feedback. ([Google, 2024](https://blog.google/products/translate/google-translate-new-languages-2024/))
6. **Add small, high value human data.** Google later released two open datasets built for this purpose. GATITOS is a lexicon of about 4,000 entries translated into 170 small languages. SMOL adds 863 sentences and 584 documents, translated by professionals, under a CC BY 4.0 licence. ([SMOL paper](https://arxiv.org/abs/2502.12301))

### What this means for us

| Google's approach | What we take from it | What we do differently |
|---|---|---|
| Web crawl plus language identification | We need a Zambian language identifier early, to sort mixed text | We ingest from known, licensed sources first and crawl second |
| One model across many languages | Start from an existing multilingual model and fine tune it | We do not train from scratch |
| Synthetic pairs with round trip filtering | Useful for growing the machine translation tier | Synthetic data never enters the verified set |
| A few thousand well chosen human translations go a long way | A core list of about 4,000 concepts per language is the right first target | Ours are reviewed by two native speakers and published with sources |
| Picks one main variety per language | Decide the standard variety per language up front | We also record town varieties such as Town Bemba and Town Nyanja as labelled variants |

One more point matters. Google's developer API lists Bemba (`bem`) and Nyanja (`ny`) and none of Tonga, Lozi, Kaonde, Lunda or Luvale. I checked the [Cloud Translation language list](https://docs.cloud.google.com/translate/docs/languages) for this. The consumer Translate app may differ, so confirm there before we quote it publicly.

## 3. What already exists for Zambian languages

Bemba and Nyanja have real data. Tonga and Lozi have some. Kaonde, Lunda and Luvale have almost nothing digital beyond Bibles and old printed dictionaries. This is a first survey from public sources, not a complete inventory. The "Can we use it" column is my reading of the published licence and needs legal confirmation.

### Speech and parallel text

| Resource | Languages | Size | Licence | Can we use it |
|---|---|---|---|---|
| [Zambezi Voice](https://github.com/unza-speech-lab/zambezi-voice) (University of Zambia) | Bemba, Nyanja, Tonga, Lozi labelled. Lunda audio only | Read speech with transcripts: Bemba 26 h, Nyanja 25 h, Tonga 22 h, Lozi 6 h. Radio audio without transcripts: Bemba 162 h, Tonga 101 h, Lunda 39 h, Lozi 30 h, Nyanja 25 h | [Paper](https://arxiv.org/abs/2306.04428) says CC BY-NC-ND 4.0. The repository shows MIT | Not commercially without permission. Ask UNZA which licence applies |
| [BIG-C](https://github.com/csikasote/bigc) | Bemba with English | 92,117 utterances, about 187 h, each transcribed and translated to English | CC BY-NC-ND 4.0 | Largest Bemba to English parallel set. Research only without a separate agreement |
| [BembaSpeech](https://github.com/csikasote/BembaSpeech) | Bemba | Over 24 h read speech | Check repository | Check |
| [ZamVoice dataset](https://huggingface.co/datasets/buumba641/Zambia-MultiLigual-ASR-Dataset) | Bemba, Nyanja, Tonga | Growing, community collected | Check dataset card | Check. Quality is self rated by users |
| [FLORES-200 / FLORES+](https://huggingface.co/datasets/openlanguagedata/flores_plus) | Bemba, Nyanja, Tumbuka | About 3,000 professionally translated sentences per language | CC BY-SA 4.0 | Yes, with attribution and share alike. Best kept as our test set |
| [Google SMOL](https://huggingface.co/datasets/google/smol) and GATITOS | To confirm per language. Tumbuka had volunteer contributors | 863 sentences, 584 documents, about 4,000 lexicon entries per covered language | CC BY 4.0 | Yes. First task is to list which Zambian languages are in it |
| [GeoPoll voice data](https://geopoll.com/voice-recordings-africa) | Bemba, Nyanja, Tonga, Lozi | Listed hours: Nyanja 13,130, Bemba 12,133, Tonga 733, Lozi 103 | Commercial | Paid. Worth a quote once we reach the speech phase |

### Models

| Model | Zambian languages | Licence | Note |
|---|---|---|---|
| [NLLB-200](https://huggingface.co/facebook/nllb-200-distilled-600M) (Meta) | Bemba, Nyanja | CC BY-NC 4.0 | Good baseline to compare against. We cannot serve it commercially |
| Google Cloud Translation | Bemba, Nyanja | Paid API | Baseline to beat. No Tonga, Lozi, Kaonde, Lunda or Luvale in the [language list](https://docs.cloud.google.com/translate/docs/languages) |

### Dictionaries and word lists

| Resource | Language | Size and form | Rights |
|---|---|---|---|
| [White Fathers' Bemba-English Dictionary](https://scholarblogs.emory.edu/bemba/?p=110) | Bemba | 1954, reprinted 1991, about 1,500 pages in print | Publisher held. Needs permission or a public domain check |
| Hoch, Bemba Pocket Dictionary | Bemba | 1960, both directions, 239 pages | Publisher held |
| [Kitwe Online English-Bemba dictionary](https://kitweonline.com/languages) | Bemba, with smaller Nyanja, Lozi, Tonga, Ngoni and Nsenga lists | Web pages and an ebook of roughly 12,900 words | Private author. Approach for a partnership |
| [Chichewa/Chinyanja Dictionary](https://www.chichewadictionary.org/history) (Steven Paas) | Nyanja | Online since 2010, free to search | Copyright Steven Paas. Approach for a licence |
| [Jalla, Dictionary of the Lozi Language](https://tsammalex.clld.org/sources/jalla1937) | Lozi | 1937. Digitised by the CBOLD project, about 24,000 items | Old enough that it may be out of copyright. Confirm |
| Horton, Luvale dictionary and grammar | Luvale | Print only, as far as I found | Unknown |

### Bibles and church material

Bibles are the largest body of carefully translated text for every one of these languages. [Find.Bible](https://find.bible/countries/ZM) lists 8 Bible resources for Nyanja, 5 for Lunda, 4 for Luvale, 3 for Bemba, 2 for Tonga, and 1 each for Lozi and Kaonde. It also lists smaller languages such as Nsenga, Tumbuka, Lamba, Mambwe-Lungu and Lenje. Texts are readable in apps such as YouVersion and Bible.is, and [the Revised Bemba Bible](https://bible.com/versions/1097) is published by the Bible Society of Zambia.

Because the verses line up one to one with English, a Bible gives about 31,000 aligned verses per language. Two cautions apply. The translations are copyrighted, mostly by the Bible Society of Zambia, so we need a written licence and must not scrape the apps. The language is formal and religious, so it teaches vocabulary well but sounds nothing like a conversation at a market.

### People and organisations

- **University of Zambia speech lab.** Built Zambezi Voice, BIG-C and BembaSpeech. Plans to cover all seven regional languages. The single most important partner.
- **BongoHive, Data for Change and Ubuntu AI.** Published the [Zambia AI Local Languages Ecosystem Report 2025](https://aimappingreport.bongohive.co.zm/), based on 123 respondents and over 43 organisations. I could only read the summary page. Get the full report before finalising this plan.
- **Bible Society of Zambia.** Rights holder for most local Bible texts.
- **Ministry of Education and the Curriculum Development Centre.** Owner of the approved orthography for the seven languages and of school readers in each.

## 4. Gaps and risks

Licensing is the biggest risk, ahead of any technical one. Most of the good data cannot legally sit behind a paid API today.

| Risk | Why it matters | What we do about it |
|---|---|---|
| Research only licences | Zambezi Voice and BIG-C are published as non commercial, no derivatives. NLLB-200 is non commercial | Ask UNZA for a partnership or a separate licence. Until then, use these for internal evaluation only |
| Copyright on dictionaries and Bibles | Zambian copyright lasts for the author's life plus 50 years under the [Copyright and Performance Rights Act](https://wipolex-res.wipo.int/edocs/lexdocs/laws/en/zm/zm070en.pdf). Most dictionaries and all modern Bible translations are still protected | Written permission per source, recorded in the source register. No scraping of apps or sites |
| Almost no text data outside Bemba and Nyanja | The public Tonga and Lozi sets are speech with transcripts. Kaonde, Lunda and Luvale have no public labelled data | Plan to create this data ourselves with volunteer native speakers, recruited through schools, churches and community groups |
| Skew toward religious and formal language | Bible text dominates what exists | Collect everyday domains on purpose: health, money, farming, government services, greetings |
| Spelling is not uniform | People write the same word several ways. An official orthography was approved in 1977 but everyday writing drifts from it | Store one standard spelling per entry plus accepted variants. Normalise on search |
| Town varieties | Town Bemba and Town Nyanja borrow heavily from English and differ from the rural forms | Tag every entry with its variety. Let callers ask for standard, town, or both |
| Bantu word structure | One verb root produces hundreds of surface forms through prefixes and suffixes. A plain word lookup will miss most of them | Store roots with noun class and plural. Build a simple stemmer per language so `/search` finds the root |
| Personal data in voice | Voice recordings are personal data. The [Data Protection Act, 2021](https://zambialii.org/akn/zm/act/2021/3/eng@2021-03-24) requires consent and restricts moving personal data out of Zambia | Written consent from every speaker, registration as a data controller, and a hosting decision made with legal advice |
| Generated filler getting in | Language models write confident Bemba that is partly wrong. If that enters the verified set, the product loses its reason to exist | Hard rule in section 7: nothing generated is ever marked verified without two human reviews |

## 5. Scope

Version 1 is a verified dictionary and phrase API for Bemba and Nyanja, paired with English. Everything else is staged after it.

### Languages

| Wave | Languages (ISO 639-3) | Why this order |
|---|---|---|
| 1 | Bemba (`bem`), Nyanja (`nya`) | Most speakers, most existing data, existing baselines to compare with |
| 2 | Tonga (`toi`), Lozi (`loz`) | Some speech data and at least one large word list exist |
| 3 | Kaonde (`kqn`), Lunda (`lun`), Luvale (`lue`) | Completes the seven regional languages. Needs the most original collection |
| Later | Tumbuka, Nsenga, Lamba, Mambwe-Lungu and others | Added when a community partner and reviewers are in place |

### Who uses it

- **App and web developers** who want their product to greet, label and guide users in a local language.
- **Government and public service teams** building portals, forms, SMS and USSD services.
- **Health, agriculture and finance organisations** that send messages to the public.
- **Schools, learners and researchers** who need a reliable dictionary.
- **Contributors and reviewers**: native speakers, teachers and linguists who add and check entries.

### In scope for version 1

- Word and phrase lookup in both directions with English.
- Ready made phrase packs for common interface text: greetings, numbers, days, months, form labels, error messages.
- Language detection across the supported languages.
- A contributor and reviewer web app.
- Public documentation, API keys and usage limits.

### Out of scope for version 1

- Free text machine translation. It arrives in phase 3 as a clearly labelled beta.
- Speech recognition and text to speech. Phase 4.
- Direct translation between two Zambian languages. Until then it goes through shared concepts, see section 9.

## 6. Functional requirements

Each requirement has an ID so we can track it. "Must" items define version 1. "Should" items follow in the phase shown.

### Lexicon

| ID | Requirement | Priority |
|---|---|---|
| LEX-1 | Look up a word in any supported language and get its meanings, part of speech and English equivalents | Must |
| LEX-2 | Look up an English word and get equivalents in one or more Zambian languages | Must |
| LEX-3 | Each entry carries its root, noun class, plural form and accepted spelling variants where they apply | Must |
| LEX-4 | Search tolerates spelling variants, missing prefixes and common typing errors | Must |
| LEX-5 | Each sense has at least one example sentence with an English translation | Should, phase 1 |
| LEX-6 | Entries carry a variety tag (standard, town, regional dialect) and a register tag (formal, everyday, religious, slang) | Must |
| LEX-7 | Entries link to a shared concept so the same idea can be found across all languages | Must |
| LEX-8 | Entries may carry an audio pronunciation recorded by a native speaker | Should, phase 2 |

### Phrases

| ID | Requirement | Priority |
|---|---|---|
| PHR-1 | Serve curated phrase packs by topic: greetings, numbers, dates, interface labels, health, money | Must |
| PHR-2 | Phrases support placeholders, for example a name or an amount, with notes on how the surrounding words change | Should, phase 2 |
| PHR-3 | Export a phrase pack as a standard localisation file (JSON, PO, XLIFF) | Should, phase 2 |

### Translation and detection

| ID | Requirement | Priority |
|---|---|---|
| TRN-1 | Detect which supported language a piece of text is written in, with a confidence score | Must |
| TRN-2 | Translate a short text by exact match against verified phrases first | Must |
| TRN-3 | When there is no verified match, optionally return a machine translation, marked `machine` with a confidence score | Should, phase 3 |
| TRN-4 | A caller can require verified results only. The API then returns nothing rather than a guess | Must |

### AI assisted candidates

| ID | Requirement | Priority |
|---|---|---|
| AI-1 | Generate candidate translations from Claude and Google Translate for every concept and phrase in the core list | Must |
| AI-2 | Store each candidate with provider, model version, date and status `machine` | Must |
| AI-3 | Compare AI candidates with each other and with licensed sources, and assign a review priority from their agreement | Must |
| AI-4 | Show reviewers all candidates side by side, without saying which provider wrote which | Must |
| AI-5 | Record the acceptance and correction rate per provider and per language | Must |
| AI-6 | Let callers request `machine` items explicitly. They are never returned by default | Should, phase 1 |
| AI-7 | Reuse reviewer corrections as examples in later generation runs and as training data in phase 3 | Should, phase 3 |

### Trust and provenance

| ID | Requirement | Priority |
|---|---|---|
| PRV-1 | Every item returned includes its verification status, source and licence | Must |
| PRV-2 | Every item has a stable ID and a version history | Must |
| PRV-3 | Any user can report an entry as wrong. Reports reach reviewers within the app | Must |
| PRV-4 | The data is published as numbered releases with a change log, so callers can pin a version | Should, phase 2 |

### Contribution and review

| ID | Requirement | Priority |
|---|---|---|
| CON-1 | Registered contributors can propose new entries, senses, examples and corrections | Must |
| CON-2 | Reviewers are approved per language and see a queue of pending items | Must |
| CON-3 | The author of an item cannot review it | Must |
| CON-4 | Disagreements between reviewers go to a senior reviewer for that language | Must |
| CON-5 | Contributors and reviewers are volunteers and are credited publicly for accepted items | Should, phase 1 |

### Developer experience

| ID | Requirement | Priority |
|---|---|---|
| DEV-1 | Self service sign up and API keys | Must |
| DEV-2 | OpenAPI description, interactive docs and copy ready examples | Must |
| DEV-3 | Free access with rate limits at launch. Paid tiers come later | Must |
| DEV-4 | Client libraries for JavaScript and Python | Should, phase 2 |
| DEV-5 | A usage dashboard per key | Should, phase 2 |

## 7. Data quality requirements

An entry is only called verified after two native speakers, working separately, both accept it. This section is what keeps the service from filling up with plausible but wrong content.

### Verification levels

| Level | Meaning | Served by default |
|---|---|---|
| `verified` | Accepted by two independent approved reviewers for that language | Yes |
| `reviewed` | Accepted by one reviewer, waiting for a second | Only if the caller asks |
| `sourced` | Taken from a licensed published source, not yet reviewed by us | Only if the caller asks |
| `machine` | Produced by a model | Only if the caller asks, always labelled |
| `disputed` | Reviewers disagree, or a user report is open | No |

### Rules

| ID | Rule |
|---|---|
| DQ-1 | No source is ingested until it has a row in the source register with its rights holder, licence, permission document and allowed uses |
| DQ-2 | Model generated text can be used to suggest candidates for reviewers. It is stored as `machine` and can only be promoted by two human reviews, the same as any other item |
| DQ-3 | Reviewers are native speakers of the language they review, confirmed by a short test set and a senior reviewer |
| DQ-4 | Reviewers see the item without the other reviewer's decision |
| DQ-5 | Each language has a written style guide: standard orthography, how to record variants, how to handle borrowed English words |
| DQ-6 | A random 5% sample of verified entries is checked again every quarter. If more than 2% fail, the batch they came from goes back to review |
| DQ-7 | Every reviewer has a tracked agreement rate. Reviewers who fall below the threshold lose approval until retrained |
| DQ-8 | A user report moves an entry to `disputed` after a reviewer confirms the report has merit, and it must be resolved within 14 days |
| DQ-9 | Machine translation quality is measured on a held out human test set per language, using chrF and human ratings, and published in the docs |
| DQ-10 | The test set is never used for training or as a source of entries |

The 5%, 2% and 14 day figures are starting proposals. Adjust them once we have real review data.

### AI translations as a source

Translations from Claude and Google Translate are part of the dataset. They are the fastest way to get a first candidate for every concept and phrase, and comparing them tells reviewers where to spend their time. Wamunyima has confirmed permission to use these outputs. File that permission in the source register under each provider, the same as any other source.

How they flow through the system:

1. **Generate.** For every concept and phrase in the core list, request a translation from Claude and from Google. Google's developer API covers Bemba and Nyanja only, so for the other five languages Claude is the only AI source.
2. **Compare.** Line up up to three signals per item: Claude, Google, and any licensed dictionary or corpus that has the item.
3. **Sort by agreement.** All signals agree: the item goes to reviewers as a quick confirm. Two agree: standard review, with the odd one shown as an alternative. No agreement, or one signal only: full review, where the reviewer writes or corrects the answer.
4. **Store.** Each candidate is saved with status `machine`, the provider, the model version and the date. Callers who ask for it can receive `machine` items, always labelled.
5. **Review.** The same blind double review as every other item. Two accepts make it `verified`. A correction replaces the AI text and keeps a record of what it said.
6. **Learn.** Reviewer corrections feed back as examples for later generation runs and as training pairs for our own model in phase 3. We track the acceptance rate per provider and per language and publish it.

Expect the AI output to be useful for Bemba and Nyanja and much weaker for Kaonde, Lunda and Luvale, where the models have seen very little text. For those languages treat it as a hint for reviewers and measure the acceptance rate on a small batch before generating at scale.

The acceptance rate is also a product result in its own right. It tells developers, with numbers, how far they can trust Claude or Google for each Zambian language.

### What we will not do

- Scrape apps or websites that have not given permission.
- Mark an entry as verified when its only support is a language model.
- Hide uncertainty. If we do not have a verified answer, the response says so.

## 8. Non-functional requirements

These are proposed targets for version 1. They are sized for a small team and a service that is mostly read traffic.

| ID | Area | Requirement |
|---|---|---|
| NFR-1 | Speed | Lookup and phrase requests answer in under 200 ms at the 95th percentile, measured at the server |
| NFR-2 | Availability | 99.5% monthly uptime for read endpoints |
| NFR-3 | Low bandwidth | Responses are compressed and cacheable. Phrase packs can be downloaded once and used offline |
| NFR-4 | Versioning | The API path carries a version (`/v1`). Breaking changes get a new version and six months of overlap |
| NFR-5 | Security | Keys are hashed at rest. All traffic uses TLS. Rate limits apply per key and per IP |
| NFR-6 | Privacy | Contributor and speaker data is handled under the Data Protection Act, 2021. Consent is recorded. People can withdraw their recordings |
| NFR-7 | Text handling | UTF-8 everywhere with Unicode NFC normalisation. Search ignores case and accent marks unless asked otherwise |
| NFR-8 | Auditability | Every change to an entry records who, when, and what changed |
| NFR-9 | Portability | The full verified dataset can be exported in open formats, so the data outlives the service |
| NFR-10 | Cost | Version 1 runs on a single modest database server and one application service, with a CDN in front |

## 9. System design

The system is a data pipeline with a review step in the middle and a read heavy API at the end. Start it as one application with clear internal modules, and split out services only when load or team size demands it.

Data enters from licensed sources or from contributors, and both paths end in the same review queue. Only the verified store feeds lookups. The translation model is a separate, labelled path into the API.

### Components

| Component | What it does | Suggested technology |
|---|---|---|
| Source register | Records every source, its rights holder, licence and permission document | Tables in the main database, with the documents in object storage |
| Ingestion jobs | Parse each source into candidate entries. One small parser per source | Go commands run as scheduled jobs |
| Normaliser | Applies the orthography rules, finds roots, merges duplicates, links to concepts | Go package shared by ingestion and search |
| Review app | Queues, blind double review, dispute handling, reviewer scores | Part of the Next.js app, on the Go backend |
| Core database | Entries, senses, translations, examples, reviews, audit log | PostgreSQL with full text and trigram search |
| Public API | Read endpoints, keys, rate limits, caching | One Go HTTP service behind a CDN |
| Release builder | Freezes the verified set into a numbered release and export files | Scheduled job |
| Translation service (phase 3) | Serves the fine tuned model, labelled as machine | Separate service so it can scale and fail on its own |
| Audio store (phase 2 onward) | Pronunciation clips and consent records | Object storage |

PostgreSQL alone is enough for search at this size. A dedicated search engine only becomes worth running past a few hundred thousand entries.

### Technology choices

The backend is written in Go and the web app in Next.js. The rule for everything else: use a free, open source library where a maintained one exists, and write our own code only for what is specific to Zambian languages.

The libraries below are candidates to evaluate in phase 0, not final picks. Check each one's licence and recent activity before adopting it.

| Need | Candidate | We write ourselves |
|---|---|---|
| HTTP routing and middleware | Go standard library `net/http`, or chi | Handlers |
| Database access | pgx with sqlc for typed queries | The queries |
| Schema migrations | golang-migrate or goose | The migrations |
| API contract and server stubs | OpenAPI file with oapi-codegen | The OpenAPI file |
| Background jobs | River, which runs on PostgreSQL | Job logic |
| Search | PostgreSQL full text and `pg_trgm` | Ranking rules |
| Rate limiting | An existing Go limiter middleware | Limits per plan |
| Sign in and sessions | Auth.js or Better Auth in the Next.js app | Roles: developer, contributor, reviewer, senior reviewer, admin |
| UI components | shadcn/ui on Tailwind | Screens |
| Docs site and tutorials | Fumadocs or Nextra, written in MDX | The content |
| Interactive API reference | Scalar or Swagger UI, generated from the OpenAPI file | Nothing |
| Object storage | Any S3 compatible store through the AWS SDK for Go | Nothing |
| Unicode normalisation | `golang.org/x/text` | Orthography rules per language |
| Language detection | No ready library covers these languages well. Evaluate existing language ID models first | Likely our own small model, trained on the verified data |
| Stemmer per language | None exists | Ours. This is the core original work |

Two honest limits on the Go choice. Model training in phase 3 will almost certainly use Python, because that is where the training tools are. Go still serves the result. And the AI candidate jobs call Claude and Google over HTTP, which Go handles without trouble.

### Web app and screens

One Next.js app with two halves: a public developer portal, and signed in tools for contributors and reviewers. Each screen below needs a UI design before it is built.

| Area | Screen | Purpose |
|---|---|---|
| Portal | Home | What the service is, supported languages, live counts of verified entries |
| Portal | Try it | Search box that calls the live API and shows the JSON beside the result |
| Portal | Endpoints | Interactive API reference with every endpoint, parameters and responses |
| Portal | Tutorial | Step by step: get a key, make a first call, add a phrase pack to an app |
| Portal | Docs | Concepts: verification levels, varieties, releases, licences |
| Portal | Examples | Copy ready snippets in JavaScript, Go, Python and curl |
| Portal | Get an API key | Sign up, create and revoke keys, see usage |
| Portal | Data and releases | Download the open dataset, read change logs |
| Portal | Contributors | Public credit page per language |
| Tools | Contribute | Propose a word, sense, example or correction |
| Tools | Review queue | Blind review, one item at a time, built for a phone |
| Tools | Disputes | Senior reviewer view of disagreements and user reports |
| Tools | Sources | Source register with licences and permission documents |
| Tools | Admin | Approve reviewers, see acceptance rates per provider and language |

### Data model

| Entity | Key fields | Notes |
|---|---|---|
| Language | code, name, own name, standard variety | ISO 639-3 codes |
| Concept | id, English gloss, domain | Language neutral meaning. The hub that links all languages |
| Entry | id, language, headword, root, part of speech, noun class, plural, variety, register | One per word or fixed expression |
| Variant | entry, spelling, type | Alternate spellings and town forms |
| Sense | id, entry, concept, definition, usage note | One entry can have several |
| Example | sense, text, English translation | Human written |
| Phrase | id, language, text, concept or pack, placeholders | For interface text and common sentences |
| Audio | entry or phrase, file, speaker, consent record | Phase 2 |
| Source | id, title, rights holder, licence, permission document, allowed uses | DQ-1 depends on this |
| Review | item, reviewer, decision, comment, time | Two accepts make an item verified |
| Release | number, date, counts, change log | What callers pin to |

### Why concepts sit in the middle

Linking every sense to a shared concept means we translate between any two languages without storing every pair. Bemba to Lozi becomes Bemba to concept to Lozi. With seven languages, that is seven sets of links to maintain where direct pairs would need twenty one. Google's GATITOS lexicon is keyed on English words, so it should map onto concepts with little work.

### Handling word structure

Bantu words carry prefixes for noun class, subject, tense and object, so users will search for forms that are not headwords. Each language gets a small rule based stemmer that strips known prefixes and suffixes and proposes candidate roots. The search endpoint looks up the typed form first, then the candidates, and tells the caller which one matched. A full morphological analyser is a later research task.

### Machine translation tier (phase 3)

Fine tune an existing multilingual model on our verified pairs plus any parallel data we hold a licence for. Choose a base model whose licence allows commercial serving, which rules out NLLB-200 as published. Measure against Google Cloud Translation and NLLB-200 on the FLORES test sentences for Bemba and Nyanja before any public launch, and publish the scores.

## 10. API design

The API is REST over HTTPS with JSON, versioned in the path, and authenticated with a key in the `Authorization` header. Every response item carries the same trust fields, so a developer learns them once.

### Endpoints

| Method and path | Purpose | Phase |
|---|---|---|
| `GET /v1/languages` | Supported languages, their codes, varieties and entry counts | 1 |
| `GET /v1/entries/search?q=&lang=` | Find entries by word, with variant and root matching | 1 |
| `GET /v1/entries/{id}` | One entry in full: senses, examples, variants, sources | 1 |
| `GET /v1/concepts/{id}` | One meaning across all languages | 1 |
| `GET /v1/lookup?q=&from=&to=` | Word or phrase equivalents between two languages | 1 |
| `GET /v1/phrases?pack=&lang=` | A curated phrase pack | 1 |
| `POST /v1/detect` | Identify the language of a text | 1 |
| `POST /v1/reports` | Report a problem with an entry | 1 |
| `GET /v1/releases` | Dataset releases and change logs | 2 |
| `GET /v1/phrases/export?pack=&lang=&format=` | Phrase pack as a localisation file | 2 |
| `POST /v1/translate` | Free text translation. Verified match first, then labelled machine output | 3 |

### Common parameters

| Parameter | Values | Default |
|---|---|---|
| `min_status` | `verified`, `reviewed`, `sourced` | `verified` |
| `variety` | `standard`, `town`, `any` | `standard` |
| `release` | A release number, or `latest` | `latest` |

### Example

Request:

```
GET /v1/lookup?q=water&from=eng&to=bem
Authorization: Bearer <key>
```

Response:

```json
{
  "query": "water",
  "from": "eng",
  "to": "bem",
  "release": "2027.1",
  "results": [
    {
      "entry_id": "bem_000412",
      "headword": "amenshi",
      "part_of_speech": "noun",
      "noun_class": 6,
      "concept_id": "c_water",
      "variety": "standard",
      "register": "everyday",
      "status": "verified",
      "reviews": 2,
      "sources": [
        { "id": "src_014", "title": "Example source", "licence": "CC BY 4.0" }
      ]
    }
  ]
}
```

The IDs, release number and source in this example are placeholders that show the shape. The word itself is real.

### Translate response

The translate endpoint always states how it got its answer, so callers can decide whether to show it.

```json
{
  "text": "...",
  "method": "machine",
  "confidence": 0.71,
  "model": "zl-mt-0.3",
  "verified_segments": []
}
```

`method` is `verified_phrase` when the whole input matched a verified phrase, and `machine` otherwise. With `verified_only=true` and no match, the API returns an empty result and a `404` style code, never a guess.

### Errors and limits

- Errors use standard HTTP status codes with a JSON body holding `code`, `message` and a link to the docs.
- Each response carries rate limit headers showing the limit, what remains and when it resets.
- Lookups are cacheable. Responses include an `ETag` tied to the release number.

## 11. Delivery plan

Five phases, each ending in something usable. Durations are estimates for a small team and should be revisited once we know who is on it.

| Phase | Rough length | What gets built | Done when |
|---|---|---|---|
| 0. Foundations | 4 to 6 weeks | Source register. Licence requests to UNZA, Bible Society of Zambia, dictionary authors. Full read of the BongoHive report. Style guides for Bemba and Nyanja. Database schema. UI designs for the portal and review tools. Recruit a first group of reviewers | At least one licensed source per language is signed. Schema and style guides are agreed |
| 1. Lexicon for Bemba and Nyanja | About 3 months | Ingestion and normaliser. Candidate generation from Claude and Google. Review app. Core concept list of about 4,000 meanings. Lookup, search, phrases and detect endpoints. Docs and keys | 4,000 verified concepts in both languages. API open to a closed group of developers |
| 2. Tonga and Lozi, plus polish | About 3 months | Two more languages through the same pipeline. Releases, exports, client libraries, usage dashboard. Pronunciation audio for the core list | Four languages live. Public launch of the dictionary and phrase API |
| 3. Machine translation beta | About 3 months | Fine tuned model for Bemba and Nyanja to and from English. Translate endpoint. Published quality scores against Google and NLLB | Our model matches or beats the baselines on the test set, and native raters accept the output for everyday sentences |
| 4. Kaonde, Lunda, Luvale and speech | Ongoing | Original data collection with volunteer speakers. Speech recognition and text to speech if licences and funding allow | All seven regional languages have a verified core list |

### First ten tasks

- Write to the UNZA speech lab about partnership and licence terms for Zambezi Voice and BIG-C
- Write to the Bible Society of Zambia about text licences for all seven languages
- Contact the Kitwe Online and Chichewa Dictionary authors about licensing their word lists
- Download Google SMOL and GATITOS and list exactly which Zambian languages they cover
- Get the full BongoHive ecosystem report and add what it lists to section 3
- Get legal advice on copyright status of the pre 1960 dictionaries and on Data Protection Act duties
- Obtain the Ministry of Education orthography for each language and draft the Bemba and Nyanja style guides
- Draft the core concept list, starting from the GATITOS word list and adding Zambian everyday terms
- Find two senior reviewers each for Bemba and Nyanja
- Build the schema and a throwaway prototype of the lookup endpoint against a few hundred hand checked entries

## 12. Decisions and open questions

Decisions made by Wamunyima on 15 November 2026, with what each one changes.

| Question | Decision | What it changes |
|---|---|---|
| Ownership | Wamunyima personally | He signs licence agreements and holds consent records, and is the data controller for any personal data |
| Code | Open source under MIT | All code is public from the first commit |
| Data | Open data, with a hosted API that can charge later | The verified set is published under an open licence with each release. Only sources that allow open redistribution can go in |
| Reviewers | Volunteers only | No payment per item. Credit is public. Review speed becomes the main schedule risk |
| Nyanja variety | Zambian Nyanja is the standard | Malawian Chichewa forms are stored as a tagged variety |
| Hosting | Decide when it is time to host | Build so it runs anywhere: containers, PostgreSQL, S3 compatible storage |
| Pricing and running costs | Free for now. Users pay once paid tiers start | Until then the costs fall on the owner, so the first deployment must be cheap |
| Backend | Go | See technology choices in section 9 |
| Web app | Next.js | One app for the developer portal and the review tools |
| Build rule | Use free, open source libraries. Write from scratch only when nothing suitable exists | Each dependency is checked for licence and upkeep before adoption |
| Base model for phase 3 | Decide at phase 3 | Nothing now |

### What open data rules out

An open data licence only works if every source inside it allows redistribution. That has three effects.

- Zambezi Voice and BIG-C, as published, are non commercial and no derivatives. They cannot be merged into an open dataset. They stay as evaluation data unless UNZA grants different terms.
- Bible and dictionary licences must explicitly allow republishing the extracted words and sentences, not only internal use. Ask for that in the first letter.
- The permission for Claude and Google output needs to cover publishing that output openly, since `machine` items will be in the public data.

Being open source and non profit in spirit should make these requests easier to win.

### What volunteer only review means

The target of 4,000 verified concepts in two languages needs 16,000 review decisions, since each item is reviewed twice per language. With volunteers that takes steady recruitment. Plan for it in phase 0.

- Recruit through Byte Builders, UNZA language students, teachers, churches and radio stations.
- Make reviewing possible on a phone in short sessions.
- Show public credit, a leaderboard per language, and a named contributors page in every release.
- The agreement sorted queue from the AI candidates matters more now. Quick confirms keep volunteers moving.

### Still open

1. **Data licence.** MIT is written for software and fits the code. The dataset needs its own licence. CC BY 4.0 is the closest match in spirit: anyone can use it, with credit. Confirm or pick another.
2. **Personal liability.** Holding the agreements and consent records personally works to start. Revisit before collecting voice recordings, since that is when Data Protection Act duties become real.
3. **Hosting.** Deferred by decision. The Data Protection Act point about storing personal data outside Zambia still has to be answered at that time.

## 13. Sources

Pages opened for this document on 15 November 2026.

- [Unlocking Zero-Resource Machine Translation to Support New Languages in Google Translate](https://research.google/blog/unlocking-zero-resource-machine-translation-to-support-new-languages-in-google-translate/), Google Research, 2022
- [110 new languages are coming to Google Translate](https://blog.google/products/translate/google-translate-new-languages-2024/), Google, 2024
- [Cloud Translation supported languages](https://docs.cloud.google.com/translate/docs/languages), Google Cloud
- [SMOL: Professionally translated parallel data for 115 under-represented languages](https://arxiv.org/abs/2502.12301), Caswell et al., 2025
- [Zambezi Voice repository](https://github.com/unza-speech-lab/zambezi-voice) and [paper](https://arxiv.org/abs/2306.04428), University of Zambia, 2023
- [BIG-C repository](https://github.com/csikasote/bigc) and [paper](https://arxiv.org/abs/2305.17202), Sikasote et al., 2023
- [BembaSpeech](https://arxiv.org/abs/2102.04889), Sikasote and Anastasopoulos
- [Zambia Multilingual ASR Dataset](https://huggingface.co/datasets/buumba641/Zambia-MultiLigual-ASR-Dataset), ZamVoice project
- [NLLB-200 model card](https://huggingface.co/facebook/nllb-200-distilled-600M), Meta
- [FLORES+](https://huggingface.co/datasets/openlanguagedata/flores_plus), Open Language Data Initiative
- [Zambia AI Local Languages Ecosystem Report 2025](https://aimappingreport.bongohive.co.zm/), BongoHive, Data for Change and Ubuntu AI (summary page only)
- [Find.Bible, Zambia](https://find.bible/countries/ZM)
- [GeoPoll voice recordings](https://geopoll.com/voice-recordings-africa)
- [Bemba bibliography](https://scholarblogs.emory.edu/bemba/?p=110), Emory University
- [Kitwe Online language resources](https://kitweonline.com/languages)
- [Chichewa Dictionary project history](https://www.chichewadictionary.org/history)
- [Jalla 1937, Lozi dictionary record](https://tsammalex.clld.org/sources/jalla1937)
- [Language data for Zambia](https://translatorswithoutborders.org/language-data-for-zambia), Translators without Borders
- [Data Protection Act, 2021](https://zambialii.org/akn/zm/act/2021/3/eng@2021-03-24), ZambiaLII
- [Copyright and Performance Rights Act](https://wipolex-res.wipo.int/edocs/lexdocs/laws/en/zm/zm070en.pdf), WIPO Lex

Some figures here come from search summaries of these pages and not a full read, notably the dictionary details, the copyright term and the Data Protection Act points. Treat licence and legal statements as leads to confirm, not as advice.
