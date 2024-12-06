package main

// #include <client.h>
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/cosmos/gogoproto/proto"
)

type SerializedProto struct {
	TypeUrl string
	Data    []byte
}

// CProtoMessageArrayToBytes converts a C proto_message_array into a [][]byte slice.
func CProtoMessageArrayToBytes(cArray *C.proto_message_array) [][]byte {
	if cArray == nil || cArray.messages == nil {
		return nil
	}

	// Create slice of byte slices.
	msgsBz := make([][]byte, cArray.num_messages)

	// Convert messages array to a byte slice.
	msgBz := (*[1 << 30]*C.serialized_proto)(unsafe.Pointer(cArray.messages))[:cArray.num_messages:cArray.num_messages]

	// Convert each message to []byte
	for i := uint64(0); i < uint64(cArray.num_messages); i++ {
		msg := msgBz[i]
		if msg != nil && msg.data != nil {
			msgsBz[i] = C.GoBytes(unsafe.Pointer(msg.data), C.int(msg.length))
		}
	}

	return msgsBz
}

// CProtoMessageArrayToGoProtoMessages converts a C proto_message_array into a []proto.Message slice.
func CProtoMessageArrayToGoProtoMessages(cArray *C.proto_message_array) (msgs []proto.Message, err error) {
	if cArray == nil || cArray.messages == nil {
		return nil, fmt.Errorf("invalid proto_message_array: %+v", cArray)
	}

	// Convert C.serialized_protos array to slice.
	cSerializedProtos := GoSliceFromCArray[C.serialized_proto, C.serialized_proto](cArray.messages, int(cArray.num_messages))

	// Convert each message to SerializedProto struct.
	for _, cSerializedProto := range cSerializedProtos {
		if cSerializedProto != nil &&
			cSerializedProto.data != nil {
			serializedProto := &SerializedProto{
				TypeUrl: C.GoString(cSerializedProto.type_url),
				Data:    C.GoBytes(unsafe.Pointer(cSerializedProto.data), C.int(cSerializedProto.length)),
			}

			msg, err := SerializedProtoToProtoMessage(serializedProto)
			if err != nil {
				return nil, err
			}

			msgs = append(msgs, msg)
		}
	}

	return msgs, nil
}

// TODO_IN_THIS_COMMIT: move & godoc...
func SerializedProtoToProtoMessage(serializedProto *SerializedProto) (proto.Message, error) {
	typeUrl := serializedProto.TypeUrl
	if !strings.HasPrefix(typeUrl, "/") {
		typeUrl = "/" + typeUrl
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

// TODO_IN_THIS_COMMIT: move & godoc...
// DEV_NOTE: ONLY USE WITH C TYPES (generic type arguments).
func GoSliceFromCArray[D, S any](cArrayPtr *S, cArrayLen int) []*D {
	// TODO_IN_THIS_COMMIT: add a  DEV_NOTE.
	return (*[1 << 30]*D)(unsafe.Pointer(cArrayPtr))[:cArrayLen:cArrayLen]
}
