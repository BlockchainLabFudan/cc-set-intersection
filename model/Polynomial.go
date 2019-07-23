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
	Length int64
	//加密前
	MCoefficients []*big.Int
	//加密后
	ECoefficients1 []*big.Int
	ECoefficients2 []*big.Int
	Ed             *ElgamalDeform
}

//初始化，配置一个加密算法
func Init(ed *ElgamalDeform) (*Polynomial, error) {
	if ed == nil {
		return nil, fmt.Errorf("ed is nil")
	}
	var p = new(Polynomial)
	p.Ed = ed
	p.Length = int64(0)
	p.IsEncrypted = false
	return p, nil
}

//多项式加密系数，加密类中的系数，为协议第二步，加密多项式
func (p *Polynomial) EncryptPolynomial() error {
	tempC1, tempC2, err := p.encrypt(p.MCoefficients)
	if err != nil {
		return err
	}
	p.ECoefficients1 = tempC1
	p.ECoefficients2 = tempC2
	p.IsEncrypted = true
	return nil
}

//协议第三步和第四步，选取多项式r1r2，计算r1fa+r2fb
func CalGcdPolynomial(fa, fb *Polynomial) (*Polynomial, error) {
	if !fa.IsEncrypted || !fb.IsEncrypted {
		return nil, fmt.Errorf("fa or fb not encrypted")
	}
	//生成随机多项式，这里调换fafb是为了长度相同
	r1, err := fb.RandomPolynomial()
	if err != nil {
		return nil, err
	}
	r2, err := fa.RandomPolynomial()
	if err != nil {
		return nil, err
	}

	//	计算整合的多项式
	left, err := fa.Mul(r1)
	if err != nil {
		return nil, err
	}
	right, err := fb.Mul(r2)
	if err != nil {
		return nil, err
	}

	res, err := left.Add(right)
	//res, err := fa.Add(fb)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//协议第五步，交换各自对整合多项式的解密分享c1sk
func (p *Polynomial) ChangeShare(c1s []*big.Int) ([]*big.Int, error) {
	var share []*big.Int
	for _, c1 := range c1s {
		c1sk, err := p.Ed.Decrypt1(c1)
		if err != nil {
			return nil, err
		}
		share = append(share, c1sk)
	}
	return share, nil
}

//协议第六步，最后解密出明文多项式
func (p *Polynomial) DecryptPolynomial(c1sk, c2s []*big.Int) (*Polynomial, error) {
	var coefficients []*big.Int
	for index, c2 := range c2s {
		coefficient, err := p.Ed.Decrypt2(c2, c1sk[index])
		if err != nil {
			return nil, err
		}
		coefficients = append(coefficients, coefficient)
	}
	res := Polynomial{false, p.Length, coefficients, nil, nil, p.Ed}
	return &res, nil
}

func (p *Polynomial) encrypt(cs []*big.Int) ([]*big.Int, []*big.Int, error) {
	var res1 []*big.Int
	var res2 []*big.Int
	fmt.Println("enc:", cs)
	for _, c := range cs {
		tempC1, tempC2, err := p.Ed.Encrypt(c)
		if err != nil {
			return nil, nil, err
		}
		res1 = append(res1, tempC1)
		res2 = append(res2, tempC2)
	}
	return res1, res2, nil
}

//多项式求值
func (p *Polynomial) CalculatePolynomial(x *big.Int) (*big.Int, error) {
	if p.IsEncrypted {
		return nil, fmt.Errorf("polynomial is encrypted")
	}

	aI := big.NewInt(1)
	up := big.NewInt(1)
	down := big.NewInt(1)
	//fmt.Println(p.MCoefficients)
	for _, coef := range p.MCoefficients {
		if aI.Sign() >= 0 {
			up.Mul(up, big.NewInt(1).Exp(coef, aI, p.Ed.P))
		} else {
			down.Mul(down, big.NewInt(1).Exp(coef, big.NewInt(0).Abs(aI), p.Ed.P))
		}
		aI.Mul(aI, x)
	}
	if down.Cmp(big.NewInt(1)) == 0 {
		return big.NewInt(0).Mod(up, p.Ed.P), nil
	}
	return big.NewInt(0).Mod(big.NewInt(0).Mul(up, big.NewInt(0).ModInverse(down, p.Ed.P)), p.Ed.P), nil
}

//多项式加法
func (p *Polynomial) Add(g *Polynomial) (*Polynomial, error) {
	//必须是加密多项式
	if !p.IsEncrypted || !g.IsEncrypted {
		return nil, fmt.Errorf("not encrypted")
	}

	var coef1 []*big.Int
	var coef2 []*big.Int
	index := int64(0)
	//密文分别相乘
	for ; index < p.Length && index < g.Length; index++ {
		c1 := big.NewInt(0).Mul(g.ECoefficients1[index], p.ECoefficients1[index])
		c2 := big.NewInt(0).Mul(g.ECoefficients2[index], p.ECoefficients2[index])
		c1.Mod(c1, p.Ed.P)
		c2.Mod(c2, p.Ed.P)
		coef1 = append(coef1, c1)
		coef2 = append(coef2, c2)
	}
	for index < p.Length {
		coef1 = append(coef1, p.ECoefficients1[index])
		coef2 = append(coef2, p.ECoefficients2[index])
	}
	for index < g.Length {
		coef1 = append(coef1, g.ECoefficients1[index])
		coef2 = append(coef2, g.ECoefficients2[index])
	}

	res := Polynomial{true, p.Length, nil, coef1, coef2, p.Ed}
	return &res, nil
}

//多项式乘法,f是加密后的,g是未加密的
func (p *Polynomial) Mul(g *Polynomial) (*Polynomial, error) {
	//必须是加密多项式和非加密多项式
	if !p.IsEncrypted || g.IsEncrypted {
		return nil, fmt.Errorf("not admitted")
	}
	//n2 := big.NewInt(0).Mul(p.Ed.P, p.Ed.P)

	var coef1 []*big.Int
	var coef2 []*big.Int

	//密文的乘方
	for i1 := int64(0); i1 < p.Length+g.Length-int64(1); i1++ {
		tempC1 := big.NewInt(1)
		tempC2 := big.NewInt(1)
		for i2 := int64(0); i2 <= i1; i2++ {
			if i2 >= p.Length || i1-i2 >= g.Length {
				continue
			}
			tempC1.Mul(tempC1, big.NewInt(0).Exp(p.ECoefficients1[i2], g.MCoefficients[i1-i2], p.Ed.P))
			tempC1.Mod(tempC1, p.Ed.P)

			tempC2.Mul(tempC2, big.NewInt(0).Exp(p.ECoefficients2[i2], g.MCoefficients[i1-i2], p.Ed.P))
			tempC2.Mod(tempC2, p.Ed.P)
		}
		coef1 = append(coef1, tempC1)
		coef2 = append(coef2, tempC2)
	}

	res := Polynomial{true, p.Length + g.Length - int64(1), nil, coef1, coef2, p.Ed}
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
	p1 := big.NewInt(2).Sub(p.Ed.P, big.NewInt(1))
	for _, v := range res {
		//fmt.Println(v)
		if len(v) != 1 || !v[0].IsInt() {
			return fmt.Errorf("unknown set error")
		}
		coef = append(coef, big.NewInt(0).Mod(v[0].Num(), p1))
	}
	p.MCoefficients = coef
	p.Length = int64(len(coef))
	return nil
}

//随机未加密多项式
func (p *Polynomial) RandomPolynomial() (*Polynomial, error) {
	var coef []*big.Int
	for i := int64(0); i < p.Length; i++ {
		temp, err := rand.Int(rand.Reader, p.Ed.P)
		if err != nil {
			return nil, err
		}
		coef = append(coef, temp)
	}
	res := Polynomial{false, p.Length, coef, nil, nil, p.Ed}
	return &res, nil
}
