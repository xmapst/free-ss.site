package utils

import (
	"crypto/tls"
	browser "github.com/EDDYCJY/fake-useragent"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetResponse(method, webUrl string, header map[string]string, data io.Reader) (response *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}

	request, err := http.NewRequest(method, webUrl, data)
	if err != nil {
		return nil, err
	}

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}
	response, err = client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return nil, err
	}
	return response, nil
}

func RespToStr(response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func GetCookie(getUrl, referer string) (result string) {
	header := make(map[string]string)
	header["Referer"] = referer
	header["User-Agent"] = browser.Chrome()

	response, err := GetResponse("GET", getUrl, header, nil)
	if err != nil {
		return ""
	}

	var strSlice []string
	for _, v := range response.Cookies() {
		value := v.Name + "=" + v.Value
		strSlice = append(strSlice, value)
	}
	result = strings.Join(strSlice, ";")
	_ = response.Body.Close()
	return
}
