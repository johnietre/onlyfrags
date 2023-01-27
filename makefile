.PHONY: clean

bin/server: server/*.go
	go build -o bin/server $^

server: bin/server

clean:
	rm -rf bin/*
