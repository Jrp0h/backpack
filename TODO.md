# TODO

What needs to be done before 1.0 release.

## Actions

- ssh-copy (scp)
- Google Drive
- Google Cloud Storage
- Azure Disk Storage

## Bugs

- `restore` command: Can't restore data when data hasn't changed since last backup. Only occurs when `--no-backup` flag is set to `false` or absent.

## Nice to have

- More consistent error messages.
- `move` action: Verify that the user has premissions to store data at the specified location.
- `s3` action: Be able to test connection.
- `backup` command: If encryption isn't enabled then wait for input unless `--no-encrypt` is set to `true`.
- `restore` command: Accept `--action` and `--file` flags to avoid interactive prompt.
