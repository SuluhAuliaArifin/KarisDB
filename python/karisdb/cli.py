import argparse
import sys
from collections.abc import Sequence

from .client import CoreClient
from .errors import KarisDBError
from .output import render


def main(argv: Sequence[str] | None = None) -> int:
    parser = argparse.ArgumentParser(prog="karisdb", add_help=False)
    parser.add_argument("-h", "--help", action="store_true", dest="show_help")
    parser.add_argument("command", nargs="?", choices=("help", "version", "status"))
    parsed = parser.parse_args(argv)
    command = "help" if parsed.show_help or parsed.command is None else parsed.command

    try:
        data = CoreClient().request(command)
    except KarisDBError as exc:
        print(f"KarisDB error [{exc.code}]: {exc.message}", file=sys.stderr)
        return 1

    print(render(command, data))
    return 0