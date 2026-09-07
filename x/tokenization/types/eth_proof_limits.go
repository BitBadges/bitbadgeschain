package types

import "fmt"

const MaxETHSignatureProofs = 100
const ETHSignatureRecoveryGas uint64 = 3000

func ValidateETHSignatureProofs(proofs []*ETHSignatureProof) error {
	if len(proofs) > MaxETHSignatureProofs {
		return fmt.Errorf("too many ETH signature proofs: maximum %d", MaxETHSignatureProofs)
	}
	for _, proof := range proofs {
		if proof == nil {
			return fmt.Errorf("ETH signature proof cannot be nil")
		}
		if len(proof.Nonce) > 256 {
			return fmt.Errorf("ETH signature nonce exceeds 256 bytes")
		}
		for _, c := range proof.Nonce {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
				return fmt.Errorf("ETH signature nonce must use ASCII letters, digits, hyphens or underscores")
			}
		}
		if len(proof.Signature) > 132 {
			return fmt.Errorf("ETH signature exceeds 65 hex-encoded bytes")
		}
	}
	return nil
}
