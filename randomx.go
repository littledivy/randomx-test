package main

//#cgo CFLAGS: -I./RandomX/src
//#cgo LDFLAGS: -lrandomx -lstdc++ -L./RandomX/ -lm
//#include <stdlib.h>
//#include "randomx.h"
import "C"
import (
    "unsafe"
    "runtime"
    "sync"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    var cache *C.randomx_cache

    flags := C.randomx_get_flags()
    print(flags)

    cache = C.randomx_alloc_cache(flags)
    if cache == nil {
        panic("Failed to allocate cache")
    }

    seed := []byte("0410591dc8b3bba89f949212982f05deeb4a1947e939c62679dfc7610c62")
    key_size := C.size_t(len(seed))
    C.randomx_init_cache(cache, unsafe.Pointer(&seed[0]), key_size)
    
    var dataset *C.randomx_dataset
	dataset = C.randomx_alloc_dataset(flags)
	if dataset == nil {
		panic("Failed to allocate dataset");
	}

    var length uint32
	length = uint32(C.randomx_dataset_item_count())
    
    var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())

    for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (length * i) / workerNum
		b := (length * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
            C.randomx_init_dataset(dataset, cache,  C.ulong(a),  C.ulong(b-a))
		}()
	}
	wg.Wait()

    var vm *C.randomx_vm
    vm = C.randomx_create_vm(flags, cache, dataset);
    if vm == nil {
        panic("Failed to create VM")
    }
    
    output := C.CBytes(make([]byte, C.RANDOMX_HASH_SIZE))
    input := []byte("58249adafb690683a800ee8d6556e2a7d25864d577afbf709ceff9e3bdd5ebae")
    C.randomx_calculate_hash(vm, C.CBytes(input), C.size_t(len(input) + 1), output)

    hash := C.GoBytes(output, C.RANDOMX_HASH_SIZE)
    print(hash)
    C.randomx_destroy_vm(vm)
}
