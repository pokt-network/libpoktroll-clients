package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#include <memory.h>
#include <protobuf.h>
*/
import "C"
import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	// NilGoRef is a Go reference which is used to indicate a nil/NULL reference.
	NilGoRef = GoRef(-1)
)

var (
	// goMemoryMapMu is a mutex used to protect goMemoryMap during concurrent usage.
	goMemoryMapMu sync.RWMutex
	// goMemoryMap is a map of Go references to values. The map is used to store
	// values which are allocated in Go and which need to be passed back to C.
	goMemoryMap = map[GoRef]any{}
	// nextGoMemRef is the next Go reference to be allocated. It is incremented
	// each time a new reference is allocated. IT is an atomic int64 to guarantee
	// refernce uniqueness in the presence of concurrent usage.
	nextGoMemRef atomic.Int64
)

// GoRef is a type alias for a monotonically increasing integer which is used as
// a key ("reference") to a value stored in the goMemoryMap.
type GoRef int64

// SetGoMem stores the given value in the goMemoryMap, protecting it from garbage
// collection, and returns the "reference" which is that value's map key. Reference
// keys are monotically increasing integers, starting at 1. -1 and/or 0 MAY be used
// to indicate a nil/NULL reference.
func SetGoMem(value any) C.go_ref {
	nextRef := GoRef(nextGoMemRef.Add(1))

	goMemoryMapMu.Lock()
	defer goMemoryMapMu.Unlock()

	goMemoryMap[nextRef] = value
	return C.go_ref(nextRef)
}

// GetGoMem returns the value associated with the given Go reference.
func GetGoMem[T any](ref C.go_ref) (T, error) {
	var zeroT T

	goMemoryMapMu.RLock()
	defer goMemoryMapMu.RUnlock()

	value, ok := goMemoryMap[GoRef(ref)]
	if !ok {
		return *new(T), fmt.Errorf("go memory reference not found: %d", ref)
	}

	valueT, ok := value.(T)
	if !ok {
		return valueT, fmt.Errorf("expected %T, got: %T", zeroT, value)
	}

	return valueT, nil
}

// FreeGoMem frees the go-allocated memory associated with the given Go reference.
//
//export FreeGoMem
func FreeGoMem(ref C.go_ref) {
	goMemoryMapMu.Lock()
	defer goMemoryMapMu.Unlock()

	delete(goMemoryMap, GoRef(ref))
}
