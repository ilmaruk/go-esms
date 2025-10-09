import math
import subprocess
import sys


def usage():
    print(f"{sys.argv[0]} [path-to-elo-file] [min-skill=14]")
    sys.exit(1)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        usage()

    with open(sys.argv[1]) as fh:
        lines = fh.readlines()

    min_skill = 14
    if len(sys.argv) > 2:
        min_skill = int(sys.argv[2])

    teams = []
    min_elo = 1_000_000
    for line in lines:
        name, code, elo = line.strip().split(",")
        elo = int(elo)
        teams.append((name, code, elo))
        if elo < min_elo:
            min_elo = elo

    for i, team in enumerate(teams):
        name, code, skill = team
        print(f"[{i+1:2d}/{len(teams)}] {name}...")

        actual_skill = round(min_skill / min_elo * skill)

        cmd = [
            "go", "run", "./cmd/manager", "roster", "create",
            "-c", code, 
            "-n", name,
            "-s", str(actual_skill),
            "-d", "./data",
        ]
        # print(" ".join(cmd))
        result = subprocess.run(cmd, capture_output=True, text=True)

        if result.returncode != 0:
            print("Output:\n", result.stdout)
            print("Errors:\n", result.stderr)
            sys.exit(2)
