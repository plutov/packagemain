How many goroutines can you create before Go gives up?

1,000? 100,000?

A million?

Ten million?

The surprising answer is that Go doesn't really give you a nice number.

Because a goroutine isn't a thread.

And once you understand what is underneath a goroutine — an execution stack, a runtime scheduler, an OS thread, a P, queues, preemption, memory, and the garbage collector — the whole “goroutines are lightweight” story becomes much more interesting.

So instead of writing another go func() {} tutorial, let's take Go apart.

We'll create goroutines until something breaks, inspect what the runtime is doing, compare them with real OS threads, and find out where the millions-of-goroutines story finally hits reality.

### Resources

- https://go.dev/src/runtime/HACKING
- https://golang.design/under-the-hood/en/part4memory/ch14stack/grow/

