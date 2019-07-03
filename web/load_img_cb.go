// +build js,wasm

package main

import (
	"bytes"
	"image"
	"reflect"
	"syscall/js"
	"time"
	"unsafe"
)

func (c *Cropper) setupOnImgLoadCb() {
	c.onImgLoadCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		reader := bytes.NewReader(c.inBuf)
		var err error
		c.sourceImg, _, err = image.Decode(reader)
		if err != nil {
			c.log(err.Error())
			return nil
		}
		c.log("Ready for operations")
		start := time.Now()
		c.updateImage(start)
		return nil
	})
}

func (c *Cropper) setupInitMemCb() {
	// The length of the image array buffer is passed.
	// Then the buf slice is initialized to that length.
	// And a pointer to that slice is passed back to the browser
	c.initMemCb = js.FuncOf(func(this js.Value, i []js.Value) interface{} {
		length := i[0].Int()
		c.console.Call("log", "length:", length)
		// make buffer by image length
		c.inBuf = make([]uint8, length)
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&c.inBuf))
		ptr := uintptr(unsafe.Pointer(hdr.Data))
		// pass pointer to JS by calling function
		js.Global().Call("gotMem", ptr)
		return nil
	})
}
