package routers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github/xmapst/free-ss.site/config"
	_ "github/xmapst/free-ss.site/docs"
	"github/xmapst/free-ss.site/utils"
	"github/xmapst/free-ss.site/utils/ip"
	"net"
	"net/http"
	"strings"
)

var (
	checker *ip.Checker
	err     error
)

const xForwardedFor = "X-Forwarded-For"

func init() {
	if !utils.InSliceString("*", config.Config.IPWhiteList) && len(config.Config.IPWhiteList) != 0 {
		checker, err = ip.NewChecker(config.Config.IPWhiteList)
	}
	if err != nil {
		logrus.Fatal(err)
	}
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(utils.Logger(), gin.Recovery(), gzip.Gzip(gzip.BestCompression))
	url := ginSwagger.URL(config.Config.PrefixUrl + "/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/", GetSubscribe)
	r.GET("/json", GetFreeSs)
	return r
}

func handlersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		render := Gin{C: c}
		if !utils.InSliceString("*", config.Config.IPWhiteList) && len(config.Config.IPWhiteList) != 0 {
			reqIPAddr := getRemoteIP(c.Request)
			reaIPaddyLenOffset := len(reqIPAddr) - 1
			for i := reaIPaddyLenOffset; i >= 0; i-- {
				err = checker.IsAuthorized(reqIPAddr[i])
				if err != nil {
					logrus.Error(err)
					render.SetError(utils.CODE_ERR_NO_PRIV, err)
					return
				}
			}
		}
		c.Next()
	}
}

func getRemoteIP(req *http.Request) []string {
	var ipList []string

	xff := req.Header.Get(xForwardedFor)
	xffs := strings.Split(xff, ",")

	for i := len(xffs) - 1; i >= 0; i-- {
		xffsTrim := strings.TrimSpace(xffs[i])

		if len(xffsTrim) > 0 {
			ipList = append(ipList, xffsTrim)
		}
	}

	hostPort, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		remoteAddrTrim := strings.TrimSpace(req.RemoteAddr)
		if len(remoteAddrTrim) > 0 {
			ipList = append(ipList, remoteAddrTrim)
		}
	} else {
		ipTrim := strings.TrimSpace(hostPort)
		if len(ipTrim) > 0 {
			ipList = append(ipList, ipTrim)
		}
	}

	return ipList
}
