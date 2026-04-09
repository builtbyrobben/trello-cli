# trello-cli

Command-line interface for Trello. Manage boards, lists, cards, members, and labels from your terminal.

## Installation

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/builtbyrobben/trello-cli/releases).

### Build from Source

```bash
git clone https://github.com/builtbyrobben/trello-cli.git
cd trello-cli
make build
```

## Configuration

trello-cli requires both a Trello API key and an API token.

**Environment variables (recommended for CI/scripts):**

```bash
export TRELLO_API_KEY="your-api-key"
export TRELLO_TOKEN="your-api-token"
```

**Keyring storage (recommended for interactive use):**

```bash
# Store API key in system keyring
trello-cli auth set-key --stdin

# Store API token in system keyring
trello-cli auth set-token --stdin
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `TRELLO_API_KEY` | API key (overrides keyring) |
| `TRELLO_TOKEN` | API token (overrides keyring) |
| `TRELLO_CLI_COLOR` | Color output: `auto`, `always`, `never` |
| `TRELLO_CLI_OUTPUT` | Default output mode: `json`, `plain` |

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output JSON to stdout (best for scripting) |
| `--plain` | Output stable, parseable text (TSV; no colors) |
| `--color` | Color output: `auto`, `always`, `never` |
| `--verbose` | Enable verbose logging |
| `--force` | Skip confirmations for destructive commands |
| `--no-input` | Never prompt; fail instead (useful for CI) |

## Commands

### auth

Manage authentication credentials.

```bash
# Store API key in system keyring
trello-cli auth set-key --stdin

# Store API token in system keyring
trello-cli auth set-token --stdin

# Check authentication status
trello-cli auth status

# Remove all stored credentials
trello-cli auth remove
```

### boards

Manage Trello boards.

```bash
# List all boards
trello-cli boards list

# Get a specific board
trello-cli boards get BOARD_ID

# Output as JSON
trello-cli boards list --json
```

### lists

Manage lists on a board.

```bash
# List all lists on a board
trello-cli lists list --board BOARD_ID

# Create a new list on a board
trello-cli lists create --board BOARD_ID --name "To Do"

# Output as JSON
trello-cli lists list --board BOARD_ID --json
```

### cards

Create, read, update, move, and delete cards.

```bash
# List all cards on a board
trello-cli cards list --board BOARD_ID

# List cards in a specific list
trello-cli cards list --board BOARD_ID --list LIST_ID

# Get a specific card
trello-cli cards get CARD_ID

# Create a new card
trello-cli cards create --list LIST_ID --name "Fix login bug"

# Create with description
trello-cli cards create --list LIST_ID --name "Fix login bug" --desc "Users can't log in with SSO"

# Update a card
trello-cli cards update CARD_ID --name "Updated name"
trello-cli cards update CARD_ID --desc "New description"

# Move a card to another list
trello-cli cards move CARD_ID --list TARGET_LIST_ID

# Delete a card
trello-cli cards delete CARD_ID

# Output as JSON
trello-cli cards list --board BOARD_ID --json
```

### members

View board members and current user info.

```bash
# List members of a board
trello-cli members list --board BOARD_ID

# Show current authenticated member
trello-cli members me

# Output as JSON
trello-cli members me --json
```

### labels

View labels on a board.

```bash
# List labels on a board
trello-cli labels list --board BOARD_ID

# Output as JSON
trello-cli labels list --board BOARD_ID --json
```

### version

Print version information.

```bash
trello-cli version
```

## License

MIT
