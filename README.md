# cursor-config

Personal Cursor configuration collections managed by **curset** CLI.

## Install curset

One-liner install (requires Go):

```bash
curl -sSL https://raw.githubusercontent.com/bilgehannal/cursor-config/main/install.sh | bash
```

### Build from source

```bash
git clone https://github.com/bilgehannal/cursor-config.git
cd cursor-config/curset
make install
```

## Usage

### List available collections

```bash
curset list
```

Shows all available collections from the remote repository, each displayed as a separate table with rounded borders.

### Install a collection

```bash
curset install devops
```

Installs the `devops` collection into the current directory's `.cursor/` folder. Creates `.cursor/rules/`, `.cursor/commands/`, and any other object type directories as needed.

- If a folder or file already exists locally, it is **skipped** (not overwritten) and a message is printed.
- New folders and files are downloaded from the remote repository and written locally.

### List local .cursor contents

```bash
curset local list
```

Scans the current directory's `.cursor/` folder and outputs its contents in `collection.json` format.

## Collections

Collections are defined in [`data/collection.json`](data/collection.json). Each collection maps object types (like `rules` and `commands`) to lists of entries:

- **rules** entries are folders containing `.mdc` rule files
- **commands** entries are individual command files

## Adding new collections

1. Add rule files under `data/.cursor/rules/<name>/`
2. Add command files under `data/.cursor/commands/`
3. Update `data/collection.json` with the new collection definition
4. Push to `main` -- `curset` always fetches the latest from GitHub
