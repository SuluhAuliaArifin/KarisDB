import json
import os
import shutil
import subprocess
import sys
import uuid
from pathlib import Path
from typing import Any

from .errors import KarisDBError

PROTOCOL_VERSION = 1
REQUEST_TIMEOUT_SECONDS = 10.0


def find_core() -> str:
    configured_path = os.environ.get("KARISDB_CORE_PATH")
    if configured_path:
        candidate = Path(configured_path).expanduser()
        if candidate.is_file():
            return str(candidate)
        raise KarisDBError("CORE_NOT_FOUND", "KARISDB_CORE_PATH does not point to a file")

    executable_name = "karisdb-core.exe" if os.name == "nt" else "karisdb-core"
    project_binary = Path(__file__).resolve().parents[2] / "dist" / executable_name
    if project_binary.is_file():
        return str(project_binary)

    sibling_binary = Path(sys.executable).resolve().parent / executable_name
    if sibling_binary.is_file():
        return str(sibling_binary)

    discovered_binary = shutil.which(executable_name)
    if discovered_binary:
        return discovered_binary
    raise KarisDBError("CORE_NOT_FOUND", "karisdb-core was not found; build it or set KARISDB_CORE_PATH")


class CoreClient:
    def __init__(self, core_path: str | None = None, timeout: float = REQUEST_TIMEOUT_SECONDS) -> None:
        self.core_path = core_path or find_core()
        self.timeout = timeout

    def request(self, command: str) -> dict[str, Any]:
        request = {
            "protocol_version": PROTOCOL_VERSION,
            "request_id": uuid.uuid4().hex,
            "command": command,
            "arguments": {},
        }
        encoded_request = json.dumps(request, separators=(",", ":")).encode("utf-8") + b"\n"

        try:
            completed = subprocess.run(
                [self.core_path],
                input=encoded_request,
                capture_output=True,
                timeout=self.timeout,
                check=False,
            )
        except subprocess.TimeoutExpired as exc:
            raise KarisDBError("CORE_TIMEOUT", "karisdb-core did not respond before the timeout") from exc
        except OSError as exc:
            raise KarisDBError("CORE_UNAVAILABLE", "karisdb-core could not be started") from exc

        try:
            response = json.loads(completed.stdout)
        except (UnicodeDecodeError, json.JSONDecodeError) as exc:
            raise KarisDBError("INVALID_CORE_RESPONSE", "karisdb-core returned invalid JSON") from exc

        if not isinstance(response, dict) or response.get("protocol_version") != PROTOCOL_VERSION:
            raise KarisDBError("INVALID_CORE_RESPONSE", "karisdb-core returned an incompatible response")
        if not all(key in response for key in ("request_id", "success", "data", "error")):
            raise KarisDBError("INVALID_CORE_RESPONSE", "karisdb-core response is missing required fields")
        if response.get("request_id") != request["request_id"]:
            raise KarisDBError("INVALID_CORE_RESPONSE", "karisdb-core returned a mismatched request ID")
        if not isinstance(response.get("success"), bool):
            raise KarisDBError("INVALID_CORE_RESPONSE", "karisdb-core response is missing its success flag")

        if not response["success"]:
            error = response.get("error")
            if not isinstance(error, dict) or not isinstance(error.get("code"), str) or not isinstance(error.get("message"), str):
                raise KarisDBError("INVALID_CORE_RESPONSE", "karisdb-core returned an invalid error")
            if response["data"] is not None or completed.returncode != 1:
                raise KarisDBError("INVALID_CORE_RESPONSE", "core exit code does not match its response")
            raise KarisDBError(error["code"], error["message"])

        if completed.returncode != 0 or response["error"] is not None:
            raise KarisDBError("INVALID_CORE_RESPONSE", "core exit code does not match its response")
        data = response["data"]
        if not isinstance(data, dict):
            raise KarisDBError("INVALID_CORE_RESPONSE", "karisdb-core returned invalid command data")
        return data