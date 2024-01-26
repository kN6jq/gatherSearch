package goqqwry

import "os"

const (
	// IndexLen 索引长度
	IndexLen = 7
	// RedirectMode1 国家的类型, 指向另一个指向
	RedirectMode1 = 0x01
	// RedirectMode2 国家的类型, 指向一个指向
	RedirectMode2 = 0x02
)

// ResultQQwry 归属地信息
type ResultQQwry struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Area    string `json:"area"`
}

type FileData struct {
	Data     []byte
	FilePath string
	Path     *os.File
	IPNum    int64
}

// QQwry 纯真ip库
type QQwry struct {
	Data   *FileData
	Offset int64
}
