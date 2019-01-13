EXE = chesser-go
all:
	mkdir -p src/
	cp -r vendor/* src/
	GOPATH=`pwd` go build -o $(EXE) .

