# Contributing to go-moneytree

Thank you for your interest in contributing to go-moneytree! This document provides guidelines and instructions for contributing to the project.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with the following information:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Environment details (Go version, OS, etc.)
- Any relevant error messages or logs

### Suggesting Features

Feature suggestions are welcome! Please create an issue with:

- A clear and descriptive title
- Detailed description of the proposed feature
- Use cases and examples
- Any potential implementation considerations

### Submitting Pull Requests

1. **Fork the repository** and create a new branch from `main`
2. **Make your changes** following the coding guidelines below
3. **Write tests** for your changes
4. **Ensure all tests pass** by running `make test`
5. **Run linting** with `make lint` and fix any issues
6. **Update documentation** if necessary
7. **Create a pull request** with a clear description of your changes

## Development Setup

### Prerequisites

- Go 1.21 or later
- `golangci-lint` (install with `make tools`)
- `octocov` (install with `make tools`)

### Development Commands

- `make test`: Run tests and measure coverage
- `make lint`: Code inspection with golangci-lint
- `make tools`: Install dependency tools
- `make coverage`: Run tests and generate coverage report with octocov

## Coding Guidelines

### General Rules

- **Working code**: Ensure that `make test` and `make lint` succeed after completing work
- **No global variables**: Do not use global variables. Manage state through function arguments and return values
- **Follow Go conventions**: Follow Golang coding rules. [Effective Go](https://go.dev/doc/effective_go) is the basic rule

### Documentation

- **Public API comments**: Comments for public functions, methods, variables, and struct fields are mandatory
- **User-friendly documentation**: Write detailed explanations and example code for public functions and methods so users can understand them easily
- **In-code comments**: Write documentation comments to help users understand how to use the code. In-code comments should explain why or why not something is done

### Error Handling

- **Never omit error handling**: Always handle errors appropriately
- **Use `errors.Is` and `errors.As`**: Use these functions for error interface equality checks

### Code Quality

- **Remove duplicate code**: After completing your work, check if you have created duplicate code and remove unnecessary code
- **Keep functions focused**: Each function should have a single, clear responsibility

### Testing

- **Test coverage**: Aim for 80% or higher coverage with unit tests
- **Readable test code**: Avoid excessive optimization (DRY) and aim for a state where it's easy to understand what tests exist
- **Clear test cases**: Create tests with `t.Run()` and clarify test case input/output. Test cases clarify test intent by explicitly showing input and expected output
- **Test descriptions**: The first argument of `t.Run()` should clearly describe the relationship between input and expected output
- **Parallel execution**: Use `t.Parallel()` to run tests in parallel whenever possible

### CHANGELOG.md Maintenance

When updating CHANGELOG.md, always include references to the relevant PR numbers and commit hashes with clickable GitHub links. This helps developers trace which specific changes were made in which PR/commit and allows them to browse the actual code changes.

Format examples:
- **Feature description ([abc1234](https://github.com/sho-hata/go-moneytree/commit/abc1234))**: Detailed explanation of the change
- **Feature description (PR #123, [abc1234](https://github.com/sho-hata/go-moneytree/commit/abc1234))**: When both PR and commit are relevant

Use `git log --oneline` and GitHub PR numbers to identify the specific changes. Always format commit hashes as clickable links: `[hash](https://github.com/sho-hata/go-moneytree/commit/hash)`

## Pull Request Process

1. Ensure your code follows all the guidelines above
2. Make sure all tests pass (`make test`)
3. Ensure linting passes (`make lint`)
4. Verify test coverage is at least 80% (`make coverage`)
5. Update README.md if you've added new APIs or changed existing functionality
6. Update CHANGELOG.md with your changes (if applicable)
7. Create a pull request with a clear title and description

## Code Review

All pull requests will be reviewed. Reviewers may request changes or ask questions. Please be responsive to feedback and work collaboratively to improve the code.

## Questions?

If you have questions about contributing, please open an issue or contact the maintainers.

Thank you for contributing to go-moneytree!

