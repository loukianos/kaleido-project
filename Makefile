.PHONY: hooks
hooks:
	pre-commit install

.PHONY: migrate
migrate:
	go run ./cmd/migrate

.PHONY: sqlc
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate

.PHONY: bindings
bindings: contracts-build
	mkdir -p contracts/generated internal/contracts
	jq '.abi' contracts/artifacts/src/LoanNote.sol/LoanNote.json > contracts/generated/LoanNote.abi
	jq -r '.bytecode' contracts/artifacts/src/LoanNote.sol/LoanNote.json > contracts/generated/LoanNote.bin
	go run github.com/ethereum/go-ethereum/cmd/abigen@v1.17.4 \
		--abi contracts/generated/LoanNote.abi \
		--bin contracts/generated/LoanNote.bin \
		--pkg contracts \
		--type LoanNote \
		--out internal/contracts/loan_note.go

.PHONY: test
test:
	go test -race ./...

.PHONY: swagger
swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 fmt -d internal/api -g server.go
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -d internal/api -g server.go -o docs

.PHONY: demo
demo:
	go run ./cmd/demo

.PHONY: lint
lint:
	golangci-lint run

.PHONY: dev-up
dev-up:
	docker compose up -d --wait postgres
	$(MAKE) migrate
	docker compose up -d --build --remove-orphans api

.PHONY: dev-down
dev-down:
	docker compose down -v --remove-orphans

.PHONY: paladin-up
paladin-up:
	kind get clusters | grep -qx paladin || kind create cluster --name paladin --config .local/paladin-kind.yaml
	helm repo add paladin https://LFDT-Paladin.github.io/paladin --force-update
	helm repo add jetstack https://charts.jetstack.io --force-update
	helm upgrade --install paladin-crds paladin/paladin-operator-crd
	helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.16.1 --set crds.enabled=true
	helm upgrade --install paladin paladin/paladin-operator -n paladin --create-namespace
	kubectl wait -n paladin --for=create statefulset/besu-node1 --timeout=180s
	kubectl wait -n paladin --for=create statefulset/besu-node2 --timeout=180s
	kubectl wait -n paladin --for=create statefulset/besu-node3 --timeout=180s
	kubectl rollout status -n paladin statefulset/besu-node1 --timeout=180s
	kubectl rollout status -n paladin statefulset/besu-node2 --timeout=180s
	kubectl rollout status -n paladin statefulset/besu-node3 --timeout=180s

.PHONY: paladin-down
paladin-down:
	kind delete cluster --name paladin

.PHONY: local-reset
local-reset: dev-down paladin-down

.PHONY: contracts-install
contracts-install:
	cd contracts && npm install

.PHONY: contracts-build
contracts-build:
	cd contracts && npm run compile

.PHONY: contracts-test
contracts-test:
	cd contracts && npm test

.PHONY: contracts-deploy
contracts-deploy:
	cd contracts && npm run deploy:besu
