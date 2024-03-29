package model

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

func copyRat(r *big.Rat) *big.Rat {
	return big.NewRat(1, 1).Set(r)
}

func SolveGaussian(eqM [][]*big.Rat, printTriangularForm bool) (res [][]*big.Rat, err error) {

	// 判断行数是否比系数多
	if len(eqM) > len(eqM[0])-1 {
		err = errors.New("the number of equations can not be greater than the number of variables")
		return
	}

	// 判断是否有重复的行
	dl, i, j := containsDuplicatesLines(eqM)
	if dl {
		err = fmt.Errorf("provided matrix contains duplicate lines (%d and %d)", i+1, j+1)
		return
	}

	for i := 0; i < len(eqM)-1; i++ {
		//fmt.Println("--------+", eqM)
		eqM = sortMatrix(eqM, i)
		//fmt.Println("---------", eqM)
		varC := big.NewRat(1, 1)
		for k := i; k < len(eqM); k++ {
			if k == i {
				varC.Set(eqM[k][i])
			} else {
				multipliedLine := make([]*big.Rat, len(eqM[i]))
				for z, zv := range eqM[i] {
					//multipliedLine[z] = zv.Multiply(eqM[k][i].Divide(varC)).MultiplyByNum(-1)
					// fmt.Println(z, zv, varC, eqM[k][i])

					// multipliedLine[z] = zv.Neg(zv.Mul(zv, big.NewRat(1, 1).Quo(eqM[k][i], varC)))

					// #### 要重新 New 一个 Rat ####
					tmp := copyRat(zv)
					tmp.Mul(tmp, big.NewRat(1, 1).Quo(eqM[k][i], varC))
					tmp.Neg(tmp)
					//fmt.Println("---------", tmp)
					multipliedLine[z] = tmp

				}
				newLine := make([]*big.Rat, len(eqM[k]))
				//fmt.Println("---------", eqM)
				//fmt.Println("---------", multipliedLine)
				for z, zv := range eqM[k] {
					newLine[z] = big.NewRat(1, 1).Add(zv, multipliedLine[z])
				}
				//fmt.Println("+++++++++++", newLine)
				eqM[k] = newLine
				//fmt.Println("+++++++++++", eqM)
			}
		}
	}
	//fmt.Println(eqM)

	// 移除为0的行，并且反转
	var resultEqM [][]*big.Rat
	for i := len(eqM) - 1; i >= 0; i-- {
		//if !rational.RationalsAreNull(eqM[i]) {
		resultEqM = append(resultEqM, eqM[i])
		//}
	}
	//fmt.Println("---------", resultEqM)

	getFirstNonZeroIndex := func(sl []*big.Rat) (index int) {
		zero := big.NewInt(0)
		for i, v := range sl {
			if v.Num().Cmp(zero) != 0 {
				index = i
				return
			}
		}
		return
	}

	// Back substitution.
	for z := 0; z < len(resultEqM)-1; z++ {
		var processIndex int
		var firstLine []*big.Rat
		for i := z; i < len(resultEqM); i++ {
			v := resultEqM[i]
			if i == z {
				//fmt.Println(v)
				processIndex = getFirstNonZeroIndex(v)
				firstLine = resultEqM[i]
			} else {
				// mult := v[processIndex].Quo(v[processIndex], firstLine[processIndex]).Mul(v[processIndex], big.NewRat(-1, 1))

				mult := copyRat(v[processIndex])
				mult.Quo(mult, firstLine[processIndex])
				mult.Mul(mult, big.NewRat(-1, 1))
				//fmt.Println("---------", resultEqM)
				//fmt.Println("---------", firstLine)
				//fmt.Println("---------", v)
				for j, jv := range v {
					resultEqM[i][j] = copyRat(firstLine[j])
					resultEqM[i][j].Mul(resultEqM[i][j], mult)
					resultEqM[i][j].Add(resultEqM[i][j], jv)
				}
				//fmt.Println("---------", resultEqM)
				//fmt.Println("---------", v)
			}
		}
	}
	//fmt.Println("---------", resultEqM)

	if printTriangularForm {
		for i := len(resultEqM) - 1; i >= 0; i-- {
			var str string
			for _, jv := range resultEqM[i] {
				temp, _ := jv.Float64()
				str += strconv.FormatFloat(temp, 'f', 2, 64) + ","
			}
			str = str[:len(str)-1]
			fmt.Println(str)
		}
	}

	//fmt.Println("---------", resultEqM)
	// Calculating variables.
	res = make([][]*big.Rat, len(eqM[0])-1)
	if getFirstNonZeroIndex(resultEqM[0]) == len(resultEqM[0])-2 {
		// All the variables have been found.
		for i, iv := range resultEqM {
			index := len(res) - 1 - i
			res[index] = append(res[index], iv[len(iv)-1].Quo(iv[len(iv)-1], iv[len(resultEqM)-1-i]))
		}
	} else {
		// Some variables remained unknown.
		var unknownStart, unknownEnd int
		for i, iv := range resultEqM {
			fnz := getFirstNonZeroIndex(iv)
			var firstRes []*big.Rat
			firstRes = append(firstRes, iv[len(iv)-1].Quo(iv[len(iv)-1], iv[fnz]))
			if i == 0 {
				unknownStart = fnz + 1
				unknownEnd = len(iv) - 2
				for j := unknownEnd; j >= unknownStart; j-- {
					res[j] = []*big.Rat{big.NewRat(0, 1)}
					firstRes = append(firstRes, iv[j].Quo(iv[j], iv[fnz]))
				}
			} else {
				for j := unknownEnd; j >= unknownStart; j-- {
					firstRes = append(firstRes, iv[j].Quo(iv[j], iv[fnz]))
				}
			}
			res[fnz] = firstRes
		}
	}

	return
}

func sortMatrix(m [][]*big.Rat, initRow int) (m2 [][]*big.Rat) {
	indexed := make(map[int]bool)

	for i := 0; i < initRow; i++ {
		m2 = append(m2, m[i])
		indexed[i] = true
	}

	greaterThanMax := func(rr1, rr2 []*big.Rat) (greater bool) {
		for i := 0; i < len(rr1); i++ {
			if big.NewRat(1, 1).Abs(rr1[i]).Cmp(big.NewRat(1, 1).Abs(rr2[i])) > 0 {
				greater = true
				return
			} else if big.NewRat(1, 1).Abs(rr1[i]).Cmp(big.NewRat(1, 1).Abs(rr2[i])) < 0 {
				return
			}
		}
		return
	}

	type maxStruct struct {
		index   int
		element []*big.Rat
	}

	for i := initRow; i < len(m); i++ {
		//fmt.Println(i, m2)
		var maxElement []*big.Rat
		for index := 0; index < len(m[i]); index++ {
			maxElement = append(maxElement, big.NewRat(0, 1))
		}
		max := maxStruct{-1, maxElement}
		var firstNotIndexed int
		for k, kv := range m {
			if !indexed[k] {
				firstNotIndexed = k
				if greaterThanMax(kv, max.element) {
					max.index = k
					max.element = kv
				}
			}
		}
		//fmt.Println(i, max.element)
		if max.index != -1 {
			m2 = append(m2, max.element)
			indexed[max.index] = true
		} else {
			m2 = append(m2, m[firstNotIndexed])
			indexed[firstNotIndexed] = true
		}
	}

	return
}

func containsDuplicatesLines(eqM [][]*big.Rat) (contains bool, l1, l2 int) {
	for i := 0; i < len(eqM); i++ {
		for j := i + 1; j < len(eqM); j++ {
			var equalElements int
			for k := 0; k < len(eqM[i]); k++ {
				if eqM[i][k].Cmp(eqM[j][k]) == 0 {
					equalElements++
				} else {
					break
				}
			}
			if equalElements == len(eqM[i]) {
				contains = true
				l1 = i
				l2 = j
				return
			}
		}
	}
	return
}
