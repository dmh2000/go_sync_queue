#!/bin/bash
echo "==================="
echo "SYNCHRONOUS TESTS  "
echo "==================="
echo "go test -v -run Test.*Sync ."
go test -v -run Test.*Sync . | tee sync.out

echo "==================="
echo "ASYNCHRONOUS TESTS "
echo "==================="
echo "go test -v -run Test.*Async ."
go test -v -run Test.*Async . | tee async.out

echo "==================="
echo "BENCHMARKS         "
echo "==================="
# uses https://github.com/cespare/prettybench
echo "go test -v -run Benchmark.* -bench . -benchmem -memprofile mem.out -cpuprofile cpu.out"
go test -v -run Benchmark.* -bench . -benchmem -memprofile mem.out -cpuprofile cpu.out | prettybench | tee bench.out