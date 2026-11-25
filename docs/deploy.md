# Deploying on one server

Vernacular runs on a single small server (NFR-10): Postgres, the API, the web app and Caddy
for HTTPS, from `deploy/compose.prod.yml`. A VPS with 2 vCPU, 4 GB of memory and 40 GB of
disk is a reasonable start; the web image build is the most memory-hungry step.

## What runs

| Service    | Image                                    | Notes                                                   |
| ---------- | ---------------------------------------- | ------------------------------------------------------- |
| `postgres` | `postgres:17-alpine`                     | Data in the `pgdata` volume                             |
| `migrate`  | `vernacular-api` (`apps/api/Dockerfile`) | Applies migrations on every start, then exits           |
| `api`      | `vernacular-api`                         | Uploaded permission documents in the `blobs` volume     |
| `worker`   | `vernacular-api`                         | AI candidate jobs; only with `--profile worker`         |
| `web`      | `vernacular-web` (`apps/web/Dockerfile`) | Next.js standalone server                               |
| `caddy`    | `caddy:2-alpine`                         | The only open ports (80, 443); gets certificates itself |

The API image also holds the `ingest` and `admin` tools, the migrations and `data/`.

## First start

1. Point two DNS names at the server, for example `vernacular.example.org` and
   `api.vernacular.example.org`.
2. Install Docker with the compose plugin, clone the repository, and open ports 80 and 443.
3. Copy `deploy/.env.production.example` to `deploy/.env.production` and fill it in. Every
   command below runs from the repository root and uses this shorthand:

   ```bash
   export COMPOSE_FILE=deploy/compose.prod.yml COMPOSE_ENV_FILES=deploy/.env.production
   ```

4. Start everything:

   ```bash
   docker compose up -d --build
   ```

5. Sign up on the site, then make yourself an admin:

   ```bash
   docker compose exec api admin -email you@example.org
   ```

6. Create the portal's key on `/account` (ADR-0011), set it as `WEB_API_KEY`, and run step 4
   again. Until then, signed-out pages that read the API show an error.
7. Import the data, as `make seed` does locally:

   ```bash
   C="docker compose exec api"
   $C ingest -all
   for f in /app/data/concepts/*.csv; do $C ingest -concepts "$f"; done
   $C ingest -phrases /app/data/phrases/en.csv
   ```

## Updating

With the shorthand from step 3 exported:

```bash
git pull
docker compose up -d --build
```

`migrate` runs before the API starts, so a new version never runs against an old schema.
Migrations must stay backward compatible for one release (the old API keeps serving while the
new images build).

## Backups

Dump the database daily and copy the dump off the server:

```bash
docker compose exec -T postgres pg_dump -U vernacular -Fc vernacular > vernacular-$(date +%F).dump
```

Back up the `blobs` volume too: it holds the source permission documents (PRV-3). Restore
with `pg_restore --clean -d vernacular` into a fresh `postgres` service.

## Latency (NFR-1)

Lookup and phrase requests must answer in under 200 ms at the 95th percentile, measured at the
server. The API records every request's duration by route in a Prometheus histogram, served on
port 9090 inside the compose network only (`METRICS_ADDR`). For a quick reading since the API
started:

```bash
docker compose exec api wget -qO- localhost:9090/metrics | scripts/p95.sh
```

```text
route                              requests  p95 (ms) under 200ms
/v1/lookup                               56         9      100.0%
/v1/phrases                              28         5      100.0%
```

With Prometheus scraping `api:9090`, the p95 over five minutes is:

```promql
histogram_quantile(0.95, sum by (le, route) (
  rate(vernacular_http_request_duration_seconds_bucket{method="GET", route=~"/v1/(lookup|phrases.*|entries.*)"}[5m])
))
```

The histogram has a bucket edge at exactly 0.2 s, so the share of requests under the target is
exact: the `le="0.2"` bucket divided by the count.

## Notes

- A CDN in front (NFR-10) is optional at launch. If you add one, proxy both domains through it
  and add its address ranges to the API's `TRUSTED_PROXIES`, or every caller shares the CDN's
  IP for rate limits.
- The images are built without the API reachable, so pages that read it render on the first
  request instead of at build time (`BUILD_WITHOUT_API`, `apps/web/lib/settle.ts`).
- Permission documents are stored in a local folder. To use an S3-compatible store instead,
  set `BLOB_URL` on the `api` service (see `.env.example`).
- Rate limits count per key and per client IP; the API trusts `X-Forwarded-For` only from
  Caddy's network.

## Translation service (optional)

Machine translation (TRN-3) needs the translation service from `ml/` (P3-03), which runs
the fine-tuned model ([ADR 0013](adr/0013-translategemma.md)). It needs a GPU, so it runs
on a separate host rather than the API's server; 4B fits a 24 GB card.

1. Train an adapter with `vernacular-train` (see `ml/README.md`) and copy its directory,
   for example `runs/2026.1`, to the GPU host.
2. Build and run the image there, with a Hugging Face token for an account that has
   accepted the Gemma Terms of Use (the base model is gated; it is cached in `/cache`):

   ```bash
   docker build -t vernacular-mt ml/
   docker run -d --gpus all --restart unless-stopped -e HF_TOKEN -e MODEL_DIR=/model \
     -v "$PWD/runs/2026.1:/model:ro" -v hf-cache:/cache -p 8090:8090 vernacular-mt
   ```

3. Keep port 8090 private to the API server (a firewall rule or a private network), then
   set `TRANSLATION_SERVICE_URL=http://<gpu-host>:8090` in `deploy/.env.production` and
   restart `api`. Leave it empty until the P3-05 benchmark shows the model is good enough.

The service answers `GET /healthz` with the model name it reports to API users.
