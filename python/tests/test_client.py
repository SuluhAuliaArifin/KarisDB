import json
import subprocess
import unittest
from unittest.mock import patch

from karisdb.client import CoreClient
from karisdb.errors import KarisDBError


class CoreClientTests(unittest.TestCase):
    @patch("karisdb.client.subprocess.run")
    def test_sends_protocol_request_and_returns_data(self, run_process):
        run_process.return_value = subprocess.CompletedProcess(
            args=["karisdb-core"],
            returncode=0,
            stdout=json.dumps(
                {
                    "protocol_version": 1,
                    "request_id": "unused",
                    "success": True,
                    "data": {"version": "0.1.0"},
                    "error": None,
                }
            ).encode(),
            stderr=b"",
        )

        client = CoreClient("karisdb-core")
        with patch("karisdb.client.uuid.uuid4") as uuid4:
            uuid4.return_value.hex = "test-request"
            run_process.return_value.stdout = json.dumps(
                {
                    "protocol_version": 1,
                    "request_id": "test-request",
                    "success": True,
                    "data": {"version": "0.1.0"},
                    "error": None,
                }
            ).encode()
            data = client.request("version")

        self.assertEqual(data, {"version": "0.1.0"})
        call = run_process.call_args
        self.assertEqual(call.args[0], ["karisdb-core"])
        self.assertEqual(call.kwargs["timeout"], 10.0)
        self.assertNotIn("shell", call.kwargs)
        sent = json.loads(call.kwargs["input"])
        self.assertEqual(sent["command"], "version")
        self.assertEqual(sent["request_id"], "test-request")

    @patch("karisdb.client.subprocess.run", side_effect=subprocess.TimeoutExpired("core", 1))
    def test_maps_timeout_to_safe_error(self, _run_process):
        with self.assertRaisesRegex(KarisDBError, "did not respond") as raised:
            CoreClient("karisdb-core", timeout=1).request("status")
        self.assertEqual(raised.exception.code, "CORE_TIMEOUT")

    @patch("karisdb.client.subprocess.run")
    def test_rejects_mismatched_response_id(self, run_process):
        run_process.return_value = subprocess.CompletedProcess(
            args=["karisdb-core"],
            returncode=0,
            stdout=b'{"protocol_version":1,"request_id":"wrong","success":true,"data":{},"error":null}',
            stderr=b"",
        )
        with patch("karisdb.client.uuid.uuid4") as uuid4:
            uuid4.return_value.hex = "expected"
            with self.assertRaisesRegex(KarisDBError, "mismatched request ID"):
                CoreClient("karisdb-core").request("status")


if __name__ == "__main__":
    unittest.main()