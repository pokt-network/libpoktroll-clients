package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
// #cgo LDFLAGS: -L${SRCDIR}/callback.c
#include <memory.h>
#include <protobuf.h>
#include <context.h>
#include <callback.h>

static void bridge_success(AsyncOperation *op, void *results) {
    if (op && op->on_success) {
        op->on_success(op->ctx, results);
    }
}

static void bridge_error(AsyncOperation *op, char *err) {
    if (op && op->on_error) {
        op->on_error(op->ctx, err);
    }
}
*/
import "C"
import (
	"context"
	"fmt"
	"strings"
	"unsafe"

	"cosmossdk.io/depinject"
	"github.com/pokt-network/poktroll/pkg/client"
	"github.com/pokt-network/poktroll/pkg/client/tx"
)

// TODO_IMPROVE: add separate constructor which supports options...

// TODO_IN_THIS_COMMIT: godoc...
//
//export NewTxClient
func NewTxClient(depsRef C.go_ref, signingKeyName *C.char, cErr **C.char) C.go_ref {
	// TODO_CONSIDERATION: Could support a version of methods which receive a go context, created elsewhere..
	ctx := context.Background()

	deps, err := GetGoMem[depinject.Config](GoRef(depsRef))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	signingKeyOpt := tx.WithSigningKeyName(C.GoString(signingKeyName))
	txClient, err := tx.NewTxClient(ctx, deps, signingKeyOpt)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	return C.go_ref(SetGoMem(txClient))
}

//export WithSigningKeyName
func WithSigningKeyName(keyName *C.char) C.go_ref {
	return C.go_ref(SetGoMem(tx.WithSigningKeyName(C.GoString(keyName))))
}

// TODO_IN_THIS_COMMIT: godoc...
//
//export TxClient_SignAndBroadcast
func TxClient_SignAndBroadcast(
	op *C.AsyncOperation,
	txClientRef C.go_ref,
	serializedProto *C.serialized_proto,
) C.go_ref {
	goCtx := context.Background()

	txClient, err := GetGoMem[client.TxClient](GoRef(txClientRef))
	if err != nil {
		err = fmt.Errorf("getting tx client ref: %s", err)
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}

	typeUrl := string(C.GoBytes(unsafe.Pointer(serializedProto.type_url), C.int(serializedProto.type_url_length)))
	if !strings.HasPrefix(typeUrl, "/") {
		typeUrl = "/" + typeUrl
	}
	if strings.HasSuffix(typeUrl, string([]byte{0x00})) {
		typeUrl = typeUrl[:len(typeUrl)-1]
	}

	msg, err := interfaceRegistry.Resolve(typeUrl)
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}

	if err = cdc.Unmarshal(C.GoBytes(unsafe.Pointer(serializedProto.data), C.int(serializedProto.data_length)), msg); err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}

	eitherAsyncErr := txClient.SignAndBroadcast(goCtx, msg)
	err, errCh := eitherAsyncErr.SyncOrAsyncError()
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}

	go func() {
		if err = <-errCh; err != nil {
			C.bridge_error(op, C.CString(err.Error()))
		} else {
			C.bridge_success(op, nil)
		}
	}()

	return C.go_ref(SetGoMem(errCh))
}

// TODO_IN_THIS_COMMIT: godoc...
//
//export TxClient_SignAndBroadcastMany
func TxClient_SignAndBroadcastMany(
	op *C.AsyncOperation,
	txClientRef C.go_ref,
	protoMessageArray *C.proto_message_array,
) C.go_ref {
	goCtx := context.Background()

	txClient, err := GetGoMem[client.TxClient](GoRef(txClientRef))
	if err != nil {
		err = fmt.Errorf("getting tx client ref: %s", err)
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}

	if protoMessageArray.num_messages == 0 {
		C.bridge_error(op, C.CString("no messages provided"))
		return C.go_ref(NilGoRef)
	}

	msgs, err := CProtoMessageArrayToGoProtoMessages(protoMessageArray)
	if err != nil {
		err = fmt.Errorf("converting C proto messages to Go: %s", err)
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}

	eitherAsyncErr := txClient.SignAndBroadcast(goCtx, msgs...)
	err, errCh := eitherAsyncErr.SyncOrAsyncError()
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}

	go func() {
		if err = <-errCh; err != nil {
			C.bridge_error(op, C.CString(err.Error()))
		} else {
			C.bridge_success(op, nil)
		}
	}()

	return C.go_ref(SetGoMem(errCh))
}
