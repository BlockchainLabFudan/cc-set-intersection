package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Paillier struct {
	P *big.Int
	Q *big.Int
	SK *big.Int
	PK *big.Int
	N2 *big.Int
}

//初始化，传入两个等长大素数
func (this *Paillier) Init(p,q *big.Int) {
	this.P = new(big.Int)
	this.Q = new(big.Int)
	this.PK = new(big.Int)
	this.SK = new(big.Int)
	this.N2 = new(big.Int)

	this.P.Set(p)
	this.Q.Set(q)
	this.PK.Mul(p, q)
	this.N2.Mul(this.PK, this.PK)

	one := big.NewInt(1)
	p.Sub(p, one)
	q.Sub(q, one)
	this.SK.Mul(p, q)
}

//加密，传入明文
func (this *Paillier) Encrypt(m *big.Int) (*big.Int, error){
	if this.PK == nil || this.N2 == nil {
		return nil, fmt.Errorf("pk not init")
	}

	res := big.NewInt(0)
	r, err := this.randomR()
	if err != nil {
		return nil, err
	}

	res.Add(big.NewInt(1), this.PK)
	res.Exp(res, m, this.N2)
	r.Exp(r, this.PK, this.N2)
	res.Mul(res, r)
	res.Mod(res, this.N2)
	return res, nil
}

//从选取随机的r
func (this *Paillier) randomR() (*big.Int, error) {
	res, err := rand.Int(rand.Reader, this.N2)
	if err != nil {
		return nil, err
	}
	//辗转相除法有点问题
	//for !euclid(res, this.N2) {
	//	res,_ = rand.Int(rand.Reader, this.N2)
	//}
	z := big.NewInt(0)
	one := big.NewInt(1)
	z.GCD(nil, nil, res, this.N2)
	for z.Cmp(one) != 0 {
		res,_ = rand.Int(rand.Reader, this.N2)
	}
	return res, nil
}

//欧几里得算法
func euclid(a,b *big.Int) bool {
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
		d3.Mod(d2, d1)
		d2.Set(d1)
		d1.Set(d3)
	}

	return d3.Cmp(one) == 0

}

//解密，传入密文
func (this *Paillier) Decrypt(c *big.Int) (*big.Int, error){
	if this.SK == nil {
		return nil, fmt.Errorf("sk not init")
	}

	m := big.NewInt(0)
	sk_1 := big.NewInt(0)
	//phi(N)的模逆
	sk_1.ModInverse(this.SK, this.PK)
	m.Exp(c, this.SK, this.N2)
	m.Sub(m, big.NewInt(1))
	m.Div(m, this.PK)
	m.Mul(m, sk_1)
	m.Mod(m, this.PK)

	return m, nil
}