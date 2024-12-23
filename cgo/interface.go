package main

import (
	"context"

	"github.com/pokt-network/poktroll/pkg/client"
	apptypes "github.com/pokt-network/poktroll/x/application/types"
	prooftypes "github.com/pokt-network/poktroll/x/proof/types"
	servicetypes "github.com/pokt-network/poktroll/x/service/types"
	sessiontypes "github.com/pokt-network/poktroll/x/session/types"
	sharedtypes "github.com/pokt-network/poktroll/x/shared/types"
)

type MultiQueryClient interface {
	// TODO_TECHDEBT(@bryanchriswhite, #2): observable-based clients.
	// client.EventsQueryClient
	// client.DelegationClient

	client.AccountQueryClient
	client.BankQueryClient
	client.BlockQueryClient
	SharedQueryClient
	ApplicationQueryClient
	SupplierQueryClient
	SessionQueryClient
	ServiceQueryClient
	ProofQueryClient

	// TODO_TECHDEBT(@bryanchriswhite): There's no gateway query client yet. 😅
	// GatewayQueryClient

	GetSharedParams(ctx context.Context) (*sharedtypes.Params, error)
	GetSessionParams(ctx context.Context) (*sessiontypes.Params, error)
	GetProofParams(ctx context.Context) (*prooftypes.Params, error)

	// TODO_BLOCKED(@bryanchriswhite poktroll#543): add once available.
	//GetApplicationParams(ctx context.Context) (*apptypes.Params, error)
	//GetSupplierParams(ctx context.Context) (*sharedtypes.Params, error)
	//GetServiceParams(ctx context.Context) (*sharedtypes.Params, error)
	//GetTokenomicsParams(ctx context.Context) (*tokenomics.Params, error)
}

// ApplicationQueryClient defines an interface that enables the querying of the
// on-chain application information
type ApplicationQueryClient interface {
	// GetApplication queries the chain for the details of the application provided
	GetApplication(ctx context.Context, appAddress string) (apptypes.Application, error)

	// GetAllApplications queries all on-chain applications
	GetAllApplications(ctx context.Context) ([]apptypes.Application, error)
}

// SupplierQueryClient defines an interface that enables the querying of the
// on-chain supplier information
type SupplierQueryClient interface {
	// GetSupplier queries the chain for the details of the supplier provided
	GetSupplier(ctx context.Context, supplierOperatorAddress string) (sharedtypes.Supplier, error)
}

// SessionQueryClient defines an interface that enables the querying of the
// on-chain session information
type SessionQueryClient interface {
	// GetSession queries the chain for the details of the session provided
	GetSession(
		ctx context.Context,
		appAddress string,
		serviceId string,
		blockHeight int64,
	) (*sessiontypes.Session, error)
}

// SharedQueryClient defines an interface that enables the querying of the
// on-chain shared module params.
type SharedQueryClient interface {
	// GetSessionGracePeriodEndHeight returns the block height at which the grace period
	// for the session that includes queryHeight elapses.
	// The grace period is the number of blocks after the session ends during which relays
	// SHOULD be included in the session which most recently ended.
	GetSessionGracePeriodEndHeight(ctx context.Context, queryHeight int64) (int64, error)
	// GetClaimWindowOpenHeight returns the block height at which the claim window of
	// the session that includes queryHeight opens.
	GetClaimWindowOpenHeight(ctx context.Context, queryHeight int64) (int64, error)
	// GetEarliestSupplierClaimCommitHeight returns the earliest block height at which a claim
	// for the session that includes queryHeight can be committed for a given supplier.
	GetEarliestSupplierClaimCommitHeight(ctx context.Context, queryHeight int64, supplierOperatorAddr string) (int64, error)
	// GetProofWindowOpenHeight returns the block height at which the proof window of
	// the session that includes queryHeight opens.
	GetProofWindowOpenHeight(ctx context.Context, queryHeight int64) (int64, error)
	// GetEarliestSupplierProofCommitHeight returns the earliest block height at which a proof
	// for the session that includes queryHeight can be committed for a given supplier.
	GetEarliestSupplierProofCommitHeight(ctx context.Context, queryHeight int64, supplierOperatorAddr string) (int64, error)
	// GetComputeUnitsToTokensMultiplier returns the multiplier used to convert compute units to tokens.
	GetComputeUnitsToTokensMultiplier(ctx context.Context) (uint64, error)
}

// ProofQueryClient defines an interface that enables the querying of the
// on-chain proof module params.
type ProofQueryClient interface {
	// TODO_BLOCKED(@bryanchriswhite poktroll#543): add once available.
	// GetProofParams queries the chain for the current shared module parameters.
	//GetProofParams(ctx context.Context) (*prooftypes.Params, error)
}

// ServiceQueryClient defines an interface that enables the querying of the
// on-chain service information
type ServiceQueryClient interface {
	// GetService queries the chain for the details of the service provided
	GetService(ctx context.Context, serviceId string) (sharedtypes.Service, error)
	GetServiceRelayDifficulty(ctx context.Context, serviceId string) (servicetypes.RelayMiningDifficulty, error)
}

// TODO_NEXT: tokenomics query client doesn't exist yet. 😅
// type TokenomicsQueryClient interface{}
