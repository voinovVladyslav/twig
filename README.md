# twig

A tiny terminal UI for cleaning up git branches.

`twig` shows every branch in your repo in an interactive list. Move with the
arrow keys (`j`, `k` also supported), select the ones you want gone, hit Enter, and they're deleted — no
need to type out branch names or run `git branch -D` one at a time.

## Why

Branches pile up. After a while `git branch` scrolls off the screen with dozens
of stale feature branches, and deleting them means copy-pasting names into
`git branch -D` over and over. `twig` turns that chore into a checklist: see
them all at once, tick the dead ones, delete in a single pass. Errors from git
(e.g. an unmerged branch) show up inline next to the branch instead of aborting
the whole thing.

## What it shows

- All local and remote branches (`git branch -a`).
- The **current** branch and **remote** branches are color-coded.
- Deletion errors appear inline, next to the branch that failed.

The current branch and `main` are protected — `twig` refuses to delete them.

## Usage

Run it from inside any git repository:

```sh
twig
```

| Key            | Action                          |
| -------------- | ------------------------------- |
| `↑`/`k`, `↓`/`j` | Move the cursor               |
| `Space`        | Select / deselect a branch      |
| `Enter`        | Delete selected branches        |
| `r`            | Refresh the branch list         |
| `q`, `Ctrl+C`  | Quit                            |

## Install

```sh
./install.sh      # builds and installs to /usr/bin/twig
./deinstall.sh    # removes it
```

Or build it yourself:

```sh
go build -o twig main.go
```

## Built with

Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea) (v2) +
[Lip Gloss](https://github.com/charmbracelet/lipgloss). It shells out to the
`git` CLI, so git must be on your `PATH`.
