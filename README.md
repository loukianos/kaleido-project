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

## Quick start

```bash
make paladin-up # start the local blockchain network
make dev-up # start database, run migrations, start the API
make demo # run a simple demo that checks readiness, deploys a contract, originates two loans, repays one, defaults the other, prints results
make dev-down # teardown database an API
make paladin-down # take down the blockchain
```

The API initializes an Ethereum client at startup, verifies `CHAIN_ID`, and logs the signer address derived from `DEPLOYER_PRIVATE_KEY`.

Generated code is committed, so you don't need to regenerate sqlc or contract bindings.

## Smart contract

The ERC-721 contract (OpenZeppelin 5.2.0, Solidity 0.8.28) lives in `contracts/`, built with Hardhat.
The local Hardhat deploy helper uses Ignition against the configured Besu JSON-RPC network.

```bash
make contracts-install # npm install
make contracts-build # compile
make contracts-deploy # deploy LoanNote to local Besu via Hardhat Ignition, not actually used for the demo, more for dev/debugging
```

## Configuration

Configuration is read from environment variables.
For hardhat and docker compose, we load env from `.env` for convenience
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

Nonce-sensitive chain writes are serialized with a short-lived DB-backed lock
keyed by chain id and signer address.
That keeps multiple API instances from submitting transactions with the same nonce while leaving the runtime itself stateless.
Nonce management is kind of tricky and could be mitigated by using one deployer key per API instance.

## Database

Postgres stores deployed contract records, loan projections, repayments, chain operation state, and coarse app locks. The compose stack exposes the dev database on `localhost:5432`.
Migrations are handled by [goose](https://github.com/pressly/goose).

## SQL generation

Database access code is generated from SQL files in `db/query` using sqlc.

```bash
make sqlc
```

In a real project we might handle this with a pre-commit hook that watches the directory that `sqlc.yaml` watches.

## Contract bindings

Go bindings for `LoanNote` are generated from the compiled Hardhat artifact using go-ethereum's `abigen`.

```bash
make bindings
```

The generated Go binding is committed under `internal/contracts` so normal Go
builds do not require Node.js.

## API

| Method | Path                      | Description                                                           |
|--------|---------------------------|-----------------------------------------------------------------------|
| `GET`  | `/`                       | Service name and build version                                        |
| `GET`  | `/healthz`                | Liveness probe (`{"status":"ok"}`)                                    |
| `GET`  | `/ready`                  | Readiness probe with start time and sub-checks                        |
| `POST` | `/admin/contracts/deploy` | Deploy and activate `LoanNote` contract                               |
| `GET`  | `/contracts/active`       | Read active contract metadata                                         |
| `POST` | `/loans`                  | Originate a loan note and mint its corresponding NFT                  |
| `GET`  | `/loans`                  | List loans by optional lender/status filters                          |
| `GET`  | `/loans/{id}`             | Read one loan                                                         |
| `POST` | `/loans/{id}/transfer`    | Transfer loan to a new lender onchain                                 |
| `POST` | `/loans/{id}/default`     | Mark an active loan defaulted on chain and in the API                 |
| `POST` | `/loans/{id}/repayments`  | Record a repayment; final payment settles the loan and burns the note |
| `GET`  | `/loans/{id}/repayments`  | List repayments for a loan                                            |
| `GET`  | `/loans/{id}/terms`       | Terms JSON target for `tokenURI`                                      |

## Testing

```bash
make test # unit tests, race detector on
```
