all: getpwuid_r_hijack.so

getpwuid_r_hijack.o: getpwuid_r_hijack.c
	gcc -fPIC -c getpwuid_r_hijack.c -o getpwuid_r_hijack.o

getpwuid_r_hijack.so: getpwuid_r_hijack.o
	gcc -shared -o getpwuid_r_hijack.so getpwuid_r_hijack.o

clean:
	rm -f *.o *.so
