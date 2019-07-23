package main

import (
	"github.com/BlockchainLabFudan/cc-set-intersection/model"
	"log"
	"math"
	"math/big"
	"testing"
)

var attrs1 = []*big.Int{big.NewInt(int64(math.Pow(2, 31))), big.NewInt(int64(math.Pow(2, 30)))}

//var attrs1 = []*big.Int{big.NewInt(2), big.NewInt(3)}
var attrs2 = []*big.Int{big.NewInt(2)}

func TestEncDec(t *testing.T) {
	//初始化两个用户
	//初始化公私钥
	ed1 := EdInit()
	ed2 := EdInit()

	//协议步骤1密钥交换，产生公钥h
	h1, _ := ed1.SetPk(ed2.PK)
	h2, _ := ed2.SetPk(ed1.PK)
	println("h is right::", h1.Cmp(h2) == 0)

	//生成多项式并加密
	user1, err := model.Init(ed1)
	if err != nil {
		log.Fatal(err)
	}
	user2, err := model.Init(ed2)
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

	share1, err := user1.ChangeShare(user1.ECoefficients1)
	if err != nil {
		log.Fatal(err)
	}

	sum2, err := user2.DecryptPolynomial(share1, user1.ECoefficients2)
	if err != nil {
		log.Fatal(err)
	}
	//初始求值
	Cal(attrs1, sum2)
	Cal(attrs2, sum2)
	//fmt.Println(big.NewInt(64).Mul(big.NewInt(16), big.NewInt(0).Exp(big.NewInt(4), ed1.P.Mod(attrs2[0], ed2.P), ed2.P)))
	//fmt.Println(big.NewInt(64).Mul(big.NewInt(16), big.NewInt(0).Exp(big.NewInt(4), attrs2[0], ed2.P)))
}

//
//func TestAdd(t *testing.T) {
//	user1Poly := NewUser(paillier)
//	err = user1Poly.SetPolynomial(attrs1)
//	if err != nil {
//		log.Fatal(err)
//	}
//	Cal(attrs1, user1Poly)
//	PrintCoef(user1Poly)
//	//	加密
//	err = user1Poly.EncryptPolynomial()
//	if err != nil {
//		log.Fatal(err)
//	}
//	//	相加
//	user1Poly, err = user1Poly.Add(user1Poly)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	Cal(attrs1, user1Poly)
//	PrintCoef(user1Poly)
//}
//
//func TestMul(t *testing.T) {
//	user1Poly := NewUser(paillier)
//	err = user1Poly.SetPolynomial(attrs1)
//	if err != nil {
//		log.Fatal(err)
//	}
//	Cal(attrs1, user1Poly)
//
//	//加密
//	err = user1Poly.EncryptPolynomial()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//随机
//	r, err := user1Poly.RandomPolynomial()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//	相乘
//	user1Poly, err = user1Poly.Mul(r)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = user1Poly.DecryptPolynomial()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	Cal(attrs1, user1Poly)
//}
