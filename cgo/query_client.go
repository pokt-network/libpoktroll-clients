package main

/*
#include <memory.h>
#include <context.h>
#include <protobuf.h>
*/
import "C"
import (
	"context"
	"math"
	"unsafe"

	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/types/query"
	gogoproto "github.com/cosmos/gogoproto/proto"
	migrationtypes "github.com/pokt-network/poktroll/x/migration/types"
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
func QueryClient_GetSessionParams(clientRef C.go_ref, cErr **C.char) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
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
func QueryClient_GetProofParams(clientRef C.go_ref, cErr **C.char) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
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

// QueryClient_GetMigrationParams queries the chain for the current migration module parameters.
//
//export QueryClient_GetMigrationParams
func QueryClient_GetMigrationParams(clientRef C.go_ref, cErr **C.char) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	ctx := context.Background()
	queryParamsRes, err := multiClient.Params(ctx, &migrationtypes.QueryParamsRequest{})
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	migrationParams := queryParamsRes.GetParams()
	cSerializedProto, err := CSerializedProtoFromGoProto(&migrationParams)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProto
}

// QueryClient_GetMorseClaimableAccounts queries the chain for the current morse claimable accounts.
//
//export QueryClient_GetMorseClaimableAccounts
func QueryClient_GetMorseClaimableAccounts(clientRef C.go_ref, cErr **C.char) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	ctx := context.Background()
	queryMorseClaimableAccoutsRes, err := multiClient.MorseClaimableAccountAll(
		ctx, &migrationtypes.QueryAllMorseClaimableAccountRequest{
			Pagination: &query.PageRequest{
				// TODO_NEXT_RELEASE: extend the pagination API to C.
				Limit: math.MaxInt,
			},
		},
	)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	morseClaimableAccounts := queryMorseClaimableAccoutsRes.GetMorseClaimableAccount()
	morseClaimableAccountProtoMsgs := make([]gogoproto.Message, len(morseClaimableAccounts))
	for i, morseClaimableAccount := range morseClaimableAccounts {
		morseClaimableAccountProtoMsgs[i] = &morseClaimableAccount
	}
	cSerializedProtoArray, err := CSerializedProtoArrayFromGoProtoMessages(morseClaimableAccountProtoMsgs)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	return cSerializedProtoArray
}

// QueryClient_GetMorseClaimableAccount queries the chain for the Morse claimable
// account with the given morse_src_address.
//
//export QueryClient_GetMorseClaimableAccount
func QueryClient_GetMorseClaimableAccount(clientRef C.go_ref, cMorseSrcAddress *C.char, cErr **C.char) unsafe.Pointer {
	multiClient, err := GetGoMem[MultiQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	ctx := context.Background()
	queryMorseClaimableAccountRes, err := multiClient.MorseClaimableAccount(
		ctx, &migrationtypes.QueryMorseClaimableAccountRequest{
			Address: C.GoString(cMorseSrcAddress),
		},
	)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.NULL
	}

	if queryMorseClaimableAccountRes == nil {
		return C.NULL
	}

	morseClaimableAccount := queryMorseClaimableAccountRes.GetMorseClaimableAccount()
	cSerializedProto, err := CSerializedProtoFromGoProto(&morseClaimableAccount)
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
