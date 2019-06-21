package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetGBKString(src string) string {
	code_rlt, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		fmt.Println("error: ", err)
		return src
	}
	return string(code_rlt)
}

func EncodeBase64(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func DecodeBase64(src string) string {
	code, err := base64.StdEncoding.DecodeString(src)
	CheckPanic(err)
	return string(code)
}
