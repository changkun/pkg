#include <time.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>

int write_all(int fd, void* buffer, size_t length) {
    while (length > 0) {
        int written = write(fd, buffer, length);
        if (written < 0)
            return -1;
        length -= written;
        buffer += written;
    }
    return length;
}

int read_call(int fd, void *buffer, size_t length) {
	return read(fd, buffer, length);
}

struct timespec timer_start(){
    struct timespec start_time;
    clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &start_time);
    return start_time;
}

long timer_end(struct timespec start_time){
    struct timespec end_time;
    clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &end_time);
    long diffInNanos = (end_time.tv_sec - start_time.tv_sec) * (long)1e9 + (end_time.tv_nsec - start_time.tv_nsec);
    return diffInNanos;
}

int main() {
    int i = 0;
    int N = 500000;
    int fds[2];
    char message[14] = "hello, world!\0";
    char buffer[14] = {0};

    socketpair(AF_UNIX, SOCK_STREAM, 0, fds);
    struct timespec vartime = timer_start();
    for(i = 0; i < N; i++) {
        write_all(fds[0], message, sizeof(message));
        read_call(fds[1], buffer, 14);
    }
    long time_elapsed_nanos = timer_end(vartime);
    printf("BenchmarkReadWritePureCCalls\t%d\t%.2ld ns/op\n", N, time_elapsed_nanos/N);
}