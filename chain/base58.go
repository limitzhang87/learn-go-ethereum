package chain

import (
	"bytes"
	"math/big"
)

// Base58编码  去掉 0 o  I l
var base58 = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58Encode Base58编码
func Base58Encode(strByte []byte) []byte { //Base58编码
	//1. 转换成ascii码对应的值
	//2. 转换十进制
	strTen := big.NewInt(0).SetBytes(strByte)
	//fmt.Println(strTen)  // 结果4612462
	//3. 取出余数
	var modSlice []byte
	for strTen.Cmp(big.NewInt(0)) > 0 {
		mod := big.NewInt(0) //余数
		strTen58 := big.NewInt(58)
		strTen.DivMod(strTen, strTen58, mod)             //取余运算
		modSlice = append(modSlice, base58[mod.Int64()]) //存储余数,并将对应值放入其中
	}
	//  处理0就是1的情况 0使用字节'1'代替
	for _, elem := range strByte {
		if elem != 0 {
			break
		} else if elem == 0 {
			modSlice = append(modSlice, byte('1'))
		}
	}
	//fmt.Println(modSlice)   //结果 [12 7 37 23] 但是要进行反转，因为求余的时候是相反的。
	//fmt.Println(string(modSlice))  //结果D8eQ
	ReverseModSlice := ReverseByteArr(modSlice)
	//fmt.Println(ReverseModSlice)  //反转[81 101 56 68]
	//fmt.Println(string(ReverseModSlice))  //结果Qe8D
	return ReverseModSlice
}

func ReverseByteArr(bytes []byte) []byte { //将字节的数组反转
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i] //前后交换
	}
	return bytes
}

// Base58Decode 就是编码的逆过程
func Base58Decode(strByte []byte) []byte { //Base58解码
	ret := big.NewInt(0)
	for _, byteElem := range strByte {
		index := bytes.IndexByte(base58, byteElem) //获取base58对应数组的下标
		ret.Mul(ret, big.NewInt(58))               //相乘回去
		ret.Add(ret, big.NewInt(int64(index)))     //相加
	}
	//fmt.Println(ret) 	// 拿到了十进制 4612462
	//fmt.Println(ret.Bytes())  //[70 97 110]
	//fmt.Println(string(ret.Bytes()))
	return ret.Bytes()
}
