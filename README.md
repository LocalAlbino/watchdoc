# Watchdoc

Watchdoc is a simple tool to help automate the creation of headers for source code files written in Python.

### Installation

Watchdoc can be installed from PyPi using pip.

```
pip install watchdoc-cli
```

### Setup

Create a `watchdoc.json` file inside your project root. Watchdoc reads this file to determine
your configurations for each extension. Each extension is its own object and
Inside each extension, there are three main parts that can optionally be configured.

1. `"comment"`: This tells Watchdoc which comment string to use for your language.
2. `"copyright"`: This is a list of strings that Watchdoc will place at the top of the file.
3. `"fields"`: This is an object that contains your defined header information.

Watchdoc also provides a few format specifiers for inserting date information:

- `%Y`: The current year
- `%m`: The current month
- `%d`: The current day
- `%D`: The current date formatted as `YYYY-mm-dd`

An example `watchdoc.json` for a project using Python and JavaScript may look like this.

```json
{
  ".py": {
    "comment": "#",
    "copyright": [
      "Copyright (c) %Y [Name Here]. All rights reserved."
    ],
    "fields": {
      "Author": "John Doe",
      "Created": "%D",
      "Description": ""
    }
  },
  ".js": {
    "comment": "//",
    "copyright": [
      "Copyright (c) %Y [Name Here]. All rights reserved."
    ]
  }
}
```

The following `watchdoc.json` will produce these headers
(with the current date in place of these example values):

When a python file is created:

```python
# Copyright (c) 2025 [Name Here]. All rights reserved.
#
# Author: John Doe
# Created: 2025-1-1
# Description:
```

When a javascript file is created:

```javascript
// Copyright (c) 2025 [Name Here]. All rights reserved.
```

Again, `"comment"`, `"copyright"`, and `"fields"` are optional so feel free to only
add the ones that you need for your project.

### Running

After you've configured your `watchdoc.json` file, just run `watchdoc` in a terminal at the root of your project.
It will detect when a file is created, and if it's part of your config, then the header will be created.

### Future plans

Here are some potential features that I might add:

- Templates that allow for prompting each contributor for their unique information.
- More format specifiers if they seem needed.

If you're interested in implementing one of these features, or something else, feel free to make a PR.
