package zzdm

import (
	"encoding/binary"
	"path/filepath"
	"hash/adler32"
	"math/rand"
	"strings"
	"time"
	"fmt"
	"os"
	"io"
)

//写入文件头
func WriteHead(file *os.File, fileName []byte, secret bool, frames int64) error {
	head := Header{frames, fileName, secret}
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, 8848)
	file.Write(bytes)
	message, err := head.Marshal()
	if err != nil {
		return err
	}
	size := uint64(len(message))
	bytes = make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, size)
	file.Write(bytes)
	file.Write(message)
	return nil
}

//写入数据帧
func WriteFrame(file *os.File, iv, data []byte, hashcode uint32) error {
	frame := Frame{
		iv,
		data,
		hashcode,
	}
	message, err := frame.Marshal()
	if err != nil {
		return err
	}
	length := uint64(len(message))
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, length)
	file.Write(bytes)
	file.Write(message)
	return nil
}

//读取文件头
func ReadHead(file *os.File) (*Header, error) {

	tag, err := ReadUInt64Value(file)
	if tag != 8848 {
		return nil, err
	}
	size, err := ReadUInt64Value(file)
	if err != nil {
		return nil, err
	}
	if size <= 0 {
		return nil, ErrorInvalidData
	}
	bytes, err := ReadBytes(file, size)
	if err != nil {
		return nil, err
	}
	if uint64(len(bytes)) != size {
		return nil, ErrorInvalidData
	}
	header := &Header{
	}
	err = header.Unmarshal(bytes)
	if err != nil {
		return nil, err
	}
	return header, nil
}

//读取指定长度的字节
func ReadBytes(file *os.File, length uint64) ([]byte, error) {
	bytes := make([]byte, length)
	num, err := file.Read(bytes)
	if err != nil {
		return nil, err
	}
	if num <= 0 {
		return nil, ErrorInvalidData
	}
	return bytes[:num], nil
}

//读取数据长度
func ReadUInt64Value(file *os.File) (uint64, error) {
	bytes := make([]byte, 8)
	num, err := file.Read(bytes)
	if err != nil {
		return 0, err
	}
	if num != 8 {
		return 0, ErrorInvalidData
	}
	tag := binary.BigEndian.Uint64(bytes[:num])
	return tag, nil
}

//读取数据帧
func ReadFrame(file *os.File) (*Frame, error) {
	length, err := ReadUInt64Value(file)
	if err != nil {
		return nil, err
	}
	if length <= 0 {
		return nil, ErrorInvalidData
	}
	bytes, err := ReadBytes(file, length)
	if err != nil {
		return nil, err
	}
	if uint64(len(bytes)) != length {
		return nil, ErrorInvalidData
	}
	frame := &Frame{

	}
	err = frame.Unmarshal(bytes)
	return frame, err
}

//解密文件
func Decrypt(input, output, password string, force bool) error {

	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := ReadHead(file)
	if err != nil {
		return err
	}
	if header == nil {
		return ErrorFileIO
	}
	ph := stringBytes(password, 32) //加密密钥
	div := defaultIv()              //加密向量
	nameBytes := header.Name
	secret := header.Secret
	fileName := ""
	if secret {
		nameBytes, err = AesDecrypt(nameBytes, ph, div)
		if err != nil {
			return err
		}
		if nameBytes == nil {
			return ErrorAES
		}
	}
	fileName = string(nameBytes)
	if len(fileName) <= 0 {
		return ErrorFileIO
	}
	fullName := decryptionName(input, output, fileName)
	//输入文件和输出文件不能相同
	//这里忽略文件大小写,对于某些系统可能会存在误判
	if strings.EqualFold(fullName, input) {
		return ErrorFileName
	}
	if Exist(fullName) {
		if force {
			err = os.Truncate(fullName, 0)
			if err != nil {
				return err
			}
		} else {
			return ErrorFileDuplicated
		}
	}
	ptr, err := Open(fullName)
	if err != nil {
		return err
	}
	defer ptr.Close()
	frameCount := header.Frames
	var index int64 = 0

	for {
		frame, err := ReadFrame(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		frameData := frame.Data
		if frameData == nil {
			return ErrorFileIO
		}
		checksum := frame.Hash
		iv := frame.Iv
		if iv == nil {
			return ErrorFileIO
		}
		ivDecrypted, err := AesDecrypt(iv, ph, div)
		if err != nil {
			return err
		}
		data, err := AesDecrypt(frameData, ph, ivDecrypted)
		if err != nil {
			return err
		}
		checksum2 := adler32.Checksum(data)
		if checksum2 != checksum {
			return ErrorChecksumMismatch
		}
		size, err := ptr.Write(data)
		if err != nil {
			return err
		}
		if size != len(data) {
			return ErrorDataMissing
		}
		index += 1
		fmt.Printf("frame{index=%d,max=%d}\n", index, frameCount)
	}
	if index != frameCount {
		return ErrorFrameMissing
	}
	return nil
}

//加密文件
func Encrypt(input, output, password string, secret, force bool) error {
	fileName := encryptionName(input, output, secret)
	//输入文件和输出文件不能相同
	//这里忽略文件大小写,对于某些系统可能会存在误判
	if strings.EqualFold(fileName, input) {
		return ErrorFileName
	}
	if Exist(fileName) {
		if force {
			err := os.Truncate(fileName, 0)
			if err != nil {
				return err
			}
		} else {
			return ErrorFileDuplicated
		}
	}

	fileSize := FileLength(input)
	left := fileSize % BUFFER
	frameCount := (fileSize - left) / BUFFER
	if left > 0 {
		frameCount ++
	}
	ptr, err := Open(fileName)
	if err != nil {
		return err
	}
	defer ptr.Close()
	ph := stringBytes(password, 32)
	div := defaultIv()
	baseName := filepath.Base(input)
	baseNameBytes := []byte(baseName)
	if secret {
		baseNameBytes, err = AesEncrypt(baseNameBytes, ph, div)
		if err != nil {
			return err
		}
	}
	err = WriteHead(ptr, baseNameBytes, secret, frameCount)
	if err != nil {
		return err
	}
	var index int64 = 0
	var offset int64 = 0
	raw, err := os.Open(input)
	if err != nil {
		return err
	}
	defer raw.Close()
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		buff, checksum, err := Read(raw, offset, BUFFER)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		byteSize := len(buff)
		offset += int64(byteSize)
		iv := randomBytes(32, 16) //32个随机字符的字符串的前16个字节
		ivEncrypt, err := AesEncrypt(iv, ph, div)
		if err != nil {
			return err
		}
		data, err := AesEncrypt(buff, ph, iv)
		if err != nil {
			return err
		}

		err = WriteFrame(ptr, ivEncrypt, data, checksum)
		if err != nil {
			return err
		}
		index ++
		fmt.Printf("frame{index=%d,max=%d,bytes=%d}\n", index, frameCount, byteSize)
		if index >= frameCount {
			break
		}
	}
	return nil
}

//加密文件保存地址
func encryptionName(input, output string, hidden bool) string {
	baseName := ""
	dir := ""
	var fileName = ""
	if !hidden {
		path := filepath.Base(input)
		extension := filepath.Ext(input)
		baseName = strings.TrimSuffix(path, extension)
	} else {
		rand.Seed(time.Now().UTC().UnixNano())
		baseName = Md5(RandomString(32))
	}
	if !IsDir(output) {
		dir = strings.TrimSuffix(input, filepath.Base(input))
	} else {
		dir = output
	}
	if strings.HasSuffix(dir, PathSeparator) {
		fileName = fmt.Sprintf("%s%s%s", dir, baseName, Extension)
	} else {
		fileName = fmt.Sprintf("%s%s%s%s", dir, PathSeparator, baseName, Extension)
	}
	return fileName
}

//解密文件保存地址
func decryptionName(input, output, realName string) string {
	dir := ""
	var fileName = ""
	if !IsDir(output) {
		dir = strings.TrimSuffix(input, filepath.Base(input))
	} else {
		dir = output
	}
	if strings.HasSuffix(dir, PathSeparator) {
		fileName = fmt.Sprintf("%s%s", dir, realName)
	} else {
		fileName = fmt.Sprintf("%s%s%s", dir, PathSeparator, realName)
	}
	return fileName
}
