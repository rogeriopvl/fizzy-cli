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
