# help info

## go tools

```bash
go get -u -v github.com/mdempsky/gocode
go get -u -v github.com/ramya-rao-a/go-outline
go get -u -v github.com/acroca/go-symbols
go get -u -v github.com/stamblerre/gocode
go get -u -v golang.org/x/lint/golint
go get -u -v golang.org/x/tools/cmd/guru
go get -u -v golang.org/x/tools/cmd/gorename
go get -u -v golang.org/x/tools/cmd/gopls
```

## bug cautions
- 使用大数big.Int或者大分数big.Rat的时候需要注意，基本上数据类型都是指针，因此对指针进行加减乘除取反取模等运算时，注意不要改变原先的值，同理在foreach的时候也容易犯错，范例如下：
```go
a := big.newInt(1)
b := big.newInt(2)
//此时a的值会变化
c := a.Add(a,b)

//考虑使用临时变量
d := big.NewInt(1)
c := d.Add(a,b)

//foreach中也需要
es := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
//注意这里不能使用var sum *big.Int，因为没有经过初始化的指针是nil，不能对其进行赋值
sum := big.NewInt(0)
for _, e := range es {
	//2*e的累加可以使用链式相加，不能使用sum.Add(sum, e.Add(e,e))
	sum.Add(sum, e).Add(sum, e)
}
```
- 对指针的数组不使用make预定义长度，因为初始值都为nil，无法进一步赋值修改数组
- 在本项目中，涉及到大数运算，因此所有中间运算临时值都不要涉及到int64，容易发生数值计算溢出、转换溢出等问题，且难以定位。一般来说小的数可以成功，而大的数发生错误基本上都是因为溢出问题
