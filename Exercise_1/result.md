
# Task 3
## Results C 

* 995166
* 1120653
* 1015499
## Results GO
* 745504
* -64869
* 379192

## Why?
because access to the variable is open to both threads, and is then non-determenistic thereby random results. 
We can relate this to a Data Race. 

# Task 4 Shared variable the proper way

## C
In this case we should use Semaphores. The main pourpose for using semaphores is to controll access to a shared resource(eg. variables) by mutliple threads.
The ownership of semaphores are not exclusive, and thus can be signaled(incremented counter) and waited(decremented counter) by different threads. 

Mutexes on the otherhand have exclusive ownership meaning that the thread that aquired the mutex is the only one that can release it. 




