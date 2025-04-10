package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
// #cgo LDFLAGS: -L${SRCDIR}/callback.c
#include <memory.h>
#include <protobuf.h>
#include <context.h>
#include <callback.h>
#include <string.h>
#include <gas.h>

static void bridge_success(AsyncOperation *op, void *results) {
    if (op && op->on_success) {
		op->ctx->completed = true;
		op->ctx->success = true;
		op->ctx->data = results;
		op->ctx->data_len = sizeof(results);

        op->on_success(op->ctx, results);
    }
}

static void bridge_error(AsyncOperation *op, char *err) {
    if (op && op->on_error) {
		op->ctx->completed = true;
		op->ctx->success = false;
		//op->ctx->error_msg = err;
		// TODO_IN_THIS_COMMIT: comment and/or debug - copy the array size - zero-initialized?
		memcpy(op->ctx->error_msg, err, 256);
		// TODO_IN_THIS_COMMIT: existing error codes?
		op->ctx->error_code = 1;

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
	"github.com/cosmos/cosmos-sdk/client/flags"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/poktroll/pkg/client"
	"github.com/pokt-network/poktroll/pkg/client/tx"
)

// TODO_IMPROVE: add separate constructor which supports options...

//// TODO_IN_THIS_COMMIT: godoc & move...
//func

// TODO_IN_THIS_COMMIT: godoc...
//
//export NewTxClient
func NewTxClient(
	depsRef C.go_ref,
	signingKeyName *C.char,
	gasSetting *C.gas_settings,
	cErr **C.char,
) C.go_ref {
	// TODO_CONSIDERATION: Could support a version of methods which receive a go context, created elsewhere..
	ctx := context.Background()

	deps, err := GetGoMem[depinject.Config](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	gasAndFeesOpts, err := getTxClientGasAndFeesOptions(gasSetting)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	opts := append(
		gasAndFeesOpts,
		tx.WithSigningKeyName(C.GoString(signingKeyName)),
	)
	txClient, err := tx.NewTxClient(ctx, deps, opts...)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	return SetGoMem(txClient)
}

// TODO_IN_THIS_COMMIT: godoc & move...
func getTxClientGasAndFeesOptions(cGasSetting *C.gas_settings) ([]client.TxClientOption, error) {
	var gasSettingTxClientOptions []client.TxClientOption

	if unsafe.Pointer(cGasSetting.fees) != C.NULL {
		fees, err := cosmostypes.ParseDecCoins(C.GoString(cGasSetting.fees))
		if err != nil {
			return nil, err
		}

		gasSettingTxClientOptions = append(
			gasSettingTxClientOptions,
			tx.WithFeeAmount(&fees),
		)
	}

	gasPrices, err := cosmostypes.ParseDecCoins(C.GoString(cGasSetting.gas_prices))
	if err != nil {
		return nil, err
	}

	gasSettingTxClientOptions = append(
		gasSettingTxClientOptions,
		tx.WithGasAdjustment(float64(cGasSetting.gas_adjustment)),
		tx.WithGasPrices(&gasPrices),
	)

	if cGasSetting.simulate {
		gasSettingTxClientOptions = append(
			gasSettingTxClientOptions,
			tx.WithGasSetting(&flags.GasSetting{
				Simulate: true,
				Gas:      uint64(cGasSetting.gas_limit),
			}),
		)
	} else {
		gasSettingTxClientOptions = append(
			gasSettingTxClientOptions,
			tx.WithGasSetting(&flags.GasSetting{
				Simulate: false,
				Gas:      uint64(cGasSetting.gas_limit),
			}),
		)
	}

	return gasSettingTxClientOptions, nil
}

//export WithSigningKeyName
func WithSigningKeyName(keyName *C.char) C.go_ref {
	return SetGoMem(tx.WithSigningKeyName(C.GoString(keyName)))
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

	txClient, err := GetGoMem[client.TxClient](txClientRef)
	if err != nil {
		err = fmt.Errorf("getting tx client ref: %s", err)
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(NilGoRef)
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
		return C.go_ref(NilGoRef)
	}

	if err = cdc.Unmarshal(C.GoBytes(unsafe.Pointer(serializedProto.data), C.int(serializedProto.data_length)), msg); err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(NilGoRef)
	}

	// TODO_IN_THIS_COMMIT: add a TxResponse data structure and return it...
	_, eitherAsyncErr := txClient.SignAndBroadcast(goCtx, msg)
	err, errCh := eitherAsyncErr.SyncOrAsyncError()
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(NilGoRef)
	}

	go func() {
		if err = <-errCh; err != nil {
			C.bridge_error(op, C.CString(err.Error()))
		} else {
			C.bridge_success(op, nil)
		}
	}()

	return SetGoMem(errCh)
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

	txClient, err := GetGoMem[client.TxClient](txClientRef)
	if err != nil {
		err = fmt.Errorf("getting tx client ref: %s", err)
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(NilGoRef)
	}

	if protoMessageArray.num_messages == 0 {
		C.bridge_error(op, C.CString("no messages provided"))
		return C.go_ref(NilGoRef)
	}

	msgs, err := CProtoMessageArrayToGoProtoMessages(protoMessageArray)
	if err != nil {
		err = fmt.Errorf("converting C proto messages to Go: %s", err)
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(NilGoRef)
	}

	// TODO_IN_THIS_COMMIT: add a TxResponse data structure and return it...
	_, eitherAsyncErr := txClient.SignAndBroadcast(goCtx, msgs...)
	err, errCh := eitherAsyncErr.SyncOrAsyncError()
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(NilGoRef)
	}

	go func() {
		if err = <-errCh; err != nil {
			C.bridge_error(op, C.CString(err.Error()))
		} else {
			C.bridge_success(op, nil)
		}
	}()

	return SetGoMem(errCh)
}
