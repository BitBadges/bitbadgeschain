package types

import "github.com/tendermint/tendermint/crypto/merkle"

func ConvertToTendermintProof(proof *Proof) *merkle.Proof {
	return &merkle.Proof{
		Total:    int64(proof.Total),
		Index:    int64(proof.Index),
		LeafHash: proof.LeafHash,
		Aunts:    proof.Aunts,
	}
}

func ConvertFromTendermintProof(proof *merkle.Proof) *Proof {
	return &Proof{
		Total:    uint64(proof.Total),
		Index:    uint64(proof.Index),
		LeafHash: proof.LeafHash,
		Aunts:    proof.Aunts,
	}
}