"""Stub LiteLLM server for the AgentHound v0.2 demo.

Implements the minimum surface the v0.2 fingerprinter and looter
exercise:
  - GET /health/liveliness → returns "I'm alive!" (matches the
    sdk/rules/builtin/fingerprints/litellm.yaml probe).
  - GET /model/info → returns three upstream providers so the looter
    emits 3 upstream Credential nodes.
  - GET /key/list → returns two virtual keys.

Both /model/info and /key/list require Bearer auth matching the
LITELLM_MASTER_KEY env var. This is not a security feature — the demo
key is "sk-DEMO-CHAIN-KEY-NOT-REAL" — but it confirms the looter
correctly forwards the operator-supplied master key.
"""
import os
from fastapi import FastAPI, Header, HTTPException

app = FastAPI()

MASTER_KEY = os.environ.get("LITELLM_MASTER_KEY", "sk-DEMO-CHAIN-KEY-NOT-REAL")


def _check_auth(authorization: str | None) -> None:
    if not authorization or not authorization.startswith("Bearer "):
        raise HTTPException(status_code=401, detail="missing bearer token")
    if authorization.removeprefix("Bearer ") != MASTER_KEY:
        raise HTTPException(status_code=401, detail="invalid master key")


@app.get("/health/liveliness")
def liveliness():
    # Plain text — the fingerprint matcher uses body_contains on
    # "I'm alive!" with HTTP 200.
    from fastapi.responses import PlainTextResponse
    return PlainTextResponse("I'm alive!")


@app.get("/model/info")
def model_info(authorization: str | None = Header(default=None)):
    _check_auth(authorization)
    return {
        "data": [
            {
                "model_name": "gpt-4",
                "litellm_params": {
                    "model": "openai/gpt-4",
                    "api_base": "https://api.openai.com/v1",
                },
                "model_info": {"litellm_provider": "openai"},
            },
            {
                "model_name": "claude-3-opus",
                "litellm_params": {
                    "model": "anthropic/claude-3-opus-20240229",
                    "api_base": "https://api.anthropic.com",
                },
                "model_info": {"litellm_provider": "anthropic"},
            },
            {
                "model_name": "bedrock-claude",
                "litellm_params": {
                    "model": "bedrock/anthropic.claude-v2",
                },
                "model_info": {"litellm_provider": "bedrock"},
            },
        ]
    }


@app.get("/key/list")
def key_list(authorization: str | None = Header(default=None)):
    _check_auth(authorization)
    return {
        "keys": [
            {
                "key_id": "vk-eng-team",
                "spend": 12.34,
                "models": ["gpt-4", "claude-3-opus"],
            },
            {
                "key_id": "vk-data-team",
                "spend": 5.67,
                "models": ["claude-3-opus"],
            },
        ]
    }
