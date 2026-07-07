# kaleido-project

[![ci](https://github.com/loukianos/kaleido-project/actions/workflows/ci.yml/badge.svg)](https://github.com/loukianos/kaleido-project/actions/workflows/ci.yml)

A Go API microservice for a loan-note business use case backed by an ERC-721 NFT.
Each token represents a lender's claim to repayment.
- Origination mints the note to the lender, who owns it outright
- Transfer assigns the claim to a new lender via owner-signed ERC-721 transfers
- Final repayment settles the loan and burns the NFT.

## Requirements

- Go 1.26
- Docker and Docker Compose v2
- Node.js 22+ (for the Hardhat smart-contract project)
- jq (for extracting ABI/bytecode when regenerating Go bindings)
- Terraform 1.15+
- Kind, kubectl 1.36.x, and Helm 4.2.x for Kubernetes/EKS workflows
- [pre-commit](https://pre-commit.com/) for the code-generation git hooks

## Quick start

To run the demo:

```bash
make bootstrap # set dummy variables for secret-like things
make paladin-up # start the local Besu network
make dev-up # start database, run migrations, start the API
make demo # authenticates all actors via Keycloak, onboards custodial lenders, deploys a contract, warehouses + sells a loan, shows anonymous/non-owner/non-onboarded calls refused, moves a note between custodial lenders under their own tokens, originates a loan on a second contract instance, then floods the servicer key pool: overflow mints come back 202 and the reconciler converges every loan to active with no client retries
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

Authorization uses OpenZeppelin `AccessControl` rather than a single owner:

- `ORIGINATOR_ROLE` may originate (mint) notes.
- `SERVICER_ROLE` may settle (burn on final repayment) and mark loans defaulted.
- `DEFAULT_ADMIN_ROLE` grants and revokes roles; the deployer starts with all three.

There is deliberately no admin transfer path.
Notes move only through standard owner-signed ERC-721 transfers, so the platform provably cannot reassign a lender's claim.
The one exception is settlement, which burns the note regardless of holder because final repayment extinguishes the claim.
Transfers are signed by whoever holds the note: a custodial lender's key when the platform custodies it, or the platform key for warehouse originations to its own address.
Externally held notes can't be moved by the API at all (409); their owners transfer on-chain directly.

Any number of contract instances can be deployed per chain; each instance is its own loan series (for example per originator or per vintage).
At most one contract is *active* — the default series for new originations — and `POST /loans` accepts a `contract_id` to originate into a specific series instead.
The chain's first contract becomes active automatically; later deploys take over the default only when the deploy request sets `activate`, and `POST /admin/contracts/{id}/activate` switches the default at any time.

```bash
make contracts-install # npm install
make contracts-build # compile
make contracts-test # hardhat contract tests
make contracts-deploy # deploy LoanNote to local Besu via Hardhat Ignition, not actually used for the demo, more for dev/debugging
```

Go bindings for `LoanNote` are generated from the compiled Hardhat artifact using [go-ethereum's `abigen`](https://geth.ethereum.org/docs/tools/abigen).

```bash
make bindings
```

## Configuration

Configuration is read from environment variables.
For hardhat and docker compose, we load env from `.env` for convenience.
`.env.example` is ready to copy to `.env`, but secret-like values are intentionally blank.
Fill them locally and do not commit `.env`.
Our application loads defaults that match `.env` but we could just as easily fail startup when the requisite env isn't present.

To generate local-only values for a new machine:

```bash
make bootstrap
```

This creates `.env` and `.local/keycloak-realm.json`, filling `DEPLOYER_PRIVATE_KEY`, `KEY_ENCRYPTION_MASTER_KEY`, `KEYCLOAK_ADMIN_PASSWORD`, `SERVICER_CLIENT_SECRET`, `ALICE_PASSWORD`, and `BOB_PASSWORD`.
The target requires `openssl` and `jq`, and refuses to overwrite existing local files.
If Keycloak is already running, recreate the local stack after changing `.local/keycloak-realm.json` so the realm import sees the new values.

| Variable               | Default                             | Description                               |
|------------------------|-------------------------------------|-------------------------------------------|
| `PORT`                 | `8080`                              | TCP port the API listens on               |
| `ETH_RPC_URL`          | `http://127.0.0.1:31545`            | Host-tool Besu JSON-RPC endpoint          |
| `API_ETH_RPC_URL`      | `http://host.docker.internal:31545` | API container Besu JSON-RPC endpoint      |
| `CHAIN_ID`             | `1337`                              | Chain id of the Besu network              |
| `DATABASE_URL`         | local Postgres                      | Postgres connection string                |
| `LOAN_BASE_URI`        | local API loans URI                 | Base URI used to build loan metadata URIs |
| `DEPLOYER_PRIVATE_KEY` | required locally / GitHub Secret in cloud | In-app transaction signer key             |
| `KEY_ENCRYPTION_MASTER_KEY` | required locally; unused by cloud KMS | AES-256 master key encrypting custodial signing keys at rest |
| `OIDC_ISSUER_URL`      | local Keycloak realm                | Expected token issuer                     |
| `OIDC_JWKS_URL`        | derived from issuer                 | Where the API fetches signing keys (compose overrides to the service-network URL) |
| `OIDC_AUDIENCE`        | `loan-notes-api`                    | Expected token audience                   |
| `SERVICER_KEY_POOL_SIZE` | `2`                               | Extra platform signing keys for concurrent servicer chain writes |
| `RECONCILE_INTERVAL_SECONDS` | `5`                           | How often the reconciler drains pending chain operations |
| `KEY_ENCRYPTOR`        | `local-aes-gcm`                     | Signing-key encryption backend: `local-aes-gcm` or `aws-kms` |
| `KMS_KEY_ID`           | —                                   | KMS key id, required when `KEY_ENCRYPTOR=aws-kms`          |
| `SERVICER_CLIENT_SECRET` | required for demo                 | Keycloak service-client secret used by `cmd/demo` |
| `ALICE_PASSWORD`       | required for demo                   | Demo lender password used by `cmd/demo` |
| `BOB_PASSWORD`         | required for demo                   | Demo lender password used by `cmd/demo` |

The API signs Ethereum transactions in-process with `DEPLOYER_PRIVATE_KEY` and submits them through `go-ethereum`'s `client.SendTransaction`, which uses `eth_sendRawTransaction` under the hood.

Nonce-sensitive chain writes are serialized with a short-lived DB-backed lock keyed by chain id and signing address.
That keeps multiple API instances from submitting transactions with the same nonce while leaving the runtime itself stateless.
Because the lock is per key, custodial identities' writes never contend with each other or with the platform key.

Platform-signed operations (originate, settle, default) go through a **servicer key pool**: `SERVICER_KEY_POOL_SIZE` extra keys, envelope-encrypted like lender keys and granted `ORIGINATOR_ROLE`/`SERVICER_ROLE` on every contract instance by the admin key.
The submitter tries each pool member's lock with a rotating start and takes the first free one, so concurrent servicer writes spread across independent nonce sequences instead of serializing; grants are reconciled at startup and after each deploy, and every chain operation records which key signed it (`signer_address`).
Pool keys sign lifecycle operations only — warehouse custody stays with the primary platform address, so they never hold assets.

## Authentication

Every endpoint except the system ones (`/`, `/healthz`, `/ready`, `/swagger`) requires an OIDC bearer token.
The compose stack runs [Keycloak](https://www.keycloak.org/) with a seeded `loan-notes` realm.
The concrete realm file is local-only and ignored by git; start from `.local/keycloak-realm.example.json`, write `.local/keycloak-realm.json`, and make its values match your `.env` demo credentials:

- a `servicer` service client (client-credentials grant) whose service account holds the `servicer` and `admin` realm roles
- two demo lender users, `alice` and `bob` (password grant via the public `loan-notes-app` client), with no roles

The API validates tokens against the issuer's JWKS (issuer + audience checks) and reads roles from the `realm_access` claim; roles are never self-assigned.
Route policy:

| Route | Who |
|-------|-----|
| `POST /admin/contracts/deploy`, `POST /admin/contracts/{id}/activate` | admin |
| `POST /lenders/onboard` | any authenticated lender (platform service accounts are refused) |
| `POST /loans`, `POST /loans/{id}/repayments`, `POST /loans/{id}/default` | servicer |
| `POST /loans/{id}/transfer` | the note's owner: the holding lender under their own token, or the servicer for platform-held warehouse notes |
| `GET` endpoints | any authenticated caller; lenders see only loans they hold |

Lender identities are keyed by the token's `(iss, sub)`.
`POST /lenders/onboard` is the explicit onboarding step: it creates the caller's identity, eagerly provisions their custodial wallet, and returns the subject + address the lender hands to the servicer.
Provisioning strictly precedes participation: naming a lender by `lender_subject` or `to_subject` requires that they onboarded first (422 otherwise), so a typo'd subject can never mint a note to an identity that doesn't exist.

Example token fetch against the dev realm:

```bash
curl -s http://localhost:8081/realms/loan-notes/protocol/openid-connect/token \
  -d "grant_type=password&client_id=loan-notes-app&username=alice&password=${ALICE_PASSWORD}"
```

## Identities and custodial keys

Lenders can be named two ways when originating or receiving a transfer: `lender_address`/`to_address` for an external wallet, or `lender_subject`/`to_subject` for a **custodial identity**.
A custodial identity gets its own secp256k1 signing key, provisioned when the lender onboards (`POST /lenders/onboard`); on a network that charges gas, onboarding is where the wallet would be funded.
The note is minted to that key's address, so the lender owns it on-chain, and API transfers of the note are signed with the lender's key — the platform key cannot move it.

Key material is envelope-encrypted at rest (AES-256-GCM under `KEY_ENCRYPTION_MASTER_KEY`) and decrypted only for the duration of a request.
The `signing_keys` schema records the encryption scheme and key version per row, so a cloud KMS implementation can slot in behind the same `keys.Encryptor` interface without a migration.

Chain writes take a DB-backed writer lock named by chain id and signing address, so writes by different identities never contend; only same-key writes serialize.

## Failure tolerance

Chain writes are journaled in `chain_operations` before submission, which fixes the error contract by failure class:

- **Futile to retry** (validation, authorization, an on-chain revert): reported immediately as 4xx (or a revert error) and the operation fails terminally.
- **Transient with intent recorded** (chain unreachable, all signer locks busy, receipt timeout): the API returns **202 Accepted** with the loan and operation id — the platform owns the retry, and the client polls `GET /loans/{id}` until the status converges. No client retry loops.
- **Intent not recorded** (database unreachable mid-request): a plain 5xx. `POST /loans` accepts an `external_ref` idempotency key so that retrying this one case can never create a duplicate loan.

A background **reconciler** (one leader across the fleet, elected via the lock manager every `RECONCILE_INTERVAL_SECONDS`) drains the journal: retryable originations, settlements, and defaults are re-driven from durable state; stale submitted transactions are resolved by receipt (mined → applied, reverted → failed, dropped → resubmitted); operations that exhaust their attempts fail terminally, failing the loan too when the operation carried its fate (mint, burn).
Transfers are deliberately never re-driven — re-signing a user-initiated custody change later is the wrong default — so their transient failures report immediately and the lender re-requests.

`GET /operations/{id}` (servicer) exposes an operation's status, attempts, last error, and transaction hash for operational debugging; lenders just poll their loan.

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

## Cloud deployment (EKS)

`deploy/terraform` provisions the production stack in AWS (profile `loukianos`, region `us-east-1`): a VPC, an EKS cluster with managed nodes and the EBS CSI addon, RDS Postgres, an ECR repository, a KMS key for signing-key encryption, and an IRSA role scoping the API pods to `kms:Encrypt`/`kms:Decrypt` on that one key.
In the cloud the API runs with `KEY_ENCRYPTOR=aws-kms`: custodial key material is sealed by KMS instead of a local master key, so the plaintext master key never exists in the cluster — the same `keys.Encryptor` seam, a stronger backend.

The Besu network runs in-cluster via the same Paladin operator install as local development, and Keycloak runs in-cluster from the same seeded realm, so the cloud environment is the local one writ large.

```bash
make cloud-up         # terraform apply: VPC, EKS, RDS, ECR, KMS (~20 minutes)
make cloud-upgrade-k8s-from-1-31 # one-time sequential upgrade for an existing 1.31 cluster to 1.36
make cloud-besu-up    # install the Paladin/Besu network on the cluster
make cloud-push       # build the linux/amd64 image and push to ECR
make cloud-deploy     # helm install: migrations, API, Keycloak; wires public LB endpoints in a second pass
make cloud-demo       # run the full demo against the cloud deployment
make cloud-ci-config  # one-time: seed GitHub variables/secrets for the deploy workflow
make cloud-down       # tear everything down
```

`cloud-deploy` runs in two phases because token issuers must match what callers see: the first install brings everything up behind LoadBalancers, and the second pass sets `OIDC_ISSUER_URL` and `LOAN_BASE_URI` to the LB hostnames once they exist.
Manual cloud deploys require `DEPLOYER_PRIVATE_KEY`, `KEYCLOAK_ADMIN_PASSWORD`, and a realm file at `.local/keycloak-realm.json` or `KEYCLOAK_REALM_FILE`.
JWKS is still fetched over the cluster network — the same issuer/JWKS split used in docker-compose.

Terraform targets new EKS clusters at Kubernetes `1.36`. Existing EKS clusters must be upgraded one minor version at a time; use `make cloud-upgrade-k8s-from-1-31` for the current production cluster before relying on the default `make cloud-up` target. The target applies `1.32`, `1.33`, `1.34`, `1.35`, then `1.36` in order so AWS does not reject a direct minor-version jump.

### Infrastructure deployment

`.github/workflows/terraform.yml` owns Terraform changes separately from application deployment. Pull requests that touch `.github/workflows/terraform.yml` or `deploy/terraform/**` run `terraform plan`; pushes to `main` for those paths run `terraform apply`. Applies detect the live EKS cluster version and advance one Kubernetes minor at a time until `KUBERNETES_VERSION` is reached.

Terraform state is stored in S3 at `s3://kaleido-project-tfstate-433484250096-us-east-1/kaleido/terraform.tfstate` with native S3 lockfiles. The state bucket has versioning, default server-side encryption, and public access blocked.

| Variable | Description |
|----------|-------------|
| `KUBERNETES_VERSION` | Optional EKS target version for the Terraform workflow; defaults to `1.31` so the first CI run is a no-op |

Set `KUBERNETES_VERSION=1.31` before merging the Terraform workflow if you want the first CI run to prove remote state access without upgrading the cluster. Change it to `1.36` when ready to perform the sequential EKS upgrade.

### Continuous deployment

`.github/workflows/deploy.yml` deploys every push to `main`: build and push to ECR, `helm upgrade`, then the full demo runs against the deployment as a smoke test.
It reads the repo variables/secrets seeded by `make cloud-ci-config`, using the `github-actions-deployer` IAM user credentials stored as repo secrets.
Secret-like deployment inputs live in GitHub Secrets, not in tracked files: `DATABASE_URL`, `DEPLOYER_PRIVATE_KEY`, `KEYCLOAK_ADMIN_PASSWORD`, `KEYCLOAK_REALM_JSON`, `SERVICER_CLIENT_SECRET`, `ALICE_PASSWORD`, and `BOB_PASSWORD`.
This workflow deploys application changes only; infrastructure changes are handled by the separate Terraform workflow.

## CI

GitHub Actions (`.github/workflows/ci.yml`) runs on pushes to `main` and on pull requests:

- **lint** — `golangci-lint` (pinned to the same version used locally)
- **test** — `make test` (race detector on)
- **generated** — regenerates sqlc, swagger, and contract bindings, then fails if the committed copies are out of sync (the CI enforcement of the [pre-commit hooks](#git-hooks))
- **contracts** — Hardhat contract tests (`make contracts-test`)
- **docker** — verifies the API image builds

Contract bytecode is compiled with `metadata.bytecodeHash: "none"` so the committed Go bindings are reproducible byte-for-byte across machines.
Without it, solc appends an IPFS metadata hash that shifts with toolchain state and would make the sync check flaky.
