package main

import (
	"context"

	"cosmossdk.io/depinject"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pokt-network/poktroll/app"
	"github.com/pokt-network/poktroll/pkg/cache/memory"
	"github.com/pokt-network/poktroll/pkg/client"
	"github.com/pokt-network/poktroll/pkg/client/query"
	"github.com/pokt-network/poktroll/pkg/client/query/cache"
	"github.com/pokt-network/poktroll/pkg/polylog/polyzero"
	apptypes "github.com/pokt-network/poktroll/x/application/types"
	prooftypes "github.com/pokt-network/poktroll/x/proof/types"
	servicetypes "github.com/pokt-network/poktroll/x/service/types"
	sessiontypes "github.com/pokt-network/poktroll/x/session/types"
	sharedtypes "github.com/pokt-network/poktroll/x/shared/types"
	suppliertypes "github.com/pokt-network/poktroll/x/supplier/types"
	"github.com/spf13/pflag"
)

var _ MultiQueryClient = (*queryClient)(nil)

// NewMultiQueryClient constructs a new MultiQueryClient and returns its Go reference to the C caller.
// Required dependencies:
//   - cosmosclient.Context (gogogrpc.ClientConn)
//   - client.BlockQueryClient
func NewMultiQueryClient(deps depinject.Config, queryNodeRPCURL string) (MultiQueryClient, error) {
	// TODO_IMPROVE: This should be parameterized.
	homedir := app.DefaultNodeHome
	clientCtx := cosmosclient.Context{}.
		WithCodec(cdc).
		WithTxConfig(TxConfig).
		WithHomeDir(homedir).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithInterfaceRegistry(InterfaceRegistry)

	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	// TODO_IMPROVE: It would be nice if the value could be set correctly based
	// on whether the test using it is running in tilt or not.
	flagSet.String(flags.FlagNode, queryNodeRPCURL, "")
	flagSet.String(flags.FlagHome, "", homedir)
	flagSet.String(flags.FlagChainID, app.Name, "use pocket chain-id")
	err := flagSet.Parse([]string{})
	if err != nil {
		return nil, err
	}

	clientCtx, err = cosmosclient.ReadPersistentCommandFlags(clientCtx, flagSet)
	if err != nil {
		return nil, err
	}

	// TODO_IN_FUTURE_RELEASE: Support and injectable logger.
	logger := polyzero.NewLogger()

	// TODO_IN_THIS_RELEASE: Refactor supply functions...
	accountCache, err := memory.NewKeyValueCache[cosmostypes.AccountI]()
	if err != nil {
		return nil, err
	}

	balancesCache, err := memory.NewKeyValueCache[query.Balance]()
	if err != nil {
		return nil, err
	}

	blockHashCache, err := memory.NewKeyValueCache[query.BlockHash]()
	if err != nil {
		return nil, err
	}

	appCache, err := memory.NewKeyValueCache[apptypes.Application]()
	if err != nil {
		return nil, err
	}

	supplierCache, err := memory.NewKeyValueCache[sharedtypes.Supplier]()
	if err != nil {
		return nil, err
	}

	sessionsCache, err := memory.NewKeyValueCache[*sessiontypes.Session]()
	if err != nil {
		return nil, err
	}

	servicesCache, err := memory.NewKeyValueCache[sharedtypes.Service]()
	if err != nil {
		return nil, err
	}

	relayMiningDifficultyCache, err := memory.NewKeyValueCache[servicetypes.RelayMiningDifficulty]()
	if err != nil {
		return nil, err
	}

	claimsCache, err := memory.NewKeyValueCache[prooftypes.Claim]()
	if err != nil {
		return nil, err
	}

	sharedParamsCache, err := cache.NewParamsCache[sharedtypes.Params]()
	if err != nil {
		return nil, err
	}

	appParamsCache, err := cache.NewParamsCache[apptypes.Params]()
	if err != nil {
		return nil, err
	}

	supplierParamsCache, err := cache.NewParamsCache[suppliertypes.Params]()
	if err != nil {
		return nil, err
	}

	sessionParamsCache, err := cache.NewParamsCache[sessiontypes.Params]()
	if err != nil {
		return nil, err
	}

	proofParamsCache, err := cache.NewParamsCache[prooftypes.Params]()
	if err != nil {
		return nil, err
	}

	deps = depinject.Configs(
		deps,
		depinject.Supply(
			clientCtx,
			logger,
			accountCache,
			balancesCache,
			blockHashCache,
			appCache,
			supplierCache,
			sessionsCache,
			servicesCache,
			claimsCache,
			relayMiningDifficultyCache,
			sharedParamsCache,
			appParamsCache,
			supplierParamsCache,
			sessionParamsCache,
			proofParamsCache,
		),
	)

	accountQuerier, err := query.NewAccountQuerier(deps)
	if err != nil {
		return nil, err
	}

	bankQuerier, err := query.NewBankQuerier(deps)
	if err != nil {
		return nil, err
	}

	blockQuerier, err := cosmosclient.NewClientFromNode(queryNodeRPCURL)
	if err != nil {
		return nil, err
	}

	sharedQuerier, err := query.NewSharedQuerier(deps)
	if err != nil {
		return nil, err
	}

	applicationQuerier, err := query.NewApplicationQuerier(deps)
	if err != nil {
		return nil, err
	}

	supplierQuerier, err := query.NewSupplierQuerier(deps)
	if err != nil {
		return nil, err
	}

	deps = depinject.Configs(deps, depinject.Supply(sharedQuerier))
	sessionQuerier, err := query.NewSessionQuerier(deps)
	if err != nil {
		return nil, err
	}

	serviceQuerier, err := query.NewServiceQuerier(deps)
	if err != nil {
		return nil, err
	}

	proofQuerier, err := query.NewProofQuerier(deps)
	if err != nil {
		return nil, err
	}

	// TODO_OPTIMIZE: lazily initialize these, so that they're only constructed when needed.
	return &queryClient{
		AccountQueryClient:     accountQuerier,
		BankQueryClient:        bankQuerier,
		BlockQueryClient:       blockQuerier,
		SharedQueryClient:      sharedQuerier,
		ApplicationQueryClient: applicationQuerier,
		SupplierQueryClient:    supplierQuerier,
		SessionQueryClient:     sessionQuerier,
		ServiceQueryClient:     serviceQuerier,
		ProofQueryClient:       proofQuerier,
	}, nil
}

// queryClient composes all pocket module query clients.
type queryClient struct {
	client.AccountQueryClient
	client.BankQueryClient
	client.BlockQueryClient
	client.SharedQueryClient
	client.ApplicationQueryClient
	client.SupplierQueryClient
	client.SessionQueryClient
	client.ServiceQueryClient
	client.ProofQueryClient
}

// GetSharedParams queries the chain for the current shared module parameters.
func (qc *queryClient) GetSharedParams(ctx context.Context) (*sharedtypes.Params, error) {
	return qc.SharedQueryClient.GetParams(ctx)
}

// GetSessionParams queries the chain for the current session module parameters.
func (qc *queryClient) GetSessionParams(ctx context.Context) (*sessiontypes.Params, error) {
	return qc.SessionQueryClient.GetParams(ctx)
}

// GetProofParams queries the chain for the current proof module parameters.
func (qc *queryClient) GetProofParams(ctx context.Context) (*prooftypes.Params, error) {
	params, err := qc.ProofQueryClient.GetParams(ctx)
	return params.(*prooftypes.Params), err
}

/* TODO_BLOCKED(@bryanchriswhite, #543): uncomment & implement once dependencies are available.

func (qc *queryClient) GetSupplierParams(ctx context.Context) (*suppliertypes.Params, error) {
	return qc.SupplierQueryClient.GetParams(ctx)
}

func (qc *queryClient) GetServiceParams(ctx context.Context) (*servicetypes.Params, error) {
	return qc.ServiceQueryClient.GetParams(ctx)
}

func (qc *queryClient) GetApplicationParams(ctx context.Context) (*apptypes.Params, error) {
	return qc.ApplicationQueryClient.GetParams(ctx)
}

*/
