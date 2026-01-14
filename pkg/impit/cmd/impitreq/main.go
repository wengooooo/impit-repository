package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"unsafe"

	"impit/pkg/impit"
)

//export ImpitRequest
func ImpitRequest(clientJSON *C.char, requestJSON *C.char) *C.char {
	var clientOpts impit.ImpitOptions
	var reqOpts impit.RequestInit
	if clientJSON != nil {
		_ = json.Unmarshal([]byte(C.GoString(clientJSON)), &clientOpts)
	}
	if requestJSON != nil {
		_ = json.Unmarshal([]byte(C.GoString(requestJSON)), &reqOpts)
	}
	c := impit.CreateClient(clientOpts)
	res := impit.HandleRequest(c, reqOpts)
	buf, _ := json.Marshal(res)
	return C.CString(string(buf))
}

//export FreeCString
func FreeCString(p *C.char) {
	if p != nil {
		C.free(unsafe.Pointer(p))
	}
}

func main() {}
