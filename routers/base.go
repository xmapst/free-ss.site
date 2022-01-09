package routers

import (
	"github.com/gin-gonic/gin"
	"github/xmapst/free-ss.site/config"
	"github/xmapst/free-ss.site/utils"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type JSONResult struct {
	Code int         `json:"code" description:"code"`
	Info *Info       `json:"info" description:"info"`
	Data interface{} `json:"data" description:"data"`
}

type Info struct {
	Ok      bool   `json:"ok" description:"status"`
	Message string `json:"message,omitempty" description:"message"`
}

func NewRes(data interface{}, err error, code int) *JSONResult {
	if code == 0 {
		code = 200
	}
	codeMsg := ""
	if config.Config.RunMode == "release" && code != 200 {
		codeMsg = utils.GetMsg(code)
	}

	return &JSONResult{
		Data: data,
		Code: code,
		Info: func() *Info {
			result := NewInfo(err)
			if _, ok := data.(string); ok && len(data.(string)) != 0 {
				result.Message = data.(string)
			}
			if codeMsg != "" {
				result.Message = codeMsg + ", Error message: " + err.Error()
			}
			return result
		}(),
	}
}

func NewInfo(err error) (info *Info) {
	info = new(Info)
	if err == nil {
		info.Ok = true
		return
	}
	info.Message = err.Error()
	return
}

// Response res
func (g *Gin) SetRes(res interface{}, err error, code int) {
	g.C.JSON(http.StatusOK, NewRes(res, err, code))
}

// Set Json
func (g *Gin) SetJson(res interface{}) {
	g.SetRes(res, nil, utils.CODE_SUCCESS)
}

// Check Error
func (g *Gin) SetError(code int, err error) {
	g.SetRes(nil, err, code)
	g.C.Abort()
}
