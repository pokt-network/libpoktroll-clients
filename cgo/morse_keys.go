package main

/*
#include <memory.h>
#include <protobuf.h>
#include <morse_keys.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/pokt-network/poktroll/x/migration/module/cmd"
	migrationtypes "github.com/pokt-network/poktroll/x/migration/types"
	sharedtypes "github.com/pokt-network/poktroll/x/shared/types"
)

//export LoadMorsePrivateKey
func LoadMorsePrivateKey(morseKeyExportPath, passphrase *C.char, cErr **C.char) C.go_ref {
	privKey, err := cmd.LoadMorsePrivateKey(
		C.GoString(morseKeyExportPath),
		C.GoString(passphrase),
		true,
	)

	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	return SetGoMem(privKey)
}

//export NewSerializedSignedMsgClaimMorseAccount
func NewSerializedSignedMsgClaimMorseAccount(cShannonDestAddr *C.char, privKeyRef C.go_ref, cErr **C.char) unsafe.Pointer {
	morsePrivKey, err := GetGoMem[ed25519.PrivKey](privKeyRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	msg, err := migrationtypes.NewMsgClaimMorseAccount(
		C.GoString(cShannonDestAddr),
		morsePrivKey,
	)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	serializedMsg, err := CSerializedProtoFromGoProto(msg)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return serializedMsg
}

//export NewSerializedSignedMsgClaimMorseApplication
func NewSerializedSignedMsgClaimMorseApplication(
	cShannonDestAddr *C.char,
	privKeyRef C.go_ref,
	serviceId *C.char,
	cErr **C.char,
) unsafe.Pointer {
	morsePrivKey, err := GetGoMem[ed25519.PrivKey](privKeyRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	msg, err := migrationtypes.NewMsgClaimMorseApplication(
		C.GoString(cShannonDestAddr),
		morsePrivKey,
		&sharedtypes.ApplicationServiceConfig{
			ServiceId: C.GoString(serviceId),
		},
	)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	serializedMsg, err := CSerializedProtoFromGoProto(msg)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return serializedMsg
}

//export NewSerializedSignedMsgClaimMorseSupplier
func NewSerializedSignedMsgClaimMorseSupplier(
	cShannonOwnerAddr *C.char,
	cShannonOperatorAddr *C.char,
	privKeyRef C.go_ref,
	cSupplierServiceConfigs *C.proto_message_array,
	cErr **C.char,
) unsafe.Pointer {
	morsePrivKey, err := GetGoMem[ed25519.PrivKey](privKeyRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	supplierServiceConfigsAny, err := CProtoMessageArrayToGoProtoMessages(cSupplierServiceConfigs)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	numSupplierServiceConfigs := int(cSupplierServiceConfigs.num_messages)
	supplierServiceConfigs := make([]*sharedtypes.SupplierServiceConfig, numSupplierServiceConfigs)
	for idx := 0; idx < numSupplierServiceConfigs; idx++ {
		supplierServiceConfigAny := supplierServiceConfigsAny[idx]
		supplierServiceConfig, ok := supplierServiceConfigAny.(*sharedtypes.SupplierServiceConfig)
		if !ok {
			*cErr = C.CString("unable to convert C proto messages to Go slice of *sharedtypes.SupplierServiceConfig")
			return C.NULL
		}

		supplierServiceConfigs[idx] = supplierServiceConfig
	}

	msg, err := migrationtypes.NewMsgClaimMorseSupplier(
		C.GoString(cShannonOwnerAddr),
		C.GoString(cShannonOperatorAddr),
		morsePrivKey,
		supplierServiceConfigs,
	)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	serializedMsg, err := CSerializedProtoFromGoProto(msg)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return serializedMsg
}

//// TODO_IN_THIS_COMMIT: godoc & move...
//type MorseClaimMessage interface {
//	GetMorseSignature() []byte
//	SignMsgClaimMorseAccount(ed25519.PrivKey) error
//}

// TODO_IN_THIS_COMMIT: godoc & move...
//
//export GetMorseAddress
func GetMorseAddress(privKeyRef C.go_ref, cErr **C.char) *C.char {
	privKey, err := GetGoMem[ed25519.PrivKey](privKeyRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.CString("")
	}

	return C.CString(privKey.PubKey().Address().String())
}

// TODO_IN_THIS_COMMIT: godoc... ONLY returns the signature, DOES NOT modify the serialized message. ... caller frees..
//
//export SignMorseClaimMsg
func SignMorseClaimMsg(
	cSerializedProto *C.serialized_proto,
	privKeyRef C.go_ref,
	cOutMorseSignature *C.uint8_t,
	cErr **C.char,
) {
	// Get the private key from the go memory reference.
	privKey, err := GetGoMem[ed25519.PrivKey](privKeyRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return
	}

	// Convert the C serialized message to a Go serialized message.
	serializedProto, err := CSerializedProtoToGoSerializedProto(cSerializedProto)
	if err != nil {
		*cErr = C.CString(err.Error())
		return
	}

	// Deserialize the message to a go protobuf message.
	protoMsg, err := SerializedProtoToProtoMessage(serializedProto)
	if err != nil {
		*cErr = C.CString(err.Error())
		return
	}

	// TODO_NEXT_RELEASE: something better...
	// Generate the signature.
	switch msg := protoMsg.(type) {
	case *migrationtypes.MsgClaimMorseAccount:
		if err = msg.SignMsgClaimMorseAccount(privKey); err != nil {
			*cErr = C.CString(err.Error())
			return
		}
	case *migrationtypes.MsgClaimMorseApplication:
		if err = msg.SignMorseSignature(privKey); err != nil {
			*cErr = C.CString(err.Error())
			return
		}
	case *migrationtypes.MsgClaimMorseSupplier:
		if err = msg.SignMorseSignature(privKey); err != nil {
			*cErr = C.CString(err.Error())
			return
		}
	default:
		*cErr = C.CString(fmt.Sprintf("unexpected message type: %T", protoMsg))
		return
	}

	// Copy the signature to a C array.
	// TODO_IN_THIS_COMMIT: clean up...
	morseSignature := protoMsg.(interface{ GetMorseSignature() []byte }).GetMorseSignature()
	for i := 0; i < C.MORSE_SIGNATURE_SIZE && i < len(morseSignature); i++ {
		// Use unsafe.Pointer arithmetic to access the array elements
		*(*C.uint8_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cOutMorseSignature)) + uintptr(i))) = C.uint8_t(morseSignature[i])
	}
	//:= C.GoBytes(unsafe.Pointer(&morseSignature[0]), C.int(C.MORSE_SIGNATURE_SIZE))

	// Return the signature as a C array.
	// Caller is responsible for freeing the memory.
	return
}
