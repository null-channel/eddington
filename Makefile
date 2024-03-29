.PHONY: build-proto vet test
build-proto: ## Build proto files.
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./proto/**/*.proto
	## Whoever sees this, plese don't be angry, it's go's fault. not mine
	## And yes, I refuse to use bazel to deal with this because... I should not need it.
	## API
	protoc --go_out=./api --go_opt=paths=source_relative \
		--go-grpc_out=./api --go-grpc_opt=paths=source_relative \
		./proto/**/*.proto
	## I'm sorry world
	## Container Builder
	protoc --go_out=./container-builder --go_opt=paths=source_relative \
		--go-grpc_out=./container-builder --go-grpc_opt=paths=source_relative \
		./proto/**/*.proto
vet: ## Run go vet.
	cd api && go vet ./...
	cd ../container-builder && go vet ./...
	cd ../null-operator && go vet ./...
test: ## Run go vet.
	cd api && go test ./...
	cd ../container-builder && go test ./...
	cd ../null-operator && go test ./...
