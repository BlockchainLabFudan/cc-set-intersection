package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Paillier struct {
	P  *big.Int
	Q  *big.Int
	SK *big.Int
	PK *big.Int
	N2 *big.Int
}

//初始化，传入两个等长大素数
func (pai *Paillier) Init(p, q *big.Int) {
	pai.P = new(big.Int)
	pai.Q = new(big.Int)
	pai.PK = new(big.Int)
	pai.SK = new(big.Int)
	pai.N2 = new(big.Int)

	pai.P.Set(p)
	pai.Q.Set(q)
	pai.PK.Mul(p, q)
	pai.N2.Mul(pai.PK, pai.PK)

	one := big.NewInt(1)
	p.Sub(p, one)
	q.Sub(q, one)
	pai.SK.Mul(p, q)
}

//加密，传入明文
func (pai *Paillier) Encrypt(m *big.Int) (*big.Int, error) {
	if pai.PK == nil || pai.N2 == nil {
		return nil, fmt.Errorf("pk not init")
	}

	res := big.NewInt(0)
	r, err := pai.randomR()
	if err != nil {
		return nil, err
	}

	res.Add(big.NewInt(1), pai.PK)
	res.Exp(res, m, pai.N2)
	r.Exp(r, pai.PK, pai.N2)
	res.Mul(res, r)
	res.Mod(res, pai.N2)
	return res, nil
}

//从选取随机的r
func (pai *Paillier) randomR() (*big.Int, error) {
	res, err := rand.Int(rand.Reader, pai.N2)
	if err != nil {
		return nil, err
	}
	//欧几里得算法
	//for !euclid(res, pai.N2) {
	//	res,_ = rand.Int(rand.Reader, pai.N2)
	//}
	z := big.NewInt(0)
	one := big.NewInt(1)
	z.GCD(nil, nil, res, pai.N2)
	for z.Cmp(one) != 0 {
		res, _ = rand.Int(rand.Reader, pai.N2)
	}
	return res, nil
}

//欧几里得算法
func euclid(a, b *big.Int) bool {
	d1 := big.NewInt(0)
	d2 := big.NewInt(0)
	zero := big.NewInt(0)
	one := big.NewInt(1)
	d3 := big.NewInt(1)

	switch a.Cmp(b) {
	case -1:
		d1.Set(a)
		d2.Set(b)
		break
	case 0:
		return false
	case 1:
		d1.Set(b)
		d2.Set(a)
	}

	for d3.Cmp(zero) != 0 {
		println("d3::", d3.String())
		println("d2::", d2.String())
		println("d1::", d1.String())
		d3.Mod(d2, d1)
		d2.Set(d1)
		d1.Set(d3)
	}

	return d2.Cmp(one) == 0
}

//解密，传入密文
func (pai *Paillier) Decrypt(c *big.Int) (*big.Int, error) {
	if pai.SK == nil {
		return nil, fmt.Errorf("sk not init")
	}

	m := big.NewInt(0)
	sk_1 := big.NewInt(0)
	//phi(N)的模逆
	sk_1.ModInverse(pai.SK, pai.PK)
	m.Exp(c, pai.SK, pai.N2)
	m.Sub(m, big.NewInt(1))
	m.Div(m, pai.PK)
	m.Mul(m, sk_1)
	m.Mod(m, pai.PK)

	return m, nil
}
