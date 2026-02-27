package main

import "C"
import (
	"context"

	"cosmossdk.io/depinject"

	"github.com/pokt-network/poktroll/pkg/client/block"
	"github.com/pokt-network/poktroll/pkg/polylog/polyzero"
)

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#include <memory.h>
*/
import "C"

//export NewBlockClient
func NewBlockClient(depsRef C.go_ref, cErr **C.char) C.go_ref {
	// TODO_CONSIDERATION: Could support a version of methods which receive a go context, created elsewhere..
	ctx := context.Background()

	deps, err := GetGoMem[depinject.Config](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	// Supply a logger; BlockClient needs polylog.Logger in deps since v0.1.23.
	logger := polyzero.NewLogger()
	fullDeps := depinject.Configs(deps, depinject.Supply(logger))

	blockClient, err := block.NewBlockClient(ctx, fullDeps)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	return SetGoMem(blockClient)
}
