package utils

import (
	"strings"
	"unicode"
)

// IsEmpty 判断字符串是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 判断字符串是否非空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// TrimSpace 去除字符串首尾空白字符
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// Contains 判断字符串是否包含子串
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// HasPrefix 判断字符串是否以指定前缀开头
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix 判断字符串是否以指定后缀结尾
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// ToLower 转换为小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper 转换为大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// IsDigit 判断字符串是否全为数字
func IsDigit(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}
