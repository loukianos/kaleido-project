# kaleido-project

A Go API microservice for a loan-note business use case backed by an ERC-721 NFT.
Each token represents a lender's claim to repayment.
- Origination mints the note
- Transfer assigns the claim to a new lender
- Final repayment settles the loan and burns the NFT.

## Requirements

- Go 1.26
- Docker and Docker Compose v2
- Node.js 22+ (for the Hardhat smart-contract project)
- jq (for extracting ABI/bytecode when regenerating Go bindings)
- Kind, kubectl, and Helm v3 for the Besu network
- [pre-commit](https://pre-commit.com/) for the code-generation git hooks

## Quick start

To run the demo:

```bash
make paladin-up # start the local Besu network
make dev-up # start database, run migrations, start the API
make demo # run a simple demo that checks readiness, deploys a contract, originates two loans, repays one, defaults the other, prints results
make dev-down # teardown database and API
make paladin-down # take down the blockchain
```

For development, a few other onetime setup commands are required

```bash
make hooks # onetime pre-commit install
make contracts-install # npm install for contract development
```

The API initializes an Ethereum client at startup, verifies `CHAIN_ID`, and logs the signer address derived from `DEPLOYER_PRIVATE_KEY`.

Generated code is regenerated via the pre-commit hooks.
By convention, we commit all generated code, so that `git clone`rs don't need to re-run generation.

## Smart contract

The ERC-721 contract (OpenZeppelin 5.2.0, Solidity 0.8.28) lives in `contracts/`, built with [Hardhat](https://hardhat.org/).

```bash
make contracts-install # npm install
make contracts-build # compile
make contracts-deploy # deploy LoanNote to local Besu via Hardhat Ignition, not actually used for the demo, more for dev/debugging
```

Go bindings for `LoanNote` are generated from the compiled Hardhat artifact using [go-ethereum's `abigen`](https://geth.ethereum.org/docs/tools/abigen).

```bash
make bindings
```

## Configuration

Configuration is read from environment variables.
For hardhat and docker compose, we load env from `.env` for convenience.
`.env.example` is ready to copy to `.env` and uses throwaway private keys for valid but unfunded addresses.
Our application loads defaults that match `.env` but we could just as easily fail startup when the requisite env isn't present.

| Variable               | Default                             | Description                               |
|------------------------|-------------------------------------|-------------------------------------------|
| `PORT`                 | `8080`                              | TCP port the API listens on               |
| `ETH_RPC_URL`          | `http://127.0.0.1:31545`            | Host-tool Besu JSON-RPC endpoint          |
| `API_ETH_RPC_URL`      | `http://host.docker.internal:31545` | API container Besu JSON-RPC endpoint      |
| `CHAIN_ID`             | `1337`                              | Chain id of the Besu network              |
| `DATABASE_URL`         | local Postgres                      | Postgres connection string                |
| `LOAN_BASE_URI`        | local API loans URI                 | Base URI used to build loan metadata URIs |
| `DEPLOYER_PRIVATE_KEY` | throwaway dev key                   | In-app transaction signer key             |

The API signs Ethereum transactions in-process with `DEPLOYER_PRIVATE_KEY` and submits them through `go-ethereum`'s `client.SendTransaction`, which uses `eth_sendRawTransaction` under the hood.

Nonce-sensitive chain writes are serialized with a short-lived DB-backed lock keyed by chain id and signer address.
That keeps multiple API instances from submitting transactions with the same nonce while leaving the runtime itself stateless.
Nonce management is kind of tricky and could be mitigated by using one deployer key per API instance.

## Database

Postgres stores deployed contract records, loan projections, repayments, chain operation state, and coarse app locks.
The compose stack exposes the dev database on `localhost:5432`.
Migrations are handled by [goose](https://github.com/pressly/goose).

## SQL generation

Database access code is generated from SQL files in `db/query` using sqlc.

```bash
make sqlc
```

The pre-commit hook reruns this automatically when staged changes touch `db/query/`, `db/migration/`, or `sqlc.yaml` (see [Git hooks](#git-hooks)).

## Swagger

The Swagger definition is generated with [swag](https://github.com/swaggo/swag) from annotations on the handlers in `internal/api` and served with [gin-swagger](https://github.com/swaggo/gin-swagger).
Browse the UI at `/swagger/index.html`, which loads the spec from `/swagger/doc.json`.
Static copies of the spec are generated and committed at `docs/swagger.json` and `docs/swagger.yaml`.

```bash
make swagger # regenerate staitc swagger pages
```

The pre-commit hook regenerates the swagger artifacts when staged changes touch `internal/api`.

## Git hooks

Generated code (sqlc, contract bindings, swagger docs) is kept in sync by [pre-commit](https://pre-commit.com/) hooks (`.pre-commit-config.yaml`).

```bash
make hooks # one-time setup, runs pre-commit install
```

When a hook regenerates files, the commit aborts so the changes can be reviewed.
Stage them and commit again.
Bypass the hooks with `git commit --no-verify`.

## Testing

```bash
make test # unit tests, race detector on
```
