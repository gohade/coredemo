package main

import (
	"coredemo/framework"
	"time"
)

func FooControllerHandler(ctx *framework.Context) error {
	time.Sleep(10 * time.Second)
	return ctx.Json(200, map[string]interface{}{
		"code": 0,
	})
}

// func Foo(request *http.Request, response http.ResponseWriter) {
// 	obj := map[string]interface{}{
// 		"errno":  50001,
// 		"errmsg": "inner error",
// 		"data":   nil,
// 	}

// 	response.Header().Set("Content-Type", "application/json")

// 	foo := request.PostFormValue("foo")
// 	if foo == "" {
// 		foo = "10"
// 	}
// 	fooInt, err := strconv.Atoi(foo)
// 	if err != nil {
// 		response.WriteHeader(500)
// 		return
// 	}
// 	obj["data"] = fooInt
// 	byt, err := json.Marshal(obj)
// 	if err != nil {
// 		response.WriteHeader(500)
// 		return
// 	}
// 	response.WriteHeader(200)
// 	response.Write(byt)
// 	return
// }

// func Foo2(ctx *framework.Context) error {
// 	obj := map[string]interface{}{
// 		"errno":  50001,
// 		"errmsg": "inner error",
// 		"data":   nil,
// 	}

// 	fooInt := ctx.FormInt("foo", 10)
// 	obj["data"] = fooInt
// 	return ctx.Json(http.StatusOK, obj)
// }

// func Foo3(ctx *framework.Context) error {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	return rdb.Set(ctx, "key", "value", 0).Err()
// }
