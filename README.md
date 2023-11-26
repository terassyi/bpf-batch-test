# Comparison of BPF_MAP_LOOKUP_ELEM and BPF_MAP_LOOKUP_BATCH in go

This repository compares the execution time and memory usage between `BPF_MAP_LOOKUP_ELEM` and `BPF_MAP_LOOKUP_BATCH` in golang.

I use [cilium/ebpf](https://pkg.go.dev/github.com/cilium/ebpf).

- [Map.Iterate](https://pkg.go.dev/github.com/cilium/ebpf#Map.Iterate)
  - [MapIterator.Next](https://pkg.go.dev/github.com/cilium/ebpf#MapIterator.Next)
- [Map.BatchLookup](https://pkg.go.dev/github.com/cilium/ebpf#Map.BatchLookup)

## How 

For test data, I define the map like below.

```c
struct
{
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(key_size, sizeof(u32));
	__uint(value_size, sizeof(u32));
	__uint(max_entries, 1024 * 1024 * 100);
} test_map SEC(".maps");
```

Before measuring, fill this map.(in [load/main.go](load/main.go))

```go
	v := uint32(0xdeadbeef)
	for i := 0; i < bpf.MAX_ENTRIES; i++ {
		k := uint32(i)
		if err := m.Update(&k, &v, ebpf.UpdateAny); err != nil {
			panic(err)
		}
	}
```

To compare, fetch all entries of `test_map` with each way.

[normal/main.go](normal/main.go) uses `bpf_map_lookup_elem`.

```go
	var k, v uint32
	count := 0

	iter := m.Iterate()
	for iter.Next(&k, &v) {
		count++
	}
```

[batch/main.go](batch/main.go) uses `bpf_map_lookup_batch`.

```go
	batchCount := 0
	count := 0

	max := m.MaxEntries()
	const chunk uint32 = 4096
	chSize := int(max / chunk)

	kout := make([]uint32, chunk)
	vout := make([]uint32, chunk)
	var k uint32
	var prev uint32

	for i := 0; i < chSize; i++ {
		c, err := m.BatchLookup(prev, &k, kout, vout, nil)
		if err != nil {
			panic(err)
		}
		count += c
		batchCount++
	}
```

## Result

Build sample code.

```console
$ make clean && make build
rm bpf/include/vmlinux.h
rm bin/*
rm pkg/bpf/xdpprog*
bpftool btf dump file /sys/kernel/btf/vmlinux format c > bpf/include/vmlinux.h
go generate ./...
Compiled /home/terassyi/workspace/bpf-batch-test/pkg/bpf/xdpprog_bpfeb.o
Stripped /home/terassyi/workspace/bpf-batch-test/pkg/bpf/xdpprog_bpfeb.o
Wrote /home/terassyi/workspace/bpf-batch-test/pkg/bpf/xdpprog_bpfeb.go
Compiled /home/terassyi/workspace/bpf-batch-test/pkg/bpf/xdpprog_bpfel.o
Stripped /home/terassyi/workspace/bpf-batch-test/pkg/bpf/xdpprog_bpfel.o
Wrote /home/terassyi/workspace/bpf-batch-test/pkg/bpf/xdpprog_bpfel.go
go build -o bin/batch batch/main.go
go build -o bin/normal normal/main.go
go build -o bin/load load/main.go
```

Setup test data.

```console
$ sudo bin/load
Load test program and map
map id is  668
Prepare test data...
Finish preparing test data
```

Run two programs.
(If you fail to run one program after running the other one, please stop `bin/load` and run again.)

`bpf_map_lookup_elem`

```console
$ multitime sudo bin/normal -mem
get from id:  678
Start to iterate bpf map
----- init -----
Alloc = 374176
HeapAlloc = 374176
TotalAlloc = 376120
Sys = 12254224
NumGC = 1
----- last -----
Alloc = 3614608
HeapAlloc = 3614608
TotalAlloc = 5453024352
Sys = 17497104
NumGC = 1510
------------------------
Finish iterating bpf map
Count 104857600 entries
===> multitime results
1: sudo bin/normal -mem
            Mean        Std.Dev.    Min         Median      Max
real        88.923      0.000       88.923      88.923      88.923
user        0.009       0.000       0.009       0.009       0.009
sys         0.004       0.000       0.004       0.004       0.004
```

`bpf_map_lookup_batch`

```console
$ multitime sudo bin/batch -mem
get from id:  754
Start batch look up bpf map
----- init -----
Alloc = 374168
HeapAlloc = 374168
TotalAlloc = 376112
Sys = 12254224
NumGC = 1
----- last -----
Alloc = 3498736
HeapAlloc = 3498736
TotalAlloc = 840618952
Sys = 17300496
NumGC = 235
------------------------
Finish batch look up bpf map
Count 104857600 entries
Batch Count is 25600
===> multitime results
1: sudo bin/batch -mem
            Mean        Std.Dev.    Min         Median      Max
real        6.973       0.000       6.973       6.973       6.973
user        0.002       0.000       0.002       0.002       0.002
sys         0.004       0.000       0.004       0.004       0.004
```
