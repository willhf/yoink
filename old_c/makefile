CC = gcc
CFLAGS = -Wall -Wextra -g
TARGET = yoink
SRC = yoink.c 

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC)

clean:
	rm -f $(TARGET)