package framework

import (
	"fmt"
	"log"
	"time"
)

func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
	return func(c *Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			fun(c)
			finish <- struct{}{}
		}()
		select {
		case p := <-panicChan:
			log.Println(p)
			c.responseWriter.WriteHeader(500)

		case <-finish:
			fmt.Println("finish")
		case <-time.After(d):
			c.responseWriter.Write([]byte("time out"))
		}
		return nil
	}
}
