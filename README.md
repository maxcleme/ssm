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
