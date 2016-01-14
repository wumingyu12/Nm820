#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <termios.h>
#include <errno.h>
#include <limits.h>

#define DATA_LEN                0xFF  

//#define DEBUG                 1


/*********************************************************************************************************
** Function name:           openSerial
** Descriptions:            open serial port at raw mod
** input paramters:         iNum        serial port which can be value at: 1, 2, 3, 4
** output paramters:        NONE
** Return value:            file descriptor
** Create by:               zhuguojun
** Create Data:             2008-05-19
**--------------------------------------------------------------------------------------------------------
** Modified by:
** Modified date:
**--------------------------------------------------------------------------------------------------------
*********************************************************************************************************/
static int openSerial(char *cSerialName)
{
    int iFd;

    struct termios opt; 

    iFd = open(cSerialName, O_RDWR | O_NOCTTY);                        
    if(iFd < 0) {
        perror(cSerialName);
        return -1;
    }

    tcgetattr(iFd, &opt);      

    //cfsetispeed(&opt, B57600);
    //cfsetospeed(&opt, B57600);

     cfsetispeed(&opt, B19200);
     cfsetospeed(&opt, B19200);

    
    /*
     * raw mode
     */
        opt.c_lflag &= ~(ECHO | ICANON | IEXTEN | ISIG);
        opt.c_iflag &= ~(BRKINT | ICRNL | INPCK | ISTRIP | IXON);
        opt.c_oflag     &= ~(OPOST);
        opt.c_cflag     &= ~(CSIZE | PARENB);
        opt.c_cflag     |=  CS8;
    /*
     * 'DATA_LEN' bytes can be read by serial
     */
    opt.c_cc[VMIN]   =   DATA_LEN;                                      
    opt.c_cc[VTIME]  =   150;

    if (tcsetattr(iFd,   TCSANOW,   &opt)<0) {
        return   -1;
    }

    tcflush(iFd,TCIOFLUSH);
    return iFd;
}

int main(void) 
{
    char tmp[1024];
    int len;
    int fd, i;

    //fd = openSerial("/dev/ttySP0");
    fd = openSerial("/dev/ttyO5");
    //char cmd[10]={0x8A,0x9B,0x00,0x01,0x05,0x00,0x00,0x48,0x88,0xfb};
    char cmd[8]={0x00,0x03,0x00,0x01,0x00,0x02,0x94,0x1a};
    
    write(fd, cmd, 8);
    int num = 0;
    char mytmp = 0x00;

    while (1) {
        len = read(fd, tmp, 0x01);
        for(i = 0; i < len; i++)
            mytmp = tmp[i];
            printf("%d--%x", num,mytmp);
            num++;
        printf("\n");
    }
}
