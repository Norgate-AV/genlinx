# Contributing ✨

Thank you for your interest in contributing to genlinx! All contributions are welcome and appreciated. 🙌

If you are proposing a new feature, please [open an issue](https://github.com/Norgate-AV/genlinx/issues/new/choose) first to confirm it aligns with the project goals before writing any code.

## Prerequisites

- [Go](https://go.dev/dl/) ≥ 1.26
- [golangci-lint](https://golangci-lint.run/welcome/install/) v2
- [lefthook](https://github.com/evilmartians/lefthook) (git hooks)
- Git

## Setup :technologist:

1. Fork this repository to your own GitHub account and clone it:

    ```bash
    git clone https://github.com/{USERNAME}/genlinx.git
    cd genlinx
    ```

2. Install the git hooks:

    ```bash
    lefthook install
    ```

3. Build and verify:

    ```bash
    make build
    make test
    ```

4. Create a new branch and start working:

    ```bash
    git checkout -b my-feature
    ```

## Development

| Command      | Description                |
| ------------ | -------------------------- |
| `make build` | Build the binary to `bin/` |
| `make test`  | Run all tests              |
| `make lint`  | Run golangci-lint          |
| `make vet`   | Run go vet                 |
| `make fmt`   | Format all Go files        |

## Conventional Commits 📝

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)

This project enforces the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification. The commit-msg hook will validate your message automatically.

**Format:** `<type>[optional scope]: <description>`

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `build`, `ci`, `revert`

**Rules:**

- Description must be ≤ 72 characters
- Use `!` after the type/scope to indicate a breaking change (e.g. `feat!: ...`)

**Examples:**

```
feat(archive): add glob pattern support
fix(cfg): resolve incorrect path handling on Windows
docs: update installation instructions
feat!: rename outputFileSuffix to outputFile
```

## Submitting a Pull Request 🚀

Before opening a PR, make sure all checks pass:

```bash
make test
make lint
```

The CI pipeline will also run tests and lint checks on Linux and macOS, so if your changes are platform-specific (particularly anything touching the Windows registry via the `studio` command), please note that in the PR description.

## Add yourself as a Contributor 🙋‍♂️

This project follows the [all-contributors](https://allcontributors.org) specification. Contribution types are listed in the [emoji key](https://allcontributors.org/docs/en/emoji-key).

Comment on your Pull Request to ask the @all-contributors bot to add you:

```
@all-contributors please add @<username> for <contributions>
```
