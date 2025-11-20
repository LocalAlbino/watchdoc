from watchdoc.cli import verbose


@verbose
def create_header_from_config(filename, config):
    print(filename, config)
