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
	"unsafe"

	"github.com/cosmos/cosmos-sdk/types"
	gogoproto "github.com/cosmos/gogoproto/proto"
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
	goMemoryMapMu.RLock()
	defer goMemoryMapMu.RUnlock()

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

// GetGoRPCProtoAsSerializedProto returns a serialized proto (C struct) corresponding
// to the given Go reference. If the referens is not found or to a non-protobuf type,
// the error string is set and NULL is returned to the C caller
//
//export GetGoProtoAsSerializedProto
func GetGoProtoAsSerializedProto(ref C.go_ref, cErr **C.char) unsafe.Pointer {
	goMemoryMapMu.RLock()
	defer goMemoryMapMu.RUnlock()

	value, ok := goMemoryMap[GoRef(ref)]
	if !ok {
		return C.NULL
	}

	proto_value, ok := value.(gogoproto.Message)
	if !ok {
		*cErr = C.CString(fmt.Sprintf("expected proto value, got: %T", value))
		return C.NULL
	}

	proto_bz, err := cdc.Marshal(proto_value)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto := C.malloc(C.size_t(unsafe.Sizeof(C.serialized_proto{})))
	*(*C.serialized_proto)(cSerializedProto) = C.serialized_proto{
		type_url:        (*C.uint8_t)(C.CBytes([]byte(types.MsgTypeURL(proto_value)))),
		type_url_length: C.size_t(len(types.MsgTypeURL(proto_value))),
		data:            (*C.uint8_t)(C.CBytes(proto_bz)),
		data_length:     C.size_t(len(proto_bz)),
	}

	return unsafe.Pointer(cSerializedProto)
}

// FreeGoMem frees the go-allocated memory associated with the given Go reference.
//
//export FreeGoMem
func FreeGoMem(ref C.go_ref) {
	goMemoryMapMu.Lock()
	defer goMemoryMapMu.Unlock()

	delete(goMemoryMap, GoRef(ref))
}
