package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	appmodule "github.com/pokt-network/poktroll/x/application/module"
	gatewaymodule "github.com/pokt-network/poktroll/x/gateway/module"
	proofmodule "github.com/pokt-network/poktroll/x/proof/module"
	servicemodule "github.com/pokt-network/poktroll/x/service/module"
	sessionmodule "github.com/pokt-network/poktroll/x/session/module"
	sharedmodule "github.com/pokt-network/poktroll/x/shared/module"
	suppliermodule "github.com/pokt-network/poktroll/x/supplier/module"
	tokenomicsmodule "github.com/pokt-network/poktroll/x/tokenomics/module"
)

var (
	interfaceRegistry = codectypes.NewInterfaceRegistry()
	cdc               = codec.NewProtoCodec(interfaceRegistry)
)

// main is a dummy function to satisfy the cgo requirements.
func main() {}

// DEV_NOTES: tl;dr, cgo IS NOT Go!
//
// 1. Functions intended to be exported to C MUST:
//   1a. have an `//export <func_name>` comment on the line preceding their declaration.
//   1b. be declared in this `main` package.
// 2. C types which are included in one package NEVER match the same C types imported from another package.
//
// For more on cgo, see: https://pkg.go.dev/cmd/cgo

func init() {
	registerAllModuleInterfaces()
}

// TODO_IN_THIS_COMMIT: godoc...
func registerAllModuleInterfaces() {
	registerCosmosModuleInterfaces()
	registerPoktrollModuleInterfaces()
}

// TODO_IN_THIS_COMMIT: godoc...
func registerCosmosModuleInterfaces() {
	banktypes.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	authz.RegisterInterfaces(interfaceRegistry)
}

// TODO_IN_THIS_COMMIT: godoc...
func registerPoktrollModuleInterfaces() {
	appmodule.NewAppModuleBasic(cdc).RegisterInterfaces(interfaceRegistry)
	gatewaymodule.NewAppModuleBasic(cdc).RegisterInterfaces(interfaceRegistry)
	proofmodule.NewAppModuleBasic(cdc).RegisterInterfaces(interfaceRegistry)
	servicemodule.NewAppModuleBasic(cdc).RegisterInterfaces(interfaceRegistry)
	sessionmodule.NewAppModuleBasic(cdc).RegisterInterfaces(interfaceRegistry)
	sharedmodule.NewAppModuleBasic(cdc).RegisterInterfaces(interfaceRegistry)
	suppliermodule.NewAppModuleBasic(cdc).RegisterInterfaces(interfaceRegistry)
	tokenomicsmodule.NewAppModuleBasic(cdc).RegisterInterfaces(interfaceRegistry)
}
