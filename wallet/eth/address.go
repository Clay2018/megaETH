package eth

import (
	"golang.org/x/exp/rand"
	"time"
)

import (
	"github.com/ethereum/go-ethereum/crypto"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

var letterRunes = []rune("0123456789abcdef")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateAddrAndPriv() (addr string, privHex string, err error) {
	var privKey = randStringRunes(64)
	a, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", "", err
	}
	address := crypto.PubkeyToAddress(a.PublicKey)
	return address.String(), privKey, nil
}
