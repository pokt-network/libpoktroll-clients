package main

/*
#include <memory.h>
*/
import "C"
import (
	"context"
	"unsafe"

	gogoproto "github.com/cosmos/gogoproto/proto"
)

//export QueryClient_GetApplication
func QueryClient_GetApplication(
	clientRef C.go_ref,
	appAddress *C.char,
	cErr **C.char,
) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	app, err := multiClient.GetApplication(context.TODO(), C.GoString(appAddress))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto, err := CSerializedProtoFromGoProto(&app)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}

//export QueryClient_GetAllApplications
func QueryClient_GetAllApplications(
	clientRef C.go_ref,
	cErr **C.char,
) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	apps, err := multiClient.GetAllApplications(context.TODO())
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	appPtrs := make([]gogoproto.Message, len(apps))
	for i, app := range apps {
		appPtrs[i] = &app
	}

	cProtoMessages, err := CSerializedProtoArrayFromGoProtoMessages(appPtrs)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cProtoMessages
}
