package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type ElgamalDeform struct {
	//阶
	P *big.Int
	//生成元
	G *big.Int
	//自己的私钥
	SK *big.Int
	//自己的公钥
	PK *big.Int
	//密钥交换得到的公钥
	H *big.Int
}

//初始化，传入素数p和生成元g，传出私钥、公钥、错误信息（若有）
func (ed *ElgamalDeform) Init(p *big.Int, g *big.Int) (*big.Int, *big.Int, error) {
	ed.P = new(big.Int)
	ed.G = new(big.Int)
	ed.PK = new(big.Int)
	ed.SK = new(big.Int)

	ed.P.Set(p)
	ed.G.Set(g)

	sk, err := rand.Int(rand.Reader, ed.P)
	if err != nil {
		return nil, nil, err
	}

	ed.SK.Set(sk)
	ed.PK.Exp(g, sk, p)

	return ed.SK, ed.PK, nil
}

//是否初始化
func (ed *ElgamalDeform) checkInitialized() bool {
	if ed.P == nil || ed.G == nil || ed.PK == nil || ed.SK == nil {
		return false
	}
	return true
}

//传入对方公钥，传出密钥交换的公钥
func (ed *ElgamalDeform) SetPk(pk *big.Int) (*big.Int, error) {
	if !ed.checkInitialized() {
		return nil, fmt.Errorf("not init")
	}
	ed.H = new(big.Int)
	ed.H.Exp(pk, ed.SK, ed.P)
	return ed.H, nil
}

//加密，传入明文，传出密文c1,c2和错误信息（如有）
func (ed *ElgamalDeform) Encrypt(m *big.Int) (*big.Int, *big.Int, error) {
	if !ed.checkInitialized() || ed.H == nil {
		return nil, nil, fmt.Errorf("not init")
	}

	c1 := big.NewInt(0)
	c2 := big.NewInt(0)

	r, err := rand.Int(rand.Reader, ed.P)
	if err != nil {
		return nil, nil, err
	}

	c1.Exp(ed.G, r, ed.P)
	c2.Mul(big.NewInt(0).Exp(ed.G, m, ed.P), big.NewInt(0).Exp(ed.H, r, ed.P))
	return c1, c2, nil
}

//解密步骤1，传入c1，传出c1^sk
func (ed *ElgamalDeform) Decrypt1(c1 *big.Int) (*big.Int, error) {
	if !ed.checkInitialized() || ed.H == nil {
		return nil, fmt.Errorf("not init")
	}

	res := big.NewInt(0).Exp(c1, ed.SK, ed.P)
	return res, nil
}

//解密步骤2，传入c2、对方的解密步骤1输出的c1^sk，传出明文m
func (ed *ElgamalDeform) Decrypt2(c2 *big.Int, c1sk *big.Int) (*big.Int, error) {
	if !ed.checkInitialized() || ed.H == nil {
		return nil, fmt.Errorf("not init")
	}

	modi := big.NewInt(0).ModInverse(big.NewInt(0).Exp(c1sk, ed.SK, ed.P), ed.P)
	res := big.NewInt(0).Mul(c2, modi)
	res.Mod(res, ed.P)
	return res, nil
}
