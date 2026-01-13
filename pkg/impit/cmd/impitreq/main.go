package main

// #include <stdlib.h>
import "C"

import (
	"bytes"
	"encoding/json"
	"unsafe"

	"impit/pkg/impit"
)

type Command struct {
	ClientOptions impit.ImpitOptions `json:"client_options"`
	Request       impit.RequestInit  `json:"request"`
}

//export ImpitHandleRequestJSON
func ImpitHandleRequestJSON(input *C.char) *C.char {
	if input == nil {
		return C.CString(`{"error":"nil input"}`)
	}

	var cmd Command
	if err := json.Unmarshal([]byte(C.GoString(input)), &cmd); err != nil {
		out, _ := json.Marshal(impit.ResponseData{Error: err.Error()})
		return C.CString(string(out))
	}

	client := impit.CreateClient(cmd.ClientOptions)
	resp := impit.HandleRequest(client, cmd.Request)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(resp); err != nil {
		out, _ := json.Marshal(impit.ResponseData{Error: err.Error()})
		return C.CString(string(out))
	}

	return C.CString(buf.String())
}

//export ImpitFreeCString
func ImpitFreeCString(p *C.char) {
	if p == nil {
		return
	}
	C.free(unsafe.Pointer(p))
}

func main() {}
