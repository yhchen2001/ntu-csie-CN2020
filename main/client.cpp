#include <iostream>
#include <stdio.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/time.h>
#include <errno.h>
#include <string>
#include "lib.hpp"

using namespace std;


int main(int argc, char *argv[])
{
    mkdir("./client_f", S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);

    char instruction[BUFF_SIZE];
    char Buf[BUFF_SIZE];
    char *ip_port = argv[1];

    // ip_port handling
    char ip[30] = "\0";
    int port = 0;
    int idx = 0;
    while (ip_port[idx] != ':')
    {
        ip[idx] = ip_port[idx];
        ip[idx + 1] = '\0';
        idx++;
    }
    idx++;
    while (ip_port[idx] != '\0')
    {
        port = port * 10 + (ip_port[idx] - '0');
        idx++;
    }
    //cout << "ip = " << ip << ", port = " << port << endl;

    bool name_entered = false;
    char name[BUFF_SIZE] = {'\0'};

    bool first_time = true;

    while(1){
        int localSocket, recved;
        localSocket = socket(AF_INET, SOCK_STREAM, 0);
        if (localSocket < 0)
        {
            perror("Create socket failed.");
            exit(1);
        }

        /* Bind to an arbitrary return address.*/
        struct sockaddr_in localAddr;
        bzero(&localAddr, sizeof(localAddr)); // #include <string.h>，將為sizeof(x)的前幾個東西清成0

        localAddr.sin_family = PF_INET;
        localAddr.sin_port = htons(port);
        localAddr.sin_addr.s_addr = inet_addr(ip);

        //Connect to the server連線
        if (connect(localSocket, (struct sockaddr *)&localAddr, sizeof(localAddr)) < 0)
        {
            perror("Connection failed.");
            exit(1);
        }

        if(name_entered == false)
        {
            memset(name, '\0', BUFF_SIZE);
            if(first_time){
                printf("input your username:\n");
                first_time = false;
            }
            scanf("%s", name);

            if(name[MAX_NAME] != '\0'){
                printf("%s is too long\n", name);
                close(localSocket);
                continue;
            }
            name[BUFF_SIZE-1] = NONAME;
            if (write(localSocket, name, BUFF_SIZE) < 0)
            {
                perror("Write name msg error");
                exit(1);
            }
            // 這邊socket會確定name有沒有被用掉，如果沒被用掉就會傳 naming_success
            char check_naming[BUFF_SIZE] = {'\0'};
            read(localSocket, check_naming, BUFF_SIZE);

            if(check_naming[0] == NAMING_FAIL){
                printf("username is in used, please try another:\n");
            }
            else if(check_naming[0] == NAMING_OK)
            {
                name_entered = true;
                printf("connect successfully\n");
            }
            close(localSocket);
            continue;
        }

        char ins[BUFF_SIZE] = {'\0'};
        //printf("Enter your instruction\n");
        scanf("%s", ins);
        //printf("Your instruction is %s\n", ins);
        ins[BUFF_SIZE-1] = NAMED;

        if(ins[0] == 'l' && ins[1] == 's' && (ins[2] == '\0' || ins[2] == ' '))
        {// ls
            ins[2] = '\0';
            int write_stat;
            if ((write_stat = write(localSocket, ins, sizeof(char) * BUFF_SIZE)) < 0)
            {
                perror("Write Instruction error");
                exit(1);
            }
            //printf("ls %d bytes written, Instruction sent successfully\n", write_stat);
            recv_ls(localSocket);
        }
        else if(ins[0] == 'g' && ins[1] == 'e' && ins[2] == 't' && (ins[3] == '\0' || ins[3] == ' '))
        {
            ins[3] = ' ';
            char filename[BUFF_SIZE] = {'\0'};
            scanf("%[^\n]", filename); // filename前面會有一個空格
            if(file_ok(filename) == false){
                printf("Command format error\n");
                close(localSocket);
                continue; 
            }
            //printf("filename is [%s]\n", filename);
            for(int off = 0; off < MAX_FILENAME; off++)
                ins[off+4] = filename[off];
            //printf("whole ins name is : %s\n", ins);
            int write_stat;
            if ((write_stat = write(localSocket, ins, sizeof(char) * BUFF_SIZE)) < 0)
            {
                perror("Write Instruction error");
                exit(1);
            }

            char check_exist[BUFF_SIZE] = {'\0'};
            if ((read(localSocket, check_exist, sizeof(char) * BUFF_SIZE)) < 0)
            {
                perror("read check_exist error");
                exit(1);
            }
            //printf("exist msg [%s]\n", check_exist);
            if(check_exist[0] == FILE_EXIST)
            {
                recv_file(localSocket, filename, CLIENT);
                printf("get %s successfully\n", filename);
            }
            else if(check_exist[0] == FILE_NOT_EXIST)
                printf("The %s doesn’t exist\n", filename);
        }

        else if(ins[0] == 'p' && ins[1] == 'u' && ins[2] == 't' && (ins[3] == '\0' || ins[3] == ' '))
        {
            ins[3] = ' ';
            char filename[BUFF_SIZE] = {'\0'};
            scanf("%[^\n]", filename); // filename前面會有一個空格
            if( file_ok(filename) == false){
                printf("Command format error\n");
                close(localSocket);
                continue; 
            }
            //printf("filename is [%s]\n", filename);
            for(int off = 0; off < MAX_FILENAME; off++)
                ins[off+4] = filename[off];
            
            //printf("Full instruction is [%s]\n", ins);
            int write_stat;
            if ((write_stat = write(localSocket, ins, sizeof(char) * BUFF_SIZE)) < 0)
            {
                perror("Write Instruction error");
                exit(1);
            }
            char exist_msg[BUFF_SIZE] = {'\0'};
            set_exist_msg(filename, exist_msg, CLIENT);
            if(exist_msg[0] == '0')
            {
                printf("The %s doesn’t exist\n", filename);
                close(localSocket);
                continue;
            }    
            send_file(filename, localSocket, CLIENT);
            printf("put %s successfully\n", filename);

        }
        else{
            printf("Command not found\n");
        }

        //printf("local socket closing\n");
        close(localSocket);
        //sleep(1);

    }
}