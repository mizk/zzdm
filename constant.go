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
	Version = "v0.1.0"
	//SKU
	SKU = "001"
)

//错误代码
var (
	ErrorFileName         = errors.New("输入和输出文件名相同")
	ErrorInvalidData      = errors.New("无效数据")
	ErrorFileIO           = errors.New("IO错误")
	ErrorAES              = errors.New("AES错误")
	ErrorInvalidFile      = errors.New(fmt.Sprintf("不是有效的%s文件", Extension))
	ErrorFileDuplicated   = errors.New("目标文件已存在且不覆盖.如果要覆盖目标文件，请指定--force")
	ErrorDataMissing      = errors.New("缺少数据")
	ErrorFrameMissing     = errors.New("缺少帧数据")
	ErrorChecksumMismatch = errors.New("数据校验错误")
)
