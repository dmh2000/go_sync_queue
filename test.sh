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
echo "go test -v -run Benchmark.* -bench . -benchmem -memprofile mem.out -cpuprofile cpu.out"
go test -v -run Benchmark.* -bench . -benchmem -memprofile mem.out -cpuprofile cpu.out | tee

echo "==================="
echo "RACE               "
echo "==================="
echo "go test -v -race *.go" 
go test -v -race *.go | tee race.out

echo "==================="
echo "PROFILE MEM      "
echo "==================="
echo "go tool pprof -text  mem.out" 
go tool pprof -top  -nodecount=15 mem.out | tee mem_prof.out

echo "==================="
echo "PROFILE CPU      "
echo "==================="
echo "go tool pprof -text  cpu.out" 
go tool pprof -top -nodecount=20  cpu.out | tee cpu_prof.out