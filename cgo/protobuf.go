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
		serializedProto, err := CSerializedProtoToGoSerializedProto(&cSerializedProto)
		if err != nil {
			return nil, err
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
func CSerializedProtoToGoSerializedProto(cSerializedProto *C.serialized_proto) (*SerializedProto, error) {
	serializedProto := &SerializedProto{
		TypeUrl: C.GoBytes(unsafe.Pointer(cSerializedProto.type_url), C.int(cSerializedProto.type_url_length)),
		Data:    C.GoBytes(unsafe.Pointer(cSerializedProto.data), C.int(cSerializedProto.data_length)),
	}

	return serializedProto, nil
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

	cSerializedProto, err := CSerializedProtoFromGoProto(value)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}

// TODO_IN_THIS_COMMIT: godoc...
func CSerializedProtoFromGoProto(value gogoproto.Message) (unsafe.Pointer, error) {
	typeURL := []byte(cosmostypes.MsgTypeURL(value))
	proto_bz, err := cdc.Marshal(value)
	if err != nil {
		return nil, err
	}

	cSerializedProto := C.malloc(C.size_t(unsafe.Sizeof(C.serialized_proto{})))
	*(*C.serialized_proto)(cSerializedProto) = C.serialized_proto{
		type_url:        (*C.uint8_t)(C.CBytes(typeURL)),
		type_url_length: C.size_t(len(typeURL)),
		data:            (*C.uint8_t)(C.CBytes(proto_bz)),
		data_length:     C.size_t(len(proto_bz)),
	}

	return cSerializedProto, nil
}

// TODO_IN_THIS_COMMIT: godoc... caller is responsible for freeing.
func CProtoMessageArrayFromGoProtoMessages(msgs []gogoproto.Message) (unsafe.Pointer, error) {
	// Allocate the main structure
	protoMessageArray := (*C.proto_message_array)(C.malloc(C.size_t(unsafe.Sizeof(C.proto_message_array{}))))

	// Set the number of messages
	protoMessageArray.num_messages = C.size_t(len(msgs))

	// Allocate array of serialized_proto structures
	sizeOfSerializedProto := unsafe.Sizeof(C.serialized_proto{})
	messagesArray := (*C.serialized_proto)(C.malloc(C.size_t(sizeOfSerializedProto) * C.size_t(len(msgs))))
	protoMessageArray.messages = messagesArray

	// Populate each message in the array
	for i, msg := range msgs {
		// Calculate pointer to current serialized_proto
		msgMemAddr := uintptr(unsafe.Pointer(messagesArray))
		msgOffset := uintptr(i) * sizeOfSerializedProto
		currentProto := (*C.serialized_proto)(unsafe.Pointer(msgMemAddr + msgOffset))

		// Populate type_url
		typeURL := []byte(cosmostypes.MsgTypeURL(msg))
		currentProto.type_url = (*C.uint8_t)(C.CBytes(typeURL))
		currentProto.type_url_length = C.size_t(len(typeURL))

		// Serialize & populate the message
		msgBz, err := cdc.Marshal(msg)
		if err != nil {
			return nil, err
		}

		currentProto.data = (*C.uint8_t)(C.CBytes(msgBz))
		currentProto.data_length = C.size_t(len(msgBz))
	}

	return unsafe.Pointer(protoMessageArray), nil
}
