#include <sys/socket.h>
#include <sys/types.h>
#include <sys/select.h>
#include <arpa/inet.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>


#define IP_listen "0.0.0.0"
#define BUFSIZE 1024
#define PORT_listen 30000



char buffer[BUFSIZE];



int main(void){
    printf("Yo");
    struct sockaddr_in fromWho;

    struct sockaddr_in address;
    socklen_t addrlen = sizeof(address);

    int recvSocket = socket(AF_INET, SOCK_DGRAM,0); 

    // if((recvSocket= Socket(AF_INET, SOCK_DGRAM,0)) < 0){
    //     perror("Socket creation faied!");
    //     exit(EXIT_FAILIURE);
    // }
    
    int rc = bind(recvSocket, &address, addrlen);

    fromWho.sin_family = AF_INET;
    fromWho.sin_port = htons(PORT_listen);
    fromWho.sin_addr.s_addr = INADDR_ANY;


    while(1){
        int n;
        socklen_t len;
        // clear buffer
        buffer[0] ='\0';
        
        n = recvfrom(recvSocket, &buffer, BUFSIZE, MSG_WAITALL,
                    (struct sockaddr *) &fromWho, &len);

        buffer[n] = '\0';
        printf("Message: ");

        for(int i = 0; i < sizeof(buffer); i++){

            printf(buffer[i]);
        } 

        
        sleep(2);
            
    }


    return 0;
}

