package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Encoder interface {
	Encode() (data []byte, contentType string, err error)
}

type HandlerFunc func(r http.ResponseWriter, w *http.Request)

//App is a wrapper for the http server.
type App struct {
	Engine *gin.Engine
}

func NewApp() *App{
	e := gin.New()
	return &App{
		Engine: e,
	}
}

func (a *App) Start(){
	err := a.Engine.Run()
	if err != nil {
		panic(err)
	}
}

func (a *App) HandlerFunc(method string, group string, path string, handlerFunc HandlerFunc, mw ...gin.HandlerFunc) {
    fullPath := group + path
    a.routeMethod(method, fullPath, handlerFunc, mw...)
}

func (a *App) HandlerFuncNoMid(method string, group string, path string, handlerFunc HandlerFunc) {
    fullPath := group + path
    a.routeMethod(method, fullPath, handlerFunc)
}

func (a *App) routeMethod(method string, path string, handlerFunc HandlerFunc, mw ...gin.HandlerFunc) {
    handler := func(c *gin.Context) {
        handlerFunc(c.Writer, c.Request)
    }

    if len(mw) > 0 {
        handler = func(c *gin.Context) {
            for _, m := range mw {
                m(c)
            }
            handlerFunc(c.Writer, c.Request)
        }
    }

    switch method {
    case "GET":
        a.Engine.GET(path, handler)
    case "POST":
        a.Engine.POST(path, handler)
    case "PUT":
        a.Engine.PUT(path, handler)
    case "DELETE":
        a.Engine.DELETE(path, handler)
    default:
        log.Printf("Unsupported method: %s", method)
    }
}
