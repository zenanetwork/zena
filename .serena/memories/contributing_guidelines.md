# Contributing Guidelines

## Pull Request Requirements

### 1. GitHub Issue Required
- **All PRs must link to an existing GitHub issue**
- PRs without issues **will NOT be reviewed**
- Issue quality requirements:
  - **Reproducibility**: For bugs, include clear reproduction steps
  - **Context**: Provide sufficient background for understanding
  - **Impact**: Explain effects on project/users:
    - Severity level
    - User scope
    - Downstream effects

### 2. Signed Commits (Mandatory)
- All commits must be GPG-signed
- Unsigned commits will be rejected
- Setup guide: https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits

### 3. Conventional Commit Format
- PR titles must use conventional commit format
- Format: `type: subject`
- Valid types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `ci`, `build`
- Reference: https://www.conventionalcommits.org/

### 4. Documentation Changes
- Only **substantial or impactful** documentation changes accepted
- Minor typo or style-only fixes will NOT be accepted

## Development Workflow

### 1. Fork and Branch
```bash
# Fork the repo on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/zena.git
cd zena

# Create feature branch from main
git checkout -b feature/your-feature-name main
```

### 2. Make Changes
- Follow code style and conventions (see `code_style_conventions.md`)
- Write tests for new functionality
- Update documentation if needed

### 3. Test Your Changes
```bash
# Format code
make format

# Run linters
make lint

# Run tests
make test-unit

# Optional: Check coverage
make test-unit-cover
```

### 4. Commit Changes
```bash
# Stage changes
git add .

# Commit with GPG signature and conventional format
git commit -S -m "feat: add new feature description"

# Example commit message formats:
# feat: add EIP-1559 support
# fix: resolve race condition in mempool
# docs: update installation guide
# test: add unit tests for vm module
# refactor: simplify fee calculation logic
```

### 5. Push and Create PR
```bash
# Push to your fork
git push origin feature/your-feature-name

# Create PR on GitHub
# - Link to the related issue
# - Use conventional commit format for PR title
# - Provide detailed description
```

## Code Review Process

### Reviewers Will Check
1. Code quality and adherence to style guide
2. Test coverage for new code
3. Documentation updates
4. Breaking changes assessment
5. Security implications
6. Performance impact

### Common Review Feedback
- Insufficient test coverage
- Missing error handling
- Code style violations
- Unclear variable/function names
- Missing documentation
- Breaking changes without migration guide

## Project Status Notes

### Current Phase: Pre-v1.0
- Version: v0.x releases
- Status: Under audit and testing
- Breaking changes may occur
- Stability features and benchmarking in progress

### Contribution Priorities
1. Bug fixes and stability improvements
2. Test coverage improvements
3. Documentation enhancements (substantial only)
4. Performance optimizations
5. New features (discuss in issue first)

## Getting Help

### Community Channels
- **Discord**: https://discord.com/invite/interchain
- **Telegram**: https://t.me/CosmosOG
- **Slack**: #Cosmos-tech channel
- **Expert Contact**: https://cosmos.network/interest-form

### Issue Guidelines
- Search existing issues before creating new ones
- Use issue templates when available
- Provide complete information:
  - Environment details (OS, Go version)
  - Steps to reproduce
  - Expected vs actual behavior
  - Relevant logs/errors

## Maintainers
- Primary: Cosmos Labs (cosmoslabs.io)
- Sponsored by: Interchain Foundation
- Key Contributors: B-Harvest, Mantra

## License
- Apache 2.0 License
- Fork of evmOS (open-sourced by Tharsis/Interchain Foundation)
