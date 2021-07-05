package framework

import (
	"log"
	"net/http"
)

// Core represent core struct
type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	panic("not implement")
}

func (c *Core) Post(url string, handler ControllerHandler) {
	panic("not implement")
}

func (c *Core) FindRouteByRequest(request *http.Request) ControllerHandler {
	return nil
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)

	router := c.FindRouteByRequest(request)
	if router == nil {
		return
	}

	ctx.SetHandler(router)

	ctx.ProcessChan = make(chan string)

	go func(ctx *Context) error {
		err := router(ctx)
		if err != nil {
			ctx.ProcessChan <- err.Error()
			return err
		}

		ctx.ProcessChan <- ""
		return nil
	}(ctx)

	select {
	case err := <-ctx.ProcessChan:
		if err != "" {
			response.Write([]byte(err))
		}
	case <-ctx.BaseContext().Done():
		log.Println("ctx timeout")
	}
}
