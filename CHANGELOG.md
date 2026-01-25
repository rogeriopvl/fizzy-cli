# Changelog

## 0.6.0 - 2026-01-25

### Features

#### Reaction Management

- `fizzy reaction list <card_number> <comment_id>` - List all reactions on a comment
- `fizzy reaction add <card_number> <comment_id> <emoji>` - Add a reaction to a comment
- `fizzy reaction remove <card_number> <comment_id> <reaction_id>` - Remove a reaction from a comment

#### Step Management

- `fizzy step create <card_number> --content <text> [--completed]` - Create a new step on a card
- `fizzy step update <card_number> <step_id> [--content <text>] [--completed]` - Update a step
- `fizzy step delete <card_number> <step_id>` - Delete a step from a card

#### Comment Management

- `fizzy comment delete <card_number> <comment_id>` - Delete a comment from a card
- `fizzy comment update <card_number> <comment_id> --body <text>` - Update a comment on a card

### Improvements

- Refactored API client logic into separate files for better organization
- Enhanced card update command to use `Flags().Changed()` for better flag handling
- Comprehensive test coverage for all new commands (168 tests total)

## 0.5.0 - 2026-01-25

### Features

#### Comment Management

- `fizzy comment list <card_number>` - List all comments on a card
- `fizzy comment show <card_number> <comment_id>` - Display a specific comment
- `fizzy comment add <card_number> <body>` - Create a new comment on a card

### Fixes

- Fixed error handling in notification read-all command
- Removed global flag usage in some commands

## 0.4.0 - 2026-01-19

### Features

#### Card Management

- `fizzy card not-now <card_number>` - Move a card to "Not Now" status
- `fizzy card untriage <card_number>` - Send a card back to triage
- `fizzy card watch <card_number>` - Subscribe to card notifications
- `fizzy card unwatch <card_number>` - Unsubscribe from card notifications
- `fizzy card golden <card_number>` - Mark a card as golden
- `fizzy card ungolden <card_number>` - Remove golden status from a card

## 0.3.0 - 2026-01-11

### Features

#### Notification Management

- `fizzy notification list` - List all notifications with optional filtering
- `fizzy notification read <notification_id>` - Mark a notification as read and display it
- `fizzy notification unread <notification_id>` - Mark a notification as unread
- `fizzy notification read-all` - Mark all unread notifications as read

#### Card Management

- `fizzy card assign <card_number> <user_id>` - Assign or unassign a user to/from a card
- `fizzy card triage <card_number> <column_id>` - Move a card from triage into a column

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
