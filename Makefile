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
	docker compose up -d --wait postgres keycloak
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

# ---------- Cloud (EKS) ----------
# All AWS access uses the loukianos profile; CI overrides AWS_PROFILE to empty for the default chain.
AWS_PROFILE ?= loukianos
KUBERNETES_VERSION ?= 1.36
KUBERNETES_UPGRADE_VERSIONS ?= 1.32 1.33 1.34 1.35 1.36
TF := AWS_PROFILE=$(AWS_PROFILE) terraform -chdir=deploy/terraform

.PHONY: cloud-up
cloud-up:
	$(TF) init -input=false
	$(TF) apply -input=false -auto-approve -var=kubernetes_version=$(KUBERNETES_VERSION)

.PHONY: cloud-upgrade-k8s-from-1-31
cloud-upgrade-k8s-from-1-31:
	$(TF) init -input=false
	@set -e; \
	for version in $(KUBERNETES_UPGRADE_VERSIONS); do \
		echo "Upgrading EKS control plane and managed node group to Kubernetes $$version"; \
		$(TF) apply -input=false -auto-approve -var=kubernetes_version=$$version; \
	done

.PHONY: cloud-kubeconfig
cloud-kubeconfig:
	aws eks update-kubeconfig --profile "$(AWS_PROFILE)" \
		--region "$$($(TF) output -raw region)" \
		--name "$$($(TF) output -raw cluster_name)"

.PHONY: cloud-besu-up
cloud-besu-up: cloud-kubeconfig
	# EKS ships gp2 without the default annotation; the Besu PVCs don't name a class, so one must be default.
	kubectl patch storageclass gp2 -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
	helm repo add paladin https://LFDT-Paladin.github.io/paladin --force-update
	helm repo add jetstack https://charts.jetstack.io --force-update
	helm upgrade --install paladin-crds paladin/paladin-operator-crd
	helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.16.1 --set crds.enabled=true
	helm upgrade --install paladin paladin/paladin-operator -n paladin --create-namespace
	kubectl wait -n paladin --for=create statefulset/besu-node1 --timeout=300s
	kubectl wait -n paladin --for=create statefulset/besu-node2 --timeout=300s
	kubectl wait -n paladin --for=create statefulset/besu-node3 --timeout=300s
	kubectl rollout status -n paladin statefulset/besu-node1 --timeout=300s
	kubectl rollout status -n paladin statefulset/besu-node2 --timeout=300s
	kubectl rollout status -n paladin statefulset/besu-node3 --timeout=300s

.PHONY: cloud-push
cloud-push:
	$(eval ECR_URL := $(shell $(TF) output -raw ecr_repository_url))
	aws ecr get-login-password --profile "$(AWS_PROFILE)" --region "$$($(TF) output -raw region)" \
		| docker login --username AWS --password-stdin "$(ECR_URL)"
	docker buildx build --platform linux/amd64 \
		--build-arg VERSION="$$(git rev-parse --short HEAD)" \
		-t "$(ECR_URL):latest" --push .

.PHONY: cloud-deploy
cloud-deploy: cloud-kubeconfig
	helm upgrade --install kaleido deploy/chart -n kaleido --create-namespace \
		--set image.repository="$$($(TF) output -raw ecr_repository_url)" \
		--set databaseUrl="$$($(TF) output -raw database_url)" \
		--set kmsKeyId="$$($(TF) output -raw kms_key_id)" \
		--set irsaRoleArn="$$($(TF) output -raw api_irsa_role_arn)" \
		--set awsRegion="$$($(TF) output -raw region)" \
		--set-file keycloak.realmJson=.local/keycloak-realm.json \
		--wait --timeout 10m
	# Phase two: the LoadBalancer hostnames exist now, so tokens can carry the public issuer and metadata URIs the public API.
	@until kubectl get svc kaleido-api -n kaleido -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' | grep -q amazonaws; do echo "waiting for api load balancer"; sleep 5; done
	@until kubectl get svc kaleido-keycloak -n kaleido -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' | grep -q amazonaws; do echo "waiting for keycloak load balancer"; sleep 5; done
	API_LB="$$(kubectl get svc kaleido-api -n kaleido -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')"; \
	KC_LB="$$(kubectl get svc kaleido-keycloak -n kaleido -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')"; \
	helm upgrade kaleido deploy/chart -n kaleido --reuse-values \
		--set oidcIssuerUrl="http://$$KC_LB/realms/loan-notes" \
		--set loanBaseUri="http://$$API_LB/loans/" \
		--wait --timeout 5m

.PHONY: cloud-demo
cloud-demo: cloud-kubeconfig
	API_URL="http://$$(kubectl get svc kaleido-api -n kaleido -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')" \
	KEYCLOAK_URL="http://$$(kubectl get svc kaleido-keycloak -n kaleido -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')" \
	go run ./cmd/demo

.PHONY: cloud-ci-config
cloud-ci-config:
	# Seeds the GitHub repo variables/secrets the deploy workflow reads, from terraform outputs. One-time after cloud-up.
	gh variable set AWS_REGION --body "$$($(TF) output -raw region)"
	gh variable set EKS_CLUSTER_NAME --body "$$($(TF) output -raw cluster_name)"
	gh variable set ECR_REPOSITORY_URL --body "$$($(TF) output -raw ecr_repository_url)"
	gh variable set KMS_KEY_ID --body "$$($(TF) output -raw kms_key_id)"
	gh variable set API_IRSA_ROLE_ARN --body "$$($(TF) output -raw api_irsa_role_arn)"
	gh secret set DATABASE_URL --body "$$($(TF) output -raw database_url)"

.PHONY: cloud-down
cloud-down: cloud-kubeconfig
	# The LoadBalancers live outside terraform; remove them first or the VPC destroy hangs.
	helm uninstall kaleido -n kaleido --wait || true
	$(TF) destroy -input=false -auto-approve

# ---------- Contracts ----------
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
