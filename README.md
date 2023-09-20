<div align="center">
    <img src=".github/grove_logo.png" alt="Grove logo" width="600"/>
    <!-- TODO Rename header -->
    <h1>Transaction HTTP DB</h1>
    <big>Implementation of the service to communicate with the transaction DB</big>
    <div>
    <br/>
        <a href="https://github.com/pokt-foundation/transaction-http-db/pulse"><img src="https://img.shields.io/github/last-commit/pokt-foundation/transaction-http-db.svg"/></a>
        <a href="https://github.com/pokt-foundation/transaction-http-db/pulls"><img src="https://img.shields.io/github/issues-pr/pokt-foundation/transaction-http-db.svg"/></a>
        <a href="https://github.com/pokt-foundation/transaction-http-db/issues"><img src="https://img.shields.io/github/issues-closed/pokt-foundation/transaction-http-db.svg"/></a>
    </div>
</div>
<br/>

  <!-- TODO Update the nelow section with development instructions (leave the pre-commit section in place) -->

# Development

## Pre-Commit Installation

Before starting development work on this repo, `pre-commit` must be installed.

In order to do so, run the command **`make init-pre-commit`** from the repository root.

Once this is done, the following checks will be performed on every commit to the repo and must pass before the commit is allowed:

### 1. Basic checks

- **check-yaml** - Checks YAML files for errors
- **check-merge-conflict** - Ensures there are no merge conflict markers
- **end-of-file-fixer** - Adds a newline to end of files
- **trailing-whitespace** - Trims trailing whitespace
- **no-commit-to-branch** - Ensures commits are not made directly to `main`

### 2. Go-specific checks

- **go-fmt** - Runs `gofmt`
- **go-imports** - Runs `goimports`
- **golangci-lint** - run `golangci-lint run ./...`
- **go-critic** - run `gocritic check ./...`
- **go-build** - run `go build`
- **go-mod-tidy** - run `go mod tidy -v`

### 3. Detect Secrets

Will detect any potential secrets or sensitive information before allowing a commit.

- Test variables that may resemble secrets (random hex strings, etc.) should be prefixed with `test_`
- The inline comment `pragma: allowlist secret` may be added to a line to force acceptance of a false positive

# Components Diagram

```mermaid
flowchart TD
    A(Client) --> |Sends sessions| B[Router]
    A --> |Sends relays| B
    A --> |Sends service records| B
    B --> |Saves sessions| E[Transaction DB]
    B --> |Sends relays| C[Relay batch]
    B --> |Sends service records| D[Service record batch]
    C --> |After X quantity of relays\nor X quantity of time\nor in case of shutdown\nsaves relays| E
    D --> |After X quantity of service records\n or X quantity of time\nor in case of shutdown\nsaves service records| E
```
Explanation of components:

- **Client**: any external service doing http requests to Transaction HTTP DB service.
- **Router**: internal component that consists of an http server routing the requests to its correct destination.
- **Relay batch**: internal component that saves in memory the relays to later be saved in the Transaction DB after some requirements are meant. This is used to not hit the transaction DB on each request.
- **Service records batch**: internal component that saves in memory the service records to later be saved in the Transaction DB after some requirements are meant. This is used to not hit the transaction DB on each request.
- **Transaction DB**: PostgreSQL database for storing sessions, relays and service records. It serves as the primary storage for transaction data.
