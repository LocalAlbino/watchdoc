import json
import sys

from watchdoc.cli import verbose


@verbose
def load_config():
    print("Loading configuration file.")
    try:
        with open("watchdoc.json", "r") as f:
            config = json.load(f)
            print("Finished loading configuration file.")
            return config
    except FileNotFoundError:
        print("Configuration file not found.")
        print("See https://github.com/LocalAlbino/watchdoc for info on how to set up watchdoc.")
        sys.exit(1)
    except json.decoder.JSONDecodeError as e:
        print("Failed to parse watchdoc.json.")
        print(f"JSON Error: {e}")
        sys.exit(1)
    except Exception as e:
        print(f"Unexpected error: {e}")
        sys.exit(1)
