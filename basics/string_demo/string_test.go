package string_demo

import (
	"strconv"
	"strings"
	"testing"
)

func TestStringIterator(t *testing.T) {
	ss := "武汉加油加油！！！"
	t.Log(len(ss), len([]rune(ss)))
	// for range遍历字符串就是rune，如果是传统for循环，i++;ss[i]的模式就是byte
	for _, s := range ss {
		t.Logf("%[1]c %[1]x", s)
	}

	t.Logf("%[1]d %[1]T", ss[0])
	t.Log(uint8('a'))
}

/* 可见遍历时，会自动转化为rune, UTF8字符
--- PASS: TestCreateAndInitStructCase1 (0.00s)
    string_test.go:8: 武 6b66
    string_test.go:8: 汉 6c49
    string_test.go:8: 加 52a0
    string_test.go:8: 油 6cb9
    string_test.go:8: 加 52a0
    string_test.go:8: 油 6cb9
    string_test.go:8: ！ ff01
    string_test.go:8: ！ ff01
    string_test.go:8: ！ ff01
*/

func TestStringSplit(t *testing.T) {
	ss := "123,456,789"
	parts := strings.Split(ss, ",")
	for _, val := range parts {
		t.Log(val)
	}
}

func TestStringJoin(t *testing.T) {
	s1 := "123"
	s2 := "456"
	s3 := "789"
	sa := [3]string{s1, s2, s3}
	ss := strings.Join(sa[:], ",")
	t.Log(ss)
}

func TestStringConvert(t *testing.T) {
	a := 3
	s := strconv.Itoa(a)
	t.Log(s)
	if b, err := strconv.Atoi(s); err == nil {
		t.Log(b)
	}
}

func TestReverseString(t *testing.T) {
	src := "你好abc李焕英"
	dst := reverse([]rune(src))
	t.Logf("%v\n", string(dst))
}

func reverse(s []rune) string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return string(s)
}

// 可以看出，字符串的一个元素就是一个uint8类型
func TestChangeString(t *testing.T) {
	s := "{hello world"
	s1 := string('a') + s[1:]
	t.Logf("%T", s[0])
	t.Log(s[0], 'h')
	t.Log(s, s[0] == 'h')
	t.Log(s1)
	var c = '}'
	var d = '{'
	t.Log(c, d)
}

func TestStringToMap(t *testing.T) {
	ss := "weitong"
	mp := map[rune]int{}
	// range时value为rune
	for index, s := range ss {
		t.Log(index, s)
		mp[s]++
	}
	for k, v := range mp {
		t.Log(k, v)
	}

	t.Log(ss[0:2], len(ss[0:0]), nil)
	t.Logf("%T", ss[0:2])
	res := findNthDigit(12)
	t.Log(res)
}

func TestStringBuilder(t *testing.T) {
	t.Log(findNthDigit(11))
	t.Log(findNthDigit(12))
	t.Log(findNthDigit(100))
}

// 1,2,3,4,5,6,7,8,9,10,11... 第11个数是0，12个是1
func findNthDigit(n int) int {
	var ss strings.Builder
	for num := 1; ; num++ {
		ss.WriteString(strconv.Itoa(num))
		lens := len(ss.String())
		if lens >= n {
			res, _ := strconv.Atoi(string((ss.String())[lens-1]))
			return res
		}
	}
}
