package main

import (
	"testing"
	"github.com/BlockchainLabFudan/cc-set-intersection/model"
	"math/big"
	"crypto/rand"
)

func TestDemo(t *testing.T) {
	paillier := new(model.Paillier)

	p1 := new(big.Int)
	p2 := new(big.Int)
	p1.SetString("242661090146032969904098483991985908921", 10) // octal
	p2.SetString("215662396313044988944834777682074105079", 10) // octal

	paillier.Init(p1, p2)
	
	testM, _ := rand.Int(rand.Reader, paillier.PK)
	
	println("testM::", testM.String())
	
	testC, _ := paillier.Encrypt(testM)

	println("testC::", testC.String())

	decM, _ := paillier.Decrypt(testC)

	println("isRight::", decM.Cmp(testM) == 0)

	r1 := big.NewInt(1561916)
	r2 := big.NewInt(1)
	println(r1.Mod(r1, r2).Int64())
}
