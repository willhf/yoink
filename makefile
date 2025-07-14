CC = gcc
CFLAGS = -Wall -Wextra -O2
TARGET = yoink
SRC = yoink.c

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC)

clean:
	rm -f $(TARGET)