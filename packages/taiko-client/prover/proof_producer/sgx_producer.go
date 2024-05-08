opinion
package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/metrics"
)

const (
	ProofTypeSgx = "sgx"
	ProofTypeCPU = "native"
)

// SGXProofProducer generates a SGX proof for the given block.
type SGXProofProducer struct {
	RaikoHostEndpoint string // a proverd RPC endpoint
	L1Endpoint        string // a L1 node RPC endpoint
	L1BeaconEndpoint  string // a L1 beacon node RPC endpoint
	L2Endpoint        string // a L2 execution engine's RPC endpoint
	ProofType         string // Proof type
// RaikoRequestProofBody represents the JSON body for requesting the proof.
type RaikoRequestProofBody struct {
    L2RPC       string                     `json:"rpc"`
    L1RPC       string                     `json:"l1_rpc"`
    L1BeaconRPC string                     `json:"beacon_rpc"`
    Block       *big.Int                   `json:"block_number"`
    Prover      string                     `json:"prover"`
    Graffiti    string                     `json:"graffiti"`
    Type        string                     `json:"proof_type"`
    SGX         *ProofParam                `json:"sgx"`
    RISC0       RISC0RequestProofBodyParam `json:"risc0"`
}

// ProofParam represents the JSON body of RaikoRequestProofBody's `sgx` field.
type ProofParam struct {
    Setup     bool `json:"setup"`
    Bootstrap bool `json:"bootstrap"`
    Prove     bool `json:"prove"`
}

// requestProof sends a RPC request to proverd to try to get the requested proof.
func (s *SGXProofProducer) requestProof(opts *ProofRequestOptions) (*RaikoHostOutput, error) {
    reqBody := RaikoRequestProofBody{
        Type:        s.ProofType,
        Block:       opts.BlockID,
        L2RPC:       s.L2Endpoint,
        L1RPC:       s.L1Endpoint,
        L1BeaconRPC: s.L1BeaconEndpoint,
        Prover:      opts.ProverAddress.Hex()[2:],
        Graffiti:    opts.Graffiti,
        SGX: &ProofParam{
            Setup:     false,
            Bootstrap: false,
            Prove:     true,
        },
    }

    jsonValue, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    res, err := http.Post(s.RaikoHostEndpoint+"/proof", "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
        return nil, err
    }

    // Handle response...

    return nil, nil
}

	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request proof, id: %d, statusCode: %d", opts.BlockID, res.StatusCode)
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var output SGXRequestProofBodyResponse
	if err := json.Unmarshal(resBytes, &output); err != nil {
		return nil, err
	}

	if output.Error != nil {
		return nil, errors.New(output.Error.Message)
	}

	return output.Result, nil
}

// Tier implements the ProofProducer interface.
func (s *SGXProofProducer) Tier() uint16 {
	return encoding.TierSgxID
}
