package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
// #cgo LDFLAGS: -L${SRCDIR}/callback.c
#include <client.h>
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

	bankv1beta1 "cosmossdk.io/api/cosmos/bank/v1beta1"
	"cosmossdk.io/math"
	"github.com/pokt-network/poktroll/api/poktroll/application"
	"github.com/pokt-network/poktroll/pkg/client/tx"
	"google.golang.org/protobuf/types/known/anypb"

	"cosmossdk.io/depinject"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosproto "github.com/cosmos/gogoproto/proto"
	"github.com/pokt-network/poktroll/pkg/client"
	apptypes "github.com/pokt-network/poktroll/x/application/types"
	sharedtypes "github.com/pokt-network/poktroll/x/shared/types"
)

// TODO_IN_THIS_COMMIT: godoc...
// TODO_IN_THIS_COMMIT: add seperate constructor which supports options...
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
// TODO_IMPROVE: support multiple msgs (if top-level JSON array).
//
//export TxClient_SignAndBroadcastAny
func TxClient_SignAndBroadcastAny(
	op *C.AsyncOperation,
	txClientRef C.go_ref,
	msgAnyJSON *C.char,
) C.go_ref {
	goCtx := context.Background()

	txClient, err := GetGoMem[client.TxClient](GoRef(txClientRef))
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(NilGoRef)
	}

	msg, err := convertAnyMsgJSON(msgAnyJSON)
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(NilGoRef)
	}

	eitherAsyncErr := txClient.SignAndBroadcast(goCtx, msg)
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

	return C.go_ref(SetGoMem(errCh))
}

//export TxClient_SignAndBroadcast
func TxClient_SignAndBroadcast(
	op *C.AsyncOperation,
	txClientRef C.go_ref,
	cTypeUrl *C.char,
	msgBz *C.uchar,
	msgBzLen C.int,
) C.go_ref {
	goCtx := context.Background()

	txClient, err := GetGoMem[client.TxClient](GoRef(txClientRef))
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}

	typeUrl := C.GoString(cTypeUrl)
	if !strings.HasPrefix(typeUrl, "/") {
		typeUrl = "/" + typeUrl
	}

	msg, err := interfaceRegistry.Resolve(typeUrl)
	if err != nil {
		C.bridge_error(op, C.CString(err.Error()))
		return C.go_ref(ZeroGoRef)
	}
	//msgType := gogoproto.MessageType(C.GoString(typeUrl))
	//if msgType == nil {
	//	C.bridge_error(op, C.CString(fmt.Sprintf("unknown message type: %s", string(C.GoString(typeUrl)))))
	//	return C.go_ref(NilGoRef)
	//}

	//msg := reflect.New(msgType.Elem()).Interface().(cosmosproto.Message)
	if err = cdc.Unmarshal(C.GoBytes(unsafe.Pointer(msgBz), msgBzLen), msg); err != nil {
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

// TODO_IN_THIS_COMMIT: move & godoc...
func convertAnyMsgJSON(msgAnyJSON *C.char) (cosmosproto.Message, error) {
	msgAnyJSONBz := []byte(C.GoString(msgAnyJSON))
	msgAny := new(anypb.Any)
	if err := cdc.UnmarshalJSON(msgAnyJSONBz, msgAny); err != nil {
		return nil, err
	}

	var resultMsg cosmosproto.Message
	switch msgAny.GetTypeUrl() {
	case "type.googleapis.com/cosmos.bank.v1beta1.MsgSend":
		apiMsg := new(bankv1beta1.MsgSend)
		if err := msgAny.UnmarshalTo(apiMsg); err != nil {
			return nil, err
		}

		msg := new(banktypes.MsgSend)

		// TODO_IN_THIS_COMMIT: automate the below via reflection...
		msg.FromAddress = apiMsg.GetFromAddress()
		msg.ToAddress = apiMsg.GetToAddress()

		coins := make(cosmostypes.Coins, len(apiMsg.GetAmount()))
		for i, apiCoin := range apiMsg.GetAmount() {
			coinAmt, ok := math.NewIntFromString(apiCoin.GetAmount())
			if !ok {
				return nil, fmt.Errorf("failed to parse coin amount %q", apiCoin.GetAmount())
			}

			coin := cosmostypes.Coin{
				Denom:  apiCoin.GetDenom(),
				Amount: coinAmt,
			}
			coins[i] = coin
		}
		msg.Amount = coins

		resultMsg = msg
	case "type.googleapis.com/poktroll.application.MsgStakeApplication":
		apiMsg := new(application.MsgStakeApplication)
		if err := msgAny.UnmarshalTo(apiMsg); err != nil {
			return nil, err
		}

		msg := new(apptypes.MsgStakeApplication)

		// TODO_IN_THIS_COMMIT: automate the below via reflection...
		msg.Address = apiMsg.GetAddress()

		stakeAmt, ok := math.NewIntFromString(apiMsg.GetStake().GetAmount())
		if !ok {
			return nil, fmt.Errorf("failed to parse stake amount %q", apiMsg.GetStake().GetAmount())
		}

		stake := &cosmostypes.Coin{
			Denom:  apiMsg.GetStake().GetDenom(),
			Amount: stakeAmt,
		}
		msg.Stake = stake

		services := make([]*sharedtypes.ApplicationServiceConfig, len(apiMsg.GetServices()))
		for i, apiService := range apiMsg.GetServices() {
			service := &sharedtypes.ApplicationServiceConfig{
				ServiceId: apiService.GetServiceId(),
			}
			services[i] = service
		}
		msg.Services = services

		resultMsg = msg
	default:
		panic("unsupported msg type")
	}

	return resultMsg, nil
}
