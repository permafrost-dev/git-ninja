# git-ninja

---

<!-- [![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/) -->
![MIT License](https://img.shields.io/badge/oss_license-MIT-blue?style=flat-square&logo=opensourceinitiative&logoColor=white)
![GitHub Release](https://img.shields.io/github/v/release/permafrost-dev/git-ninja?include_prereleases&sort=semver&display_name=tag&style=flat-square&logo=github&color=F9DC3E)
![GitHub Release Date](https://img.shields.io/github/release-date/permafrost-dev/git-ninja?display_date=published_at&style=flat-square&logo=github)
![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability/permafrost-dev/git-ninja?style=flat-square&logo=codeclimate)

A powerful command-line tool designed to enhance your Git workflow with advanced commands tailored for developers. 
It simplifies complex Git operations, making branch management and navigation more efficient and intuitive.

## Screenshots

Recently used branches:
![image](https://github.com/user-attachments/assets/31142cc5-1f1c-4f07-bc9c-2b50f6701b43)

Frequently used branches:
![image](https://github.com/user-attachments/assets/8f3ddea1-24c4-41cc-93db-1f3e938b5dec)

## Available Commands

- `branch:current` - Work with the current branch
- `branch:exists` - Check if the specified branch name exists
- `branch:freq` - List branches frequently checked out
- `branch:last` - Work with the last checked out branch
- `branch:recent` - List branches recently checked out
- `branch:search` - Search branch names for a substring or regex match
- `checkout` - Check out a branch

## Examples

Check out, then `git pull` the `main` branch:

```bash
git-ninja checkout main --pull
git-ninja co main -p
```

List recently checked out branches:

```bash
git-ninja branch:recent
```

List recently checked out branches, limit to 5 results and exclude 'develop' and 'main' from the list:

```bash
git-ninja branch:recent -c 5 -e 'develop|main'
```

Show the last checked out branch name:

```bash
git-ninja branch:last
```

Switch to the last checked out branch:

```bash
git checkout main
git checkout feature/my-feature

# switch from feature/my-feature to main:
git-ninja branch:last --checkout 
```

Search for branches containing "fix":

```bash
git-ninja branch:search fix
```

Search for branches matching a regex pattern (e.g., all branches starting with `GN-12`):

```bash
git-ninja branch:search -r "GN-12.+"
```

Search for a substring in branch names, and check out the first result:

```bash
git checkout main
git-ninja branch:search some-fix -o
# our active branch is now "feature/some-fix" (assuming that was the first result)
```

### Git Aliases - Configuration

Add the following aliases to your `.gitconfig` file to use `git-ninja` commands as Git aliases:

```ini
[alias]
    co = "!f() { git-ninja checkout $@; }; f"
    # list recently checked out branches
    lrb = "!f() { git-ninja branch:recent $@; }; f"
    # list frequently checked out branches
    lfb = "!f() { git-ninja branch:freq $@; }; f"
    # search branches
    sb = "!f() { git-ninja branch:search $@; }; f"
    # push current branch
    pcb = "!f() { git-ninja branch:current --push; }; f"
    # switch to the last checked out branch
    co-last = "!f() { git-ninja branch:last --checkout; }; f"
```

### Git Aliases - Examples

List recently checked out branches:

```bash
git lrb
```

Search branch names:

```bash
# find branches containing "fix"
git sb fix
# find branches matching a regex pattern
git sb -r "fix.+"
# find branches containing "fix" and check out the first result
git sb fix -o
# or checkout a branch by ticket number:
git sb -o 1123 # checks out 'GN-1123-my-feature-branch'
```

Switch to the last checked out branch:

```bash
git checkout main
git checkout feature/my-feature
git co-last # switch from feature/my-feature to main
```

## JIRA Integration
The `branch:recent` command can be run with the `--jira` flag to refine the ordering of the results using live data from JIRA. When your branch names include a JIRA issue key, branches tied to active issues you've updated recently appear closer to the top.

### Configuration

1. **Create an API token**
   - Visit <https://id.atlassian.com/manage-profile/security/api-tokens>.
   - Click **Create API token** and copy the generated value.

   # Recommended: Use a secure credential manager or encrypted .env file
   # Store these credentials in an encrypted .env file, not in your shell profile
   export JIRA_API_TOKEN="your-token"
   export JIRA_SUBDOMAIN="acme"        # for 
   export JIRA_EMAIL_ADDRESS="you@example.com"
   export JIRA_EMAIL_ADDRESS="you@example.com"
   ```

   Add these to your shell profile or `.env` file so `git-ninja` can authenticate with JIRA.

3. **Run `branch:recent` with the `--jira` flag**

   ```bash
   git-ninja branch:recent --jira
   ```

Results are cached for five minutes to avoid repeated API calls. Tickets with higher numbers that have been updated recently are ranked above older or inactive issues.

## Development Setup

```bash
go mod tidy
```

### Building the project

`git-ninja` uses the [task](https://github.com/go-task/task) build tool. To build the project, run the following command:

```bash
task build
```

---

## Changelog

Please see [CHANGELOG](CHANGELOG.md) for more information on what has changed recently.

## Contributing

Please see [CONTRIBUTING](.github/CONTRIBUTING.md) for details.

## Security Vulnerabilities

Please review [our security policy](../../security/policy) on how to report security vulnerabilities.

## Credits

- [Patrick Organ](https://github.com/patinthehat)
- [All Contributors](../../contributors)

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
