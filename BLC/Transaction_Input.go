package BLC

import "bytes"

type TXInput struct {
	TxHash    []byte //交易hash
	Vout      int    //out索引
	Signature []byte //数字签名
	PublicKey []byte //公钥
}

//判断当前的消费者是谁
func (txInput *TXInput) UnLockRipmed(ripmedHash []byte) bool {

	publicKey := PublickKeyHash(txInput.PublicKey)
	return bytes.Compare(publicKey, ripmedHash) == 0
}
