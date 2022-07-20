package pow

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
)

const maxNonce = math.MaxInt64

type ProofOfWork struct {
	Data       []byte
	Nonce      uint64
	TargetBits uint16
	target     *big.Int
}

func NewProofOfWork(in []byte, targetBits uint16) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{
		Data:       in,
		TargetBits: targetBits,
		target:     target,
	}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	data := make([]byte, 8+len(pow.Data))
	binary.LittleEndian.PutUint64(data, nonce)
	copy(data[8:], pow.Data)
	return data
}

func (pow *ProofOfWork) Run(ctx context.Context) uint64 {
	var (
		hashInt big.Int
		hash    [32]byte
		nonce   uint64
	)

	for nonce < maxNonce && ctx.Err() == nil {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce
}

func Validate(in []byte, targetBits uint16, nonce uint64) bool {
	var hashInt big.Int

	pow := NewProofOfWork(in, targetBits)
	pow.Nonce = nonce

	data := pow.prepareData(pow.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

func Calc(ctx context.Context, in []byte, targetBits uint16) (nonce uint64) {
	p := NewProofOfWork(in, targetBits)
	nonce = p.Run(ctx)
	return
}
