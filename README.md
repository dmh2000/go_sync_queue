Bounded Queue : Channels vs Condition Variable

#### test

<pre>
$ go test -v .
=== RUN   TestChannel
--- PASS: TestChannel (0.00s)
=== RUN   TestList
--- PASS: TestList (0.00s)
=== RUN   TestCircular
--- PASS: TestCircular (0.00s)
=== RUN   TestQueueNative
--- PASS: TestQueueNative (0.00s)
PASS
ok      dmh2000.xyz/queue       (cached)
</pre>

#### benchmark

<pre>
$ go test -v -bench .
=== RUN   TestChannel
--- PASS: TestChannel (0.00s)
=== RUN   TestList
--- PASS: TestList (0.00s)
=== RUN   TestCircular
--- PASS: TestCircular (0.00s)
=== RUN   TestQueueNative
--- PASS: TestQueueNative (0.00s)
goos: linux
goarch: amd64
pkg: dmh2000.xyz/queue
BenchmarkQueueChannel
BenchmarkQueueChannel-4          2176928               547 ns/op
BenchmarkList
BenchmarkList-4                  1399761               854 ns/op
BenchmarkCircular
BenchmarkCircular-4              2579527               440 ns/op
BenchmarkQueueNative
BenchmarkQueueNative-4           3204572               373 ns/op
PASS
ok      dmh2000.xyz/queue       7.005s
</pre>

#### benchmark memory

<pre>
$ go test -v -bench . -benchmem -memprofile mem.out
=== RUN   TestChannel
--- PASS: TestChannel (0.00s)
=== RUN   TestList
--- PASS: TestList (0.00s)
=== RUN   TestCircular
--- PASS: TestCircular (0.00s)
=== RUN   TestQueueNative
--- PASS: TestQueueNative (0.00s)
goos: linux
goarch: amd64
pkg: dmh2000.xyz/queue
cpu: Intel(R) Core(TM) i5-3470 CPU @ 3.20GHz
BenchmarkQueueChannel
BenchmarkQueueChannel-4          2203729               525.5 ns/op             0 B/op          0 allocs/op
BenchmarkList
BenchmarkList-4                  1384489               857.8 ns/op           384 B/op          8 allocs/op
BenchmarkCircular
BenchmarkCircular-4              2714256               439.2 ns/op             0 B/op          0 allocs/op
BenchmarkQueueNative
BenchmarkQueueNative-4           3204912               427.2 ns/op             0 B/op          0 allocs/op
PASS
ok      dmh2000.xyz/queue       7.174s
$ go tool pprof mem.out
File: queue.test
Type: alloc_space
Time: Mar 10, 2021 at 10:04am (PST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 837.04MB, 99.88% of 838.05MB total
Dropped 15 nodes (cum <= 4.19MB)
      flat  flat%   sum%        cum   cum%
  837.04MB 99.88% 99.88%   837.04MB 99.88%  container/list.(*List).insertValue (inline)
         0     0% 99.88%   837.04MB 99.88%  container/list.(*List).PushBack (inline)
         0     0% 99.88%   837.04MB 99.88%  dmh2000.xyz/queue.(*List).Put
         0     0% 99.88%   837.04MB 99.88%  dmh2000.xyz/queue.BenchmarkList
         0     0% 99.88%   837.04MB 99.88%  dmh2000.xyz/queue.b1
         0     0% 99.88%   837.04MB 99.88%  testing.(*B).launch
         0     0% 99.88%   837.55MB 99.94%  testing.(*B).runN
(pprof)
</pre>

#### CPU profile

<pre>
$ go test -bench .  -cpuprofile cpu.out
goos: linux
goarch: amd64
pkg: dmh2000.xyz/queue
cpu: Intel(R) Core(TM) i5-3470 CPU @ 3.20GHz
BenchmarkQueueChannel-4          2238103               530.7 ns/op
BenchmarkList-4                  1000000              1045 ns/op
BenchmarkCircular-4              2707912               441.6 ns/op
BenchmarkQueueNative-4           3199048               372.1 ns/op
PASS
ok      dmh2000.xyz/queue       6.122s
$ go tool pprof cpu.out
File: queue.test
Type: cpu
Time: Mar 10, 2021 at 10:03am (PST)
Duration: 6.12s, Total samples = 6.11s (99.90%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 4410ms, 72.18% of 6110ms total
Dropped 66 nodes (cum <= 30.55ms)
Showing top 10 nodes out of 42
      flat  flat%   sum%        cum   cum%
     950ms 15.55% 15.55%      950ms 15.55%  sync.(*Mutex).Lock
     920ms 15.06% 30.61%      920ms 15.06%  sync.(*Mutex).Unlock
     460ms  7.53% 38.13%      460ms  7.53%  runtime.unlock2
     400ms  6.55% 44.68%     4470ms 73.16%  dmh2000.xyz/queue.b1
     390ms  6.38% 51.06%      780ms 12.77%  dmh2000.xyz/queue.(*NativeInt).Put
     370ms  6.06% 57.12%      370ms  6.06%  runtime.lock2
     320ms  5.24% 62.36%      730ms 11.95%  dmh2000.xyz/queue.(*NativeInt).Get
     260ms  4.26% 66.61%      630ms 10.31%  dmh2000.xyz/queue.(*Circular).Get
     210ms  3.44% 70.05%      710ms 11.62%  dmh2000.xyz/queue.(*Circular).Put
     130ms  2.13% 72.18%      420ms  6.87%  runtime.mallocgc
(pprof)
</pre>

#### escape analysis

- elided : ./queue_channel.go:29:19: &errors.errorString literal escapes to heap

<pre>
# redirect stderr to stdout
  2>&1 go build -gcflags="-m"  | grep escapes
./queue_circular.go:119:18: make([]interface {}, size, size) escapes to heap
./queue_circular.go:125:24: &sync.Cond literal escapes to heap
./queue_list.go:36:20: &list.Element literal escapes to heap
./queue_list.go:111:22: new(list.List) escapes to heap
./queue_list.go:113:24: &sync.Cond literal escapes to heap
./queue_native.go:119:18: make([]int, size, size) escapes to heap
./queue_native.go:125:24: &sync.Cond literal escapes to heap
</pre>
