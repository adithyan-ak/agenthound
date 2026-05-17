# Tiny FastAPI app that mocks the three LiteLLM endpoints the v0.2
# fingerprinter + looter use. Avoids the real LiteLLM image's Postgres
# dependency and need for actual provider keys.
FROM python:3.12-slim

WORKDIR /app
RUN pip install --no-cache-dir fastapi==0.115.0 uvicorn==0.32.0 && \
    apt-get update && apt-get install -y --no-install-recommends curl && \
    rm -rf /var/lib/apt/lists/*

COPY app.py .

EXPOSE 4000
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "4000"]
