package router

import (
	"fmt"
	"github.com/Serendipity-sw/gutil"
	"github.com/gin-gonic/gin"
	"github.com/swgloomy/gutil/glog"
	"net/http"
	"os"
	"rtsp-monitoring-server/controller"
	"rtsp-monitoring-server/service"
	"rtsp-monitoring-server/store"
	"strings"
)

func Init(debug bool) {
	gin.SetMode(gutil.If(debug, gin.DebugMode, gin.ReleaseMode).(string))
	rt := gin.Default()
	rt.Use(cors())
	routerInit(rt)
	go func() {
		err := rt.Run(store.ListenPort)
		if err != nil {
			fmt.Printf("rt run err! err: %s \n", err.Error())
			os.Exit(0)
		}
	}()
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Headers", "token,redirect,content-type,x-requested-with")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")

		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}

		c.Next()
	}
}

func routerInit(r *gin.Engine) {
	g := &r.RouterGroup
	if store.RootPrefix != "" {
		g = r.Group(store.RootPrefix)
	}
	{
		g.GET("/ok", func(c *gin.Context) { c.String(http.StatusOK, "server start") })
		g.GET("/", func(c *gin.Context) { c.File(fmt.Sprintf("%s/index.html", store.Template)) })
		g.GET("/asset/:params", getAsset)
		api := g.Group("/api")
		{
			api.GET("/open", controller.Open)
			api.GET("/close", controller.Close)
		}
	}
	r.NoRoute(noRoute)
}

func getAsset(c *gin.Context) {
	urlPath := c.Param("params")
	urlPath = fmt.Sprintf("%s/%s", store.Asset, urlPath)
	bo, err := gutil.PathExists(urlPath)
	if err != nil {
		glog.Error("router getAsset path exists run err! urlPath: %s err: %+v \n", urlPath, err)
		return
	}
	if bo {
		go service.SetMonitorStartTime(strings.TrimSpace(c.Query("channel")))
		c.File(urlPath)
	} else {
		noRoute(c)
	}
}

func noRoute(c *gin.Context) {
	urlPath := c.Request.URL.Path
	if urlPath != "" {
		filePath := fmt.Sprintf("%s%s", store.Template, urlPath)
		bo, err := gutil.PathExists(filePath)
		if err != nil {
			glog.Error("router noRoute path exists run err! urlPath: %s filePath: %s err: %+v \n", urlPath, filePath, err)
			return
		}
		if bo {
			c.File(filePath)
		}
	}
}
