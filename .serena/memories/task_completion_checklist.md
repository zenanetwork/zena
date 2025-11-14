# Task Completion Checklist

## Before Submitting Any Code Change

### 1. Format Code
```bash
make format
```
- Runs gofumpt for Go files
- Runs black + isort for Python files
- Runs shfmt for shell scripts

### 2. Lint Code
```bash
make lint
```
- Runs golangci-lint for Go (15m timeout)
- Runs pylint + flake8 for Python
- Runs solhint for Solidity contracts
- **Fix all linting errors before proceeding**

### 3. Run Tests
```bash
make test-unit
```
- Runs all unit tests (excludes E2E and simulation)
- Timeout: 15 minutes
- Must pass with no failures

### 4. Coverage Check (Optional but Recommended)
```bash
make test-unit-cover
```
- Generates `coverage.txt` with detailed report
- Review coverage for new code
- Aim for high coverage on critical paths

### 5. Protobuf Changes (If Applicable)
```bash
make proto-format
make proto-lint
make proto-gen
```
- Format proto files
- Check proto linting
- Regenerate Go code from protos

### 6. Git Commit Requirements

#### Sign Your Commits
```bash
git commit -S -m "type: description"
```
- **All commits MUST be GPG-signed**
- Unsigned commits will be rejected

#### Conventional Commit Format
```
type: subject line

body (optional)

footer (optional)
```

**Valid types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `ci`, `build`

**Example**:
```
feat: add EIP-1559 fee market support

Implements dynamic base fee calculation based on block gas usage.
Includes unit tests and integration tests.

Closes #123
```

### 7. Link to GitHub Issue
- Every PR must link to an existing GitHub issue
- PRs without issues **will not be reviewed**
- Ensure issue has:
  - Clear problem description
  - Reproduction steps (if bug)
  - Context and impact explanation

### 8. Pre-PR Checklist
- [ ] Code formatted (`make format`)
- [ ] All linters pass (`make lint`)
- [ ] All unit tests pass (`make test-unit`)
- [ ] Commits are GPG-signed
- [ ] Commit messages follow conventional format
- [ ] Linked to GitHub issue
- [ ] Documentation updated (if needed)
- [ ] CHANGELOG.md updated (if needed)

## Quick Command Sequence
```bash
# Full quality check before commit
make format && make lint && make test-unit

# If all pass, commit with signature
git add .
git commit -S -m "type: your message"
```

## Additional Checks for Specific Changes

### Solidity Contract Changes
```bash
make test-solidity
make lint-contracts
```

### Python Script Changes
```bash
make test-scripts
make lint-python
```

### Breaking Changes
```bash
make proto-check-breaking
```
- Check for breaking proto changes against main branch

### Security Audit
```bash
make vulncheck
```
- Check for known vulnerabilities in dependencies

## Darwin (macOS) Specific Notes
- Ensure CGO is enabled for builds
- Fuzz tests use old Apple linker (xcode workaround)
- All standard Unix commands available
