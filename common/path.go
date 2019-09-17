package common

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
)

//文件内容遍历
func SeekFile(filepath string) (func() (string, bool), error) {
	fobj, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		fobj.Close()
		return nil, err
	}
	scanner := bufio.NewScanner(fobj)

	return func() (string, bool) {
		for scanner.Scan() {
			code := scanner.Text()
			return string(code), true
		}
		if err := scanner.Err(); err != nil {
			LogPanic("SeekFile %s Error: %v", filepath, err)
		}
		fobj.Close()
		return "", false
	}, nil
}

//读取完整文件内容
func ReadFile(filepath string) (string, error) {
	fobj, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		fobj.Close()
		return "", err
	}
	defer fobj.Close()

	buf, err := ioutil.ReadAll(fobj)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

//写入内容到文件
func WriteFile(filepath string, data string) error {
	LogDebug("WriteFile: %s", filepath)
	dir := path.Dir(filepath)
	CreateDir(dir)

	fobj, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fobj.Close()

	fileWritor := bufio.NewWriter(fobj)
	_, err = fileWritor.WriteString(data)
	if err != nil {
		return err
	}
	err = fileWritor.Flush()
	if err != nil {
		return err
	}
	return nil
}

//复制文件
func CopyFile(fromPath string, tarPath string) error {
	fobj, err := os.OpenFile(fromPath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer fobj.Close()
	target, err := os.OpenFile(tarPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer target.Close()

	dir := path.Dir(tarPath)
	CreateDir(dir)
	_, err = io.Copy(target, fobj)
	return err
}

//移动文件
func MoveFile(fromPath string, tarPath string) error {
	dir := path.Dir(tarPath)
	CreateDir(dir)
	return os.Rename(fromPath, tarPath)
}

/*
判断文件或文件夹是否存在
如果返回的错误为nil,说明文件或文件夹存在
如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
如果返回的错误为其它类型,则不确定是否在存在
*/
func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CreateDir(dirPath string) error {
	if IsPathExists(dirPath) {
		return nil
	}
	err := os.MkdirAll(dirPath, 0755)
	return err
}

func GetFilePath(dir string, filename string) string {
	return path.Join(dir, filename)
}
