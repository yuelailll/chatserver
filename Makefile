all: server

server: 
	go build main.go && && chmod 755 && ./main

clean:
	rm -fr ./main
