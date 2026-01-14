Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?

> Concurrency is about multiple things (threads) happening in order with the shared result. The threads is said to be "intertwined"
>
> Parallelism is about  multiple things (threads) happening simultaneously, independent of each other.

What is the difference between a *race condition* and a *data race*?

> A race condition is when something (like a variable) changes at an unfortunate timing leading to an action behaving differently than how it was supposed to, like an if statement.
>A data race is like the shared variable in part 3. It happens when multiple threads have access to the same variable, but there is no lock to control access. 

*Very* roughly - what does a *scheduler* do, and how does it do it?

> A scheduler decides when and how long a task should run. Looks at ready tasks -> Saves state of currently running task -> starts the next task. 

### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?

> *Your answer here*

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?

> *Your answer here*

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?

> *Your answer here*

What do you think is best - *shared variables* or *message passing*?

> *Your answer here*
