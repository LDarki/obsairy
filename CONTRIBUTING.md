# Contributing to Obsairy

Thanks for your interest in contributing!

Obsairy aims to stay **fast and maintainable**.  
Please keep changes **focused and minimal**.

---

## Development Setup

### Requirements

* **Bun** ≥ 1.3.6 (JS/TS runtime & package manager)  
* **Go** ≥ 1.25.6 (backend)

### Clone and Build

```bash
# Clone the repo
git clone https://github.com/LDarki/obsairy.git
cd obsairy

# Frontend setup
cd frontend
bun install
bun run dev 

## Backend setup
cd ..
go run cmd/server/main.go
```

## Coding Guidelines 

### JavaScript/TypeScript (React + Vite)

- Prefer simple, readable code over clever abstractions
- Avoid unnecessary dependencies
- Keep performance and scalability in mind
- Document exported components and hooks
- Follow Bun/Prettier formatting:

```bash
bun format
```

### Go (Backend)

- Document exported functions
- Use `go fmt ./...` to keep code style consistent
- Avoid unnecessary dependencies
- Keep performance and scalability in mind
- Prioritize security

## Commit Style

Use clear, descriptive commit messages:

```
feat(chat): optimize websocket connection
fix: prevent race condition in log collector
docs: update readme
```

Small commits are preferred over large, mixed changes.

## Issues

Before opening a new issue:

* Check if it already exists
* Provide a clear description
* Include logs or reproduction steps when relevant

Feature requests should explain the use case.

## Pull Requests

1. Fork the repository
2. Create a feature branch
3. Implement your change
4. Open a pull request

Please keep pull requests focused on a single change.
Large refactors should be discussed in an issue first.

## Philosophy

Obsairy prioritizes:

- Performance
- Security
- Flexibility
- Community

If a change conflicts with these goals, it may not be accepted.

Thanks for contributing!