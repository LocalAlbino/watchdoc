import json
import sys

from watchdoc.cli import verbose


@verbose
def load_config():
    print("loading configuration file.")
    try:
        with open("../../watchdoc.json", "r") as f:
            config = json.load(f)
            print("finished loading configuration file.")
            return config
    except FileNotFoundError:
        print("configuration file not found.")
        print("see https://github.com/LocalAlbino/watchdoc for info on how to set up watchdoc.")
        sys.exit(1)
