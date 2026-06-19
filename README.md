# pit

A small TUI for writing standup notes without starting from a blank page.

`pit` opens a focused TUI for writing a daily standup entry, saves entries locally as JSON, renders entries as Markdown, and can prefill your "Yesterday" field from Git commits across tracked repositories.
<p align="center">
  <img width="800" alt="pit" src="https://github.com/user-attachments/assets/a1985fe6-c0eb-4bf6-ab70-e4f1c81a310f" />
</p>


## Why

Standups often start with a small panic: what did I actually do yesterday, what did I say I would do next, and is anything blocking me?

`pit` is for the moment before standup, an async update, or a personal work log when you need to answer:

- What did I do?
- What is blocking me?
- What am I doing next?

Instead of starting from a blank note or digging through Git logs, tickets, and yesterday's scratchpad, `pit` keeps a small local trail and gives you a starting point. The goal is to reduce the mental load of remembering the shape of your work, then let you edit the update in your own words.

## Features

- Fast Bubble Tea TUI for daily standup notes
- Three focused fields: `Yesterday`, `Blocked`, and `Today`
- Local JSON persistence
- History browser with filtering
- Markdown preview and clipboard copy
- Git commit prefill from tracked repositories
- Per-repo or global Git author email support
- Previous workday carry-over from yesterday's `Today`
- `--days-back` support for missed standups

## Install

With Go:

```sh
go install github.com/zacaguas/pit@latest
```

Or from a local checkout:

```sh
go build -o pit
./pit
```

## Usage

```sh
pit
```

Look back more than one workday:

```sh
pit --days-back=2
```

This is useful if you missed a day. For example, running `pit --days-back=2` on a Friday will seed from Wednesday's entry and query commits since Wednesday.

## Keybindings

### Today

| Key | Action |
| --- | --- |
| `i`, `enter` | edit focused field |
| `esc` | leave edit mode |
| `j`, `down`, `tab` | next field |
| `k`, `up`, `shift+tab` | previous field |
| `1` | focus Yesterday |
| `2` | focus Blocked |
| `3` | focus Today |
| `s` | save entry |
| `c` | copy current entry as Markdown |
| `v` | preview current entry as Markdown |
| `b` | bulletize focused field |
| `h` | open history |
| `a` | track current repo, when prompted |
| `q`, `esc` | quit |

### History

| Key | Action |
| --- | --- |
| `enter` | open selected entry |
| `/` | filter entries |
| `esc` | clear/cancel filter |
| `q`, `esc` | back to today, when not filtering |

### Detail Preview

| Key | Action |
| --- | --- |
| `j`, `down` | scroll down |
| `k`, `up` | scroll up |
| `c` | copy entry as Markdown |
| `q`, `esc` | back |

## Git Prefill

When there is no saved entry for today, `pit` can prefill `Yesterday` from:

1. The previous workday entry's `Today` field
2. Git commits from configured repositories since the previous workday

Commit lines are appended as Markdown bullets.

If you launch `pit` inside a Git repo that is not tracked yet, the today view shows a prompt. Press `a` to add it to your config. Once tracked, `pit` immediately tries to load commits for that repo.

## Config

`pit` stores config in your user config directory:

```text
pit/config.toml
```

The exact base directory depends on your OS. On macOS this is usually:

```text
~/Library/Application Support/pit/config.toml
```

Example:

```toml
global_email = "you@example.com"

[[repos]]
path = "/Users/you/dev/project-a"

[[repos]]
path = "/Users/you/dev/project-b"
email = "work@example.com"
```

Email precedence for Git queries:

1. Repo-specific `email` in `config.toml`
2. Repo-local `git config user.email`
3. Global `git config --global user.email`

## Data

Entries are stored as JSON files under:

```text
pit/entries/YYYY-MM-DD.json
```

Example entry:

```json
{
  "date": "2026-06-19",
  "did": "Finished git prefill",
  "blocked": "",
  "tomorrow": "Polish the README"
}
```

## Development

Run tests:

```sh
go test ./...
```

Run locally:

```sh
go run .
```

Enable Bubble Tea debug logging:

```sh
DEBUG=1 go run .
```

## Status

`pit` is small, local-only, and intentionally boring. It is built for people who want a quick terminal workflow for daily updates without sending their notes to a service.

## License

MIT
