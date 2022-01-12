CC = g++
CLIENT = client.cpp
SERVER = server.cpp
LIB = lib.cpp
CLI = client
SER = server

all: server client
server: $(SERVER)
	$(CC) $(LIB) $(SERVER) -o $(SER)
client: $(CLIENT)
	$(CC) $(LIB) $(CLIENT) -o $(CLI)

.PHONY: clean

clean:
	rm $(CLI) $(SER)

