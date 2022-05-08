# TODO

What needs to be done before 1.0 release.

## Actions

- ssh-copy (scp)
- Google Drive
- Google Cloud Storage
- Azure Disk Storage
- backblaze B2

## Bugs

- Can't zip data when using absolute paths. Potential fix: Strip all directories and only keep the deepest one.
- Hash fails if prev_hash file has newline
- Wrong descriptions on flags

## Config

- Add `rsa` encryption.

## Nice to have

- More consistent error messages.
- `move` action: Verify that the user has premissions to store data at the specified location. Allow chown and chmod.
- `s3` action: Be able to test connection.
- Updates on what's going on
