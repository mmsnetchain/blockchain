package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"

	"crypto/rand"
	"crypto/x509"
	"encoding/base64"

	"encoding/pem"
	"fmt"
	"github.com/prestonTao/libp2parea/nodeStore"
	"mmschainnewaccount/config"
	"os"

	"encoding/hex"
	"github.com/prestonTao/utils"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type A struct {
	Id   string `"json:"id"`
	Name []byte `json:"name"`
}

type X struct {
	E []byte
}

func (x *X) Ts(b *[]byte, c string) {
	x.Tss(b)
}
func (x *X) Tss(e *[]byte) {
	x.E = []byte("aaa")
}
func EncodePubkey(pubkey *ecdsa.PublicKey) ([]byte, error) {
	pubKey, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(pubKey))
	buf := bytes.NewBuffer(nil)
	err = pem.Encode(buf, &pem.Block{Type: config.Core_addr_puk_type, Bytes: pubKey})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func DecodePubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	b, _ := pem.Decode(pubkey)
	fmt.Println(b)
	prk, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return prk.(*ecdsa.PublicKey), nil

}
func Verify() error {
	curve := btcec.S256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	pub, _ := EncodePubkey(&private.PublicKey)
	fmt.Printf("2%s\n", pub)

	x, err := DecodePubkey(pub)
	fmt.Printf("%+v %v\n", x, err)
	sign, err := utils.Sign(private, []byte("aaa"))
	if err != nil {
		fmt.Println(x, err)
	}
	fmt.Printf("sign:%v\n", sign)
	ok, _ := utils.Verify(pub, []byte("aaa"), hex.EncodeToString(*sign))
	if ok {
		fmt.Println("verify success")
	} else {
		fmt.Println("verify fail")
	}
	return nil
}

func VerifyS256(text string, pub []byte, sign string) bool {

	pubKeyBytes := pub
	fmt.Printf("%x\n", pubKeyBytes)

	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		fmt.Println("pub", err)
		return false
	}

	sigBytes, err := hex.DecodeString(sign)

	if err != nil {
		fmt.Println("sign", err)
		return false
	}

	signature, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Println("parse", err)
		return false
	}

	message := text
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := signature.Verify(messageHash, pubKey)
	fmt.Println("Signature Verified?", verified)

	return verified
}
func main() {

	js := `{
			"idinfo": {
				"id": "CPHfTv5LLyAZJY0tvl7/eovt8Jn1q2V83CU+3LPFPMQ=",
				"sign": "3045022100b869bb326cbc418cd6714f1247852a40ca5940a8304ffb0d9a8e9264d4defc4c02205e0082a9f2278e829f3e1231c6cf2df3dd72b8b8416e5abacdef0f197df061f0",
				"ctype": "s256",
				"puk": "BN8YWizuPl9mgeLX+D3kpMaV8YhmilRxjI4qZES2jZEwzoNRWuCr331gNhe63n8xXx+GoA9q0CeXqv3hH7ybfN4="
			},
			"issuper": false,
			"tcpport": 49537
		}`
	js = `{"idinfo":{"id":"QjRRZVN1S0NvekJvZ2tDRkxqb1NrTjI3RTdZSjFGYlJqTkc2a005dGNDOU4=","puk":"AhKZo9QM3ewdv8f87dKHBct+0gphM3Ys5yleZywLC5Df","sign":"3044022050be7307f5a266893b52f5d75dfe7387280ca844fec21067535f8d0a72092dd4022032fe51d57156d2f899ca259a55ebc746e396b1f1d5e8d4b30a47c8fe07e8ab5c"},"issuper":false,"tcpport":19981,"addr":"192.168.2.194"}`
	ss, ess := nodeStore.ParseNode([]byte(js))
	fmt.Printf("%+v %v\n", ss, ess)
	fmt.Println(ss.IdInfo.Id.B58String())

	fmt.Println(nodeStore.CheckSafeAddr(ss.IdInfo.Puk))
	fmt.Println(ss.IdInfo.Id.B58String(), "bvB8b8GMc2i5ybEpV44SbEeuA8BZHCmUZXVH68sZgsu")
	fmt.Println(VerifyS256(ss.IdInfo.Id.B58String(), ss.IdInfo.Puk, ss.IdInfo.Sign))
	fmt.Println(utils.VerifyS256(*ss.IdInfo.Id, ss.IdInfo.Puk, ss.IdInfo.Sign))
	os.Exit(0)
	pubtext, puberr := base64.StdEncoding.DecodeString("BN8YWizuPl9mgeLX+D3kpMaV8YhmilRxjI4qZES2jZEwzoNRWuCr331gNhe63n8xXx+GoA9q0CeXqv3hH7ybfN4=")
	fmt.Println(nodeStore.CheckSafeAddr(pubtext), puberr)
	os.Exit(0)
	text, err := base64.StdEncoding.DecodeString("YnZCOGI4R01jMmk1eWJFcFY0NFNiRWV1QThCWkhDbVVaWFZINjhzWmdzdQ==")
	fmt.Println(string(text))

	ecc1, err1 := nodeStore.GetKeyPair()
	pub1, _ := ecc1.GetPukBytes()
	sha := sha256.Sum256(pub1)
	fmt.Println(sha, err1)
	os.Exit(0)

	st := time.Now()
	ecc, err := nodeStore.GetKeyPair()
	fmt.Println(ecc, err)
	et := time.Now().Sub(st)
	fmt.Println(et)

	Verify()
	os.Exit(0)
}
