import random
import subprocess
import sys


def usage():
    print(f"{sys.argv[0]} [path-to-elo-file]")
    sys.exit(1)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        usage()

    with open(sys.argv[1]) as fh:
        lines = fh.readlines()

    teams = []
    for line in lines:
        name, code, elo = line.strip().split(",")
        teams.append(code)

    random.shuffle(teams)

    cmd = [
        "go", "run", "./cmd/manager", "fixtures", "create",
        "-d", ".",
        "-t", ",".join(teams)
    ]

    # print(" ".join(cmd))
    result = subprocess.run(cmd, capture_output=True, text=True)
    print("Output:\n", result.stdout)

    if result.returncode != 0:
        print("Errors:\n", result.stderr)
        sys.exit(2)
