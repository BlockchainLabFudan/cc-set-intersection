package main

import (
	"log"
	"math/big"
	"testing"
)

var paillier = PaillierInit()
var err error

func TestTwoUsers(t *testing.T) {

	//初始化两个用户的初始多项式
	user1Poly := NewUser(paillier)
	user2Poly := NewUser(paillier)

	//构造属性集
	attrs1 := []*big.Int{big.NewInt(789789789), big.NewInt(5204002545), big.NewInt(123123123)}
	attrs2 := []*big.Int{big.NewInt(123123123), big.NewInt(5), big.NewInt(789456123)}

	//生成多项式并加密
	err = user1Poly.SetPolynomial(attrs1)
	if err != nil {
		log.Fatal(err)
	}

	err = user1Poly.EncryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	err = user2Poly.SetPolynomial(attrs2)
	if err != nil {
		log.Fatal(err)
	}

	err = user2Poly.EncryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	//用户选择随机多项式r00,r01
	r00, err := user1Poly.RandomPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	r01, err := user1Poly.RandomPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	//计算最终加密多项式
	r00E, err := user1Poly.Mul(r00)
	if err != nil {
		log.Fatal(err)
	}

	r01E, err := user2Poly.Mul(r01)
	if err != nil {
		log.Fatal(err)
	}

	sum1, err := r00E.Add(r01E)
	if err != nil {
		log.Fatal(err)
	}

	//解密
	err = sum1.DecryptPolynomial()
	if err != nil {
		log.Fatal(err)
	}

	//求值
	Cal(attrs1, sum1)

}
