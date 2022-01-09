package routers

import (
	"encoding/base64"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/anaskhan96/soup"
	"github.com/bitly/go-simplejson"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/unknwon/com"
	"github/xmapst/free-ss.site/utils"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
)

// GetFreeSs
// @Summary 获取 free-ss.site
// @Tags SS/SSR
// @Produce  json
// @Param php query string false "data.php页面" default(data2.php)
// @Param v query int false "版本" default(1)
// @Success 200 {object} JSONResult
// @Failure 500 {object} JSONResult
// @Router /json [get]
func GetFreeSs(c *gin.Context) {
	render := Gin{C: c}
	//data2.php
	dataPhp, ok := c.GetQuery("php")
	if !ok || dataPhp == "undefined" {
		dataPhp = "data2.php"
	}
	// version
	version := c.GetInt("v")
	if version != 2 {
		version = 1
	}
	source, err := freeSs(dataPhp, version)
	if err != nil {
		logrus.Error(err)
		render.SetError(utils.CODE_ERR_MSG, err)
		return
	}
	if source == nil {
		render.SetError(utils.CODE_ERR_MSG, fmt.Errorf("get failed"))
		return
	}
	render.SetJson(source)
}

// GetSubscribe
// @Summary 获取 ssr 订阅
// @Tags SS/SSR
// @Produce  json
// @Param php query string false "data.php页面" default(data2.php)
// @Param v query int false "版本" default(1)
// @Success 200 {object} JSONResult
// @Failure 500 {object} JSONResult
// @Router / [get]
func GetSubscribe(c *gin.Context) {
	render := Gin{C: c}
	//data2.php
	dataPhp, ok := c.GetQuery("php")
	if !ok || dataPhp == "undefined" {
		dataPhp = "data2.php"
	}
	// version
	version := c.GetInt("v")
	if version != 2 {
		version = 1
	}
	source, err := freeSs(dataPhp, version)
	if err != nil {
		logrus.Error(err)
		render.SetError(utils.CODE_ERR_MSG, err)
		return
	}
	if source == nil {
		render.SetError(utils.CODE_ERR_MSG, fmt.Errorf("get failed"))
		return
	}
	var listUrl []string
	for _, list := range source {
		listUrl = append(listUrl, freeSsSubscribe(list))
	}
	subscribeStr := utils.Bs64EnStr(strings.Join(listUrl, "\n\r"))
	c.String(200, "%s", subscribeStr)
}

type freeSsRes struct {
	Url string
	Res string
	Err error
}

type ServerJson struct {
	Remarks    string `json:"remarks"`
	Server     string `json:"server"`
	ServerPort string `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	Speed      string `json:"speed"`
	Uptime     string `json:"uptime"`
	SsUrl      string `json:"ss_url"`
}

var vm = goja.New()

func freeSs(dataPhp string, version int) ([]ServerJson, error) {
	webUrl := "https://free-ss.site/"
	enCcUrl := fmt.Sprintf("%sajax/libs/encc/0.0.0/encc.min.js", webUrl)
	cryptoUrl := fmt.Sprintf("%sajax/libs/crypto-js/3.1.9-1/crypto-js.min.js", webUrl)

	header := make(map[string]string)
	header["User-Agent"] = browser.Chrome()
	header["Connection"] = "keep-alive"
	header["Upgrade-Insecure-Requests"] = "1"
	header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	header["Accept-Language"] = "en-US,en;q=0.9"
	header["Referer"] = webUrl
	header["Origin"] = webUrl

	ch := make(chan freeSsRes)
	urls := []string{webUrl, enCcUrl, cryptoUrl}
	var html, enCcJs, cryptoJs string
	var err error
	for _, u := range urls {
		go getResStr("GET", u, header, ch)
	}

	for range urls {
		reqs := <-ch
		if reqs.Err != nil {
			return nil, reqs.Err
		}
		switch reqs.Url {
		case webUrl:
			html = reqs.Res
		case enCcUrl:
			enCcJs = reqs.Res
		case cryptoUrl:
			cryptoJs = reqs.Res
		}
	}
	close(ch)

	doc := soup.HTMLParse(html).FindAll("script")
	if doc == nil {
		return nil, fmt.Errorf("element `script` with attributes not found")
	}
	var script string
	for k, src := range doc {
		if strings.Contains(src.Text(), "grecaptcha") {
			script = src.Text()
			break
		}
		if k == len(doc)-1 {
			return nil, fmt.Errorf("element `script` with attributes `grecaptcha` not found")
		}
	}

	valueDict := make(map[string]string)
	valueDict["a"], err = getValue(script, "a")
	if err != nil {
		return nil, err
	}
	valueDict["b"], err = getValue(script, "b")
	if err != nil {
		return nil, err
	}
	valueDict["c"], err = getValue(script, "c")
	if err != nil {
		return nil, err
	}
	cValue, err := vm.RunString(fmt.Sprintf("%sencc('%s')", enCcJs, valueDict["c"]))

	if err != nil {
		return nil, err
	}
	valueDict["c"] = cValue.String()

	v := make(url.Values)
	v.Set("a", valueDict["a"])
	v.Set("b", valueDict["b"])
	v.Set("c", valueDict["c"])

	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	header["Content-Type"] = "application/x-www-form-urlencoded"
	response, err := utils.GetResponse("POST", fmt.Sprintf("%s%s", webUrl, dataPhp), header, body)
	if err != nil {
		return nil, err
	}
	d := utils.RespToStr(response)
	if len(d) == 0 {
		return nil, fmt.Errorf("request for encrypted data failed")
	}

	evJsFunc, err := getJsFun(script)
	if err != nil {
		return nil, err
	}
	evJs, err := vm.RunString(evJsFunc)
	if err != nil {
		return nil, err
	}
	ev := fmt.Sprintf("%sdec.toString(CryptoJS.enc.Base64);", evJs.String())

	variables := fmt.Sprintf(`;var a = "%s";var d = "%s";var b = "%s"`,
		valueDict["a"], d, valueDict["b"])

	xy := ";var x=CryptoJS.enc.Latin1.parse(a);var y=CryptoJS.enc.Latin1.parse(b);"
	base64Data, err := vm.RunString(fmt.Sprintf("%s%s%s%s", cryptoJs, variables, xy, ev))
	if err != nil {
		return nil, err
	}

	jsonData, err := base64.StdEncoding.DecodeString(base64Data.String())
	if err != nil {
		return nil, err
	}
	return servers(jsonData, version)
}

func getResStr(method, u string, header map[string]string, ch chan<- freeSsRes) {
	response, err := utils.GetResponse(method, u, header, nil)
	if err != nil {
		ch <- freeSsRes{
			Url: u,
			Res: "",
			Err: err,
		}
		return
	}

	resStr := utils.RespToStr(response)
	_ = response.Body.Close()
	ch <- freeSsRes{
		Url: u,
		Res: resStr,
		Err: nil,
	}
}

func servers(str []byte, version int) (result []ServerJson, err error) {
	jsons, err := simplejson.NewJson(str)
	if err != nil {
		return nil, err
	}
	getData := jsons.Get("data").Interface()
	for _, v := range getData.([]interface{}) {
		server := ServerJson{
			Speed:      com.ToStr(v.([]interface{})[0]),
			Server:     com.ToStr(v.([]interface{})[1]),
			ServerPort: com.ToStr(v.([]interface{})[2]),
			Method:     com.ToStr(v.([]interface{})[3]),
			Password:   com.ToStr(v.([]interface{})[4]),
			Uptime:     com.ToStr(v.([]interface{})[5]),
			Remarks:    com.ToStr(v.([]interface{})[6]),
		}
		if version == 2 {
			server.Password, server.Method = server.Method, server.Password
		}
		server.SsUrl = getSSUrl(server)
		result = append(result, server)
	}

	return
}

func getSSUrl(server ServerJson) string {
	//method:password@server:port
	ssUrl := fmt.Sprintf("%s:%s@%s:%s",
		server.Method,
		server.Password,
		server.Server,
		server.ServerPort)
	bs64Url := utils.Bs64EnStr(ssUrl)
	return "ss://" + bs64Url
}

func getJsFun(javascript string) (string, error) {
	re := `eval\(function.*`
	match, err := regexp.MatchString(re, javascript)
	if err != nil {
		return "", err
	}
	if match {
		jsExp, err := regexp.Compile(re)
		if err != nil {
			return "", err
		}
		params := jsExp.FindAllString(javascript, -1)
		return string([]rune(params[1])[4:len([]rune(params[1]))]), nil
	}
	return "", fmt.Errorf("not match")
}

func getValue(content, char string) (string, error) {
	re := fmt.Sprintf(`%s\s*='\s*[^;]+`, char)
	match, err := regexp.MatchString(re, content)
	if err != nil {
		return "", err
	}
	if !match {
		return "", fmt.Errorf("not match")
	}
	valueExp, err := regexp.Compile(re)
	if err != nil {
		return "", err
	}
	params := valueExp.FindAllString(content, -1)
	for _, p := range params {
		if ok := strings.Contains(p, "1"); ok {
			return string([]rune(p)[3 : len([]rune(p))-1]), nil
		}
	}
	return "", nil
}

func freeSsSubscribe(server ServerJson) (result string) {
	// host:port:origin:aes-256-cfb:plain:/?remarks=&group=
	str := fmt.Sprintf("%s:%s:origin:%s:plain:%s/?remarks=%s&group=%s",
		server.Server,
		server.ServerPort,
		server.Method,
		utils.Bs64RawEnStr(server.Password),
		utils.Bs64RawEnStr("YFDou"+"-"+server.Remarks),
		utils.Bs64RawEnStr("YFDouSS"),
	)
	return fmt.Sprintf("ssr://%s", utils.Bs64RawEnStr(str))
}
