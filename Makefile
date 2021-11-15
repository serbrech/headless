
.PHONY: setup
setup:
	@kind create cluster --name kind
	@kubectl config set-context kind-kind

.PHONY: run
run:
	skaffold dev

