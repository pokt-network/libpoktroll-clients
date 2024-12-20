package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#include <memory.h>
#include <stdint.h>
#include <errno.h>
*/
import "C"
import (
	"context"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/pokt-network/poktroll/pkg/client"
)

//export NewBlockQueryClient
func NewBlockQueryClient(cometWebsocketURL *C.char, cErr **C.char) C.go_ref {
	// TODO_TECHDEBT: support opts args.
	blockQueryClient, err := sdkclient.NewClientFromNode(C.GoString(cometWebsocketURL))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	return SetGoMem(blockQueryClient)
}

//export BlockQueryClient_Block
func BlockQueryClient_Block(clientRef C.go_ref, cHeight *C.int64_t, cErr **C.char) C.go_ref {
	var height *int64
	if cHeight != nil {
		*height = int64(*cHeight)
	}

	// TODO_CONSIDERATION: Could support a version of methods which receive a go context, created elsewhere..
	ctx := context.Background()

	blockQueryClient, err :=
		GetGoMem[client.BlockQueryClient](clientRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	resultBlock, err := blockQueryClient.Block(ctx, height)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	blockProto, err := resultBlock.Block.ToProto()
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	// TODO_IN_THIS_COMMIT: return C-native struct.
	return SetGoMem(blockProto)
}
