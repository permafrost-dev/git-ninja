# git-ninja

---

A command-line tool that provides advanced git commands for developers.

## Usage

```bash
git-ninja [command] [flags] [args]
```

### Available Commands

- `branch:current` - Work with the current branch or return the current branch name
- `branch:exists` - Check if the given branch name exists
- `branch:recent` - Show recently checked out branch names
- `branch:last` - Show the last checked out branch name
- `branch:search` - Search branch names for matching substrings or a regex pattern

## Development Setup

```bash
go mod tidy
```

### Building the project

`git-ninja` uses the [task](https://task.dev) build tool. To build the project, run the following command:

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
