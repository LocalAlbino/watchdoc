import argparse

_cli_args = None


def cli_args():
    global _cli_args
    if _cli_args is not None:
        return _cli_args

    parser = argparse.ArgumentParser(prog="watchdoc", usage="%(prog)s [options]",
                                     description="Automates creation of file headers for source code.",
                                     epilog="For more usage information, see https://github.com/LocalAlbino/watchdoc")

    parser.add_argument("-v", "--verbose", action="store_true", help="enables verbose logging")
    _cli_args = parser.parse_args()
    return _cli_args


def verbose(func):
    def wrapper(*args, **kwargs):
        if _cli_args.verbose:
            print("[verbose] ", func.__name__, args, kwargs)
        return func(*args, **kwargs)

    return wrapper
