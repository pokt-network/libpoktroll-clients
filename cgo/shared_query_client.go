package main

/*
#include <memory.h>
*/
import "C"
import (
	"context"
	"fmt"
)

// QueryClient_GetSharedParams queries the chain for the current shared module parameters.
//
//export QueryClient_GetSharedParams
func QueryClient_GetSharedParams(depsRef C.go_ref, cErr **C.char) C.go_ref {
	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	sharedParams, err := multiClient.GetSharedParams(context.TODO())
	if err != nil {
		*cErr = C.CString(err.Error())
		return C.go_ref(NilGoRef)
	}

	return SetGoMem(sharedParams)
}

// QueryClient_GetSessionGracePeriodEndHeight returns the block height at which
// the grace period for the session that includes queryHeight elapses.
// The grace period is the number of blocks after the session ends during which relays
// SHOULD be included in the session which most recently ended.
//
//export QueryClient_GetSessionGracePeriodEndHeight
func QueryClient_GetSessionGracePeriodEndHeight(depsRef C.go_ref, queryHeight C.int64_t, cErr **C.char) C.int64_t {
	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	gracePeriodEndHeight, err := multiClient.GetSessionGracePeriodEndHeight(context.TODO(), int64(queryHeight))
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	return C.int64_t(gracePeriodEndHeight)
}

// QueryClient_GetClaimWindowOpenHeight returns the block height at which the claim window of
// the session that includes queryHeight opens.
//
//export QueryClient_GetClaimWindowOpenHeight
func QueryClient_GetClaimWindowOpenHeight(depsRef C.go_ref, queryHeight C.int64_t, cErr **C.char) C.int64_t {
	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	claimWindowOpenHeight, err := multiClient.GetClaimWindowOpenHeight(context.TODO(), int64(queryHeight))
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	return C.int64_t(claimWindowOpenHeight)
}

// QueryClient_GetEarliestSupplierClaimCommitHeight returns the earliest block height at which a claim
// for the session that includes queryHeight can be committed for a given supplier.
//
//export QueryClient_GetEarliestSupplierClaimCommitHeight
func QueryClient_GetEarliestSupplierClaimCommitHeight(
	depsRef C.go_ref,
	queryHeight C.int64_t,
	supplierOperatorAddr *C.char,
	cErr **C.char,
) C.int64_t {
	fmt.Printf(">>> Go queryHeight: %d\n", queryHeight)

	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	claimWindowOpenHeight, err := multiClient.GetEarliestSupplierClaimCommitHeight(
		context.TODO(),
		int64(queryHeight),
		C.GoString(supplierOperatorAddr),
	)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	return C.int64_t(claimWindowOpenHeight)
}

// QueryClient_GetProofWindowOpenHeight returns the block height at which the proof window of
// the session that includes queryHeight opens.
//
//export QueryClient_GetProofWindowOpenHeight
func QueryClient_GetProofWindowOpenHeight(depsRef C.go_ref, queryHeight C.int64_t, cErr **C.char) C.int64_t {
	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	proofWindowOpenHeight, err := multiClient.GetProofWindowOpenHeight(context.TODO(), int64(queryHeight))
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	return C.int64_t(proofWindowOpenHeight)
}

// QueryClient_GetEarliestSupplierProofCommitHeight returns the earliest block height at which a proof
// for the session that includes queryHeight can be committed for a given supplier.
//
//export QueryClient_GetEarliestSupplierProofCommitHeight
func QueryClient_GetEarliestSupplierProofCommitHeight(
	depsRef C.go_ref,
	queryHeight C.int64_t,
	supplierOperatorAddr *C.char,
	cErr **C.char,
) C.int64_t {
	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	proofWindowOpenHeight, err := multiClient.GetEarliestSupplierProofCommitHeight(
		context.TODO(),
		int64(queryHeight),
		C.GoString(supplierOperatorAddr),
	)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	return C.int64_t(proofWindowOpenHeight)
}

// QueryClient_GetComputeUnitsToTokensMultiplier returns the multiplier used to convert compute units to tokens.
//
//export QueryClient_GetComputeUnitsToTokensMultiplier
func QueryClient_GetComputeUnitsToTokensMultiplier(depsRef C.go_ref, cErr **C.char) C.uint64_t {
	multiClient, err := GetGoMem[MultiQueryClient](depsRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	computeUnitsToTokensMultiplier, err := multiClient.GetComputeUnitsToTokensMultiplier(context.TODO())
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	return C.uint64_t(computeUnitsToTokensMultiplier)
}
