package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//填充字符串（末尾）
func PaddingText1(str []byte, blockSize int) []byte {
	//需要填充的数据长度
	paddingCount := blockSize - len(str)%blockSize
	//填充数据为：paddingCount ,填充的值为：paddingCount
	paddingStr := bytes.Repeat([]byte{byte(paddingCount)}, paddingCount)
	newPaddingStr := append(str, paddingStr...)
	//fmt.Println(newPaddingStr)
	return newPaddingStr
}

//去掉字符（末尾）
func UnPaddingText1(str []byte) []byte {
	n := len(str)
	count := int(str[n-1])
	newPaddingText := str[:n-count]
	return newPaddingText
}

//---------------DES加密  解密--------------------
func EncyptogAES(src, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(nil)
		return nil
	}
	src = PaddingText1(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return src

}
func DecrptogAES(src, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(nil)
		return nil
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)
	src = UnPaddingText1(src)
	return src
}

// func main() {
// 	str := `{"orderSerialId":"sz5ua7ffud21255bu85498942","partnerOrderId":"sz5ua7ffud21255bu85498942","consumeDate":"2020-04-29 23:59:08","tickets":1}}`
// 	fmt.Println("编码的数据为：", str)
// 	key := []byte("4MTU1KBG")
// 	src := EncyptogAES([]byte(str), key)
// 	//DecrptogAES(src, key)
// 	fmt.Println("解码之后的数据为：", string(src))

// }

func f1() (result int) {
	defer func() {
		result++
	}()
	return 0
}

func f2() (t int) {
	t = 5
	defer func() {
		t = t + 5
		fmt.Println("t:", t)
	}()
	return t
}

func f3() (r int) {
	defer func(r int) {

		r = r + 5

	}(r)
	return 1
}

type Test struct {
	Max int
}

func NewTest(max int) *Test {
	return &Test{
		Max: max,
	}
}
func (t *Test) Println() {
	fmt.Println(t.Max)
}

func deferExec(f func()) {
	f()
}

func call() {
	var t *Test
	//defer deferExec(t.Println)
	//t = new(Test)
	t = NewTest(6)
	t.Println()
}

// func main() {
// 	// for i := 0; i < 5; i++ {
// 	// 	func(idx int) {
// 	// 		fmt.Println(idx)
// 	// 	}(i) // 传入的 i，会立即被求值保存为 idx
// 	// }
// 	fmt.Println("f1:", f1())
// 	fmt.Println("f2:", f2())
// 	fmt.Println("f3:", f3())
// 	call()
// }
// func main() {
// 	c := boring("Joe")
// 	timeout := time.After(5 * time.Second)
// 	for {
// 		select {
// 		case s := <-c:
// 			fmt.Println(s)
// 		case <-timeout:
// 			fmt.Println("You talk too much.")
// 			return
// 		}
// 	}
// 	fmt.Printf()
// }
