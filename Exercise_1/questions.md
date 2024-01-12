Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> *Concurrency is when multiple tasks starts, runs and complete in overlapping time-periods, and parallelism is when multiple tasks executes at the same time (in parallell - this requires hardware that is able to run several tasks at the same time). In the concurrency-case the tasks never runs at the exact same time, but the CPU executes several tasks at the same time by switching between the tasks. Concurrency is an alternative to parallelism if you don't have hardware to execute several tasks at the same time.*

What is the difference between a *race condition* and a *data race*? 
> *When a task is dependent on the sequence or timing of other events it is called a race condition. A data race occurs when: two or more threads in a single process access the same memory location concurrently, and at least one of the accesses is for writing, and the threads are not using any exclusive locks to control their accesses to that memory.* 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> *A component in an operative system that "plans" the tasks running and is responsible for optimizing the use of system recourses by administrating the time-use of the CPU for the different tasks. It will receive tasks-request in a queue and choose which one that is to be performed next based on different criterias for said task. It may decide based on priority or by the use of algorithms.* 


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> *Let's use the elevator-project as an example. Threads will allow us to receive a request from the second floor while still processing a former request from first floor. Threading will make the elevators more efficient, especially if we are running multiple elevators. One elevator can handle one request while one or more of the other elevators process other requests. We can also handle unforseen problems better with threads. If one elevator crashes for instance.*

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> *Green threads are administrated on a user-level in a runtime-environment or a programming-language in opposition to directly by the OS. They are "simpler" to manage and yields faster performance and reduced latency. Since regular threads consumes more recourses on the CPU, green threads is preferred if possible. Also preferred because of simpler management. In high-level programming the green threads takes care of the distribution between them, so this is not something the programmer has to worry about.*

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> *Depending on the task, but the programming can easily get more complex for "small" tasks and you have to pay more attention to details to avoid deadlocks etc. But some complex tasks may also become easier. Dynamic code may become easier, and when the problem is large, for instance, multiple elevators working together, you can significantly shorten the code.*

What do you think is best - *shared variables* or *message passing*?
> *In the elevator-project we think message passing would be the best solution as there are no main server. All the elevators have their own CPU, and message passing allows for a secure way of communicating wihtout sharing a global memory. It will also make the system more scalable and modular which is a criteria for the project.*

7: Thinking about elevators
---------------------------

The main problem of the project is to ensure that no orders are lost. 
 - What sub-problems do you think this consists of?

 **
 - What will you have to make in order to solve these problems?

Maybe try thinking about the happy case of the system:
 - If we push the button one place, how do we make (preferably only) one elevator start moving?
 - Once an elevator arrives, how do we inform the others that it is safe to clear that order?

Maybe try thinking about the worst-case (http://xkcd.com/748/) behavior of the system:
 - What if the software controlling one of the elevators suddenly crashes?
 - What if it doesn't crash, but hangs?
 - What if a message between machines is lost?
 - What if the network cable is suddenly disconnected? Then re-connected?
 - What if the elevator car never arrives at its destination?

