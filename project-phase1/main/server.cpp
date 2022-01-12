#include <iostream>
#include <stdio.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h> 
#include <string.h>
#include <stdlib.h>
#include <dirent.h>
#include <sys/stat.h>
#include <sys/time.h> 
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
    mkdir("./server_dir", S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);

    char Buf[BUFF_SIZE];//Buf用來存放read到的資料
    char name_map[MAX_CLIENT][MAX_NAME] = {'\0'};
    FILE *fp;

    int mainSocket, clientSocket;                          
    
    string port_stg = argv[1];
    int port = atoi(port_stg.c_str());
    printf("Get port number: %d\n",port);
    struct  sockaddr_in mainAddr,remoteAddr;
    int addrLen = sizeof(struct sockaddr_in);  
    
    if ((mainSocket = socket(AF_INET , SOCK_STREAM , 0)) < 0){
        perror("socket create error!");
        exit(1);
    }
    printf("Server socket creation successful.\n");

    mainAddr.sin_family = AF_INET;
    mainAddr.sin_port = htons(port);
    mainAddr.sin_addr.s_addr = INADDR_ANY;

    if(bind(mainSocket,(struct sockaddr *)&mainAddr , sizeof(mainAddr)) < 0) {
        perror("Bindint error!");
        exit(1);
    }

    if(listen(mainSocket , 3) < 0){
    	perror("Listening error");
    	exit(1);
    }

    int client_fd[MAX_CLIENT];
    fd_set readfds;
    for(int i = 0; i < MAX_CLIENT; i++)
        client_fd[i] = 0;

    while(1){ 

        FD_ZERO(&readfds);
        FD_SET(mainSocket, &readfds);
        int max_fd = mainSocket;
        //printf("max_fd = %d\n", max_fd);

        for(int i = 0;  i< MAX_CLIENT; i++){
            if(client_fd[i] > 0)
                FD_SET(client_fd[i], &readfds);
            if(client_fd[i] > max_fd)
                max_fd = client_fd[i];
        }

        int select_stat = select(max_fd + 1, &readfds, NULL, NULL, NULL);
        if ((select_stat  < 0) && (errno != EINTR)) { 
            printf("select error\n"); 
        } 

        if (FD_ISSET(mainSocket, &readfds)){ 
            if ((clientSocket = accept(mainSocket, (struct sockaddr *)&remoteAddr, (socklen_t*)&addrLen)) < 0) {
                perror("Accept error");
                exit(1);
            }
            printf("Accept ok.\n");

            for(int i = 0; i < MAX_CLIENT; i++){
                if(client_fd[i] == 0){
                    client_fd[i] = clientSocket;
                    break;
                }
            }  
        } 

        for(int i = 0; i < MAX_CLIENT; i++){
            int crr_fd = client_fd[i];
            if(FD_ISSET(crr_fd, &readfds) ){
                memset(Buf, '\0', BUFF_SIZE); 
                if ((read(crr_fd, Buf, BUFF_SIZE)) < 0) {
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

                    if(check_name(name_map, new_name, client_fd, crr_fd) == REPEAT){
                        printf("same name\n");
                        if ((write(crr_fd, naming_fail, BUFF_SIZE)) < 0) {
                            perror("writing naming fail error.");
                            exit (1);
                        }
                    }
                    else{
                        printf("naming success\n");
                        if ((write(crr_fd, naming_ok, BUFF_SIZE)) < 0) {
                            perror("writing naming ok error.");
                            exit (1);
                        }
                        for(int off = 0; off < NAME_MAX; off++)
                            name_map[i][off] = new_name[off];
                    }
                    close(crr_fd);
                    client_fd[i] = 0;
                }
                else{
                    printf("Instruction: %s\n", Buf);
                    // ls
                    if(Buf[0] == 'l' && Buf[1] == 's')
                    { 
                        list_file(crr_fd);
                    }
                    else if(Buf[0] == 'g' && Buf[1] == 'e' && Buf[2] == 't')
                    {
                        char filename[BUFF_SIZE] = {'\0'};
                        for(int i = 0; i < MAX_FILENAME; i++)
                            filename[i] = Buf[i+4];
                        printf("file name = [%s]\n", filename);

                        char exist_msg[BUFF_SIZE] = {'\0'};
                        set_exist_msg(filename, exist_msg, SERVER);

                        write(crr_fd, exist_msg, BUFF_SIZE);  
                        if(exist_msg[0] == '1')
                            send_file(filename, crr_fd, SERVER);        
                    }
                    else if(Buf[0] == 'p' && Buf[1] == 'u' && Buf[2] == 't')
                    {
                        char filename[BUFF_SIZE] = {'\0'};
                        for(int i = 0; i < MAX_FILENAME; i++)
                            filename[i] = Buf[i+4];
                        printf("file name = [%s]\n", filename);
                        recv_file(crr_fd, filename, SERVER);
                    }

                    close(crr_fd);
                    client_fd[i] = 0;
                } 
            }
        }
	}
    close(mainSocket);
    return 0;
}
