package zzdm

import (
	"os"
	"errors"
	"fmt"
)

const (
	//路径分隔符
	PathSeparator = string(os.PathSeparator)
	//扩展名
	Extension = ".scc"
	//4kb
	BUFFER int64 = 4096
	//autor
	Author = "mizk.chen@gmail.com"
	//version
	Version = "v0.1.1"
	//SKU
	SKU = "1806262316"
)

//错误代码
var (
	ErrorFileName         = errors.New("the input file name is same with the output")
	ErrorInvalidData      = errors.New("invalid bytes")
	ErrorFileIO           = errors.New("io error")
	ErrorAES              = errors.New("aes error")
	ErrorInvalidFile      = errors.New(fmt.Sprintf("not a valid %s file", Extension))
	ErrorFileDuplicated   = errors.New("output file already exists,to overwrite it,specify the flag --force")
	ErrorDataMissing      = errors.New("no more bytes to read")
	ErrorFrameMissing     = errors.New("mssing frames")
	ErrorChecksumMismatch = errors.New("checksum mismatch")
)
