package main

/*
#include <memory.h>
#include <context.h>
#include <protobuf.h>
*/
import "C"
import (
	"context"
	"unsafe"

	"cosmossdk.io/depinject"
)

// NewQueryClient constructs a new MultiQueryClient and returns its Go reference to the C caller.
//
//export NewQueryClient
func NewQueryClient(depsRef C.go_ref, queryNodeRPCURL *C.char, cErr **C.char) C.go_ref {
	deps, err := GetGoMem[depinject.Config](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	multiClient, err := NewMultiQueryClient(deps, C.GoString(queryNodeRPCURL))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	return SetGoMem(multiClient)
}

// QueryClient_GetSessionParams queries the chain for the current session module parameters.
//
//export QueryClient_GetSessionParams
func QueryClient_GetSessionParams(depsRef C.go_ref, cErr **C.char) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	sessionParams, err := multiClient.GetSessionParams(context.TODO())
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto, err := CSerializedProtoFromGoProto(sessionParams)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}

// QueryClient_GetProofParams queries the chain for the current proof module parameters.
//
//export QueryClient_GetProofParams
func QueryClient_GetProofParams(depsRef C.go_ref, cErr **C.char) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	proofParams, err := multiClient.GetProofParams(context.TODO())
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	cSerializedProto, err := CSerializedProtoFromGoProto(proofParams)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}

/* TODO_BLOCKED(@bryanchriswhite, #543): uncomment & implement once dependencies are available.

//export QueryClient_GetServiceParams
func QueryClient_GetServiceParams(depsRef C.go_ref, op *C.AsyncOperation) {}

//export QueryClient_GetApplicationParams
func QueryClient_GetApplicationParams(depsRef C.go_ref, op *C.AsyncOperation) {}

//export QueryClient_GetGetewayParams
func QueryClient_GetGetewayParams(depsRef C.go_ref, op *C.AsyncOperation) {}

//export QueryClient_GetSupplierParams
func QueryClient_GetSupplierParams(depsRef C.go_ref, op *C.AsyncOperation) {}

//export QueryClient_GetServiceParams
func QueryClient_GetServiceParams(depsRef C.go_ref, op *C.AsyncOperation) {}

//export QueryClient_GetTokenomicsParams
func QueryClient_GetTokenomicsParams(depsRef C.go_ref, op *C.AsyncOperation) {}

*/
