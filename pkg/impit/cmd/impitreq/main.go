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
func ImpitHandleRequestJSON(input *C.char, output *C.char, maxLen C.int) C.int {
	if input == nil {
		return -1
	}

	var cmd Command
	if err := json.Unmarshal([]byte(C.GoString(input)), &cmd); err != nil {
		out, _ := json.Marshal(impit.ResponseData{Error: err.Error()})
		return copyToBuffer(out, output, maxLen)
	}

	client := impit.CreateClient(cmd.ClientOptions)
	resp := impit.HandleRequest(client, cmd.Request)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(resp); err != nil {
		out, _ := json.Marshal(impit.ResponseData{Error: err.Error()})
		return copyToBuffer(out, output, maxLen)
	}

	return copyToBuffer(buf.Bytes(), output, maxLen)
}

func copyToBuffer(data []byte, output *C.char, maxLen C.int) C.int {
	if len(data) >= int(maxLen) {
		// Buffer too small, return required size (negative)
		return C.int(-len(data))
	}

	// Copy data to C buffer
	cBuf := unsafe.Slice((*byte)(unsafe.Pointer(output)), int(maxLen))
	copy(cBuf, data)
	cBuf[len(data)] = 0 // Null terminator

	return C.int(len(data))
}

//export ImpitFreeCString
func ImpitFreeCString(p *C.char) {
	if p == nil {
		return
	}
	C.free(unsafe.Pointer(p))
}

func main() {}
