package main

/*
#include <memory.h>
*/
import "C"
import (
	"context"
	"unsafe"
)

// TODO_IN_THIS_COMMIT: godoc...
//
//export QueryClient_GetService
func QueryClient_GetService(
	clientRef C.go_ref,
	serviceId *C.char,
	cErr **C.char,
) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	service, err := multiClient.GetService(context.TODO(), C.GoString(serviceId))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto, err := CSerializedProtoFromGoProto(&service)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}

// TODO_IN_THIS_COMMIT: godoc...
//
//export QueryClient_GetServiceRelayDifficulty
func QueryClient_GetServiceRelayDifficulty(
	clientRef C.go_ref,
	serviceId *C.char,
	cErr **C.char,
) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	relayDifficulty, err := multiClient.GetServiceRelayDifficulty(context.TODO(), C.GoString(serviceId))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto, err := CSerializedProtoFromGoProto(&relayDifficulty)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}
