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
#include <vector>

using namespace std;

bool check_name(char name_map[][MAX_NAME], char new_name[], int client_socket[], int sd)
{
    bool tot_same = false;
    for(int i = 0; i < MAX_CLIENT; i++)
    {
        if(client_socket[i] == 0 || client_socket[i] == sd)
            continue;
        bool crr_same = true;
        for(int j = 0; j < MAX_NAME; j++)
            if(name_map[i][j] != new_name[j])
                crr_same = false ;
        //只要有一個就會是named
        tot_same = tot_same | crr_same;
    }
    if(tot_same == true)
        return REPEAT;
    else  
        return NOREPEAT;
}

void list_file(int remoteSocket){
    int sent;
    DIR *d = opendir("./server_f");// opendir() returns a pointer of DIR type
    struct dirent *dir;// Pointer for directory entry
    vector<string> files;
    if(d){
        while((dir = readdir(d)) != NULL){
            char Buf[BUFF_SIZE] = {'\0'};
            strcpy(Buf, dir->d_name);
            if(Buf[0] == '.')  
                continue;
            Buf[strlen(Buf)] = '\0';
            string str(Buf);
            files.push_back(str);
        }
        closedir(d);
    }

    sort(files.begin(), files.end());

    for(auto s: files){
        char Buf[BUFF_SIZE] = {'\0'};
        strcpy(Buf, s.c_str());
        cout << s << endl;
        printf("Buf = %s\n", Buf);

        if(write(remoteSocket, Buf, sizeof(char) * BUFF_SIZE) < 0)
        {
            perror("ls fail\n");
            exit(1);
        }
    }
}


void set_exist_msg(char filenm[], char exist_msg[], int type){
    DIR *d; 
    if(type == SERVER)
        d = opendir("./server_f");
    else if(type == CLIENT)
        d = opendir("./client_f");
    struct dirent *dir;
    if(d){
        while((dir = readdir(d)) != NULL){
            if(strcmp(filenm, dir->d_name) == 0){
                //printf("Finded out the file, file name: %s\n", filenm);
                exist_msg[0] = '1';
                return;
            }
        }
        closedir(d);
    }

    exist_msg[0] = '0';
    //printf("exist[0] = %c, Finding file failed.\n",exist_msg[0]);
}

void send_file(char filename[], int remoteSocket, int type){
    FILE *fp;
    string str_1;

    if(type == CLIENT)
        str_1 = "./client_f/";
    else if(type == SERVER)
        str_1 = "./server_f/";

    string str_2 = filename;
    string str_3 = str_1 + str_2;
    fp = fopen(str_3.c_str(), "rb");
    
    while(!feof(fp)){
        char Buf[BUFF_SIZE] = {'\0'};
        int numbytes = fread(Buf, sizeof(char), sizeof(Buf), fp);
        //printf("fread %d bytes, ", numbytes);
        numbytes = write(remoteSocket, Buf, BUFF_SIZE);
        //printf("Sending %d bytes\n",numbytes);
    }
}

void recv_file(int remoteSocket,char filename[], int type){
    FILE *fp;
    string str_1;

    if(type == CLIENT)
        str_1 = "./client_f/";
    else if(type == SERVER)
        str_1 = "./server_f/";

    string str_2 = filename;
    string str_3 = str_1 + str_2;
    if ((fp = fopen(str_3.c_str(), "wb")) == NULL){
        perror("fopen error");
        exit(1);
    }
    while(1){
        char Buf[BUFF_SIZE] = {'\0'};
        int numbytes = read(remoteSocket, Buf, BUFF_SIZE);
        //printf("read %d bytes, ", numbytes);
        if(numbytes <= 0){
            break;
        }
        numbytes = fwrite(Buf, sizeof(char), numbytes, fp);
        //printf("fwrite %d bytes\n", numbytes);
    }
    fclose(fp);
}

void recv_ls(int localSocket){
    int rec;
    while(1){
        char Buf[BUFF_SIZE] = {'\0'};
        rec = read(localSocket, Buf, BUFF_SIZE);
        if(rec <= 0){
            break;
        }
        printf("%s\n", Buf);
    }
    //printf("end ls\n");
}

bool file_ok(char filename[])
{
    if(filename[MAX_FILENAME] != '\0'){
        //printf("filename too long\n");
        return false;
    }
    if(filename[0] == '\0'){
        //printf("no file name\n");
        return false;
    }
    bool spaced = false;
    for(int i = 0; i < MAX_FILENAME; i++)
    {
        if(spaced && filename[i] == ' '){
            //printf("more then one file\n");
            return false;
        }
        if(filename[i] == ' ')
            spaced = true;
    } 
    if(filename[0] == ' ')
    {
        for(int i = 0; i < MAX_FILENAME; i++)
            filename[i] = filename[i+1];
    }
    return true;      
}