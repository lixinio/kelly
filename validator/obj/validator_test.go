package obj

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func checkFail(t *testing.T, obj interface{}) {
	v := NewValidator()
	err := v.Validate(obj)
	_, ok := err.(validator.ValidationErrors)
	if err == nil || !ok {
		t.Error(err)
	}
}

func checkOk(t *testing.T, obj interface{}) {
	v := NewValidator()
	err := v.Validate(obj)
	if err != nil {
		t.Error(err)
	}
}

func toStr(str string) *string {
	return &str
}

func toInt(v int) *int {
	return &v
}

func TestRequired(t *testing.T) {
	type A struct {
		B string  `validate:"required"`
		C *string `validate:"required"`
	}

	checkFail(t, &A{})
	checkFail(t, &A{B: "abc"})
	checkFail(t, &A{B: "", C: toStr("")})
	checkOk(t, &A{B: "abc", C: toStr("")})
}

func TestLength(t *testing.T) {
	// 忽略空值
	type A struct {
		B *string `validate:"omitempty,max=5,min=1"`
		C *int    `validate:"omitempty,gte=0,lte=100"`
		D string  `validate:"omitempty,len=6"`
	}

	checkOk(t, &A{})
	checkOk(t, &A{B: toStr("abcde")})
	checkOk(t, &A{C: toInt(100)})
	checkFail(t, &A{B: toStr("abcdef")})
	checkFail(t, &A{C: toInt(101)})
	checkFail(t, &A{B: toStr("abcdef")})
	checkFail(t, &A{B: toStr("abcdef"), C: toInt(100)})
	checkFail(t, &A{B: toStr("abcde"), C: toInt(101)})
	checkFail(t, &A{B: toStr("abcde"), C: toInt(100), D: "123"})
	checkOk(t, &A{B: toStr("abcde"), C: toInt(100), D: "234567"})
	checkOk(t, &A{B: toStr("a"), C: toInt(0), D: "123456"})
}

func TestFormat(t *testing.T) {
	// 忽略空值
	type A struct {
		B string `validate:"email"`
		C string `validate:"ipv4"`
		D string `validate:"url"`
	}

	checkOk(t, &A{
		B: "a@b.com",
		C: "192.168.1.1",
		D: "http://www.baidu.com",
	})
	checkFail(t, &A{
		B: "a@b",
		C: "192.168.1.1",
		D: "http://www.baidu.com",
	})
	checkFail(t, &A{
		B: "a@b.com",
		C: "192.168.1.256",
		D: "http://www.baidu.com",
	})
	checkFail(t, &A{
		B: "a@b.com",
		C: "192.168.1.1",
		D: "www.baidu.com",
	})
}

func TestDate(t *testing.T) {
	// 忽略空值
	type A struct {
		B string `validate:"date"`
		C string `validate:"datetime"`
	}

	checkOk(t, &A{
		B: "2020-11-22",
		C: "2020-11-22 11:22:33",
	})
	// checkFail(t, &A{
	// 	B: "2021-02-30",
	// 	C: "2020-11-22 11:22:33",
	// })
	checkFail(t, &A{
		B: "2021-02",
		C: "2020-11-22 11:22:33",
	})
	checkFail(t, &A{
		B: "2020-11-22",
		C: "2020-11-22 11:22",
	})
}

func TestMisc(t *testing.T) {
	type A struct {
		B string `validate:"omitempty,lowercase,startswith=abc"`
		C string `validate:"omitempty,uppercase,endswith=efg"`
		D string `validate:"omitempty,alpha"`
		E string `validate:"omitempty,numeric"`
	}

	checkFail(t, &A{B: "Abc"})
	checkFail(t, &A{B: "abef"})
	checkOk(t, &A{B: "abcd"})

	checkFail(t, &A{C: "Abc"})
	checkFail(t, &A{C: "abcdef"})
	checkOk(t, &A{B: "abcdefg"})

	checkFail(t, &A{D: "中国"})
	checkFail(t, &A{D: "123"})
	checkFail(t, &A{D: ",.-%"})
	checkOk(t, &A{D: "efg"})

	checkFail(t, &A{E: "中国"})
	checkFail(t, &A{E: "efg"})
	checkFail(t, &A{E: ",.-%"})
	checkOk(t, &A{E: "123"})
	checkOk(t, &A{E: "-123"})
}

func TestEnum(t *testing.T) {
	type A struct {
		B string `validate:"omitempty,oneof=ab cd ef 'gh i'"`
		C int    `validate:"omitempty,oneof=12 34 56"`
	}

	checkFail(t, &A{B: "中国"})
	checkFail(t, &A{B: "bc"})
	checkFail(t, &A{B: "gh"})
	checkOk(t, &A{B: "cd"})
	checkOk(t, &A{B: "gh i"})

	checkFail(t, &A{C: 23})
	checkFail(t, &A{C: -12})
	checkOk(t, &A{C: 56})
}

func TestArray(t *testing.T) {
	// D 数组里面的字符串只能是枚举中的值
	type A struct {
		B []string `validate:"omitempty,unique,len=3"`
		C []string `validate:"omitempty,max=3,min=2"`
		D []string `validate:"omitempty,dive,oneof=ab cd ef 'gh i'"`
	}

	checkFail(t, &A{B: []string{"a", "b", "a"}})
	checkFail(t, &A{B: []string{"a", "b"}})
	checkFail(t, &A{B: []string{"a", "b", "c", "d"}})
	checkOk(t, &A{B: []string{"a", "b", "c"}})

	checkFail(t, &A{C: []string{"ab"}})
	checkFail(t, &A{C: []string{"ab", "bc", "cd", "ef"}})
	checkOk(t, &A{C: []string{"ab", "cd"}})
	checkOk(t, &A{C: []string{"ab", "cd", "ef"}})

	checkFail(t, &A{D: []string{"ac"}})
	checkFail(t, &A{D: []string{"ab", "cd", "gh"}})
	checkOk(t, &A{D: []string{"ab"}})
	checkOk(t, &A{D: []string{"ab", "cd", "gh i"}})
}

func TestRequiredIf(t *testing.T) {
	// 如果C=="test", B必须存在
	type A struct {
		B string `validate:"required_if=C test"`
		C string `validate:"-"`
	}

	checkFail(t, &A{C: "test"})
	checkOk(t, &A{C: "test2"})
	checkOk(t, &A{C: "test1"})
}

func TestRequiredUnless(t *testing.T) {
	// 如果C!="test", B必须存在
	type A struct {
		B string `validate:"required_unless=C test"`
		C string `validate:"-"`
	}

	checkFail(t, &A{C: "test2"})
	checkFail(t, &A{C: "test1"})
	checkOk(t, &A{C: "test"})
}

func TestRequiredWith(t *testing.T) {
	// B/C要不都存在，要不都不存在
	type A struct {
		B *string `validate:"required_with=C"`
		C *string `validate:"required_with=B"`
	}

	checkFail(t, &A{B: toStr("x")})
	checkFail(t, &A{B: toStr("y")})
	checkFail(t, &A{C: toStr("x")})
	checkFail(t, &A{C: toStr("y")})
	checkOk(t, &A{B: toStr("x"), C: toStr("y")})
	checkOk(t, &A{B: toStr("y"), C: toStr("y")})
}

func TestRequiredWithout(t *testing.T) {
	// B/C 至少一个存在
	type A struct {
		B *string `validate:"required_without=C"`
		C *string `validate:"required_without=B"`
	}

	checkFail(t, &A{})
	checkOk(t, &A{B: toStr("x")})
	checkOk(t, &A{B: toStr("y")})
	checkOk(t, &A{C: toStr("x")})
	checkOk(t, &A{C: toStr("y")})
	checkOk(t, &A{B: toStr("x"), C: toStr("y")})
	checkOk(t, &A{B: toStr("y"), C: toStr("y")})
}

// gtfield / gtefield / ltefield / ltfield / nefield
func TestEqfield(t *testing.T) {
	// B 必须等于 C
	type A struct {
		B string `validate:"eqfield=C"`
		C string
	}

	checkFail(t, &A{B: "test"})
	checkFail(t, &A{B: "test", C: "test2"})
	checkOk(t, &A{B: "test", C: "test"})
	checkOk(t, &A{B: "test2", C: "test2"})
}

func TestNefield(t *testing.T) {
	// B 必须不等于 C
	type A struct {
		B string `validate:"nefield=C"`
		C string
	}

	checkFail(t, &A{B: "test", C: "test"})
	checkFail(t, &A{B: "test2", C: "test2"})
	checkOk(t, &A{B: "test"})
	checkOk(t, &A{B: "test", C: "test2"})
}

func TestEqcsfield(t *testing.T) {
	type C struct {
		D string
	}

	// B 必须等于 E.D （跨结构）
	type A struct {
		B string `validate:"eqcsfield=E.D"`
		E *C
	}

	checkFail(t, &A{B: "test"})
	checkFail(t, &A{B: "test", E: &C{D: "test2"}})
	checkOk(t, &A{B: "test", E: &C{D: "test"}})
	checkOk(t, &A{B: "test2", E: &C{D: "test2"}})
}
