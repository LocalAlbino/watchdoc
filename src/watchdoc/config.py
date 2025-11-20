import json

from watchdoc.cli import verbose


@verbose
def load_config():
    print("loading configuration file.")
    with open("watchdoc.json", "r") as f:
        config = json.load(f)
        print("finished loading configuration file.")
        return config
