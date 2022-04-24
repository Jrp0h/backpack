# TODO

What needs to be done before 1.0 release.

## Actions

- ssh-copy (scp)
- Google Drive
- Google Cloud Storage
- Azure Disk Storage

## Config

- Add `rsa` encryption.
- Add name format.

## Bugs

- `restore` command: Can't restore data when data hasn't changed since last backup. Only occurs when `--no-backup` flag is set to `false` or absent.
- `backup` command: Don't write new hash to file until backup has succeded.

## Nice to have

- More consistent error messages.
- `move` action: Verify that the user has premissions to store data at the specified location.
- `s3` action: Be able to test connection.
- `restore` command: Accept `--action` and `--file` flags to avoid interactive prompt.
