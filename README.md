# git-ninja

---

A powerful command-line tool designed to enhance your Git workflow with advanced commands tailored for developers. 
It simplifies complex Git operations, making branch management and navigation more efficient and intuitive.

## Usage

```bash
git-ninja [command] [flags] [args]
```

### Available Commands

- `branch:current` - Work with the current branch or return the current branch name
- `branch:exists` - Check if the given branch name exists
- `branch:recent` - Show recently checked out branches
- `branch:freq` - Show frequently checked out branches
- `branch:last` - Work with the last checked out branch
- `branch:search` - Search branch names for matching substrings or a regex pattern

### Examples

Check if a branch exists:

```bash
git-ninja branch:exists feature/new-ui
```

List recently checked out branches:

```bash
git-ninja branch:recent
```

List frequently checked out branches, limit to 5 results while excluding 'develop' and 'main' from the list:

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

### Git Aliases - Configuration

Add the following aliases to your `.gitconfig` file to use `git-ninja` commands as Git aliases:

```ini
[alias]
    # list recently checked out branches
    lrb = "!f() { git-ninja branch:recent; }; f"
    # list frequently checked out branches
    lfb = "!f() { git-ninja branch:freq; }; f"
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
```

Switch to the last checked out branch:

```bash
git checkout main
git checkout feature/my-feature
git co-last # switch from feature/my-feature to main
```

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
