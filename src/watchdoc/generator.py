from datetime import datetime

from watchdoc.cli import verbose

_formats = {
    "%Y": datetime.now().strftime("%Y"),
    "%m": datetime.now().strftime("%m"),
    "%d": datetime.now().strftime("%d"),
    "%D": datetime.now().strftime("%Y-%m-%d")
}


@verbose
def replace_format_strings(text):
    for key in _formats:
        if key in text:
            text = text.replace(key, _formats[key])

    return text


@verbose
def update_timestamps():
    # Do this once per file to make sure that these are up to date
    global _formats
    _formats["%Y"] = datetime.now().strftime("%Y")
    _formats["%m"] = datetime.now().strftime("%m")
    _formats["%d"] = datetime.now().strftime("%d")
    _formats["%D"] = datetime.now().strftime("%Y-%m-%d")


@verbose
def create_header_from_config(filename, config, tries):
    if tries >= 5:
        print("Failed to create header several times.")
        print(f"Skipping file: {filename}")
        return

    if "comment" not in config or not isinstance(config["comment"], str):
        config["comment"] = ""

    update_timestamps()
    try:
        with open(filename, "w") as f:
            f.seek(0)  # Prevents overwriting any text that may already be in the file at this point
            if "copyright" in config and isinstance(config["copyright"], list):
                for line in config["copyright"]:
                    updated = f"{config["comment"]} {replace_format_strings(line)}\n"
                    f.write(updated)
                f.write(f"{config["comment"]}\n")

            if "fields" in config and isinstance(config["fields"], dict):
                for key, value in config["fields"].items():
                    f.write(f"{config["comment"]} {key}: {replace_format_strings(value)}\n")

            f.write('\n')
    except FileNotFoundError:
        print(f"Unable to open file: {filename}")
        print("Retrying...")
        create_header_from_config(filename, config, tries + 1)
    except PermissionError:
        print(f"Unable to open file: {filename}")
        print("Retrying...")
        create_header_from_config(filename, config, tries + 1)
    except Exception as e:
        print(f"Unexpected error: {e}")
        print(f"Skipping file: {filename}")
        return
