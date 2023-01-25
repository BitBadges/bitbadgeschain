package types

import "github.com/tendermint/tendermint/crypto/merkle"

func ConvertToTendermintProof(proof *Proof) *merkle.Proof {
	aunts := make([][]byte, len(proof.Aunts))
	for i, v := range proof.Aunts {
		aunts[i] = []byte(v)
	}
	return &merkle.Proof{
		Total:    int64(proof.Total),
		Index:    int64(proof.Index),
		LeafHash: []byte(proof.LeafHash),
		Aunts:    aunts,
	}
}

func ConvertFromTendermintProof(proof *merkle.Proof) *Proof {
	aunts := make([]string, len(proof.Aunts))
	for i, v := range proof.Aunts {
		aunts[i] = string(v)
	}

	return &Proof{
		Total:    uint64(proof.Total),
		Index:    uint64(proof.Index),
		LeafHash: string(proof.LeafHash),
		Aunts:    aunts,
	}
}