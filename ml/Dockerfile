# The translation service (P3-03), for a GPU host apart from the API's server.
# Build and run it as in docs/deploy.md, "Translation service".
FROM python:3.12-slim
COPY --from=ghcr.io/astral-sh/uv:0.9 /uv /bin/uv
ENV UV_COMPILE_BYTECODE=1 UV_LINK_MODE=copy
WORKDIR /app
COPY pyproject.toml uv.lock README.md ./
RUN --mount=type=cache,target=/root/.cache/uv uv sync --locked --no-dev --group serve --no-install-project
COPY vernacular_ml ./vernacular_ml
RUN --mount=type=cache,target=/root/.cache/uv uv sync --locked --no-dev --group serve
# The gated base model is downloaded into /cache once; it must be writable by the service.
RUN useradd --system --uid 10001 --no-create-home vernacular && mkdir /cache && chown vernacular /cache
ENV PATH="/app/.venv/bin:$PATH" HF_HOME=/cache
VOLUME /cache
EXPOSE 8090
USER vernacular
CMD ["uvicorn", "--factory", "vernacular_ml.service:from_env", "--host", "0.0.0.0", "--port", "8090"]
