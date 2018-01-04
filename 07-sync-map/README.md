In Go 1.9 there is one more exported type sync.Map.
sync.Map provides atomic versions of most of the usual map operations plus LoadOrStore method.

First of all — why does this even exist and when should you use it, as you know there was always built-in `map` type.

Optimized to slve cache contention.

If cache contention is not a problem for you, RW mutexes provide better performance and better type safety.

costs:
overhead
type safety
limited API (no Len...)

2 maps:
1 read-only map, and 1 read-write map.

As Bryan Mills says: sync.Map is in the stlib so taht we can use it in stdlib. And it's not the best concurrent-map solution.

Map is a concurrent map with amortized-constant-time loads, stores, and deletes. It is safe for multiple goroutines to call a Map's methods concurrently.