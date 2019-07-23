package main

import (
	"crypto/rand"
	"testing"
)

func TestPaillier(t *testing.T) {
	paillier := PaillierInit()

	testM, _ := rand.Int(rand.Reader, paillier.PK)

	println("testM::", testM.String())

	testC, _ := paillier.Encrypt(testM)

	println("testC::", testC.String())

	decM, _ := paillier.Decrypt(testC)

	println("isRight::", decM.Cmp(testM) == 0)
}
