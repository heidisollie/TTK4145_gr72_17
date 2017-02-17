
#include <pthread.h>
#include <stdio.h>
#include <unistd.h>


int i = 0;



pthread_mutex_t mutex;



void* thread_1(){
	for ( int j = 0; j <= 100; j++ ){
		pthread_mutex_lock(&mutex);
		i++;
		printf("A: %d\n", i);
		pthread_mutex_unlock(&mutex);

	}
	

	return NULL;
}

void* thread_2(){
	for ( int n = 0; n <= 99; n++ ){
		pthread_mutex_lock(&mutex);
		i--;
		printf("B: %d\n", i);
		pthread_mutex_unlock(&mutex);
	}

	return NULL;
}


int main(){

	if (pthread_mutex_init(&mutex, NULL) != 0){
		printf("\nmutex init failed\n");
		return 1;}
		
	
	pthread_t thread1, thread2;
	pthread_create(&thread1, NULL, thread_1, NULL);
	pthread_create(&thread2, NULL, thread_2, NULL);
	pthread_join( thread1, NULL);
	pthread_join( thread2, NULL);

	sleep(0.01);
	
	printf("%i\n", i);
	
	pthread_mutex_destroy(&mutex);
	return 0;
}

