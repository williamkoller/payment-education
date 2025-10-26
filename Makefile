test:
	go test -v ./...

cover:
	go test -buildvcs=false -coverpkg=./... -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html