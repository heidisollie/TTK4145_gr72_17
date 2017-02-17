
#include <pthread.h>
#include <stdio.h>


int i = 0;



void* thread_1(){
	for ( int j = 0; j <= 1000000; j++ ){
		i++;
	}
	return NULL;
}

void* thread_2(){
	for ( int n = 0; n <= 1000000; n++ ){
		i--;
	}
	return NULL;
}


int main(){
	pthread_t thread1, thread2;
	pthread_create(&thread1, NULL, thread_1, NULL);
	pthread_join( thread1, NULL);
	
	pthread_create(&thread2, NULL, thread_2, NULL);
	pthread_join( thread2, NULL);

	printf("%i\n", i);
	return 0;
}

