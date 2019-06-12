package main

import (
	"log"
	"math"
	"math/big"
	"testing"
)

var attrs1 = []*big.Int{big.NewInt(int64(math.Pow(2, 31))), big.NewInt(int64(math.Pow(2, 31)))}

func TestEncDec(t *testing.T) {
	//初始化用户的初始多项式
	user1Poly := NewUser(paillier)

	//生成多项式并加密
	err = user1Poly.SetPolynomial(attrs1)
	if err != nil {
		log.Fatal(err)
	}

	//初始求值
	Cal(attrs1, user1Poly)

	err = user1Poly.EncryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	//加密后解密求值
	Cal(attrs1, user1Poly)
}

func TestAdd(t *testing.T) {
	user1Poly := NewUser(paillier)
	err = user1Poly.SetPolynomial(attrs1)
	if err != nil {
		log.Fatal(err)
	}
	Cal(attrs1, user1Poly)
	PrintCoef(user1Poly)
	//	加密
	err = user1Poly.EncryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}
	//	相加
	user1Poly, err = user1Poly.Add(user1Poly)
	if err != nil {
		log.Fatal(err)
	}

	Cal(attrs1, user1Poly)
	PrintCoef(user1Poly)
}

func TestMul(t *testing.T) {
	user1Poly := NewUser(paillier)
	err = user1Poly.SetPolynomial(attrs1)
	if err != nil {
		log.Fatal(err)
	}
	Cal(attrs1, user1Poly)

	//加密
	err = user1Poly.EncryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	//随机
	r, err := user1Poly.RandomPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	//	相乘
	user1Poly, err = user1Poly.Mul(r)
	if err != nil {
		log.Fatal(err)
	}
	err = user1Poly.DecryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	Cal(attrs1, user1Poly)
}
