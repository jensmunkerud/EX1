// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t lockIn;

// Note the return type: void*
void* incrementingThreadFunction(){
	// TODO: increment i 1_000_000 times
	for (int k = 0; k < 1000000; k++) {
		pthread_mutex_lock(&lockIn);
		i++;
		pthread_mutex_unlock(&lockIn);
	}
	return NULL;
}

void* decrementingThreadFunction(){
	// TODO: decrement i 1_000_000 times
	for (int k = 0; k < 999999; k++) {
		pthread_mutex_lock(&lockIn);
		i--;
		pthread_mutex_unlock(&lockIn);
	}
	return NULL;
}


int main(){
	// TODO: 
	// start the two functions as their own threads using `pthread_create`
	// Hint: search the web! Maybe try "pthread_create example"?
	
	pthread_mutex_init(&lockIn, NULL);

	// Creates the two thread variables we will populate in the next two lines with pthread_create
	pthread_t thread1;
	pthread_t thread2;
	
	// Links the two threads to each function (and runs it automatically?)
	// thread structure ptr, function parameters ptr, function to run, function to start the thread..
	int s = pthread_create(&thread1, NULL, incrementingThreadFunction, NULL);
	int t = pthread_create(&thread2, NULL, decrementingThreadFunction, NULL);
	// s & t are return error codes, 0 is success

	// TODO:
	// wait for the two threads to be done before printing the final result
	// Hint: Use `pthread_join`    

	// These two lines pauses the main thread until both thread1 and thread2 is done
	pthread_join(thread1, NULL);
	pthread_join(thread2, NULL);
	
	pthread_mutex_destroy(&lockIn);
	// The value varies, both negatively and positively of order 10^6
	printf("The magic number is: %d\n", i);
	return 0;
}
