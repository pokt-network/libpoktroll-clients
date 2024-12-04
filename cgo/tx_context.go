package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#include <client.h>
*/
import "C"
import (
	"os"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cosmostx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/pflag"

	"github.com/pokt-network/poktroll/app"
	"github.com/pokt-network/poktroll/cmd/poktrolld/cmd"
	"github.com/pokt-network/poktroll/pkg/client/tx"
	txtypes "github.com/pokt-network/poktroll/pkg/client/tx/types"
)

var (
	// TxConfig provided by app.AppConfig(), intended as a convenience for use in tests.
	TxConfig client.TxConfig
	// Marshaler provided by app.AppConfig(), intended as a convenience for use in tests.
	Marshaler codec.Codec
	// InterfaceRegistry provided by app.AppConfig(), intended as a convenience for use in tests.
	InterfaceRegistry codectypes.InterfaceRegistry
)

func init() {
	cmd.InitSDKConfig()

	deps := depinject.Configs(
		// TODO_TECHDEBT: Avoid importing the entire app package - bloats the .so file.
		app.AppConfig(),
		depinject.Supply(
			log.NewLogger(os.Stderr),
		),
	)

	// Ensure that the global variables are initialized.
	if err := depinject.Inject(
		deps,
		&TxConfig,
		&Marshaler,
		&InterfaceRegistry,
	); err != nil {
		panic(err)
	}

	//// If VALIDATOR_RPC_ENDPOINT environment variable is set, use it to override the default localnet endpoint.
	//if endpoint := os.Getenv("VALIDATOR_RPC_ENDPOINT"); endpoint != "" {
	//	CometLocalTCPURL = fmt.Sprintf("tcp://%s", endpoint)
	//	CometLocalWebsocketURL = fmt.Sprintf("ws://%s/websocket", endpoint)
	//}
}

// TODO_IN_THIS_COMMIT: godoc...
// TODO_IN_THIS_COMMIT: add seperate constructor which supports deps...
// func NewTxContext(depsRef C.go_ref, cErr **C.char) C.go_ref {
//
//export NewTxContext
func NewTxContext(tcpURL *C.char, cErr **C.char) C.go_ref {
	//deps, err := GetGoMem[depinject.Config](depsRef)
	//if err != nil {
	//	*cErr = C.CString(err.Error())
	//	return 0
	//}

	flagSet, err := newFlagSet(C.GoString(tcpURL))
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	clientCtx, err := newClientCtx(flagSet)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	txFactory, err := cosmostx.NewFactoryCLI(clientCtx, flagSet)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	deps := depinject.Supply(
		txtypes.Context(clientCtx),
		txFactory,
	)

	txCtx, err := tx.NewTxContext(deps)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	return C.go_ref(SetGoMem(txCtx))
}

// TODO_IN_THIS_COMMIT: godoc...
func newFlagSet(tcpURL string) (*pflag.FlagSet, error) {
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// TODO_IN_THIS_COMMIT: parameterize these
	flagSet.String(flags.FlagHome, app.DefaultNodeHome, "")
	flagSet.String(flags.FlagKeyringBackend, "test", "")
	// ---

	flagSet.String(flags.FlagNode, tcpURL, "")
	flagSet.String(flags.FlagChainID, app.Name, "use poktroll chain-id")
	if err := flagSet.Parse([]string{}); err != nil {
		return nil, err
	}

	return flagSet, nil
}

// TODO_IN_THIS_COMMIT: godoc...
func newClientCtx(flagSet *pflag.FlagSet) (client.Context, error) {
	homedir, err := flagSet.GetString(flags.FlagHome)
	if err != nil {
		return client.Context{}, err
	}

	clientCtx := client.Context{}.
		WithCodec(Marshaler).
		WithTxConfig(TxConfig).
		WithHomeDir(homedir).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithInterfaceRegistry(InterfaceRegistry)

	return client.ReadPersistentCommandFlags(clientCtx, flagSet)
}
