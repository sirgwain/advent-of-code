build: generate
	go build -o advent-of-code main.go

generate: 
	go generate ./...
