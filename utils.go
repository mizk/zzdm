package zzdm

import (
	"path/filepath"
	"hash/adler32"
	"math/rand"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"io"
	"os"
	"regexp"
	"fmt"
)

//固定长度的随机字符串
func RandomString(length int) string {
	builder := strings.Builder{}
	for index := 0; index < length; index++ {
		num := rand.Intn(3)
		if num == 0 { //0-9
			builder.WriteRune(rune(rand.Intn(10) + 48))
		} else if num == 1 { //a-z
			builder.WriteRune(rune(rand.Intn(26) + 97))
		} else { //A-Z
			builder.WriteRune(rune(rand.Intn(26) + 65))
		}
	}
	return builder.String()
}

func randomBytes(length, capacity int) []byte {
	return stringBytes(RandomString(length), capacity)
}

func stringBytes(content string, capacity int) []byte {
	if capacity < 0 {
		capacity = 0
	}
	if len(content) == 0 {
		return make([]byte, capacity)
	}
	contentBytes := []byte(content)
	length := len(contentBytes)
	if length >= capacity {
		return contentBytes[0:capacity]
	} else {
		left := capacity - length
		empty := make([]byte, left)
		contentBytes = append(contentBytes, empty...)
	}

	return contentBytes
}

//文本的MD5值
func Md5(content string) string {
	h := md5.New()
	io.WriteString(h, content)
	return hex.EncodeToString(h.Sum(nil))
}

//打开文件
func Open(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

//给定路径的目录名
func Dir(path string) (string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		return path, nil
	}
	baseName := filepath.Base(path)
	return strings.TrimSuffix(path, baseName), nil
}

//文件是否存在
func Exist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func IsDir(filename string) bool {
	stat, err := os.Stat(filename)
	if err == nil {
		return stat.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//从文件中读取固定长度的字节
func Read(file *os.File, offset, len int64) ([]byte, uint32, error) {
	empty := make([]byte, 0, 0)
	var hashcode uint32 = 0
	buffer := make([]byte, len)
	_, err := file.Seek(offset, 0)
	if err != nil {
		return empty, hashcode, err
	}
	num, err := file.Read(buffer)
	if err != nil {
		return empty, hashcode, err
	}
	bytes := buffer[:num]
	hashcode = adler32.Checksum(bytes)
	return bytes, hashcode, nil
}

//文件长度
func FileLength(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return stat.Size()
}

func defaultIv() []byte {
	return stringBytes(fmt.Sprintf("%s", Author), 16)
}

func PasswordLevel(password string) int {
	/*
	密码建议
	1.至少有一个大写字母
	2.至少有一个小写字母
	3.至少有一个数字
	4.长度至少8位
	5.应该包含特殊字符
	*/
	regex, err := regexp.Compile("[A-Z]+")
	if err != nil {
		return -6
	}
	if !regex.MatchString(password) {
		return -1
	}
	regex, err = regexp.Compile("[a-z]+")
	if err != nil {
		return -6
	}
	if !regex.MatchString(password) {
		return -2
	}
	regex, err = regexp.Compile("[0-9]+")
	if err != nil {
		return -6
	}
	if !regex.MatchString(password) {
		return -3
	}
	if strings.Count(password, "")-1 < 8 {
		return -4
	}
	pattern := "~`!@#$%^&*()_+-=[]{}|\\<,>.?/;:\"'"
	if !strings.ContainsAny(password, pattern) {
		return -5
	}
	return 0
}
