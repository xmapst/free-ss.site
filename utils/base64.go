package utils

import "encoding/base64"

func Bs64EnStr(src string) string {
	return base64.URLEncoding.EncodeToString([]byte(src))
}

func Bs64DeStr(src string) string {
	byteSrc, _ := base64.URLEncoding.DecodeString(src)
	return string(byteSrc)
}

func Bs64RawEnStr(src string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(src))
}

func Bs64RawDeStr(src string) string {
	byteSrc, _ := base64.RawURLEncoding.DecodeString(src)
	return string(byteSrc)
}
