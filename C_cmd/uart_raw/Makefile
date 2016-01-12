EXEC = uart_raw
OBJS = uart_raw.o
SRC  = uart_raw.c

CC = arm-none-linux-gnueabi-gcc
CFLAGS += -O2 -Wall 
LDFLAGS += 

all:$(EXEC)

$(EXEC):$(OBJS)
	$(CC) $(LDFLAGS) -o $@ $(OBJS)

%.o:%.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	@rm -vf $(EXEC) *.o *~
