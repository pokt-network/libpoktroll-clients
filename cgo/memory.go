package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#include <memory.h>
*/
import "C"
import "fmt"

const (
	// TODO_IN_THIS_COMMIT: godoc...
	NilGoRef = GoRef(-1)
	// TODO_IN_THIS_COMMIT: godoc...
	ZeroGoRef = GoRef(0)
)

var (
	goMemoryMap  = map[GoRef]any{}
	nextGoMemRef = GoRef(0)
)

type GoRef int64

// TODO_IN_THIS_COMMIT: godoc...
func SetGoMem(value any) C.go_ref {
	nextGoMemRef++
	goMemoryMap[nextGoMemRef] = value
	return C.go_ref(nextGoMemRef)
}

// TODO_IN_THIS_COMMIT: godoc...
func GetGoMem[T any](ref C.go_ref) (T, error) {
	value, ok := goMemoryMap[GoRef(ref)]
	if !ok {
		return *new(T), fmt.Errorf("go memory reference not found: %d", ref)
	}

	valueT, ok := value.(T)
	if !ok {
		return valueT, fmt.Errorf("expected %T, got: %T", *new(T), value)
	}

	return valueT, nil
}

// TODO_IN_THIS_COMMIT: godoc...
//
//export FreeGoMem
func FreeGoMem(ref C.go_ref) {
	delete(goMemoryMap, GoRef(ref))
}
