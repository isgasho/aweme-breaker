package aweme

// #include "aweme.h"
// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"github.com/anonymous5l/console"
	"github.com/valyala/fasthttp"
	"os"
	"unsafe"
)

const (
	AwemeLibPath = "libs/AwemeDylib"
)

func signHandler(ctx *fasthttp.RequestCtx) {
	var req map[string]interface{}

	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	var url,cookie,stub *C.char

	goUrl, ok := req["URL"].(string)
	if !ok {
		ctx.Response.SetStatusCode(400)
		return
	} else {
		url = C.CString(goUrl)
	}

	goCookie, ok := req["Cookie"].(string)
	if ok {
		cookie = C.CString(goCookie)
	}

	goStub, ok := req["Stub"].(string)
	if ok {
		stub = C.CString(goStub)
	}

	defer func() {
		if url != nil {
			C.free(unsafe.Pointer(url))
		}

		if cookie != nil {
			C.free(unsafe.Pointer(cookie))
		}

		if stub != nil {
			C.free(unsafe.Pointer(stub))
		}
	}()

	gorgon := make([]byte, 54)
	timestamp := make([]byte, 12)

	ret := C.Signature(url, stub, cookie, unsafe.Pointer(&gorgon[0]), unsafe.Pointer(&timestamp[0]))
	if ret != 0 {
		console.Err("Generate ret: %d", ret)
		ctx.Response.SetStatusCode(500)
		return
	}

	m := map[string]interface{}{
		"X-Gorgon": string(gorgon[:52]),
		"X-Khronos": string(timestamp[:10]),
	}

	b, err := json.Marshal(m)
	if err != nil {
		console.Err("json: %s", err)
		ctx.Response.SetStatusCode(500)
		return
	}

	ctx.Response.SetStatusCode(200)
	ctx.SetBody(b)
}

func InitAwemeRouter() map[string]fasthttp.RequestHandler {
	lp := AwemeLibPath
	if _, err := os.Stat(lp); err != nil {
		lp = "/var/lib/ibreaker/" + AwemeLibPath
		if _, err := os.Stat(lp); err != nil {
			console.Err("Can't found aweme signature engine")
			return nil
		}
	}
	path := C.CString(lp)
	code := C.Init(path)
	C.free(unsafe.Pointer(path))

	if code != 0 {
		console.Err("Load aweme signature engine error %d", code)
		return nil
	}

	console.Ok("Aweme signature engine load succeed")

	return map[string]fasthttp.RequestHandler{
		"sign": signHandler,
	}
}
