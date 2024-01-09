all: build

build:
	go build -0 dist/kabrepl cmd/repl/*.go
	go build -o dist/kabc cmd/compiler/*.go 
	go build -o dist/kabv cmd/vm/*.go 

clean:
	rm -rf dist

test:
	go test ./pkg/lexer/