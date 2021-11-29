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
#include <sys/types.h>
#include <sys/socket.h>

using namespace std;

bool check_name(char name_map[][MAX_NAME], char new_name[], int client_socket[], int sd)
{
    bool tot_same = false;
    for(int i = 0; i < MAX_CLIENT; i++)
    {
        int optval;
        socklen_t optlen = sizeof(optval);
        int get_socket_stat = getsockopt(sd, SOL_SOCKET, SO_TYPE, &optval, &optlen);

        if(client_socket[i] == 0 || client_socket[i] == sd || get_socket_stat == -1)
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
    DIR *d = opendir("./server_dir");
    struct dirent *file_under_dir;
    vector<string> files;
    if(d){
        while((file_under_dir = readdir(d)) != NULL){
            char Buf[BUFF_SIZE] = {'\0'};
            strcpy(Buf, file_under_dir->d_name);
            if(Buf[0] == '.')  
                continue;
            Buf[strlen(Buf)] = '\0';
            string str(Buf);
            files.push_back(str);
        }
        closedir(d);
    }

    sort(files.begin(), files.end());

    for(int i = 0; i < files.size(); i++){
        string s = files[i];
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


void set_exist_msg(char filename[], char exist_msg[], int type){
    DIR *d; 
    if(type == SERVER)
        d = opendir("./server_dir");
    else if(type == CLIENT)
        d = opendir("./client_dir");
    struct dirent *dir;
    if(d){
        while((dir = readdir(d)) != NULL){
            if(strcmp(filename, dir->d_name) == 0){
                exist_msg[0] = '1';
                return;
            }
        }
        closedir(d);
    }

    exist_msg[0] = '0';
}

void send_file(char filename[], int remoteSocket, int type){
    string dir_name;
    if(type == CLIENT)
        dir_name = "./client_dir/";
    else if(type == SERVER)
        dir_name = "./server_dir/";

    string filename_str = filename;
    string path_name = dir_name + filename_str;
    FILE *fp; 

    if ((fp = fopen(path_name.c_str(), "rb")) == NULL){
        perror("fopen error");
        exit(1);
    }
    
    while(!feof(fp)){
        char Buf[BUFF_SIZE] = {'\0'};
        int numbytes = fread(Buf, sizeof(char), sizeof(Buf), fp);
        numbytes = write(remoteSocket, Buf, BUFF_SIZE);
    }
}

void recv_file(int remoteSocket,char filename[], int type){

    string dir_name;
    if(type == CLIENT)
        dir_name = "./client_dir/";
    else if(type == SERVER)
        dir_name = "./server_dir/";
    string filename_str = filename;
    string path_name = dir_name + filename_str;

    FILE *fp;
    if ((fp = fopen(path_name.c_str(), "wb")) == NULL){
        perror("fopen error");
        exit(1);
    }
    while(1){
        char Buf[BUFF_SIZE] = {'\0'};
        int numbytes = read(remoteSocket, Buf, BUFF_SIZE);
        if(numbytes <= 0){
            break;
        }
        numbytes = fwrite(Buf, sizeof(char), numbytes, fp);
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
    for(int i = 0; i < MAX_FILENAME; i++)
    {
        if(filename[i] == ' '){
            //printf("more then one file\n");
            return false;
        }
    } 
    if(filename[0] == ' ')
    {
        for(int i = 0; i < MAX_FILENAME; i++)
            filename[i] = filename[i+1];
    }
    return true;      
}