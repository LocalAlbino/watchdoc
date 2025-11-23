import os

from watchdog.events import FileSystemEventHandler
from watchdog.observers import Observer

from watchdoc.cli import verbose
from watchdoc.generator import create_header_from_config


@verbose
def create_configured_observer(config):
    observer = Observer()
    watcher = FileWatcher(config)
    observer.schedule(watcher, '.', recursive=True)
    return observer


class FileWatcher(FileSystemEventHandler):
    def __init__(self, config):
        self.config = config

    @verbose
    def on_created(self, event):
        if event.is_directory:
            return
        extension = os.path.splitext(event.src_path)[1]
        if extension in self.config:
            create_header_from_config(event.src_path, self.config[extension], 0)
        else:
            print(f"Skipping {event.src_path}, extension {extension} not configured.")
