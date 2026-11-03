import json
from typing import Any


def render(command: str, data: dict[str, Any]) -> str:
    if command == "help":
        return "Commands: " + ", ".join(data["commands"])
    if command == "version":
        return f"KarisDB {data['version']}"
    return json.dumps(data, indent=2, sort_keys=True)