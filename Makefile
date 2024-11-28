AGENT_BINARY_NAME=agent
SERVER_BINARY_NAME=server


agent-build:
	cd cmd/agent &&  go build -ldflags "-X main.buildVersion=1.0.0 -X main.buildDate=$(date +%Y-%m-%d) -X main.buildCommit=$(git rev-parse HEAD)" -o ${AGENT_BINARY_NAME}

agent-build-platforms:
	cd cmd/agent && GOARCH=amd64 GOOS=darwin go build -o ${AGENT_BINARY_NAME}-darwin main.go
	cd cmd/agent && GOARCH=amd64 GOOS=linux go build -o ${AGENT_BINARY_NAME}-linux main.go
	cd cmd/agent && GOARCH=amd64 GOOS=windows go build -o ${AGENT_BINARY_NAME}-windows main.go

agent-run: agent-build
	./cmd/agent/${AGENT_BINARY_NAME}


server-build:
	cd cmd/server &&  go build -ldflags "-X main.buildVersion=1.0.0 -X main.buildDate=$(date +%Y-%m-%d) -X main.buildCommit=$(git rev-parse HEAD)" -o ${SERVER_BINARY_NAME}

server-build-platforms:
	cd cmd/server && GOARCH=amd64 GOOS=darwin go build -o ${SERVER_BINARY_NAME}-darwin main.go
	cd cmd/server && GOARCH=amd64 GOOS=linux go build -o ${SERVER_BINARY_NAME}-linux main.go
	cd cmd/server && GOARCH=amd64 GOOS=windows go build -o ${SERVER_BINARY_NAME}-windows main.go

server-run: server-build
	./cmd/server/${SERVER_BINARY_NAME}


test:
	go test ./...

test_coverage:
	go test ./... -coverprofile cover.out && go tool cover -func cover.out


vet:
	go vet

