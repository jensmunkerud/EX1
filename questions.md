Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?

> Concurrency is about multiple things (threads) happening in order with the shared result. The threads is said to be "intertwined"
>
> Parallelism is about  multiple things (threads) happening simultaneously, independent of each other.

What is the difference between a *race condition* and a *data race*?

> A race condition is when something (like a variable) changes at an unfortunate timing leading to an action behaving differently than how it was supposed to, like an if statement.
> A data race is like the shared variable in part 3. It happens when multiple threads have access to the same variable, but there is no lock to control access.

*Very* roughly - what does a *scheduler* do, and how does it do it?

> A scheduler decides when and how long a task should run. Looks at ready tasks -> Saves state of currently running task -> starts the next task.

### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?

> Multiple threads can increase performance, making a problem faster to solve. It can lower the amount of time a system needs to perform a task, useful for realtime systems. It can have lower overhead.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?

> Coroutines allow resources to be shared between multiple "lighter weight" threads. It suspends a resources rather than blocking it.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?

> Both. It adds complexity to problems through dependencies. We may be able to solve problems we otherwise could not. It will usually solve a problem faster, by handling two problems at the same time.

What do you think is best - *shared variables* or *message passing*?

> It depends. Message passing might be simpler to understand and maintain. Shared variables might have less overhead for simpler tasks, but may be less scalable than message passing.
