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
//export QueryClient_GetSupplier
func QueryClient_GetSupplier(
	clientRef C.go_ref,
	supplierAddress *C.char,
	cErr **C.char,
) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	supplier, err := multiClient.GetSupplier(context.TODO(), C.GoString(supplierAddress))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto, err := CSerializedProtoFromGoProto(&supplier)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}

//// TODO_IN_THIS_COMMIT: godoc...
////
////export QueryClient_GetAllSuppliers
//func QueryClient_GetAllSuppliers(
//	clientRef C.go_ref,
//	cErr **C.char,
//) unsafe.Pointer {
//	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	suppliers, err := multiClient.GetAllSuppliers(context.TODO())
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	supplierPtrs := make([]gogoproto.Message, len(suppliers))
//	for i, supplier := range suppliers {
//		supplierPtrs[i] = &supplier
//	}
//
//	cProtoMessages, err := CProtoMessageArrayFromGoProtoMessages(supplierPtrs)
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	return cProtoMessages
//}
