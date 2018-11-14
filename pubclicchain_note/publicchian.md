# day01



## Block

###定义区块结构体

```go
type Block struct{
    Version int //版本号
    Height int //区块高度
    Timestamp int //时间戳
    PreBlockHash []byte //上一个区块hash
 	Merkroot []byte //梅克尔根
    Nonce int //随机数
    Hash []byte //当前区块hash
}
```

### 创建区块

```go
func NewBlock(Merkroot []byte,height int,prevBlockHash []byte) *Blcok{
 	//创建区块
    block := &Block{版本号，区块高度+1，时间戳，上一个区块哈希，梅克尔根，0，[]byte{}}
    //运行pow进行挖矿
    pow:=ProofWork(block)
    nonce,hash := pow.Run()
    //设置当前区块哈希
    block.Hash = hash
    //设置随机谁
    block.Nonce = nonce
    //校验pow有效性
    if !pow.Validata(){
        
    }
    
    return block
}
```

### 创建创世区块

```go
func CreateGenisBlock(Merkroot []byte)*Block{
    return NewBlock(Merkroot,1,[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
```

### 序列化和反序列化

```go
func (block *Block)Serialize()[]byte{
    var buffer bytes.Buffer
    encoder := gob.NewEncoder(&buffer)//创建编码器
    err := encoder.Encode(block)//编码
    if err != nil{
        log.panic(err)
    }
    return buffer.Bytes()
}

func DeSerialize(data []byte)*Block{
    var block Block
    decoder:=gob.NewDecoder(bytes.NewReader(data))//创建解码器
    err:=decoder.Decode(&block)//解码
    
    return &block
}
```

## ProofWork

### 定义难度值以及随机数最大数

```go
const (
	MaxNonce = math.MaxInt64
)
const targetBit = 16
```

### 定义pow结构体

```go
type ProofWork struct{
    block *Block//用于验证区块
    targetBit *big.Int//难度值
}
```

### 创建pow

```go
func (block *Block)*ProofWork(){
    target := big.NewInt(1)//创建目标值
    target.Lsh(target,unint(256-tarfetBit))//左移
    pow := &ProofWork{block,target}
    
    return pow
}
```

### 拼接字符串

```go
func (pow *ProofWork)PreareDate(nonce int)[]byte{
    data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.SerialzeTransation(),
		IntToHex(int64(pow.block.Timestamp)),
		IntToHex(int64(targetBit)),
		IntToHex(int64(nonce)),
		IntToHex(pow.block.Height)},
		[]byte{})

	return data
}
```

### 运行pow

```go
func (pow *ProofWork)Run(){
    var hashInt big.int
    var hash [32]byte
    nonce := 0
    for pow.block.Nonce < MaxNonce{
        data:=pow.PrepareData(nonce)
        hash = sha256.sum(data)
        //打印hash值
        hashInt.setbytes(hash[:])
        //比较hash和难度
        if hashInt.Cmp(pow.target) == -1{
            break
        }else {
            nonce++
        }
    }   
}
```

### 证明pow有效性

```go
func (pow *ProofWork)IsProofWorkVaild()bool{
    var hashInt big.int
    data:=pow.PreareData(pow.Nonce)
    hash := sha256.sum(data)
    hashInt.setbytes(hash[:])
    return hashInt.Cmp(pow.Hash) == -1
}
```



# day02



## BlockChain

### 定义数据库名，表名，最近区块hashkey

```go
const dbName = "blockchain.db"
const blockTableName = "blocks"
const LastHashKey = "l"
```

### 定义blockchain结构体

```go
type BlockChain struct{
	LastHash []byte //最新区块hash1
    DB *bolt.DB //数据库句柄
}
```

### 创建创世区块链

```go
func NewBlockchain() *BlockChain{
    //判断数据库是否存在
    if dbExists(){
        os.Exit(1)
    }
    
    var lastHash []byte
    //创建打开数据库
    db,_:=bolt.Open(dbName,0600,nil)
    //更新数据库
    db.Update(func (tx *bolt.Tx)error{
        //校验表是否存在
        b:=tx.Bucket(表名)
        if b != nil{
            //创建表
            tx.CreateBucket(表名)
            //创建创世区块
            block:=CreateGenisBlock(数据)
            //将创世区块存储到数据库中
            b.Put(创世区块hash，创世区块序列化)
            
           //存储最新区块hash
            b.put(LastHashKey,创世区块hash)
            //更新最新hash
            lastHash = block.Hash
        }
        
        return nil
    })
    
    return &BlockChain{lastHash,db}
}
```

### 进行交易挖矿

```go
func (blockchain *BlockChain)MineBlock(){
    //校验数据库是否存在
    if !dbExsits(){
        os.Exit()
    }
    
    //打开数据库获得最新交易hash和区块
    var block *Block
    blockchain.DB.View(func (tx *bolt.Tx)error{
        b := tx.Bucket(表名)
        if b != nil{
            hash := b.Get(LastHashKey)
            blockBytes := b.Get(hash)
            //反序列化hash
            block:=Deserialize(blockBytes)
        }
        
        return nil
    })
    //建立新区块
    block = NewBlock(梅克尔根，block.Height+1,block.hash)
    //将新区块存储到数据库中
    blockchain.DB.Update(func (tx *bolt.Tx)error{
        bucket := tx.Bucket(表名)
        if bucket != nil{
            bucket.put(block.hash,block.Serialize())
            //更新最新hash值
             b.put(LastHashKey,创世区块hash)
            //更新最新hash
            lastHash = block.Hash
        }
    })
}
```

### 打印区块链

```go
func (blockchain *Blockchain)PrintChain(){
    //获取迭代器
    it := blockchian.Iterator()
    //循环打印链
    for{
        it.Next()
        //打印
        。。。。。。。
        。。。、。。。。
    }
}
```

### 判断数据库是否存在

```go
func dbExsits()bool{
    if err:=os.stat(数据库名);os.isNotExist(err){
        return false
    }
    
    return true
}
```



## blockIterator

###定义迭代器

```go
type Iterator struct{
    CurrentHash []byte //当前hash
    DB *bolt.db //数据库句柄
}
```

### 创建迭代器

```go
func (blockchain *BlockChain)Iterator()*Iterator{
    return &Iterator{blockchain.LastHash,blockchain.DB}
}
```

### Next

```go
func (it *Iterator)Next()*Block{
    //校验数据库是否存在
    dbExsits()
    //查看数据库
    db,_:=it.DB.View(func (tx *Tx)error{
        //打开表
        bucket:=db.Bucket(表名)
        if bucket != nil {
            //根据当前hash找到block
            blockBytes:=bucket.Get(it.CurrentHash)
            block = Deserialize(blockBytes)
            //更新当前hash
            it.CurrentHash = block.PrevHash
        }
        
        return nil
    })
    return block
}
```



## 命令行

### 定义命令行

```go
type CLI struct{}
```

### 运行命令行

```go
func (cli *CLI)Run(){
    switch os.args[2:]
    case:
    case:
	default:
}
```





# day03

## 定义当前Transaction

```go
type Transacton struct{
    TxID []byte //当前交易hash
    TxInputs []*TxInput
    TxOutputs []*TxOutput
}
```



## 定义Input和Output

```go
type TxInput struct{
    TxID []byte //交易hsah
    Vout int //未花费输出按交易索引
    Signature string //数字签名
}
```

```go
type TxOutput struct{
    Value int //交易额
    PubKeyHash []byte //交易接受者公钥hash
}
```



## 获取UTXO

```go
func (bc *Blockchain)FindMyUtxos(address string)[]TxOutput{
    var UTXOs []TxOutput
    spentUTXOs := make(map[string][]*Transaction)//定义已经花费输出
    //遍历账本
    it := bc.Iterator
    for {
        block := it.Next()
        //遍历交易
        for _,tx := range bc.Transactions{
       		//遍历output,过滤已经花费输出
            OUTPUT:
            for index,output := range tx.TxOutputs{
                if address == output.PubkeyHash{
                    key := string(tx.TxID)
                    if len(spentUTXOs[key]) != 0{
                        for _,vout := range spentUTXOs[key]{
                            if index == vout {
                                continue OUTPUT
                            }
                        }
                    }
                    
                    UTXOs := append(UTXOs,output)
                }
            }
            
            //遍历input获得已经花费输出
            for _,input := range tx.TxInputs{
                if address == input.Signature{
                    key := string(inputs.TxId)
                    spentUTXOs[key] = append(spentUTXOs[key],input.Vout)
                }
            }
        }
       
        if len(bc前hash) == 0{
            break
        }
    }
    //遍历交易
    //找到和address地址相同的output
    //过滤掉已花费输出
}
```

## 获取余额函数

```go
func (bc *Blockchain)GetBalance(){
    //找到所有属于address的utxo的数组
	utxos := bc.FindMyUtxos(address)

	//总额
	total := 0.0

	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("%s 的余额为: %f\n", address, total)
}
```

## Tranasction

```go
func (tx *Transaction) SetTxID() {
	//使用gob编码，生成交易的哈希

	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}

//创建挖矿交易CoinbaseTx
//比特币挖矿的人可以有权利填充zhege Sig字段， 中本聪的创世语就是写在这里

func NewCoinbaseTx(miner, data string) *Transaction {
	//挖矿交易的特点， 没有输入， 只有输出
	input := TXInput{nil, -1, data}
	output := TXOutput{12.5, miner}

	tx := Transaction{nil, []TXInput{input}, []TXOutput{output}}
	tx.SetTxID()

	return &tx
}

//1. 找到所有合适utxos的集合
//
//1. 如果找到总额小于转账金额，转账失败
//
//1. 将所有的utxo转成input
//2. 创建output
//3. 如果有剩余，找零

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {

	//spentUTXOs := make(map[string][]int64)

	//1. 找到所有合适utxos的集合, 并且返回来
	//spentUTXOs[0x222] = []int64{0}
	//spentUTXOs[0x333] = []int64{0} //中间值
	//spentUTXOs[0x333] = []int64{0, 1}
	spentUTXOs, calcMoney := bc.FindNeedUTXOs(from, amount)

	if calcMoney < amount {
		fmt.Printf("余额不足，交易失败\n")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//创建input
	for txid, indexArray := range spentUTXOs {
		//key 是2222， 3333
		//value 0,     0, 1
		for _, i := range indexArray {
			//每一个output都要创建一个input
			input := TXInput{[]byte(txid), i, from}
			inputs = append(inputs, input)
		}
	}

	//创建output
	output := TXOutput{amount, to}
	outputs = append(outputs, output)

	//找零
	if calcMoney > amount {
		outputs = append(outputs, TXOutput{calcMoney - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTxID()
	return &tx
}
```



# day04

## 定义秘钥对

```go
type KeyPair struct{
    PrivateKey *ecdsa.PrivateKey
    Pubkey []byte
}
```

## 创建秘钥对

```go
//初始化，创建私钥
//通过私钥创建公钥
//pubKey := append(pubKeyRaw.X.Bytes(), pubKeyRaw.Y.Bytes()...) 真正传递公钥的形式
```

## 创建钱包地址

```go
//将公钥进行一次sha256哈希，再进行一次ripemd160哈希
//将版本号添加到公钥hash前一个字节
//进行两次sha256哈希
//取前四位作为校验位，添加到两次256哈希之后组成25字节数组
//进行一次base58生成地址
```

## 数字签名和验证

### SignTransaction

```go
func (blockchain *BlockChain)SignTransaction(tx *Transaction,privateKey *ecdsa.PrivateKey)bool{
    prevTxs := make(map[string]*Transaction)
    
    for _,input := range tx.TxInputs{
        tx := blockchain.FindTransaction(input.TxId)
        if tx == nil{
            return false
        }
        prevTxs[string(input.TxId)] = tx
    }
    
    return tx.Sign(privateKey,prevTxs)
}
```



### VerifyTransaction

```go
func (blockchain *Blockchain)VerifyTransaction(tx *transaction)bool{
     prevTxs := make(map[string]*Transaction)
    
    for _,input := range tx.TxInputs{
        tx := blockchain.FindTransaction(input.TxId)
        if tx == nil{
            return false
        }
        prevTxs[string(input.TxId)] = tx
    }
    
    return tx.Verify(prevTxs)
}
```



###FindTransaction

```go
func (blockchain *BlockChain)FindTransaction(Id []byte)*Transaction{
    iterator := blockchain.Iterator()
    
    for {
        block := iterator.next()
        
        for _,tx := range block.Transactions{
            if btytes.Equal(tx.TxId,Id){
                return tx
            }
        }
        
        if len(block.PrevBlockHash) == 0{
            fmt.Println("Not Find!")
            return nil
        }
    }
    
    return nil
}
```

### TrimmerCopy

```go
func (tx *Transaction)TrimmerCopy()*Transaction{
    var Inputs []TxInput
    
    for _,input := range tx.Inputs{
        Inputs = append(Inputs,TxInput{input.TxId,input.Index,nil.nil})   
    }
    
    return Inputs
}
```



### 签名

```go
func (tx *Transation) Sign(privkey ecdsa.PrivateKey, prevTxs map[string]Transation) {
   if tx.IsCoinbaseTransaction() {
      return
   }

   for _, vin := range tx.Vins {
      if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
         log.Panic("ERROE:Previoux transaction is not found!")
      }
   }
   txCopy := tx.TimemedCopy()

   for inHash, vin := range txCopy.Vins {
      prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
      txCopy.Vins[inHash].Signature = nil
      txCopy.Vins[inHash].PublicKey = prevTx.Vouts[vin.Vout].Pubkey
      txCopy.TxHash = txCopy.Hash()
      txCopy.Vins[inHash].PublicKey = nil

      //签名代码
      r, s, err := ecdsa.Sign(rand.Reader, &privkey, txCopy.TxHash)
      if err != nil {
         log.Panic(err)
      }
      signature := append(r.Bytes(), s.Bytes()...)
      tx.Vins[inHash].Signature = signature
   }
}

```

```go
//验证签名
func (tx *Transation) Verify(prevTxs map[string]Transation) bool {
   if tx.IsCoinbaseTransaction() {
      return true
   }

   for _, vin := range tx.Vins {
      if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
         log.Panic("ERROE:Previoux transaction is not found!")
      }
   }

   txCopy := tx.TimemedCopy()
   curve := elliptic.P256()

   for inHash, vin := range tx.Vins {
      prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
      txCopy.Vins[inHash].Signature = nil
      txCopy.Vins[inHash].PublicKey = prevTx.Vouts[vin.Vout].Pubkey
      txCopy.TxHash = txCopy.Hash()
      txCopy.Vins[inHash].PublicKey = nil

      //私钥ID
      r := big.Int{}
      s := big.Int{}
      sigLen := len(vin.Signature)
      r.SetBytes(vin.Signature[:(sigLen / 2)])
      s.SetBytes(vin.Signature[(sigLen / 2):])

      x := big.Int{}
      y := big.Int{}
      keyLen := len(vin.PublicKey)
      x.SetBytes(vin.PublicKey[:(keyLen / 2)])
      y.SetBytes(vin.PublicKey[(keyLen / 2):])

      rawPubkey := ecdsa.PublicKey{curve, &x, &y}
      if ecdsa.Verify(&rawPubkey, txCopy.TxHash, &r, &s) == false {
         return false
      }
   }

   return true
}
```

