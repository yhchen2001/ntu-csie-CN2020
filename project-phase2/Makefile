build:
	go build -o ./server.o ./server/server.go ./server/sign.go ./server/friend.go ./server/chat.go
	go build -o client.o ./client/client.go ./client/friend.go ./client/chat.go ./client/sign.go
	go build -o webclient.o ./web/webclient.go ./web/sign_handler.go ./web/friend_handler.go ./web/chat_handler.go

clean:
	rm server.o
	rm client.o
	rm webclient.o