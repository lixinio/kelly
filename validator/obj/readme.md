# 规则
``` c
-        //忽略字段，告诉验证跳过这个struct字段;这对于忽略嵌入式结构的验证尤其方便。 （用法： - ）
|        //这是'or'运算符，允许使用和接受多个验证器。 （用法：rbg | rgba）< - 这将允许接受rgb或rgba颜色。这也可以与'and'结合使用（例如：用法：omitempty，rgb | rgba）
structonly    //当遇到嵌套结构的字段并包含此标志时，将运行嵌套结构上的任何验证，但不会验证任何嵌套结构字段。如果您在程序内部知道结构有效，但需要验证它是否已分配，这非常有用。注意：结构本身只能使用“required”和“omitempty”。
nostructlevel    //与structonly标记相同，但不会运行任何结构级别验证。
omitempty    //允许条件验证，例如，如果字段未设置值（由“required”验证器确定），则其他验证（如min或max）将不会运行，但如果设置了值，则验证将运行。
dive    //这告诉验证者潜入切片，数组或映射，并使用后面的验证标记验证切片，数组或映射的该级别。还支持多维嵌套，您希望dive的每个级别都需要另一个dive标签。dive有一些子标签，'keys'和'endkeys'，请参阅下面的keys和endkeys部分
```
``` c
required    //这将验证该值不是数据类型的默认零值。数字不为０，字符串不为 " ", slices, maps, pointers, interfaces, channels and functions 不为 nil
isdefault    //这验证了该值是默认值，几乎与所需值相反。
len=10    //对于数字，长度将确保该值等于给定的参数。对于字符串，它会检查字符串长度是否与字符数完全相同。对于切片，数组和map，验证元素个数。
max=10    //对于数字，max将确保该值小于或等于给定的参数。对于字符串，它会检查字符串长度是否最多为该字符数。对于切片，数组和map，验证元素个数。
min=10
eq=10    //对于字符串和数字，eq将确保该值等于给定的参数。对于切片，数组和map，验证元素个数。
ne=10    //和eq相反
oneof=red green (oneof=5 7 9)    //对于字符串，整数和uint，oneof将确保该值是参数中的值之一。参数应该是由空格分隔的值列表。值可以是字符串或数字。
gt=10    //对于数字，这将确保该值大于给定的参数。对于字符串，它会检查字符串长度是否大于该字符数。对于切片，数组和map，它会验证元素个数。
gt    //对于time.Time确保时间值大于time.Now.UTC（）
gte=10    //大于等于
gte    //对于time.Time确保时间值大于或等于time.Now.UTC（）
lt=10    //小于
lt    //对于time.Time确保时间值小于time.Now.UTC（）
lte=10    //小于等于
lte    //对于time.Time确保时间值小于等于time.Now.UTC（）
```

``` c
unique    //对于数组和切片，unique将确保没有重复项。对于map，unique将确保没有重复值。
alpha    //这将验证字符串值是否仅包含ASCII字母字符
alphanum    //这将验证字符串值是否仅包含ASCII字母数字字符
alphaunicode    //这将验证字符串值是否仅包含unicode字符
alphanumunicode    //这将验证字符串值是否仅包含unicode字母数字字符
numeric    //这将验证字符串值是否包含基本数值。基本排除指数等...对于整数或浮点数，它返回true。
hexadecimal    //这将验证字符串值是否包含有效的十六进制
hexcolor    //这验证字符串值包含有效的十六进制颜色，包括＃标签（＃）
rgb    //这将验证字符串值是否包含有效的rgb颜色
rgba    //这将验证字符串值是否包含有效的rgba颜色
hsl    //这将验证字符串值是否包含有效的hsl颜色
hsla    //这将验证字符串值是否包含有效的hsla颜色
email    //这验证字符串值包含有效的电子邮件这可能不符合任何rfc标准的所有可能性，但任何电子邮件提供商都不接受所有可能性
file    //这将验证字符串值是否包含有效的文件路径，并且该文件存在于计算机上。这是使用os.Stat完成的，它是一个独立于平台的函数。
url    //这会验证字符串值是否包含有效的url这将接受golang请求uri接受的任何url，但必须包含一个模式，例如http：//或rtmp：//
uri    //这验证了字符串值包含有效的uri。这将接受uri接受的golang请求的任何uri
base64    //这将验证字符串值是否包含有效的base64值。虽然空字符串是有效的base64，但这会将空字符串报告为错误，如果您希望接受空字符串作为有效字符，则可以将此字符串与omitempty标记一起使用。
base64url    //这会根据RFC4648规范验证字符串值是否包含有效的base64 URL安全值。尽管空字符串是有效的base64 URL安全值，但这会将空字符串报告为错误，如果您希望接受空字符串作为有效字符，则可以将此字符串与omitempty标记一起使用。
btc_addr    //这将验证字符串值是否包含有效的比特币地址。检查字符串的格式以确保它匹配P2PKH，P2SH三种格式之一并执行校验和验证
btc_addr_bech32    //这验证了字符串值包含bip-0173定义的有效比特币Bech32地址（https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki）特别感谢Pieter Wuille提供的参考实现。
eth_addr    //这将验证字符串值是否包含有效的以太坊地址。检查字符串的格式以确保它符合标准的以太坊地址格式完全验证被https://github.com/golang/crypto/pull/28阻止
contains=@    //这将验证字符串值是否包含子字符串值
containsany=!@#?    //这将验证字符串值是否包含子字符串值中的任何Unicode code points。
containsrune=@    //这将验证字符串值是否包含提供的符文值。
excludes=@    //这验证字符串值不包含子字符串值。
excludesall=!@#?    //这将验证字符串值在子字符串值中是否包含任何Unicode code points。
excludesrune=@    //这将验证字符串值是否包含提供的符文值。
```

``` c
isbn    //这将验证字符串值是否包含有效的isbn10或isbn13值。
isbn10    //这将验证字符串值是否包含有效的isbn10值。
isbn13    //这将验证字符串值是否包含有效的isbn13值。
uuid    //这将验证字符串值是否包含有效的UUID。
uuid3    //这将验证字符串值是否包含有效的版本3 UUID。
uuid4    //这将验证字符串值是否包含有效的版本4 UUID。
uuid5    //这将验证字符串值是否包含有效的版本5 UUID。
```

``` c
ascii    //这将验证字符串值是否仅包含ASCII字符。注意：如果字符串为空，则验证为true
printascii    //这将验证字符串值是否仅包含可打印的ASCII字符。注意：如果字符串为空，则验证为true。
multibyte    //这将验证字符串值是否包含一个或多个多字节字符。注意：如果字符串为空，则验证为true
datauri    //这将验证字符串值是否包含有效的DataURI。注意：这也将验证数据部分是否有效base64
latitude    //这将验证字符串值是否包含有效的纬度。
longitude    //这将验证字符串值是否包含有效经度。
ssn    //这将验证字符串值是否包含有效的美国社会安全号码。
ip    //这将验证字符串值是否包含有效的IP地址
ipv4    //这将验证字符串值是否包含有效的v4 IP地址
ipv6    //这将验证字符串值是否包含有效的v6 IP地址
cidr    //这将验证字符串值是否包含有效的CIDR地址
cidrv4    //这将验证字符串值是否包含有效的v4 CIDR地址
cidrv5    //这将验证字符串值是否包含有效的v5 CIDR地址
tcp_addr    //这将验证字符串值是否包含有效的可解析TCP地址
tcp4_addr    //这将验证字符串值是否包含有效的可解析v4 TCP地址
tcp6_addr    //这将验证字符串值是否包含有效的可解析v6 TCP地址
udp_addr    //这将验证字符串值是否包含有效的可解析UDP地址
udp4_addr    //这将验证字符串值是否包含有效的可解析v4 UDP地址
udp6_addr    //这将验证字符串值是否包含有效的可解析v6 UDP地址
ip_addr    //这将验证字符串值是否包含有效的可解析IP地址
ip4_addr    //这将验证字符串值是否包含有效的可解析v4 IP地址
ip6_addr    //这将验证字符串值是否包含有效的可解析v6 IP地址
unix_addr    //这将验证字符串值是否包含有效的Unix地址
```
``` c
mac    //这将验证字符串值是否包含有效的MAC地址
//注意：有关可接受的格式和类型，请参阅Go的ParseMAC: http://golang.org/src/net/mac.go?s=866:918#L29
hostname    //根据RFC 952 https://tools.ietf.org/html/rfc952验证字符串值是否为有效主机名
hostname_rfc1123 or if you want to continue to use 'hostname' in your tags, create an alias    //根据RFC 1123 https://tools.ietf.org/html/rfc1123验证字符串值是否为有效主机名
fqdn    //这将验证字符串值是否包含有效的FQDN (完全合格的有效域名)，Full Qualified Domain Name (FQDN)
html    //这将验证字符串值是否为HTML元素标记，包括https://developer.mozilla.org/en-US/docs/Web/HTML/Element中描述的标记。
html_encoded    //这将验证字符串值是十进制或十六进制格式的正确字符引用
url_encoded    //这验证了根据https://tools.ietf.org/html/rfc3986#section-2.1对字符串值进行了百分比编码（URL编码）
```
# 参考
+ <https://github.com/go-playground/validator>
+ <https://www.cnblogs.com/zhzhlong/p/10033234.html>
+ <https://www.cnblogs.com/jiujuan/p/13823864.html>