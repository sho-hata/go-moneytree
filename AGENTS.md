# go-moneytree

go-moneytree is a Go HTTP API Client library for the [Moneytree LINK API](https://docs.link.getmoneytree.com/docs/product-and-tech-overview).

## Codebase Information
### Development Commands
- `make test`: Run tests and measure coverage
- `make lint`: Code inspection with golangci-lint
- `make tools`: Install dependency tools


## Development Rules
- Working code: Ensure that `make test` and `make lint` succeed after completing work.
- User-friendly documentation comments: Write detailed explanations and example code for public functions and methods so users can understand them easily.

## Coding Guidelines
- No global variables: Do not use global variables. Manage state through function arguments and return values.
- Coding rules: Follow Golang coding rules. [Effective Go](https://go.dev/doc/effective_go) is the basic rule.
- Comments for public functions, methods, variables, and struct fields are mandatory: When visibility is public, always write comments.
- Remove duplicate code: After completing your work, check if you have created duplicate code and remove unnecessary code.
- Error handling: Use `errors.Is` and `errors.As` for error interface equality checks. Never omit error handling.
- Documentation comments: Write documentation comments to help users understand how to use the code. In-code comments should explain why or why not something is done.
- CHANGELOG.md maintenance: When updating CHANGELOG.md, always include references to the relevant PR numbers and commit hashes with clickable GitHub links. This helps developers trace which specific changes were made in which PR/commit and allows them to browse the actual code changes. Format examples:
  - **Feature description ([abc1234](https://github.com/sho-hata/go-moneytree/commit/abc1234))**: Detailed explanation of the change
  - **Feature description (PR #123, [abc1234](https://github.com/sho-hata/go-moneytree/commit/abc1234))**: When both PR and commit are relevant
  - Use `git log --oneline` and GitHub PR numbers to identify the specific changes
  - Always format commit hashes as clickable links: `[hash](https://github.com/sho-hata/go-moneytree/commit/hash)`
  - This improves traceability and allows developers to browse code changes directly in their browser
  - Users want to see the actual implementation, so always provide GitHub links

- Readable test code: Avoid excessive optimization (DRY) and aim for a state where it's easy to understand what tests exist.
- Clear input/output: Create tests with `t.Run()` and clarify test case input/output. Test cases clarify test intent by explicitly showing input and expected output.
- Test descriptions: The first argument of `t.Run()` should clearly describe the relationship between input and expected output.
- Test granularity: Aim for 80% or higher coverage with unit tests.
- Parallel test execution: Use `t.Parallel()` to run tests in parallel whenever possible.
- Using `octocov`: Run `octocov` after `make test` to confirm test coverage exceeds 80%.
