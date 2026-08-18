# Contributing to FTN-AI

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR-USERNAME/FTN-AI.git`
3. Add upstream: `git remote add upstream https://github.com/beparykamrul-dev/FTN-AI.git`
4. Create a feature branch: `git checkout -b feature/amazing-feature`

## Development Setup

```bash
# Install dependencies
go mod download
npm install

# Start development environment
make dev

# Run tests
make test
```

## Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Run `go vet` before committing
- Add tests for new features

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `style:` Code style
- `refactor:` Code refactoring
- `test:` Tests
- `chore:` Chores

### Examples

```
feat(auth): add two-factor authentication
fix(api): resolve database connection leak
docs(readme): update installation instructions
```

## Pull Request Process

1. Update documentation
2. Add/update tests
3. Run `make test` and ensure all tests pass
4. Commit with conventional messages
5. Push to your fork
6. Create Pull Request with:
   - Clear title
   - Description of changes
   - Link to related issues
   - Screenshots if UI changes

## Testing

- Write tests for all new features
- Ensure coverage ≥ 80%
- Run: `make test-coverage`

## Code Review

- Address review comments promptly
- Request changes if needed
- Keep conversations professional

## Reporting Bugs

1. Check if bug is already reported
2. Create issue with:
   - Clear title
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details

## Requesting Features

1. Check existing issues and discussions
2. Create issue with:
   - Clear title
   - Feature description
   - Use cases
   - Proposed solution (optional)

## License

By contributing, you agree to license your work under the MIT License.