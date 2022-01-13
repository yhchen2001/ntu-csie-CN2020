build:
	go build -o ./server.o ./server/server.go ./server/sign.go ./server/friend.go ./server/chat.go
	go build -o client.o ./client/client.go ./client/friend.go ./client/chat.go ./client/sign.go

clean:
	rm server.o
	rm client.o