#include <iostream>
#include <stdio.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h> 
#include <string.h>
#include <stdlib.h>
#include <dirent.h>
#include <sys/stat.h>
#include <sys/time.h> //FD_SET, FD_ISSET, FD_ZERO macros 
#include <errno.h> 
#include <string>
#include "lib.hpp"
#include <iostream>
#include <string>
#include <signal.h>

using namespace std;

int main(int argc , char *argv[]){
    struct sigaction sa;
    sa.sa_handler = SIG_IGN;
    sa.sa_flags = 0;
    if (sigaction(SIGPIPE, &sa, 0) == -1) {
        perror("sigaction");
        exit(1);
    }
    mkdir("./server_f", S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);

    char Buf[BUFF_SIZE];//Buf用來存放read到的資料
    char filenm[BUFF_SIZE];//fike name buffer
    char name_map[MAX_CLIENT][MAX_NAME] = {'\0'};
    FILE *fp;

    int localSocket, remoteSocket;// local: server, remote: client                            
    
    string port_stg = argv[1];
    int port = atoi(port_stg.c_str());
    printf("Get port number: %d\n",port);
    //int port = 8888;
    struct  sockaddr_in localAddr,remoteAddr;
    int addrLen = sizeof(struct sockaddr_in);  
    
    //(localSocket = socket(AF_INET , SOCK_STREAM , 0));
    if ((localSocket = socket(AF_INET , SOCK_STREAM , 0)) < 0){
        perror("socket create error!");
        exit(1);
    }
    printf("Server socket creation successful.\n");

    localAddr.sin_family = AF_INET;
    localAddr.sin_port = htons(port);
    localAddr.sin_addr.s_addr = INADDR_ANY;

    if(bind(localSocket,(struct sockaddr *)&localAddr , sizeof(localAddr)) < 0) {
        perror("Bind error!");
        exit(1);
    }
    printf("Binding successful.\n");

    if(listen(localSocket , 3) < 0){
    	perror("Listen error");
    	exit(1);
    }

    // select_1: start
    int client_socket[MAX_CLIENT];
    fd_set readfds;
    for(int i = 0; i < MAX_CLIENT; i++)
        client_socket[i] = 0;
    int sd;


    while(1){ 
        printf("Waiting for connections...\n"); 
        printf("Server Port: %d\n", port);

        FD_ZERO(&readfds);
        FD_SET(localSocket, &readfds);
        int max_sd = localSocket;
        printf("max_sd = %d\n", max_sd);

        for(int i = 0;  i< MAX_CLIENT; i++){
            if(client_socket[i] > 0)
                FD_SET(client_socket[i], &readfds);
            if(client_socket[i] > max_sd)
                max_sd = client_socket[i];
        }

        int select_stat = select(max_sd + 1, &readfds, NULL, NULL, NULL);
        if ((select_stat  < 0) && (errno != EINTR)) { 
            printf("select error\n"); 
        } 
        //If something happened on the master socket , 
        //then its an incoming connection 
        if (FD_ISSET(localSocket, &readfds)){ 
            //remoteSocket = accept(localSocket, (struct sockaddr *)&remoteAddr, (socklen_t*)&addrLen);  
            if ((remoteSocket = accept(localSocket, (struct sockaddr *)&remoteAddr, (socklen_t*)&addrLen)) < 0) {
                perror("Accept error");
                exit(1);
            }
            printf("Connection accepted.\n");
            //add new socket to array of sockets 
            for(int i = 0; i < MAX_CLIENT; i++){
                if(client_socket[i] == 0){
                    client_socket[i] = remoteSocket;
                    printf("Adding to socket list : %d\n" , i);
                    break;
                }
            }  
        } 

        for(int i = 0; i < MAX_CLIENT; i++){
            sd = client_socket[i];
            if(FD_ISSET(sd, &readfds) ){
                bzero(Buf, sizeof(char)* BUFF_SIZE); 
                printf("Reading instruction...\n");

                if ((read(sd, Buf, BUFF_SIZE)) < 0) {
                    perror("Read instruction error.");
                    exit (1);
                }                
                printf("checking = %c\n", Buf[BUFF_SIZE-1]);
                char last_bit = Buf[BUFF_SIZE-1];
                if(last_bit == NONAME){
                    char new_name[BUFF_SIZE];
                    for(int off = 0; off < NAME_MAX; off++)
                        new_name[off] = Buf[off];
                    printf("Name = %s\n", new_name);

                    if(check_name(name_map, new_name, client_socket, sd) == REPEAT){
                        printf("same name\n");
                        if ((write(sd, naming_fail, BUFF_SIZE)) < 0) {
                            perror("writing naming fail error.");
                            exit (1);
                        }
                    }
                    else{
                        printf("naming success\n");
                        if ((write(sd, naming_ok, BUFF_SIZE)) < 0) {
                            perror("writing naming ok error.");
                            exit (1);
                        }
                        for(int off = 0; off < NAME_MAX; off++)
                            name_map[i][off] = new_name[off];
                    }
                    close(sd);
                    client_socket[i] = 0;
                }
                else{
                    printf("Instruction: %s\n", Buf);
                    // ls
                    if(Buf[0] == 'l' && Buf[1] == 's')
                    { 
                        list_file(sd);
                    }
                    else if(Buf[0] == 'g' && Buf[1] == 'e' && Buf[2] == 't')
                    {
                        char filename[BUFF_SIZE] = {'\0'};
                        for(int i = 0; i < MAX_FILENAME; i++)
                            filename[i] = Buf[i+4];
                        printf("file name = [%s]\n", filename);

                        char exist_msg[BUFF_SIZE] = {'\0'};
                        set_exist_msg(filename, exist_msg, SERVER);

                        write(sd, exist_msg, BUFF_SIZE);  
                        if(exist_msg[0] == '1')
                            send_file(filename, sd, SERVER);        
                    }
                    else if(Buf[0] == 'p' && Buf[1] == 'u' && Buf[2] == 't')
                    {
                        char filename[BUFF_SIZE] = {'\0'};
                        for(int i = 0; i < MAX_FILENAME; i++)
                            filename[i] = Buf[i+4];
                        printf("file name = [%s]\n", filename);
                        recv_file(sd, filename, SERVER);
                    }

                    close(sd);
                    client_socket[i] = 0;
                } 
            }
        }
	}
    close(localSocket);
    return 0;
}
