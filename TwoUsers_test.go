package main

import (
	"crypto/rand"
	"fmt"
	"github.com/BlockchainLabFudan/cc-set-intersection/model"
	"log"
	"math/big"
	"testing"
)

var err error

func TestTwoUsers(t *testing.T) {

	//初始化两个用户
	//初始化公私钥
	ed1 := EdInit()
	ed2 := EdInit()

	//协议步骤1密钥交换，产生公钥h
	h1, _ := ed1.SetPk(ed2.PK)
	h2, _ := ed2.SetPk(ed1.PK)
	println("h is right::", h1.Cmp(h2) == 0)

	//构造属性集
	var m = 50
	var n = 10
	var attrs1, attrs2 []*big.Int
	for i := 0; i < m; i++ {
		a, _ := rand.Int(rand.Reader, ed1.P)
		attrs1 = append(attrs1, big.NewInt(1).Set(a))
		attrs2 = append(attrs2, big.NewInt(1).Set(a))
	}
	for i := 0; i < n; i++ {
		a, _ := rand.Int(rand.Reader, ed1.P)
		attrs1 = append(attrs1, big.NewInt(1).Set(a))
		b, _ := rand.Int(rand.Reader, ed1.P)
		attrs2 = append(attrs2, big.NewInt(1).Set(b))
	}
	//给user1一些多的数
	for i := 0; i < 10; i++ {
		a, _ := rand.Int(rand.Reader, ed1.P)
		attrs1 = append(attrs1, big.NewInt(1).Set(a))
	}
	fmt.Println("随机", m+n, "个数，前", m, "个相同，后", n, "个不同")
	fmt.Println(attrs1)
	fmt.Println(attrs2)

	//协议步骤2，生成多项式并加密
	user1, err := model.Init(ed1)
	if err != nil {
		log.Fatal(err)
	}
	err = user1.SetPolynomial(attrs1)
	if err != nil {
		log.Fatal(err)
	}
	err = user1.EncryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	user2, err := model.Init(ed2)
	if err != nil {
		log.Fatal(err)
	}
	err = user2.SetPolynomial(attrs2)
	if err != nil {
		log.Fatal(err)
	}
	err = user2.EncryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	//协议步骤3、4，密文给公信力机构计算整合的多项式
	//mergePolynomial := new(model.Polynomial)
	mergePolynomial, err := model.CalGcdPolynomial(user1, user2)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(mergePolynomial)

	//协议步骤5，交换解密参数
	share1, err := user1.ChangeShare(mergePolynomial.ECoefficients1)
	if err != nil {
		log.Fatal(err)
	}

	share2, err := user2.ChangeShare(mergePolynomial.ECoefficients1)
	if err != nil {
		log.Fatal(err)
	}

	//协议步骤6，解密并求值
	sum1, err := user1.DecryptPolynomial(share2, mergePolynomial.ECoefficients2)
	if err != nil {
		log.Fatal(err)
	}
	sum2, err := user2.DecryptPolynomial(share1, mergePolynomial.ECoefficients2)
	if err != nil {
		log.Fatal(err)
	}

	//求值
	Cal(attrs1, sum1)
	Cal(attrs2, sum2)

}
