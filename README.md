# Fizzy CLI

[![Tests](https://github.com/rogeriopvl/fizzy-cli/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/rogeriopvl/fizzy-cli/actions/workflows/tests.yml)
[![Commits per month](https://img.shields.io/github/commit-activity/m/rogeriopvl/fizzy-cli)](https://github.com/rogeriopvl/fizzy-cli/commits)
[![Last commit](https://img.shields.io/github/last-commit/rogeriopvl/fizzy-cli)](https://github.com/rogeriopvl/fizzy-cli/commits)

This is a command-line interface for https://fizzy.do

## Install

You have multiple options to install:

### Homebrew (Mac/Linux)

```bash
brew tap rogeriopvl/tap

brew install fizzy-cli
```

### NPM (Mac/Linux/Windows)

```bash
npm i -g fizzy-cli
```

Or...

```bash
npx -y fizzy-cli@latest
```

## Setup

Before you start using `fizzy-cli`, you need to authenticate it with your Fizzy
account:

```bash
fizzy login
```

The first time you run this command it will print out the instructions for you
to follow. These instructions will guide you through creating an access token on
the Fizzy web app.

After you get the token, setup your shell to export the `FIZZY_ACCESS_TOKEN`
environment variable. Edit your `~/.bashrc` or `~/.zshrc` and add the following
line:

```bash
export FIZZY_ACCESS_TOKEN=your_access_token
```

Reload the shell, and run the `fizzy login` command again to confirm that you're
authenticated. If you have only one Fizzy account, `fizzy-cli` will select it
automatically. Otherwise you will be able to select which one you want to use.

### Board selection

You can choose which board you want to pre-select for all your commands with:

```bash
fizzy use --board <board_name>
```

To get the board name just type:

```bash
fizzy board list
```

The `use` command also supports selecting a different account:

```bash
fizzy use --account <account_slug>
```

## Commands

Top-level commands, grouped by what they do. Run `fizzy <command> --help` (or `fizzy <command> <subcommand> --help`) for subcommands and flags.

**Boards & navigation**

- `fizzy board` — create, list, show, update, delete, publish, and manage access for boards
- `fizzy column` — manage columns and list a column's cards
- `fizzy use` — select the active board or account
- `fizzy whoami` — show the current user and accessible accounts

**Cards**

- `fizzy card` — create, list, show, update, delete, assign, tag, watch, pin, triage, close/reopen, mark golden, postpone, and remove a card's image
- `fizzy comment` — create, list, show, update, delete card comments
- `fizzy step` — manage checklist items on a card
- `fizzy reaction` — manage emoji reactions on comments
- `fizzy tag` — list account tags
- `fizzy pin` — list pinned cards

**Notifications & activity**

- `fizzy notification` — list notifications, mark read/unread, manage notification settings
- `fizzy activity` — view the account-wide activity feed

**Accounts & users**

- `fizzy account` — show settings, manage entropy (auto-postpone), join codes
- `fizzy user` — list, show, update, deactivate users; manage avatars and email-change flow
- `fizzy export` — create and view account or user-data exports

**Integrations & auth**

- `fizzy webhook` — create, list, show, update, delete, activate webhooks and view delivery logs
- `fizzy token` — create personal access tokens
- `fizzy login` / `fizzy logout` — authenticate or destroy the session

## Development

### Tests

```bash
make test
```

### Run

```bash
make install
```

```bash
fizzy --help
```
