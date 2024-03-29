package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

//多项式
type Polynomial struct {
	IsEncrypted bool
	//k阶
	Length       int64
	Coefficients []*big.Int
	Paillier     *Paillier
}

func (p *Polynomial) Init(paillier *Paillier) error {
	if paillier == nil {
		return fmt.Errorf("paillier is nil")
	}
	p.Paillier = paillier
	p.Length = int64(0)
	p.IsEncrypted = false
	return nil
}

//多项式加密系数，加密类中的系数
func (p *Polynomial) EncryptPolynomial() error {
	temp, err := p.encrypt(p.Coefficients)
	if err != nil {
		return err
	}
	p.Coefficients = temp
	p.IsEncrypted = true
	return nil
}

//多项式解密系数，解密类中的系数
func (p *Polynomial) DecryptPolynomial() error {
	temp, err := p.decrypt(p.Coefficients)
	if err != nil {
		return err
	}
	p.Coefficients = temp
	p.IsEncrypted = false
	return nil
}

//多项式加密系数，加密外来系数
func (p *Polynomial) EncryptPolynomialForOthers(cs []*big.Int) ([]*big.Int, error) {
	return p.encrypt(cs)
}

func (p *Polynomial) encrypt(cs []*big.Int) ([]*big.Int, error) {
	var res []*big.Int
	for _, c := range cs {
		temp, err := p.Paillier.Encrypt(c)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}

func (p *Polynomial) decrypt(cs []*big.Int) ([]*big.Int, error) {
	var res []*big.Int
	for _, c := range cs {
		temp, err := p.Paillier.Decrypt(c)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}

//多项式求值,加密函数
func (p *Polynomial) CalculateForEnc(x *big.Int) (*big.Int, error) {
	if !p.IsEncrypted {
		err := p.EncryptPolynomial()
		if err != nil {
			return nil, err
		}
	}

	aI := big.NewInt(1)
	res := big.NewInt(1)
	temp := big.NewInt(1)
	for _, coef := range p.Coefficients {
		res.Mul(res, temp.Exp(coef, aI, p.Paillier.PK))
		aI.Mul(aI, x)
	}
	return res.Mod(res, p.Paillier.PK), nil
}

//多项式求值,解密函数
func (p *Polynomial) CalculateForDec(x *big.Int) (*big.Int, error) {
	if p.IsEncrypted {
		err := p.DecryptPolynomial()
		if err != nil {
			return nil, err
		}
	}

	aI := big.NewInt(1)
	res := big.NewInt(0)
	temp := big.NewInt(1)
	for _, coef := range p.Coefficients {
		res.Add(res, temp.Mul(coef, aI))
		aI.Mul(aI, x)
	}
	return res.Mod(res, p.Paillier.PK), nil
}

//多项式加法
func (p *Polynomial) Add(g *Polynomial) (*Polynomial, error) {
	if p.Length != g.Length {
		return nil, fmt.Errorf("length not adjust")
	}

	//必须是加密多项式
	if !p.IsEncrypted || !g.IsEncrypted {
		return nil, fmt.Errorf("not encrypted")
	}

	var coef []*big.Int
	for index, x := range p.Coefficients {
		temp := big.NewInt(0)
		coef = append(coef, temp.Mul(x, g.Coefficients[index]))
	}

	res := Polynomial{true, p.Length, coef, p.Paillier}
	return &res, nil
}

//多项式乘法,f是加密后的,g是未加密的
func (p *Polynomial) Mul(g *Polynomial) (*Polynomial, error) {
	if p.Length != g.Length {
		return nil, fmt.Errorf("length not adjust")
	}

	//必须是加密多项式和非加密多项式
	if !p.IsEncrypted || g.IsEncrypted {
		return nil, fmt.Errorf("not admitted")
	}

	var coef []*big.Int

	for i1 := int64(0); i1 < p.Length+g.Length-int64(1); i1++ {
		temp1 := big.NewInt(1)
		for i2 := int64(0); i2 <= i1; i2++ {
			if i2 >= p.Length || i1-i2 >= g.Length {
				continue
			}
			temp2 := big.NewInt(0)
			temp2.Exp(p.Coefficients[i2], g.Coefficients[i1-i2], p.Paillier.N2)

			temp1.Mul(temp1, temp2)
			temp1.Mod(temp1, p.Paillier.N2)
		}
		coef = append(coef, temp1.Mod(temp1, p.Paillier.N2))
	}

	res := Polynomial{true, p.Length, coef, p.Paillier}
	return &res, nil
}

//多项式生成,使用高斯消元法得到系数
func (p *Polynomial) SetPolynomial(attrs []*big.Int) error {
	nr := func(i int64) *big.Rat {
		return big.NewRat(i, 1)
	}

	var equations [][]*big.Rat

	for index := int64(0); index <= int64(len(attrs)); index++ {
		var temp []*big.Rat
		coef := big.NewInt(index)
		prod := big.NewInt(index)
		//第一个参数为1
		temp = append(temp, nr(int64(1)))

		now := big.NewInt(index)
		value := big.NewInt(1)

		for _, attr := range attrs {
			value.Mul(value, big.NewInt(1).Sub(now, attr))
			temp = append(temp, big.NewRat(1, 1).SetInt(prod))
			prod.Mul(prod, coef)
		}
		temp = append(temp, big.NewRat(1, 1).SetInt(value))

		equations = append(equations, temp)
	}

	res, gausErr := SolveGaussian(equations, false)
	if gausErr != nil {
		return gausErr
	}

	var coef []*big.Int
	for _, v := range res {
		//fmt.Println(v)
		if len(v) != 1 || !v[0].IsInt() {
			return fmt.Errorf("unknown set error")
		}
		coef = append(coef, big.NewInt(0).Mod(v[0].Num(), p.Paillier.PK))
	}
	p.Coefficients = coef
	p.Length = int64(len(coef))
	return nil
}

//随机未加密多项式
func (p *Polynomial) RandomPolynomial() (*Polynomial, error) {
	var coef []*big.Int
	for i := int64(0); i < p.Length; i++ {
		temp, err := rand.Int(rand.Reader, p.Paillier.PK)
		if err != nil {
			return nil, err
		}
		coef = append(coef, temp)
	}
	res := Polynomial{false, p.Length, coef, p.Paillier}
	return &res, nil
}
