# strings 常用函数
## 判断字符串与子串关系
```
func EqualFold(s, t string) bool // 判断两个utf-8编码字符串，大小写不敏感
func HasPrefix(s, prefix string) bool // 判断s是否有前缀字符串prefix
func Contains(s, substr string) bool // 判断字符串s是否包含子串substr
func ContainsAny(s, chars string) bool // 判断字符串s是否包含字符串chars中的任一字符
func Count(s, sep string) int // 返回字符串s中有几个不重复的sep子串
```

## 获取字符串中子串位置
```
func Index(s, sep string) int // 子串sep在字符串s中第一次出现的位置，不存在则返回-1
func IndexByte(s string, c byte) int // 字符c在s中第一次出现的位置，不存在则返回-
func IndexAny(s, chars string) int // 字符串chars中的任一utf-8码值在s中第一次出现的位置，如果不存在或者chars为空字符串则返回-1
func IndexFunc(s string, f func(rune) bool) int // s中第一个满足函数f的位置i（该处的utf-8码值r满足f(r)==true），不存在则返回-1
func LastIndex(s, sep string) int // 子串sep在字符串s中最后一次出现的位置，不存在则返回-1
```

## 字符串中字符处理
```
func Title(s string) string // 返回s中每个单词的首字母都改为标题格式的字符串拷贝
func ToLower(s string) string // 返回将所有字母都转为对应的小写版本的拷贝
func ToUpper(s string) string // 返回将所有字母都转为对应的大写版本的拷贝
func Repeat(s string, count int) string // 返回count个s串联的字符串
func Replace(s, old, new string, n int) string // 返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串
func Map(mapping func(rune) rune, s string) string // 将s的每一个unicode码值r都替换为mapping(r)，返回这些新码值组成的字符串拷贝。如果mapping返回一个负值，将会丢弃该码值而不会被替换
```

## 字符串前后端处理
```
func Trim(s string, cutset string) string // 返回将s前后端所有cutset包含的utf-8码值都去掉的字符串
func TrimSpace(s string) string // 返回将s前后端所有空白（unicode.IsSpace指定）都去掉的字符串
func TrimFunc(s string, f func(rune) bool) string // 返回将s前后端所有满足f的unicode码值都去掉的字符串
```

## 字符串分割与合并
```
func Fields(s string) []string // 返回将字符串按照空白（通过unicode.IsSpace判断，可以是一到多个连续的空白字符）分割的多个字符串
func Split(s, sep string) []string // 用去掉s中出现的sep的方式进行分割，会分割到结尾，并返回生成的所有片段组成的切片
func Join(a []string, sep string) string // 将一系列字符串连接为一个字符串，之间用sep来分隔
```

## 字符串比较
```
strings.Compare(str1, str2)
```