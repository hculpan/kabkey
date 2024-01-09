# kabkey
A programming language and virtual machine written in Go, based on "Writing an Interpreter in Go" by Thorsten Ball.

This project is based on Ball's Monkey programming language, but will have some differences, mostly minor.

# Building

To build executables:
```make build```

This will produce 3 executable files, ```kabr``` (REPL), ```kabc``` (compiler), and ```kabv``` (virtural machine). These files may be found in the ```dist``` directory.

# Testing

To run tests: ```make test```

# Run without building

Run REPL: ```go run cmd/repl/*.go```  
Run Compiler: ```go run cmd/compiler/*.go <source file>```  
Run VM: ```go run cmd/vm/*.go <exe file>```