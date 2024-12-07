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

	"github.com/cosmos/gogoproto/proto"
)

type SerializedProto struct {
	TypeUrl []byte
	Data    []byte
}

// CProtoMessageArrayToGoProtoMessages converts a C proto_message_array into a []proto.Message slice.
func CProtoMessageArrayToGoProtoMessages(cArray *C.proto_message_array) (msgs []proto.Message, err error) {
	if cArray == nil || cArray.messages == nil {
		return nil, fmt.Errorf("invalid proto_message_array: %+v", cArray)
	}

	// Convert C.serialized_protos array to slice.
	cSerializedProtos := GoSliceFromCArray[C.serialized_proto](cArray.messages, int(cArray.num_messages))

	// Convert each message to SerializedProto struct.
	for _, cSerializedProto := range cSerializedProtos {
		//fmt.Printf(">>> msg.type_url: %s\n", C.GoBytes(unsafe.Pointer(cSerializedProto.type_url), C.int(cSerializedProto.type_url_length)))
		//fmt.Printf(">>> msg.type_url: %x\n", C.GoBytes(unsafe.Pointer(cSerializedProto.type_url), C.int(cSerializedProto.type_url_length)))
		//fmt.Printf(">>> msg.data: %s\n", C.GoBytes(unsafe.Pointer(cSerializedProto.data), C.int(cSerializedProto.data_length)))
		//fmt.Printf(">>> msg.type_url_length: %d\n", cSerializedProto.type_url_length)
		//typeUrlBytes := C.GoBytes(unsafe.Pointer(cSerializedProto.type_url), C.int(cSerializedProto.type_url_length))
		//fmt.Printf(">>> actual bytes read: %d\n", len(typeUrlBytes))

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
func SerializedProtoToProtoMessage(serializedProto *SerializedProto) (proto.Message, error) {
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
