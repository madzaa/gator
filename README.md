# Gator CLI

Gator is a command-line tool for managing RSS feeds and posts, built in Go. This project uses PostgreSQL as its database backend.

## Prerequisites

- **Go** (version 1.20 or newer recommended)
- **PostgreSQL** (version 13 or newer recommended)

## Installation

You can install the Gator CLI using the following command:

```
go install github.com/madzaa/gator@latest
```

## Configuration

Before running the program, you need to set up a configuration file. By default, Gator looks for a config file at `~/.gatorconfig.json`.

Example `gatorconfig.json`:

```json
{"db_url":"postgres://localhost:5432/gator?sslmode=disable","current_user_name":"user"}                                                                                                                                                   
```

Make sure your PostgreSQL server is running and the database/user exist.

## Database Setup

Run the SQL migrations in `sql/schema/` to set up the database tables:

```
psql -U your_db_user -d gator -f sql/schema/001_users.sql
psql -U your_db_user -d gator -f sql/schema/002_feeds.sql
# ...and so on for each file in sql/schema/
```

## Usage

After installation and configuration, you can run the CLI:

```
gator <command> [flags]
```

Some example commands:

- `gator register` — Register a new user
- `gator add-feed <feed-url>` — Add a new RSS feed to follow
- `gator list-feeds` — List all followed feeds
- `gator fetch` — Fetch new posts from followed feeds
- `gator list-posts` — List posts from your feeds

Run `gator help` to see all available commands and options.

## Development

To build and run locally:

```
go build -o gator
./gator <command>
```

## License

MIT License

