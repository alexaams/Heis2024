// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>
#include <semaphore.h>

int i = 0;
sem_t sem;  // create instance of semaphore


// Note the return type: void*
void* incrementingThreadFunction(){
    // TODO: increment i 1_000_000 times
    for(int x = 0; x < 1000; x++){
        // wait
        sem_wait(&sem);
        i++;
        //signal
        sem_post(&sem);
    }
    return NULL;
}

void* decrementingThreadFunction(){
    // TODO: decrement i 1_000_000 times

    for(int x = 0; x < 1000; x++){
        // wait
        sem_wait(&sem);
        i--;
        // signal
        sem_post(&sem);
    }
    return NULL;
}


int main(){
    // TODO: 
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?

    // intialize semaphore  
    sem_init(&sem,0,1);

    pthread_t inc_thread;
    pthread_create(&inc_thread, NULL, incrementingThreadFunction, &i );

    pthread_t dec_thread;
    pthread_create(&dec_thread, NULL, decrementingThreadFunction, &i );

    
    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`    

    pthread_join(inc_thread, NULL);
    pthread_join(dec_thread, NULL);

    // done with semaphore, destroy
    sem_destroy(&sem);
    
    printf("The magic number is: %d\n", i);
    return 0;
}
