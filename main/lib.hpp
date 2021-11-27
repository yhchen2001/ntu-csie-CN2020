#define BUFF_SIZE 1024
#define MAX_NAME 8
#define MAX_FILENAME 32
#define EXIST 1
#define NOTEXIST 0
#define MAX_CLIENT 5

const char NAMED = '1';
const char NONAME  = '0';
const char NAMING_OK  = '1';
const char NAMING_FAIL  = '0';
#define REPEAT 1
#define NOREPEAT 0
const char named_msg[BUFF_SIZE] = "1\0";
const char noname_msg[BUFF_SIZE] = "0\0";
const char naming_ok[BUFF_SIZE] = "1\0";
const char naming_fail[BUFF_SIZE] = "0\0";

#define SERVER 1
#define CLIENT 0

const char FILE_EXIST = '1';
const char FILE_NOT_EXIST = '0';

bool check_name(char name_map[][MAX_NAME], char new_name[], int client_socket[], int sd);
void list_file(int remoteSocket);
void set_exist_msg(char filenm[], char exist_msg[], int type);
void send_file(char filename[], int remoteSocket, int type);
void recv_file(int remoteSocket,char filename[], int type);

void recv_ls(int localSocket);
bool file_ok(char filename[]);

