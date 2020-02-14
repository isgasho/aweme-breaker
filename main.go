package main

import (
	"github.com/anonymous5l/console"
	"github.com/valyala/fasthttp"
	"ibreaker/aweme"
	"strings"
)

var router map[string]map[string]fasthttp.RequestHandler

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.RequestURI())
	paths := strings.Split(path, "/")[1:]

	if len(paths) < 2 {
		ctx.Response.SetStatusCode(400)
		return
	}

	// main
	console.Log("Request: %#v", paths)

	if first, ok := router[paths[0]]; !ok {
		ctx.Response.SetStatusCode(404)
		return
	} else if first == nil {
		ctx.Response.SetStatusCode(500)
		return
	} else if second, ok := first[paths[1]]; !ok {
		ctx.Response.SetStatusCode(404)
		return
	} else {
		second(ctx)
	}
}

func main() {
	router = make(map[string]map[string]fasthttp.RequestHandler)
	router["aweme"] = aweme.InitAwemeRouter()

	if err := fasthttp.ListenAndServe("127.0.0.1:25583", requestHandler); err != nil {
		console.Err("fasthttp: %s", err)
	}
}
