package interpreter

import "fmt"

// import "fmt"

var MEMSIZE = 1024 //1kb
var HEAP = make([]byte, MEMSIZE)
var reserved = make([]bool, MEMSIZE)

var pointers = make(map[int]*Pointer, 0)

func malloc(size int, scid string, temp bool) *Pointer {

	cap := 0
	for i := len(HEAP) - 1; i >= 0; i-- {
		if HEAP[i] == 0 && !reserved[i] {
			cap++
		} else {
			cap = 0
		}
		if cap == size {
			//reserve [i:i+size] and return pointer to i
			for j := i; j < i+size; j++ {
				reserved[j] = true
			}
			p := Pointer{i, size, scid, temp}
			pointers[i] = &p
			return &p
		}
	}
	//todo: try to defragment memory and try again

	panic("out of memory")
}

func freePtr(ptr *Pointer) {
	delete(pointers, ptr.address)
	for i := ptr.address; i < ptr.address+ptr.size; i++ {
		reserved[i] = false
		HEAP[i] = 0
	}
}

func freeAll() {
	for _, ptr := range pointers {
		freePtr(ptr)
	}
}

// frees all temp pointers
func gc() {
	for _, ptr := range pointers {
		if ptr.temp {
			freePtr(ptr)
		}
	}
}

func validatePointer(ptr Pointer) {
	if !reserved[ptr.address] {
		panic("invalid pointer " + fmt.Sprint(ptr))
	}
}

func printMemoryStats() {
	rvd := 0
	for _, v := range reserved {
		if v {
			rvd = rvd + 1
		}
	}
	debug_info("Occupied", rvd, "/", MEMSIZE, "bytes")
}
