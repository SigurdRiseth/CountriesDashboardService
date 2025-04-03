test:
	go test ./... -v -coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out