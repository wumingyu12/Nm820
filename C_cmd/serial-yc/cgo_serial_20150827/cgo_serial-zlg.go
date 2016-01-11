//短接扬创串口1，J1的4,5脚可以看到自发自收

package main

import (
	"fmt"
	//"reflect"
	"unsafe"
)

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <termios.h>
#include <errno.h>
#include <limits.h>
#include <asm/ioctls.h>
#include <time.h>
#include <pthread.h>

#define DATA_LEN                0xFF

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


     cfsetispeed(&opt, B4800);
     cfsetospeed(&opt, B4800);

    opt.c_lflag   &=   ~(ECHO   |   ICANON   |   IEXTEN   |   ISIG);
    opt.c_iflag   &=   ~(BRKINT   |   ICRNL   |   INPCK   |   ISTRIP   |   IXON);
    opt.c_oflag   &=   ~(OPOST);
    opt.c_cflag   &=   ~(CSIZE   |   PARENB);
    opt.c_cflag   |=   CS8;

    opt.c_cc[VMIN]   =   DATA_LEN;
    opt.c_cc[VTIME]  =   150;

    if (tcsetattr(iFd,   TCSANOW,   &opt)<0) {
        return   -1;
    }


    return iFd;
}

char* sss(void)
{
    char cmd[10]={0x8A,0x9B,0x00,0x01,0x05,0x02,0x00,0x09,0x00,0x36};
    char tmp[1024];
    int len;
    int fd, i;

    fd = openSerial("/dev/ttyO2");

    write(fd, cmd, 10);

    len = read(fd, tmp, 0x64);
    printf("%x\n",*tmp);
    for(i = 0; i < len; i++){
        printf("%d--",i);
        printf("%x",tmp[i]);
        printf("\n");
    }
    return &tmp;

}

*/
import "C"

func main() {

	//var t = make([]byte, 100)

	var p *C.char = C.sss()
	/*
		    hdr := reflect.SliceHeader{
				Data: uintptr(unsafe.Pointer(p)),
				Len:  8,
				Cap:  8,
			}
	*/
	gobyte := C.GoBytes(unsafe.Pointer(p), C.int(100))
	//t[0] = *(unsafe.Pointer(C.sss()))
	fmt.Printf("t:%x", gobyte)

}
