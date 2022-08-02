// Copyright 2017 David Lazar. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package accumulator_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/trevormil/bitbadgeschain/accumulator"
)

func TestAccumulator(t *testing.T) {
	publicKey := accumulator.UniversalPublicKey
	privateKey := accumulator.UniversalPrivateKey

	items := makeItems(100)
	acc, witnesses := publicKey.Accumulate(items...)
	accPriv, witnessesPriv := privateKey.Accumulate(items...)
	if acc.Cmp(accPriv) != 0 {
		t.Fatalf("public/private accumulators differ: %s != %s", acc, accPriv)
	}

	for i := range witnesses {
		if witnesses[i].Cmp(witnessesPriv[i]) != 0 {
			t.Fatalf("witness %d: public/private witnesses differ: %s != %s", i, witnesses[i], witnessesPriv[i])
		}
	}

	badItem := make([]byte, 32)
	rand.Read(badItem)

	for i, w := range witnesses {
		if !publicKey.Verify(acc, w, items[i]) {
			t.Fatal("item not found")
		}

		if publicKey.Verify(acc, w, badItem) {
			t.Fatal("bad item was verified")
		}

		if publicKey.Verify(acc, w, items[(i+1)%len(items)]) {
			t.Fatal("bad item was verified")
		}
	}
}

func makeItems(count int) [][]byte {
	items := make([][]byte, count)
	for i := range items {
		items[i] = make([]byte, 32)
		rand.Read(items[i])
	}
	return items
}

func benchAccumulate(b *testing.B, count int) {
	privateKey := accumulator.UniversalPrivateKey
	items := makeItems(count)
	b.ResetTimer()
	b.SetBytes(int64(count) * 32)

	acc := new(big.Int)
	for i := 0; i < b.N; i++ {
		acc, _ = privateKey.Accumulate(items...)
	}

	_ = acc
}

func benchAccumulatePublic(b *testing.B, count int) {
	publicKey := accumulator.UniversalPublicKey
	items := makeItems(count)
	b.ResetTimer()
	b.SetBytes(int64(count) * 32)
	acc := new(big.Int)
	for i := 0; i < b.N; i++ {
		acc, _ = publicKey.Accumulate(items...)
	}

	_ = acc
}

var verify_items = makeItems(1)

func benchVerify(b *testing.B, count int) {
	for i := 0; i < count; i++ {
		accumulator.UniversalPublicKey.Verify(accumulator.PrivateKeyPBigInt, accumulator.PrivateKeyPBigInt, verify_items[0])
	}
	//1000 4.432 s
	//1 0.017 s
	//100
	//10000 17.082 s
	// _ = acc
}

func BenchmarkVerify(b *testing.B) {
	benchVerify(b, 1000)
}

func BenchmarkAccumulate100(b *testing.B) {
	benchAccumulate(b, 100)
}

func BenchmarkAccumulate1000(b *testing.B) {
	benchAccumulate(b, 1000)
}

func BenchmarkAccumulate10000(b *testing.B) {
	benchAccumulate(b, 10000)
}

//Public Key is pretty slow comparitively
func BenchmarkAccumulate100Public(b *testing.B) {
	benchAccumulatePublic(b, 100)
}

func BenchmarkAccumulate1000Public(b *testing.B) {
	benchAccumulatePublic(b, 1000)
}

func BenchmarkAccumulate10000Public(b *testing.B) {
	benchAccumulatePublic(b, 10000)
}
