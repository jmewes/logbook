# Logbook [![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

This project allows creating a digital engineering logbook with a command-line program.

## Setup

If a computer has the [go command](https://go.dev/dl/) and [Git](https://git-scm.com) installed, the Logbook can be
installed by cloning its Git repository and then running the `go install` command.

```sh
git clone git@github.com:jmewes/logbook.git && cd ./logbook
go install
```

Then the program can be executed with the `logbook` command:

```sh
logbook
```

In the `~/.config/logbook/config.yaml` file it can be configured what directories are used for reading and writing log entries.

The following snippet shows the configuration options with their default values:

```yaml
# The directory where new logbook entries are added.
logDirectory: ~/Logs

# The directory where logbook entries are moved when they are archived.
archiveDirectory: ~/Archive
```

## Usage

### Add logbook entry

```sh
# Add logbook entry
logbook add "${TITLE}"

# Add logbook entry and open its root directory in a text editor
${EDITOR} $(logbook add "${TITLE}")
```

### Search logbook entries

```sh
logbook search "${SEARCH_TERM}"
```

### Archive logbook entries

```sh
# Archive single logbook entry
logbook archive "${PATH}"

# Archive multiple logbook entries
logbook archive $(logbook search --output-format list "${SEARCH_TERM}")
```

### Remove logbook entries

```sh
# Remove single logbook entry
logbook remove "${PATH}"

# Remove multiple logbook entries
logbook remove $(logbook search --output-format list "${SEARCH_TERM}")
```

### Customization

User-specific utilities may be defined with shell features, e.g., these Bash alias and functions on a macOS computer that has VS Code installed:

```sh
alias log=logbook

# Creates logbook entry with title "Scratch Note" and opens it in VS Code.
function note() {
  local LOGBOOK_ENTRY_TITLE="$@"
  if [[ -z "$LOGBOOK_ENTRY_TITLE" ]]; then
    LOGBOOK_ENTRY_TITLE="Scratch Note"
  fi
  LOGBOOK_ENTRY=$(log add "$LOGBOOK_ENTRY_TITLE")
  code "$LOGBOOK_ENTRY"
  code "$LOGBOOK_ENTRY"/*.md
}

# Create logbook entry with architecture decision record
function adr() {
  local SCOPE="$@"
  local LOGBOOK_ENTRY_DIR=$(log add "ADR: $SCOPE")
  cp ~/Vorlagen/adr.md $LOGBOOK_ENTRY_DIR/
  perl -pi -e "s/SCOPE/${SCOPE}/g" $LOGBOOK_ENTRY_DIR/adr.md
  local TODAY=$(date '+%Y-%m-%d')
  perl -pi -e "s/DATE/${TODAY}/g" $LOGBOOK_ENTRY_DIR/adr.md
  code $LOGBOOK_ENTRY_DIR
}
```

– `~/.bash_profile`

## Testing

### Component test

```sh
# Run all tests
go test ./...
```

```sh
# Run all tests with coverage checks
go test ./... -coverprofile=./cov.out
```

With the help of the [gremlins](https://gremlins.dev/) program, the tests can be executed with mutations:

```sh
# Run mutation tests
gremlins unleash
```

### Component integration test

```sh
go run main.go
go run main.go search
go run main.go add "Just a test"
go run main.go archive /path/to/2026/01/11/17.28_wip
go run main.go search -a 
```

## Maintenance

### Static code analysis

https://sonarcloud.io/project/overview?id=jmewes_logbook

## Alternative projects

- [Paper-based engineering logbook](https://github.com/jmewes/logbook/wiki/Paper%E2%80%90based-engineering-logbook)
- [QOwnNote](https://www.qownnotes.org)
- [Emacs OrgMode](https://orgmode.org)
- [Evernote](https://evernote.com)
- [Roam Research](https://roamresearch.com)
- [Quiver](https://yliansoft.com/)
- [Notion](https://www.notion.so/product)
- [Obsidian](https://obsidian.md/)
- [Joplin](https://joplinapp.org/)
- [Zettelkasten](https://zettelkasten.de/)
- [Dendron](https://wiki.dendron.so/)
