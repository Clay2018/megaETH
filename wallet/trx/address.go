package trx

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/coins/tron"
	"golang.org/x/exp/rand"
	"os/exec"
	"time"
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

func GenerateAddrAndPriv() (addr string, pub string, privHex string, err error) {
	var privKey = randStringRunes(64)
	privKey = "0x" + privKey

	cmd := exec.Command("node", "GenerateAddrAndPriv.js", privKey)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", "", err
	}

	type TmpS struct {
		Addr string `json:"address"`
		Pub  string `json:"publicKey"`
	}
	var tmp TmpS
	err = json.Unmarshal(output, &tmp)
	if err != nil {
		return "", "", "", err
	}

	addr = tmp.Addr
	pub = tmp.Pub
	privHex = privKey
	return
}
func GenerateAddrAndPriv2() (addr string, pubHex string, privHex string, err error) {
	priv, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", "", err
	}
	privHex = hex.EncodeToString(priv.Serialize())

	pub := priv.PubKey().SerializeCompressed()
	pubHex = hex.EncodeToString(pub)

	//{
	//	fmt.Println("-----")
	//	addr2, err := tron.GetAddressByPublicKey(pubHex)
	//	if err != nil {
	//		fmt.Println("addr2:", err.Error())
	//	} else {
	//		fmt.Println("addr2:", addr2)
	//	}
	//	privHex2 := ""
	//	_, pub := btcec.PrivKeyFromBytes(hexutils.HexToBytes(privHex2))
	//	addr3 := tron.GetAddress(pub)
	//	fmt.Println("addr3:", addr3)
	//	fmt.Println("-----")
	//}
	addr = tron.GetAddress(priv.PubKey())
	return
}
