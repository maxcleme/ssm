# SSM

A minimalistic CLI tool for running commands with AWS SSM parameter substitution.

Heavily inspired by [1Password CLI's `op run` command](https://developer.1password.com/docs/cli/reference/commands/run).

## Installation

```bash
go install github.com/maxcleme/ssm@latest
```

## Usage

Set environment variables with `ssm://` prefixes:

```bash
export API_KEY="ssm:///myapp/api/key"
```

Run commands with SSM parameter substitution:

```bash
ssm run -- printenv API_KEY
ssm run -- docker run -e API_KEY myapp
```

### Using environment files

Load environment variables from a file using `--env-file`:

```bash
ssm run --env-file .env.prod -- node server.js
```

The file format supports:

- `KEY=value` pairs (one per line)
- SSM references (`KEY=ssm:///path/to/secret`)
- Comments (lines starting with `#`)

Variables from the file are merged with your current environment, with file values taking precedence.
