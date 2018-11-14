package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//版本
const version = byte(0x00)

//截取长度
const addressChecksumLen = 4

type Wallet struct {
	//私钥
	PrivateKey ecdsa.PrivateKey
	//公钥
	PublicKey []byte
}

func NewWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()

	return &Wallet{privateKey, publicKey}
}

//生成公私钥
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey
}

//生成地址
func (w *Wallet) GetAddress() string {
	//公钥hash256，160
	publicKeyHash := PublickKeyHash(w.PublicKey)
	//拼接版本号
	publicKeyHash2 := append([]byte{version}, publicKeyHash...)
	//获得校验位
	checkSumBytes := CheckSum(publicKeyHash2)
	//再次拼接校验位
	bytes := append(publicKeyHash2, checkSumBytes...)
	//再次进行Base58
	address := Base58Encode(bytes)

	return string(address)
}

//公钥hash256，160
func  PublickKeyHash(publicKey []byte) []byte {
	//256
	hash := sha256.New()
	hash.Write(publicKey)
	hash256 := hash.Sum(nil)

	//160
	hash160 := ripemd160.New()
	hash160.Write(hash256)
	publicKeyHash := hash160.Sum(nil)

	return publicKeyHash
}

//获得校验位
func CheckSum(b []byte) []byte {
	hash1 := sha256.Sum256(b)
	hash2 := sha256.Sum256(hash1[:])

	return hash2[:addressChecksumLen]
}

//校验地址有效性
func IsValidForAddress(address string)  bool{
	version_key_check := Base58Decode([]byte(address))
	version_key := version_key_check[:len(version_key_check)-addressChecksumLen]
	checkbytes := version_key_check[len(version_key_check)-addressChecksumLen:]
	checksum := CheckSum(version_key)

	if bytes.Compare(checksum,checkbytes) == 0{
		return true
	}

	return false
}
