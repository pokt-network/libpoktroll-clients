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

	"github.com/pokt-network/poktroll/pkg/client/events"

	"github.com/pokt-network/poktroll/pkg/client"
)

//export NewEventsQueryClient
func NewEventsQueryClient(cometWebsocketURLCString *C.char) C.go_ref {
	// TODO_TECHDEBT: support opts args.
	cometWebsocketURL := C.GoString(cometWebsocketURLCString)
	eventsQueryClient := events.NewEventsQueryClient(cometWebsocketURL)

	return C.go_ref(SetGoMem(eventsQueryClient))
}

//export EventsQueryClientEventsBytes
func EventsQueryClientEventsBytes(clientRef C.go_ref, query *C.char, cErr **C.char) C.go_ref {
	// TODO_CONSIDERATION: Could support a version of methods which receive a go context, created elsewhere..
	ctx := context.Background()

	eventsQueryClient, err :=
		GetGoMem[client.EventsQueryClient](GoRef(clientRef))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	eventsBytesObs, err := eventsQueryClient.EventsBytes(ctx, C.GoString(query))
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	return C.go_ref(SetGoMem(eventsBytesObs))
}
