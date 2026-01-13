package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"strconv"
	"unsafe"
)

//export ImpitHandleRequestJSON
func ImpitHandleRequestJSON(req *C.char) *C.char {
	in := C.GoString(req)
	out := `{"ok":true,"echo":` + strconv.Quote(in) + `}`
	return C.CString(out)
}

//export ImpitFree
func ImpitFree(p *C.char) {
	C.free(unsafe.Pointer(p))
}

func main() {}
