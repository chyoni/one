package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/chiwon99881/one/utils"
)

type wallet struct {
	privateKey ecdsa.PrivateKey
	Address    string `json:"address"`
}

const (
	walletFile string = "one.wallet"
)

var w *wallet
var once sync.Once

func bigIntToBytes(aBig, bBig big.Int) string {
	aBytes := aBig.Bytes()
	bBytes := bBig.Bytes()
	aBytes = append(aBytes, bBytes...)
	bytesAsHex := fmt.Sprintf("%x", aBytes)
	return bytesAsHex
}

func hasPrivateKeyFile() bool {
	_, err := os.ReadFile(walletFile)
	fileExist := os.IsNotExist(err)
	return !fileExist
}

func parsePrivateKey(pk []byte) *ecdsa.PrivateKey {
	privateKey, err := x509.ParseECPrivateKey(pk)
	utils.HandleErr(err)
	return privateKey
}

func marshalPrivateKey(pk *ecdsa.PrivateKey) []byte {
	pkAsBytes, err := x509.MarshalECPrivateKey(pk)
	utils.HandleErr(err)
	return pkAsBytes
}

func Wallet() *wallet {
	if w == nil {
		once.Do(func() {
			w = &wallet{}
			if hasPrivateKeyFile() {
				pkAsBytes, err := os.ReadFile(walletFile)
				utils.HandleErr(err)
				w.privateKey = *parsePrivateKey(pkAsBytes)
			} else {
				privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				utils.HandleErr(err)
				_, err = os.OpenFile(walletFile, os.O_CREATE, 0644)
				utils.HandleErr(err)
				pkAsBytes := marshalPrivateKey(privateKey)
				err = os.WriteFile(walletFile, pkAsBytes, 0644)
				utils.HandleErr(err)

				w.privateKey = *privateKey
			}
			w.Address = bigIntToBytes(*w.privateKey.X, *w.privateKey.Y)
		})
	}
	return w
}
