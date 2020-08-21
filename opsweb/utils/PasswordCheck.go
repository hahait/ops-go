package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

// 加密密码
func EncryptionPassword(password string) (string) {
	if pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err == nil {
		return string(pass)
	} else {
		fmt.Println("生成加密密码失败, 错误信息: ", err)
		return ""
	}
}

// 校验密码是否与数据库中一致
func CheckPassword(pass_in, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(pass_in)); err == nil {
		return true
	} else {
		return false
	}
}

// 密码强度校验
func CheckPasswordComplex(min_level int, password string) error {
	// min_level 声明密码强度的最低要求, 就是对包含字符种类（大写字母、小写字母、数字、特殊字符）的要求, 值越高包含的种类越多, 因此强度也越高;
	level := 0
	charList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, char := range charList {
		match, _ := regexp.MatchString(char, password)
		if match {
			level++
		}
	}

	if level < min_level {
		return fmt.Errorf("密码复杂度不合规")
	}
	return nil
}
