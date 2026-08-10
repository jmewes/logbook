# Logbook [![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

The **Logbook** projects provides a command-line program that supports keeping a digital engineering logbook.
Each logbook entry has a directory with a structure similar to blog posts, with a Markdownfile with an ISO 8601 timestamp that defines the logbook entry.
This allows to keep notes in the logbook entry file and open the logbook entry directory in a text editor to created arbitray further files and directories, depending on the work in progress.
The target group are engineers that have a terminal window open most of the time (e.g., software developers, devops engineers, testers).

| :warning: WARNING          |
|:---------------------------|
| Currently the program gets only tested on macOS. Probably it works on Linux. Windows is currently not supported.      |

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

With the `archive` command, a logbook entry can be moved into an archive directory, to make it disappearch from search results, unless explicitly requested:

```sh
# Archive single logbook entry
logbook archive "${PATH}"

# Archive multiple logbook entries
logbook archive $(logbook search --output-format list "${SEARCH_TERM}")
```

### Remove logbook entries

With the `remove` command, the logbook entry at the given path will be moved in the operating system trash bin:

```sh
# Remove single logbook entry
logbook remove "${PATH}"

# Remove multiple logbook entries
logbook remove $(logbook search --output-format list "${SEARCH_TERM}")
```

**Also see**:

- https://github.com/laurent22/go-trash

### Customization

User-specific utilities may be defined with shell features, e.g., these Bash alias and functions on a macOS computer that has VS Code installed:

```sh
alias log=logbook

# ==============================================================================
# FUNCTIONS SECTION
# ==============================================================================

# NAME
#   note
#
# SYNOPSIS
#   note [<logbook_entry_title>]
#
# DESCRIPTION
#   Creates a logbook entry and opens it in VS Code.
#
# PARAMETERS
#   $@ - logbook_entry_title (String): The title that should be used for the logbook entry (Default: Scratch Note).
# ==============================================================================
note() {
  local logbook_entry_title="$@"
  local logbook_entry
  if [[ -z "$logbook_entry_title" ]]; then
    logbook_entry_title="Scratch Note"
  fi
  logbook_entry=$(log add "$logbook_entry_title")
  code "$logbook_entry"
  code "$logbook_entry"/*.md
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
