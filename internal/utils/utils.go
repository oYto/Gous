package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// Contains 检测性别合法性
func Contains(source []string, tg string) bool {
	for _, s := range source {
		if s == tg {
			return true
		}
	}
	return false
}

// Md5String 这个函数使用MD5算法对输入字符串进行哈希处理，并返回哈希值的16进制表示。
// 首先，将输入字符串转换为字节数组，然后使用md5.New()函数创建一个MD5哈希器。
// 接下来，使用Write()方法将字节数组写入哈希器，计算出哈希值。
// h.Sum 返回的是哈希值 []byte 类型
// 最后，使用hex.EncodeToString()函数将哈希值转换为16进制字符串返回。
func Md5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

// GenerateSession 生成一个会话ID
func GenerateSession(userName string) string {
	return Md5String(fmt.Sprintf("%s:%s", userName, "session"))
}
