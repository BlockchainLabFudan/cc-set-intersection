package main

import (
	"crypto/rand"
	"testing"
)

func TestElgamalDeform(t *testing.T) {
	//初始化
	ed1 := EdInit()
	ed2 := EdInit()

	//密钥交换，产生公钥h
	h1, _ := ed1.SetPk(ed2.PK)
	h2, _ := ed2.SetPk(ed1.PK)
	println("h is right::", h1.Cmp(h2) == 0)

	//明文
	testM, _ := rand.Int(rand.Reader, ed1.P)

	println("testM::", testM.String())

	//加密
	testC1, testC2, _ := ed1.Encrypt(testM)

	println("testC1::", testC1.String())
	println("testC2::", testC2.String())

	//解密
	c1sk, _ := ed2.Decrypt1(testC1)
	decM, _ := ed1.Decrypt2(testC2, c1sk)

	println("decM::", decM.String())
	println("isRight::", decM.Cmp(testM.Exp(ed1.G, testM, ed1.P)) == 0)
}
