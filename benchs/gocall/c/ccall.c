#include <time.h>
#include <stdio.h>

// call this function to start a nanosecond-resolution timer
struct timespec timer_start(){
    struct timespec start_time;
    clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &start_time);
    return start_time;
}

// call this function to end a timer, returning nanoseconds elapsed as a long
long timer_end(struct timespec start_time){
    struct timespec end_time;
    clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &end_time);
    long diffInNanos = (end_time.tv_sec - start_time.tv_sec) * (long)1e9 + (end_time.tv_nsec - start_time.tv_nsec);
    return diffInNanos;
}

void empty() {}

int main() {
    int i = 0;
    long N = 1000000000;
    struct timespec vartime = timer_start();
    for(i = 0; i < N; i++) {
        empty();
    }
    long time_elapsed_nanos = timer_end(vartime);
    printf("BenchmarkEmptyCCalls\t%ld\t%.2f ns/op\n", N, ((double)time_elapsed_nanos)/N);
}