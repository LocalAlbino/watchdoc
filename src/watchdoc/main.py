import time

from watchdoc.cli import cli_args
from watchdoc.config import load_config
from watchdoc.filewatcher import create_configured_observer


def main():
    cli_args()
    config = load_config()
    observer = create_configured_observer(config)
    observer.start()

    print("Watching for new files...")
    print("Press Ctrl-C to stop.")
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
    observer.join()


if __name__ == "__main__":
    main()
