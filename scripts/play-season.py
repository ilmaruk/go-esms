import subprocess
import sys


if __name__ == "__main__":
    fixtures = []
    with open("./data/fixtures.txt") as fh:
        fixtures = fh.readlines()

    for fixture in fixtures:
        home, away = fixture.strip().split(" ")
        print(home, away)
        cmd = [
            "go", "run", "./cmd/manager", "play",
            "-d", ".",
            "--home", home.lower(),
            "--away", away.lower(),
        ]
        # print(" ".join(cmd))
        result = subprocess.run(cmd, capture_output=True, text=True)
        print(result.stdout)

        if result.returncode != 0:
            print("Errors:\n", result.stderr)
            sys.exit(2)
