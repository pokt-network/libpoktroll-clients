package main

/*
#include <memory.h>
#include <protobuf.h>
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	gogoproto "github.com/cosmos/gogoproto/proto"
)

type SerializedProto struct {
	TypeUrl []byte
	Data    []byte
}

// CProtoMessageArrayToGoProtoMessages converts a C proto_message_array into a []proto.Message slice.
func CProtoMessageArrayToGoProtoMessages(cArray *C.proto_message_array) (msgs []gogoproto.Message, err error) {
	if cArray == nil || cArray.messages == nil {
		return nil, fmt.Errorf("invalid proto_message_array: %+v", cArray)
	}

	// Convert C.serialized_protos array to slice.
	cSerializedProtos := GoSliceFromCArray[C.serialized_proto](cArray.messages, int(cArray.num_messages))

	// Convert each message to SerializedProto struct.
	for _, cSerializedProto := range cSerializedProtos {
		serializedProto := &SerializedProto{
			TypeUrl: C.GoBytes(unsafe.Pointer(cSerializedProto.type_url), C.int(cSerializedProto.type_url_length)),
			Data:    C.GoBytes(unsafe.Pointer(cSerializedProto.data), C.int(cSerializedProto.data_length)),
		}

		msg, err := SerializedProtoToProtoMessage(serializedProto)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

// TODO_IN_THIS_COMMIT: move & godoc...
func SerializedProtoToProtoMessage(serializedProto *SerializedProto) (gogoproto.Message, error) {
	typeUrl := string(serializedProto.TypeUrl)
	if !strings.HasPrefix(typeUrl, "/") {
		typeUrl = "/" + typeUrl
	}
	if strings.HasSuffix(typeUrl, string([]byte{0x00})) {
		typeUrl = typeUrl[:len(typeUrl)-1]
	}

	msg, err := interfaceRegistry.Resolve(typeUrl)
	if err != nil {
		return nil, err
	}

	if err = cdc.Unmarshal(serializedProto.Data, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// GoSliceFromCArray converts a C array pointer and length into a Go slice.
// Warning: The caller must ensure the C array remains valid for the lifetime of the returned slice.
func GoSliceFromCArray[T any](cArrayPtr *T, cArrayLen int) []T {
	if cArrayLen < 0 {
		panic("negative length in GoSliceFromCArray")
	}
	if cArrayPtr == nil && cArrayLen > 0 {
		panic("nil pointer with non-zero length in GoSliceFromCArray")
	}

	// Convert to a slice without allocating new memory
	return unsafe.Slice(cArrayPtr, cArrayLen)
}

// GetGoRPCProtoAsSerializedProto returns a serialized proto (C struct) corresponding
// to the given Go reference. If the referens is not found or to a non-protobuf type,
// the error string is set and NULL is returned to the C caller
//
//export GetGoProtoAsSerializedProto
func GetGoProtoAsSerializedProto(ref C.go_ref, cErr **C.char) unsafe.Pointer {
	goMemoryMapMu.RLock()
	defer goMemoryMapMu.RUnlock()

	value, err := GetGoMem[gogoproto.Message](ref)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	proto_bz, err := cdc.Marshal(value)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto := C.malloc(C.size_t(unsafe.Sizeof(C.serialized_proto{})))
	*(*C.serialized_proto)(cSerializedProto) = C.serialized_proto{
		type_url:        (*C.uint8_t)(C.CBytes([]byte(cosmostypes.MsgTypeURL(value)))),
		type_url_length: C.size_t(len(cosmostypes.MsgTypeURL(value))),
		data:            (*C.uint8_t)(C.CBytes(proto_bz)),
		data_length:     C.size_t(len(proto_bz)),
	}

	return unsafe.Pointer(cSerializedProto)
}
