# Changelog

## [0.2.1] - 2025-12-20

### Fixes

- NPM package publishing script

## [0.2.0] - 2025-12-20

### Features

#### Card Management

- `fizzy card update <card_number>` - Update card properties (title, description, status, tags)
- `fizzy card delete <card_number>` - Delete a card permanently
- `fizzy card close <card_number>` - Close an existing card (already existed, now documented)
- `fizzy card reopen <card_number>` - Reopen a closed card (already existed, now documented)

#### Account Management

- `fizzy account list` - List all accounts you have access to

#### Improvements

- `fizzy board` without arguments now displays the currently selected board
- Added `--version` flag to display CLI version
- Fixed HTTP client leak in API requests
- Updated API specification to latest version

## [0.1.0] - 2025-12-16

### Initial Release

The first stable release of Fizzy CLI with core functionality for managing
boards, cards, and columns.

### Features

#### Authentication

- `fizzy login` - Authenticate with Fizzy API using access tokens

#### Board Management

- `fizzy board list` - List all boards
- `fizzy board create` - Create a new board

#### Card Management

- `fizzy card list` - List all cards in the selected board
- `fizzy card create` - Create a new card
- `fizzy card show <card_id>` - Display details for a specific card

#### Column Management

- `fizzy column list` - List all columns in the selected board
- `fizzy column create` - Create a new column

#### Configuration

- `fizzy use --board <name>` - Set the active board for subsequent commands
- `fizzy use --account <slug>` - Set the active account for subsequent commands

### Distribution

- Multi-platform support (macOS, Linux, Windows)
- Multi-architecture binaries (x64, arm64)
- Distributed via NPM with automatic binary download on install
