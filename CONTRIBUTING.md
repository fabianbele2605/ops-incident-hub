# Contributing to Ops Incident Hub

First off, thank you for considering contributing to Ops Incident Hub! It's people like you that make this project better.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Setup](#development-setup)
4. [How to Contribute](#how-to-contribute)
5. [Coding Standards](#coding-standards)
6. [Commit Convention](#commit-convention)
7. [Pull Request Process](#pull-request-process)
8. [Testing](#testing)
9. [Documentation](#documentation)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to fabian.bele@example.com.

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker 20.10+
- Docker Compose 2.0+
- Git
- Make (optional, for convenience)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/YOUR_USERNAME/ops-incident-hub.git
cd ops-incident-hub
```

3. Add the upstream repository:

```bash
git remote add upstream https://github.com/fabianbele2605/ops-incident-hub.git
```

## Development Setup

### 1. Install Dependencies

```bash
cd backend
go mod download
```

### 2. Start Development Environment

```bash
# Start PostgreSQL
docker-compose up -d postgres

# Run migrations
go run cmd/api/main.go migrate up

# Start the API
go run cmd/api/main.go
```

### 3. Verify Setup

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","timestamp":"...","checks":{"database":"healthy"}}
```

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the behavior
- **Expected behavior**
- **Actual behavior**
- **Environment** (OS, Go version, Docker version)
- **Logs** if applicable

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear title and description**
- **Use case** - why is this enhancement useful?
- **Proposed solution** - how would you implement it?
- **Alternatives considered**

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:

- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `documentation` - Documentation improvements

## Coding Standards

### Go Style Guide

We follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

**Key points:**

- Use `gofmt` to format your code
- Use `golangci-lint` for linting
- Follow Clean Architecture principles
- Keep functions small and focused
- Write self-documenting code

### Project Structure

```
backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/           # Business entities (no external dependencies)
│   ├── usecase/          # Business logic orchestration
│   ├── infrastructure/   # External implementations (DB, etc.)
│   └── api/              # HTTP layer (handlers, middleware)
├── migrations/           # Database migrations
└── tests/                # Test files
```

**Rules:**

- `domain/` must NOT import from other internal packages
- `usecase/` can only import from `domain/`
- `infrastructure/` implements interfaces from `domain/`
- `api/` orchestrates everything

### Code Quality

**Before submitting:**

```bash
# Format code
gofmt -w .

# Run linter
golangci-lint run

# Run tests
go test ./... -v

# Check coverage
go test ./... -cover
```

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```bash
# Feature
feat(domain): add incident priority field

# Bug fix
fix(api): handle nil pointer in incident handler

# Documentation
docs(readme): update quick start guide

# Test
test(usecase): add tests for assign incident use case

# Refactor
refactor(infrastructure): optimize database connection pool
```

### Scope

Common scopes:
- `domain` - Domain layer
- `usecase` - Use case layer
- `infrastructure` - Infrastructure layer
- `api` - API layer
- `ci` - CI/CD changes
- `docs` - Documentation

## Pull Request Process

### 1. Create a Feature Branch

```bash
# Update your fork
git fetch upstream
git checkout develop
git merge upstream/develop

# Create feature branch
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

- Write clean, readable code
- Follow coding standards
- Add tests for new functionality
- Update documentation if needed

### 3. Commit Your Changes

```bash
# Stage changes
git add .

# Commit with conventional commit message
git commit -m "feat(domain): add incident priority field"
```

### 4. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 5. Open a Pull Request

1. Go to the original repository on GitHub
2. Click "New Pull Request"
3. Select your fork and branch
4. Fill in the PR template:

```markdown
## Summary
Brief description of changes

## Changes
- [ ] Added feature X
- [ ] Fixed bug Y
- [ ] Updated documentation

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Tests pass locally
- [ ] Documentation updated
- [ ] Conventional commit messages used
```

### 6. Code Review

- Address review comments promptly
- Push additional commits to the same branch
- Keep the PR focused (one feature/fix per PR)

### 7. Merge

Once approved:
- Squash commits if requested
- Maintainer will merge the PR
- Delete your feature branch

## Testing

### Running Tests

```bash
# All tests
go test ./... -v

# Specific package
go test ./internal/domain/... -v

# With coverage
go test ./... -cover

# Integration tests (requires PostgreSQL)
go test ./internal/infrastructure/postgres/... -v
```

### Writing Tests

**Unit Tests:**

```go
func TestNewIncident(t *testing.T) {
    // Arrange
    title := "Test Incident"
    description := "Test Description"
    severity := "high"

    // Act
    incident, err := domain.NewIncident(title, description, severity)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, incident)
    assert.Equal(t, title, incident.Title)
}
```

**Integration Tests:**

```go
func TestIncidentRepository_Create(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer teardownTestDB(t, db)

    // Test implementation
    repo := postgres.NewIncidentRepository(db)
    incident := createTestIncident()

    err := repo.Create(context.Background(), incident)
    assert.NoError(t, err)
}
```

### Test Coverage

Aim for:
- **Domain layer:** 90%+ coverage
- **Use case layer:** 85%+ coverage
- **Infrastructure layer:** 75%+ coverage
- **API layer:** 70%+ coverage

## Documentation

### Code Documentation

```go
// NewIncident creates a new incident with validation.
// Returns an error if title is empty or severity is invalid.
func NewIncident(title, description, severity string) (*Incident, error) {
    // Implementation
}
```

### README Updates

If your change affects:
- Installation process
- Configuration
- API endpoints
- Features

Update the README.md accordingly.

### Technical Documentation

For significant changes, update:
- `docs/arquitectura-consolidada.md` - Architecture decisions
- `docs/guia-operacion.md` - Operational procedures
- API documentation (if applicable)

## Development Workflow

### Daily Workflow

```bash
# 1. Update your fork
git fetch upstream
git checkout develop
git merge upstream/develop

# 2. Create feature branch
git checkout -b feature/my-feature

# 3. Make changes, commit, push
git add .
git commit -m "feat: add my feature"
git push origin feature/my-feature

# 4. Open PR on GitHub

# 5. Address review comments

# 6. After merge, cleanup
git checkout develop
git pull upstream develop
git branch -d feature/my-feature
```

### Branch Naming

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation
- `refactor/` - Code refactoring
- `test/` - Test additions

Examples:
- `feature/add-incident-priority`
- `fix/nil-pointer-in-handler`
- `docs/update-deployment-guide`

## Questions?

- **General questions:** Open a GitHub Discussion
- **Bug reports:** Open a GitHub Issue
- **Security issues:** Email fabian.bele@example.com (do NOT open public issue)

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes (for significant contributions)
- README acknowledgments section

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Ops Incident Hub! 🎉
