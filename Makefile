all:
	go build -o server cmd/cli/cli.go
clean:
	rm server
