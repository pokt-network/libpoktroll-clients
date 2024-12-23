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
//export QueryClient_GetSession
func QueryClient_GetSession(
	clientRef C.go_ref,
	appAddress *C.char,
	serviceId *C.char,
	blockHeight C.int64_t,
	cErr **C.char,
) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	session, err := multiClient.GetSession(
		context.TODO(),
		C.GoString(appAddress),
		C.GoString(serviceId),
		int64(blockHeight),
	)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto, err := CSerializedProtoFromGoProto(session)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}
