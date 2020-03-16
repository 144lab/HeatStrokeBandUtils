// +build go1.12

package frontend

import (
	"syscall/js"
)

func Uint8ArrayToSlice(value js.Value) []byte {
	// Note that TypedArrayOf cannot work correcly on Wasm.
	// See https://github.com/golang/go/issues/31980

	s := make([]byte, value.Get("byteLength").Int())
	a := js.TypedArrayOf(s)
	a.Call("set", value)
	a.Release()
	return s
}
