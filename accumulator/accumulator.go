// Code obtained from https://github.com/davidlazar/accumulator/blob/master/accumulator.go
package accumulator

import (
	"crypto/rand"
	"io"
	"math/big"

	"golang.org/x/crypto/sha3"
)

// PrivateKey is the private key for an RSA accumulator.
// It is not needed for typical uses of an accumulator.
type PrivateKey struct {
	P, Q    *big.Int
	N       *big.Int // N = P*Q
	Totient *big.Int // Totient = (P-1)*(Q-1)
}

type PublicKey struct {
	N *big.Int
}

var PublicKeyStringN = "24227546953215951287727064035368461575332419432574266440193853859851156467664419848319735726803096008502036921842235773798551074887424749877810043783572420752424771456226887126107794826910408176885682945396048976015452308387318178817965681783714591465119598005486388111532169088413164420536426401315398145976447600034060291329210192469428497195160618668747544165520958683818616747500209886954403848192999574018478546092368072430079078603526245172959641259911082301202277612458044027094045890847346955590842481758180690186727822068510061120090891623461608219913733354804152451660057735596605570214940623968701132696087"
var PrivateKeyStringN = "24227546953215951287727064035368461575332419432574266440193853859851156467664419848319735726803096008502036921842235773798551074887424749877810043783572420752424771456226887126107794826910408176885682945396048976015452308387318178817965681783714591465119598005486388111532169088413164420536426401315398145976447600034060291329210192469428497195160618668747544165520958683818616747500209886954403848192999574018478546092368072430079078603526245172959641259911082301202277612458044027094045890847346955590842481758180690186727822068510061120090891623461608219913733354804152451660057735596605570214940623968701132696087"
var PrivateKeyStringQ = "143091116592808346196252835841636052064524216086100358146619759996150434952009528817929370094366813089212289652591815721762895861412165592283130027522353208904497376636170195479251230073303448816477448895691260479132426722365263333794917311459096958205028718233797527814315531559685036355277255463236590992149"
var PrivateKeyStringP = "143091116592808346196252835841636052064524216086100358146619759996150434952009528817929370094366813089212289652591815721762895861412165592283130027522353208904497376636170195479251230073303448816477448895691260479132426722365263333794917311459096958205028718233797527814315531559685036355277255463236590992149"
var PrivateKeyStringTotient = "24227546953215951287727064035368461575332419432574266440193853859851156467664419848319735726803096008502036921842235773798551074887424749877810043783572420752424771456226887126107794826910408176885682945396048976015452308387318178817965681783714591465119598005486388111532169088413164420536426401315398145976135193396489616702961436578090662093628779713469876434027444352229940042117854419213444757937998514442227287190492803194086301972193266149451664751090577694716105603695511133200909494414773299224742781197879849970221142303148697943845613396196599797210073227711158602051065285880611686542175447397220293149576"
var PublicKeyNBigInt, _ = new(big.Int).SetString(PublicKeyStringN, 10)
var PrivateKeyNBigInt, _ = new(big.Int).SetString(PrivateKeyStringN, 10)
var PrivateKeyQBigInt, _ = new(big.Int).SetString(PrivateKeyStringQ, 10)
var PrivateKeyPBigInt, _ = new(big.Int).SetString(PrivateKeyStringP, 10)
var PrivateKeyTotientBigInt, _ = new(big.Int).SetString(PrivateKeyStringTotient, 10)

var UniversalPublicKey = &PublicKey{
	N: PublicKeyNBigInt,
}

var UniversalPrivateKey = &PrivateKey{
	N:       PrivateKeyNBigInt,
	Q:       PrivateKeyQBigInt,
	P:       PrivateKeyPBigInt,
	Totient: PrivateKeyTotientBigInt,
}

func HashToPrime(data []byte) *big.Int {
	// Unclear if this is a good hash function.
	h := sha3.NewShake256()
	h.Write(data)
	p, err := rand.Prime(h, 256)
	if err != nil {
		panic(err)
	}
	return p
}

var base = big.NewInt(65537)
var bigOne = big.NewInt(1)
var bigTwo = big.NewInt(2)

// GenerateKey generates an RSA accumulator keypair. The private key
// is mostly used for debugging and should usually be destroyed
// as part of a trusted setup phase.

// This is mainly still here for testing purposes. We replace this with the Universal Ones above
func GenerateKey(random io.Reader) (*PublicKey, *PrivateKey, error) {
	for {
		p, err := rand.Prime(random, 1024)
		if err != nil {
			return nil, nil, err
		}

		q, err := rand.Prime(random, 1024)
		if err != nil {
			return nil, nil, err
		}

		pminus1 := new(big.Int).Sub(p, bigOne)
		qminus1 := new(big.Int).Sub(q, bigOne)
		totient := new(big.Int).Mul(pminus1, qminus1)

		g := new(big.Int).GCD(nil, nil, base, totient)
		if g.Cmp(bigOne) == 0 {
			privateKey := &PrivateKey{
				P:       p,
				Q:       q,
				N:       new(big.Int).Mul(p, q),
				Totient: totient,
			}
			publicKey := &PublicKey{
				N: new(big.Int).Set(privateKey.N),
			}
			return publicKey, privateKey, nil
		}
	}
}

func (key *PrivateKey) Accumulate(items ...[]byte) (acc *big.Int, witnesses []*big.Int) {
	primes := make([]*big.Int, len(items))
	for i := range items {
		primes[i] = HashToPrime(items[i])
	}

	exp := big.NewInt(1)
	for i := range primes {
		exp.Mul(exp, primes[i])
		exp.Mod(exp, key.Totient)
	}
	acc = new(big.Int).Exp(base, exp, key.N)

	witnesses = make([]*big.Int, len(items))
	for i := range items {
		inv := new(big.Int).ModInverse(primes[i], key.Totient)
		inv.Mul(exp, inv)
		inv.Mod(inv, key.Totient)
		witnesses[i] = new(big.Int).Exp(base, inv, key.N)
	}

	return
}

func (key *PublicKey) Accumulate(items ...[]byte) (acc *big.Int, witnesses []*big.Int) {
	primes := make([]*big.Int, len(items))
	for i := range items {
		primes[i] = HashToPrime(items[i])
	}

	acc = new(big.Int).Set(base)
	for i := range primes {
		acc.Exp(acc, primes[i], key.N)
	}

	witnesses = make([]*big.Int, len(items))
	for i := range items {
		wit := new(big.Int).Set(base)
		for j := range primes {
			if j != i {
				wit.Exp(wit, primes[j], key.N)
			}
		}
		witnesses[i] = wit
	}

	return
}

func (key *PublicKey) Verify(acc *big.Int, witness *big.Int, item []byte) bool {
	c := HashToPrime(item)
	v := new(big.Int).Exp(witness, c, key.N)
	return acc.Cmp(v) == 0
}
