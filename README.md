# boot-gator

boot-gator is a small CLI RSS/Atom aggregator written in Go. It stores feeds, follows and posts in PostgreSQL and provides a compact command-line interface to manage feeds and browse posts.

This README covers prerequisites, installation, configuration, available commands, and basic usage examples to get you started.

## Quick summary

- Language: Go
- Database: PostgreSQL
- Install: `go install github.com/shavits/boot-gator@latest`

Run the tool directly with `go run . <command>` while in the repo, or after installing run the binary `boot-gator <command>`.

## Prerequisites

- Go (module-enabled). The project uses Go modules (see `go.mod`). A recent Go toolchain is recommended (Go 1.20+).
- PostgreSQL database. You need a reachable Postgres instance and a connection URL.

## Install

To build and install the CLI on your machine:

```bash
# from any directory
go install github.com/shavits/boot-gator@latest

# run the installed binary (binary name is `boot-gator`)
boot-gator <command> [args]
```

If you prefer to run in-place while developing:

```bash
go run . <command> [args]
```

## Configuration

The app expects a JSON config file at the current user's home directory named `.gatorconfig.json`. The code constructs the path using `os.UserHomeDir()`.

Config structure:

```json
{
  "db_url": "<postgres connection url>",
  "current_user_name": "<optional current username>"
}
```

Example `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://postgres:password@localhost:5432/bootgator?sslmode=disable",
  "current_user_name": "alice"
}
```

- `db_url` is required. Use a standard lib/pq/pgx connection string.
- `current_user_name` is optional; using `login` or `register` will update it.

The config loader and writer live in `internal/config`.

## Database / Migrations

Schema and migrations live in the `sql/schema` folder. Files use goose-style markers (`-- +goose Up` / `-- +goose Down`) so you can run them with your preferred migration tool or apply the SQL files directly.

Example (tool-specific):

```bash
# with goose (example):
# goose postgres "$DB_URL" up

# or apply SQL files directly with psql:
# psql "$DB_URL" -f sql/schema/001_users.sql
```

If you'd like, you can run Postgres locally during development (docker-compose or similar) — see the "Where to look next" section.

## Commands (handlers)

The CLI registers the following commands (see `main.go`):

- `login <username>` — set the current user for commands that require login.
- `register <username>` — create a new user and set it as current.
- `reset` — reset users (DB helper used by the project).
- `users` — list known users (current user is marked).
- `agg <duration>` — run the aggregator that periodically fetches feeds; `<duration>` is a Go duration like `30s` or `1m`.
- `addfeed <name> <url>` — add a new feed (requires logged-in user). Also automatically creates a follow for the current user.
- `feeds` — list feeds.
- `follow <feed-url>` — follow a feed (requires logged-in user).
- `following` — list the feeds the current user is following.
- `unfollow <feed-url>` — unfollow a feed (requires logged-in user).
- `browse [limit]` — show recent posts from feeds the current user follows (default limit 2). Requires logged-in user.

Commands that require a logged-in user are wrapped with middleware that loads `current_user_name` from the config file and looks up the corresponding user in the database.

## Usage examples

1. Create `~/.gatorconfig.json` and point `db_url` to your Postgres DB.
2. Run migrations (see SQL folder) to create tables.
3. Create/register a user:

```bash
go run . register alice
# or, after installing:
boot-gator register alice
```

4. Add a feed and follow it:

```bash
boot-gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
```

5. Run the aggregator (fetch every minute):

```bash
boot-gator agg 1m
```

6. Browse recent posts:

```bash
boot-gator browse 10
```

## Implementation notes & tips

- Duplicate posts are prevented by database constraints; the scraper intentionally ignores Postgres unique-violation errors (SQLSTATE `23505`) when inserting posts.
- The aggregator parses RSS date strings using a set of common layouts (see `parsePubDate` in `handler_agg.go`). If you encounter an unrecognized format, add the layout there or consider using a permissive parser like `github.com/araddon/dateparse`.
- `GetNextFeedToFetch` chooses feeds with `last_fetched_at` NULL first, then the oldest `last_fetched_at` value.

## Where to look next in the repo

- `main.go` — command registration and program entry.
- `handler_*.go` — command handlers for feeds, follows, posts and user management.
- `rss.go` — RSS fetch + XML parsing and HTML unescape utility.
- `internal/config` — config loading/writing logic (`~/.gatorconfig.json`).
- `sql/schema` and `sql/queries` — DB schema and SQLC query files.

## Development helpers (optional)

If you want a local Postgres for development, a simple `docker-compose.yml` with a Postgres service is a convenient option. Add migration tooling (goose, migrate) to your workflow to apply schema changes reliably.


