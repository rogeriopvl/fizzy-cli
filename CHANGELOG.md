# Changelog

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
