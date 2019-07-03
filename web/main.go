// +build js,wasm

package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	sh := New()
	sh.Start()
}
