import json
from datetime import datetime
import subprocess
import sys

OUTPUT_TEMPLATE = """package {}

import (
    "runtime"
    "time"
)

const (
    Version = "{}"
    BuildDate = "{}"
    LinesOfCode = "{}"
    NumFiles = "{}"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
"""

cmd = "gocloc --output-type=json".split(" ")
cmd.append(sys.argv[1])

cloc = None
for x in json.loads(subprocess.check_output(cmd).decode().strip())["languages"]:
    if x["name"] == "Go":
        cloc = x
        break

lines_of_code = cloc["code"]
num_files = cloc["files"]
current_date = datetime.now().strftime("%d/%m/%Y at %H:%M:%S")

with open(sys.argv[2], "w") as f:
    f.write(OUTPUT_TEMPLATE.format(sys.argv[3], sys.argv[4], current_date, lines_of_code, num_files))  
