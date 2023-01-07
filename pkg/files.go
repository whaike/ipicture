package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// FileMD5 对于大文件计算m5d只需很小的内存
func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
