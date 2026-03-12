package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#include <memory.h>
#include <protobuf.h>
#include <stdlib.h>
*/
import "C"
import (
	"context"
	"fmt"
	"unsafe"

	"cosmossdk.io/depinject"
	ring_secp256k1 "github.com/pokt-network/go-dleq/secp256k1"
	"github.com/pokt-network/poktroll/pkg/crypto"
	"github.com/pokt-network/poktroll/pkg/crypto/rings"
	"github.com/pokt-network/poktroll/pkg/polylog/polyzero"
	"github.com/pokt-network/poktroll/pkg/signer"
	gogoproto "github.com/cosmos/gogoproto/proto"
	servicetypes "github.com/pokt-network/poktroll/x/service/types"
)

// cgoRingCurve is the secp256k1 curve used for ring signature operations.
var cgoRingCurve = ring_secp256k1.NewCurve()

// NewRingClient constructs a new RingClient using the query clients from a
// MultiQueryClient reference. The RingClient can then be used to construct
// rings and sign relay requests.
//
// Parameters:
//   - queryClientRef: go_ref to a MultiQueryClient (from NewQueryClient)
//   - cErr: error output parameter
//
// Returns: go_ref to a crypto.RingClient
//
//export NewRingClient
func NewRingClient(queryClientRef C.go_ref, cErr **C.char) C.go_ref {
	multiClient, err := GetGoMem[MultiQueryClient](queryClientRef)
	if err != nil {
		*cErr = C.CString(fmt.Sprintf("getting query client ref: %s", err))
		return C.go_ref(NilGoRef)
	}

	// The rings.NewRingClient requires these dependencies via depinject:
	// - polylog.Logger
	// - client.ApplicationQueryClient
	// - client.AccountQueryClient
	// - client.SharedQueryClient
	//
	// MultiQueryClient embeds all three query client interfaces.
	logger := polyzero.NewLogger()
	deps := depinject.Supply(
		logger,
		multiClient.GetAccountQueryClient(),
		multiClient.GetApplicationQueryClient(),
		multiClient.GetSharedQueryClient(),
	)

	ringClient, err := rings.NewRingClient(deps)
	if err != nil {
		*cErr = C.CString(fmt.Sprintf("constructing ring client: %s", err))
		return C.go_ref(NilGoRef)
	}

	return SetGoMem(ringClient)
}

// RingClient_SignRelayRequest signs a serialized RelayRequest using a ring signature.
//
// It uses the RingClient to look up the application's ring (delegated gateways)
// and signs the relay request's signable bytes hash with the provided private key.
//
// Parameters:
//   - ringClientRef: go_ref to a crypto.RingClient (from NewRingClient)
//   - cPrivKeyBz: raw secp256k1 private key bytes (32 bytes)
//   - cPrivKeyBzLen: length of the private key bytes
//   - cRelayRequestBz: serialized RelayRequest protobuf bytes
//   - cRelayRequestBzLen: length of the relay request bytes
//   - cOutSigBz: pointer to output buffer for the signature bytes (caller must free)
//   - cOutSigBzLen: pointer to output the signature length
//   - cErr: error output parameter
//
//export RingClient_SignRelayRequest
func RingClient_SignRelayRequest(
	ringClientRef C.go_ref,
	cPrivKeyBz *C.uint8_t,
	cPrivKeyBzLen C.size_t,
	cRelayRequestBz *C.uint8_t,
	cRelayRequestBzLen C.size_t,
	cOutSigBz **C.uint8_t,
	cOutSigBzLen *C.size_t,
	cErr **C.char,
) {
	// Get the RingClient from the stored reference.
	ringClient, err := GetGoMem[crypto.RingClient](ringClientRef)
	if err != nil {
		*cErr = C.CString(fmt.Sprintf("getting ring client ref: %s", err))
		return
	}

	// Convert C bytes to Go bytes.
	privKeyBz := C.GoBytes(unsafe.Pointer(cPrivKeyBz), C.int(cPrivKeyBzLen))
	relayRequestBz := C.GoBytes(unsafe.Pointer(cRelayRequestBz), C.int(cRelayRequestBzLen))

	// Deserialize the RelayRequest.
	relayRequest := new(servicetypes.RelayRequest)
	if err := gogoproto.Unmarshal(relayRequestBz, relayRequest); err != nil {
		*cErr = C.CString(fmt.Sprintf("unmarshaling relay request: %s", err))
		return
	}

	// Extract the app address and session end height from the relay request.
	meta := relayRequest.GetMeta()
	sessionHeader := meta.GetSessionHeader()
	if sessionHeader == nil {
		*cErr = C.CString("relay request meta has no session header")
		return
	}

	appAddress := sessionHeader.GetApplicationAddress()
	sessionEndBlockHeight := sessionHeader.GetSessionEndBlockHeight()

	// Get the ring for this application at the session end block height.
	ctx := context.Background()
	appRing, err := ringClient.GetRingForAddressAtHeight(ctx, appAddress, sessionEndBlockHeight)
	if err != nil {
		*cErr = C.CString(fmt.Sprintf("getting ring for %s at height %d: %s", appAddress, sessionEndBlockHeight, err))
		return
	}

	// Decode the private key to a secp256k1 scalar for ring signing.
	privScalar, err := cgoRingCurve.DecodeToScalar(privKeyBz)
	if err != nil {
		*cErr = C.CString(fmt.Sprintf("decoding private key to scalar: %s", err))
		return
	}

	// Get the signable bytes hash of the relay request.
	signableBzHash, err := relayRequest.GetSignableBytesHash()
	if err != nil {
		*cErr = C.CString(fmt.Sprintf("getting relay request signable bytes hash: %s", err))
		return
	}

	// Sign with the ring using the private key.
	ringSigner := signer.NewRingSigner(appRing, privScalar)
	signatureBz, err := ringSigner.Sign(signableBzHash)
	if err != nil {
		*cErr = C.CString(fmt.Sprintf("ring signing relay request: %s", err))
		return
	}

	// Allocate C memory for the signature output.
	sigLen := len(signatureBz)
	cSigBz := C.malloc(C.size_t(sigLen))
	if cSigBz == nil {
		*cErr = C.CString("failed to allocate memory for signature")
		return
	}

	// Copy Go bytes to C memory.
	cSigSlice := unsafe.Slice((*byte)(cSigBz), sigLen)
	copy(cSigSlice, signatureBz)

	*cOutSigBz = (*C.uint8_t)(cSigBz)
	*cOutSigBzLen = C.size_t(sigLen)
}

// FreeCBytes frees memory that was allocated via C.malloc (e.g., signature output
// from RingClient_SignRelayRequest). Callers MUST call this to avoid memory leaks.
//
//export FreeCBytes
func FreeCBytes(cBz *C.uint8_t) {
	if cBz != nil {
		C.free(unsafe.Pointer(cBz))
	}
}
