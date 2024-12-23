package main

///*
//#include <memory.h>
//*/
//import "C"
//import (
//	"context"
//	"unsafe"
//
//	gogoproto "github.com/cosmos/gogoproto/proto"
//)
//
//// TODO_IN_THIS_COMMIT: godoc...
////
////export QueryClient_GetGateway(
//func QueryClient_GetGateawy(
//	clientRef C.go_ref,
//	gatewayAddress *C.char,
//	cErr **C.char,
//) unsafe.Pointer {
//	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	gateway, err := multiClient.GetGateway(context.TODO(), C.GoString(gatewayAddress))
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	cSerializedProto, err := CSerializedProtoFromGoProto(&gateway)
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	return cSerializedProto
//}
//
//// TODO_IN_THIS_COMMIT: godoc...
////
////export QueryClient_GetAllGateways
//func QueryClient_GetAllGateways(
//	clientRef C.go_ref,
//	cErr **C.char,
//) unsafe.Pointer {
//	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	gateways, err := multiClient.GetAllGateways(context.TODO())
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	gatewayPtrs := make([]gogoproto.Message, len(gateways))
//	for i, gateway := range gateways {
//		gatewayPtrs[i] = &gateway
//	}
//
//	cProtoMessages, err := CProtoMessageArrayFromGoProtoMessages(gatewayPtrs)
//	if err != nil {
//		*cErr = C.CString(err.Error())
//		return C.NULL
//	}
//
//	return cProtoMessages
//}
