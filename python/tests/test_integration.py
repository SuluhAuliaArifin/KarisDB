import json
import os
import subprocess
import sys
import unittest
from pathlib import Path

from karisdb.client import CoreClient


CORE_PATH = os.environ.get("KARISDB_CORE_PATH")
GO_CLI_PATH = os.environ.get("KARISDB_GO_CLI_PATH")
INTEGRATION_READY = bool(
    CORE_PATH
    and GO_CLI_PATH
    and Path(CORE_PATH).is_file()
    and Path(GO_CLI_PATH).is_file()
)


@unittest.skipUnless(INTEGRATION_READY, "build both Go binaries and set their paths to run integration tests")
class CLIIntegrationTests(unittest.TestCase):
    def setUp(self) -> None:
        self.environment = os.environ.copy()
        self.environment["KARISDB_CORE_PATH"] = str(CORE_PATH)

    def run_go_cli(self, command: str) -> str:
        result = subprocess.run(
            [str(GO_CLI_PATH), command],
            capture_output=True,
            text=True,
            timeout=10,
            check=False,
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        return result.stdout.strip()

    def run_python_cli(self, command: str) -> str:
        result = subprocess.run(
            [sys.executable, "-m", "karisdb", command],
            capture_output=True,
            text=True,
            env=self.environment,
            timeout=15,
            check=False,
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        return result.stdout.strip()

    def test_help_and_version_match_between_frontends(self) -> None:
        for command in ("help", "version"):
            with self.subTest(command=command):
                self.assertEqual(self.run_python_cli(command), self.run_go_cli(command))

    def test_status_matches_between_frontends(self) -> None:
        python_status = json.loads(self.run_python_cli("status"))
        go_status = json.loads(self.run_go_cli("status"))
        self.assertEqual(python_status, go_status)

    def test_python_client_calls_go_core_protocol(self) -> None:
        response = CoreClient(str(CORE_PATH)).request("version")
        self.assertEqual(response, {"version": "0.1.0"})


if __name__ == "__main__":
    unittest.main()