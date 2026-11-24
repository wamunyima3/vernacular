# The web app as a Next.js standalone server (P1-OPS-01). Build from the repo root:
#   docker build -f apps/web/Dockerfile -t vernacular-web .
FROM node:22-alpine AS build
RUN corepack enable
WORKDIR /src
COPY pnpm-lock.yaml pnpm-workspace.yaml package.json ./
COPY apps/web/package.json apps/web/
COPY packages packages
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile --filter @vernacular/web... --ignore-scripts
COPY apps/web apps/web
COPY contracts contracts
COPY examples examples
# Only the build sees these; the real values are set at run time. The API is not
# reachable while building, so pages read it per request (lib/settle.ts).
ENV NEXT_TELEMETRY_DISABLED=1 \
    BUILD_WITHOUT_API=1 \
    BETTER_AUTH_SECRET=build-time-placeholder-not-used-at-run-time \
    DATABASE_URL=postgres://build:build@localhost:5432/build
RUN pnpm --filter @vernacular/web exec fumadocs-mdx && pnpm --filter @vernacular/web build

FROM node:22-alpine
RUN adduser -D -H -u 10001 vernacular
WORKDIR /app
COPY --from=build /src/apps/web/.next/standalone ./
COPY --from=build /src/apps/web/.next/static apps/web/.next/static
COPY --from=build /src/apps/web/public apps/web/public
# Read from disk by the docs (API reference, tutorial examples).
COPY contracts contracts
COPY examples examples
WORKDIR /app/apps/web
ENV NODE_ENV=production NEXT_TELEMETRY_DISABLED=1 PORT=3000 HOSTNAME=0.0.0.0
USER vernacular
EXPOSE 3000
CMD ["node", "server.js"]
