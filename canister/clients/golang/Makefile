build:
	go mod tidy
	go build .

# Generate client from candid
# https://github.com/aviate-labs/agent-go/blob/main/cmd/goic/README.md
generate_client:
	#goic fetch $(BACKEND_CANISTER_ID) --output=backend/backend.did
	goic generate did backend/backend.did client --output=backend/backend.go --packageName=backend

run: build
	./icp_vault_client

install-goic:
	go install github.com/aviate-labs/agent-go/cmd/goic@v0.3.3