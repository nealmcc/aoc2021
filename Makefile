test:
	go test ./... -covermode count -coverpkg ./... -coverprofile cover.out ./... -vet
	go tool cover -html cover.out
