package main

import (
	"github.com/BlockchainLabFudan/cc-set-intersection/model"
	"log"
	"math/big"
)

var zero = big.NewInt(0)

func PaillierInit() *model.Paillier {
	//paillier初始化
	paillier := new(model.Paillier)

	p1 := new(big.Int)
	p2 := new(big.Int)
	p1.SetString("242661090146032969904098483991985908921", 10) // octal
	p2.SetString("215662396313044988944834777682074105079", 10) // octal

	paillier.Init(p1, p2)
	return paillier
}

func EdInit() *model.ElgamalDeform {

	ed := new(model.ElgamalDeform)

	p := new(big.Int)
	g := new(big.Int)
	p.SetString("242661090146032969904098483991985908921", 10) // octal
	g.SetString("75419874865198741652318945128641254566", 10)  // octal
	//p2.SetString("215662396313044988944834777682074105079", 10) // octal
	//g, _ := rand.Int(rand.Reader, p)

	_, _, _ = ed.Init(p, g)
	return ed
}

func PrintCoef(p *model.Polynomial) {
	println("----------------")
	for _, c := range p.Coefficients {
		println(c.Int64())
	}
}

func NewUser(p *model.Paillier) *model.Polynomial {
	user := new(model.Polynomial)
	err := user.Init(p)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func Cal(attrs []*big.Int, p *model.Polynomial) {
	i := 0
	for _, attr := range attrs {
		ok, err := p.CalculateForDec(attr)
		if err != nil {
			log.Fatal(err)
		}
		ookk := ok.Cmp(zero) == 0
		if ookk {
			i++
		}
		log.Println("测试交集：attr::", attr.Int64(), " intersection::", ookk, " value:", ok)
	}
	log.Printf("%d个相同,%d个不同", i, len(attrs)-i)
}

func Equal(p *model.Polynomial, q *model.Polynomial) bool {
	if p.Length != q.Length {
		return false
	}
	for i := int64(0); i < p.Length; i++ {
		if p.Coefficients[i].Cmp(q.Coefficients[i]) != 0 {
			return false
		}
	}
	return true
}
