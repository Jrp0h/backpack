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

## backpackd

A daemon for automatic backups.

Might not be a separate program but instead the command `backpack deamon`.

When it's run it will automatically run the specified tasks, keep track of previous hashes so
the `hash_path` field in the normal config file isn't needed 
(hashes will be stored at `/var/lib/backpackd/<taskname>.hash`).

### Config File

It will have a config file, default at /etc/backpackd/config.yml

Example:
```yml
actions:
  - move_to_tmp:
      type: move
      dir: /tmp/test

tasks:
  - mywebsite-database: 
      config: /home/user1/mywebsite/backpack.yml
      interval: 0 0 * * 1,5
      only:
        - move_to_tmp
```

interval should accept 'hourly', 'daily', 'weekly' and 'monthly', with optional 'bi'-prefix to make it every other week.

hourly: `0 * * * *` At minute 0.
bihourly: `0 */2 * * *` At minute 0 past every 2nd hour.
daily: `0 0 * * *` At 00:00.
bidaily: `0 0 */2 * *` At 00:00 on every 2nd day-of-month.
weekly: `0 0 * * 0` At 00:00 on Sunday.
biweekly: `0 0 1-7,15-21 * 0` At 00:00 on every day-of-month from 1 through 7 and every day-of-month from 15 through 21 and on Sunday.
monthly: `0 0 1 * *` At 00:00 on day-of-month 1.
bimonthly: `0 0 1 */2 *` At 00:00 on day-of-month 1 in every 2nd month.
