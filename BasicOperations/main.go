package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"crypto/hmac"
	"crypto/aes"
	"crypto/cipher"
    "encoding/base64"
    "github.com/ethereum/go-ethereum/crypto/ecies" 
	"os"
	"io"
	"time"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"crypto/sha256"
	"math/big"
	"fmt"
)
//生成ECC椭圆曲线密钥对，保存到文件
func GenerateECCKey() {
	//生成密钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//生成文件
	privatefile, err := os.Create("eccprivate.pem")
	if err != nil {
		panic(err)
	}
	//x509编码
	eccPrivateKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	//pem编码
	privateBlock := pem.Block{
		Type:  "ecc private key",
		Bytes: eccPrivateKey,
	}
	pem.Encode(privatefile, &privateBlock)
	//保存公钥
	publicKey := privateKey.PublicKey
	//创建文件
	publicfile, err := os.Create("eccpublic.pem")
	//x509编码
	eccPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem编码
	block := pem.Block{Type: "ecc public key", Bytes: eccPublicKey}
	pem.Encode(publicfile, &block)
}

//取得ECC私钥
func GetECCPrivateKey(path string) *ecdsa.PrivateKey {
	//读取私钥
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return privateKey
}

//取得ECC公钥
func GetECCPublicKey(path string) *ecdsa.PublicKey {
	//读取公钥
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解密
	block, _ := pem.Decode(buf)
	//x509解密
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey := publicInterface.(*ecdsa.PublicKey)
	return publicKey
}


func Hash256(msg string) string{	
	sum := sha256.Sum256([]byte(msg))
	Hstring := hex.EncodeToString(sum[:])
	fmt.Println(Hstring)
	return Hstring
}

func Hash256_1(msg []byte) string {
	hash := sha256.New()
	//填入数据
	hash.Write(msg)
	bytes := hash.Sum(nil)
	res := hex.EncodeToString(bytes)
	return res
}

func genHMAC(key []byte, msg string)string {
	mac := hmac.New(sha256.New, key)
	io.WriteString(mac, msg)
	sum := mac.Sum(nil)
 	fmt.Println(hex.EncodeToString(sum[:]))
	return hex.EncodeToString(sum[:])
}

func CheckMAC(message []byte, messageMAC []byte, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func AesEncrypt(orig string, key string) string {
    // 转成字节数组
    origData := []byte(orig)
    k := []byte(key)
    // 分组秘钥
    // NewCipher该函数限制了输入k的长度必须为16, 24或者32
    block, _ := aes.NewCipher(k)
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 补全码
    origData = PKCS7Padding(origData, blockSize)
    // 加密模式
    blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
    // 创建数组
    cryted := make([]byte, len(origData))
    // 加密
    blockMode.CryptBlocks(cryted, origData)
    return base64.StdEncoding.EncodeToString(cryted)
}
func AesDecrypt(cryted string, key string) string {
    // 转成字节数组
    crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
    k := []byte(key)
    // 分组秘钥
    block, _ := aes.NewCipher(k)
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 加密模式
    blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
    // 创建数组
    orig := make([]byte, len(crytedByte))
    // 解密
    blockMode.CryptBlocks(orig, crytedByte)
    // 去补全码
    orig = PKCS7UnPadding(orig)
    return string(orig)
}
//补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
    padding := blocksize - len(ciphertext)%blocksize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}
//去码
func PKCS7UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func AESTest() []byte {
	raw := "guoshuaijieffggggggg"
	key := "huyanyan87654321"
	encryptByte := AesEncrypt(raw, key)
	decryptByte := AesDecrypt(encryptByte, key)
		
	fmt.Println(string(decryptByte))
		
	
	return nil
}

func AESTest1(){
	start := time.Now().UnixNano()
	
    //orig := "helloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldAESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1AESTest1"
    orig:="大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大大"
    key := "0123456789012345"
    fmt.Println("原文：", orig)
    encryptCode := AesEncrypt(orig, key)
    fmt.Println("密文：" , encryptCode)
    decryptCode := AesDecrypt(encryptCode, key)
    fmt.Println("解密结果：", decryptCode)
    
    end := time.Now().UnixNano()
    fmt.Printf("runtime: %v (ns)\n", end-start)
}

func genPrivateKey() (*ecies.PrivateKey, error) {
	pubkeyCurve := elliptic.P256() //初始化椭圆曲线
	//随机挑选基点，生成私钥
	p, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) //用golang标准库生成公私钥
	if err != nil {
		return nil, err
	} else {
		return ecies.ImportECDSA(p), nil //转换成以太坊的公私钥对
	}
}

//ECCEncrypt 椭圆曲线加密
func ECCEncrypt(plain string, pubKey *ecies.PublicKey) ([]byte, error) {
	src := []byte(plain)
	return ecies.Encrypt(rand.Reader, pubKey, src, nil, nil)
}

//ECCDecrypt 椭圆曲线解密
func ECCDecrypt(cipher []byte, prvKey *ecies.PrivateKey) (string, error) {
	if src, err := prvKey.Decrypt(cipher, nil, nil); err != nil {
		return "", err
	} else {
		return string(src), nil
	}
}

func ECCTest() {
	prvKey, err := genPrivateKey()
	if err != nil {
		fmt.Println(err)
	}
	pubKey := prvKey.PublicKey
	plain := "我们没什么不同"
	cipher, err := ECCEncrypt(plain, &pubKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("密文：%v\n", cipher)
	plain, err = ECCDecrypt(cipher, prvKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("明文：%s\n", plain)
}




//对消息的散列值生成数字签名
func SignECC(msg []byte, path string)([]byte,[]byte) {
	//取得私钥
	privateKey := GetECCPrivateKey(path)
	//计算哈希值
	hash := sha256.New()
	//填入数据
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//对哈希值生成数字签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, bytes)
	if err != nil {
		panic(err)
	}
	rtext, _ := r.MarshalText()
	stext, _ := s.MarshalText()
	return rtext, stext
}

//验证数字签名
func VerifySignECC(msg []byte,rtext,stext []byte,path string) bool{
	//读取公钥
	publicKey:=GetECCPublicKey(path)
	//计算哈希值
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//验证数字签名
	var r,s big.Int
	r.UnmarshalText(rtext)
	s.UnmarshalText(stext)
	verify := ecdsa.Verify(publicKey, bytes, &r, &s)
	return verify
}
//测试
func main() {
	//生成ECC密钥对文件
	//start := time.Now()
	GenerateECCKey()
    //elapsed := time.Since(start)
    //fmt.Printf("Time taken: %s", elapsed)
    
	//模拟发送者
	//要发送的消息
	//start := time.Now()
	msg:=[]byte("hello world")
	//生成数字签名
	rtext,stext:=SignECC(msg,"eccprivate.pem")
	//elapsed := time.Since(start)
	
	//fmt.Printf("Time taken: %s", elapsed)
	//模拟接受者
	//接受到的消息
	acceptmsg:=[]byte("hello world")
	//接收到的签名
	acceptrtext:=rtext
	acceptstext:=stext
	//验证签名
	verifySignECC := VerifySignECC(acceptmsg, acceptrtext, acceptstext, "eccpublic.pem")
	fmt.Println("验证结果：",verifySignECC)
	
	
	Hash256("hello world")
	
	
	
	genHMAC([]byte("12345678"),"helloworld")
	
	AESTest1()
	
	
    start := time.Now()
    ECCTest()
    elapsed := time.Since(start)
    fmt.Printf("Time taken: %s", elapsed)
}
