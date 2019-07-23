package main

import (
	"github.com/BlockchainLabFudan/cc-set-intersection/model"
	"log"
	"math/big"
)

var one = big.NewInt(1)

func EdInit() *model.ElgamalDeform {

	ed := new(model.ElgamalDeform)

	p := new(big.Int)
	g := new(big.Int)
	p.SetString("242661090146032969904098483991985908921", 10) // octal
	//p.SetString("67", 10) // octal
	g.SetString("75419874865198741652318945128641254566", 10) // octal
	//g.SetString("5964294", 10) // octal
	//p2.SetString("215662396313044988944834777682074105079", 10) // octal
	//g, _ := rand.Int(rand.Reader, p)

	_, _, _ = ed.Init(p, g)
	return ed
}

func PrintCoef(p *model.Polynomial) {
	println("----------------")
	for _, c := range p.MCoefficients {
		println(c.Int64())
	}
}

func Cal(attrs []*big.Int, p *model.Polynomial) {
	i := 0
	for _, attr := range attrs {
		ok, err := p.CalculatePolynomial(attr)
		if err != nil {
			log.Fatal(err)
		}
		ookk := ok.Cmp(one) == 0
		if ookk {
			i++
		}
		log.Println("测试交集：attr::", attr.String(), " intersection::", ookk, " value:", ok)
	}
	log.Printf("%d个相同,%d个不同", i, len(attrs)-i)
}
