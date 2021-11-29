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
    mkdir("./client_dir", S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);

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

    bool name_entered = false;
    char name[BUFF_SIZE] = {'\0'};

    bool first_time = true;

    while(1){
        int mainSocket;
        mainSocket = socket(AF_INET, SOCK_STREAM, 0);
        if (mainSocket < 0)
        {
            perror("Create socket failed.");
            exit(1);
        }

        struct sockaddr_in mainAddr;
        bzero(&mainAddr, sizeof(mainAddr));

        mainAddr.sin_family = PF_INET;
        mainAddr.sin_port = htons(port);
        mainAddr.sin_addr.s_addr = inet_addr(ip);

        if (connect(mainSocket, (struct sockaddr *)&mainAddr, sizeof(mainAddr)) < 0)
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
                close(mainSocket);
                continue;
            }
            name[BUFF_SIZE-1] = NONAME;
            if (write(mainSocket, name, BUFF_SIZE) < 0)
            {
                perror("Write name msg error");
                exit(1);
            }
            // 這邊socket會確定name有沒有被用掉，如果沒被用掉就會傳 naming_success
            char check_naming[BUFF_SIZE] = {'\0'};
            read(mainSocket, check_naming, BUFF_SIZE);

            if(check_naming[0] == NAMING_FAIL){
                printf("username is in used, please try another:\n");
            }
            else if(check_naming[0] == NAMING_OK)
            {
                name_entered = true;
                printf("connect successfully\n");
            }
            close(mainSocket);
            continue;
        }

        char ins[BUFF_SIZE] = {'\0'};
        //printf("Enter your instruction\n");
        cin.getline(ins, sizeof(char) * BUFF_SIZE);
        //printf("ins[0] == '0' == %d\n", ins[0] == '\0');
        if(ins[0] == '\0')
        {
            close(mainSocket);
            continue;
        }
        //printf("Your instruction is %s\n", ins);
        ins[BUFF_SIZE-1] = NAMED;

        if(ins[0] == 'l' && ins[1] == 's' && (ins[2] == '\0' || ins[2] == ' '))
        {// ls
            ins[2] = '\0';
            for(int i = 3; i < MAX_FILENAME; i++)
            {
                if(ins[i] != ' ' && ins[i] != '\0' && ins[i] != '\n')
                {
                    printf("Command format error\n");
                    goto end;
                }
            }
            int write_stat;
            if ((write_stat = write(mainSocket, ins, sizeof(char) * BUFF_SIZE)) < 0)
            {
                perror("Write Instruction error");
                exit(1);
            }
            //printf("ls %d bytes written, Instruction sent successfully\n", write_stat);
            recv_ls(mainSocket);
        }
        else if(ins[0] == 'g' && ins[1] == 'e' && ins[2] == 't' && (ins[3] == '\0' || ins[3] == ' '))
        {
            ins[3] = ' ';
            char filename[BUFF_SIZE] = {'\0'};
            for(int off = 0; off <= MAX_FILENAME; off++)
                filename[off] = ins[off+4];

            if(file_ok(filename) == false){
                printf("Command format error\n");
                close(mainSocket);
                continue; 
            }
            int write_stat;
            if ((write_stat = write(mainSocket, ins, sizeof(char) * BUFF_SIZE)) < 0)
            {
                perror("Write Instruction error");
                exit(1);
            }

            char check_exist[BUFF_SIZE] = {'\0'};
            if ((read(mainSocket, check_exist, sizeof(char) * BUFF_SIZE)) < 0)
            {
                perror("read check_exist error");
                exit(1);
            }
            //printf("exist msg [%s]\n", check_exist);
            if(check_exist[0] == FILE_EXIST)
            {
                recv_file(mainSocket, filename, CLIENT);
                printf("get %s successfully\n", filename);
            }
            else if(check_exist[0] == FILE_NOT_EXIST)
                printf("The %s doesn’t exist\n", filename);
        }

        else if(ins[0] == 'p' && ins[1] == 'u' && ins[2] == 't' && (ins[3] == '\0' || ins[3] == ' '))
        {
            ins[3] = ' ';
            char filename[BUFF_SIZE] = {'\0'};
            for(int off = 0; off <= MAX_FILENAME; off++)
                filename[off] = ins[off+4];

            if(file_ok(filename) == false){
                printf("Command format error\n");
                close(mainSocket);
                continue; 
            }
            int write_stat;
            if ((write_stat = write(mainSocket, ins, sizeof(char) * BUFF_SIZE)) < 0)
            {
                perror("Write Instruction error");
                exit(1);
            }
            char exist_msg[BUFF_SIZE] = {'\0'};
            set_exist_msg(filename, exist_msg, CLIENT);
            if(exist_msg[0] == '0')
            {
                printf("The %s doesn’t exist\n", filename);
                close(mainSocket);
                continue;
            }    
            send_file(filename, mainSocket, CLIENT);
            printf("put %s successfully\n", filename);

        }
        else{
            printf("Command not found\n");
        }
end:
        close(mainSocket);
    }
}