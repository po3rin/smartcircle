// +build js,wasm

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"reflect"
	"syscall/js"
	"time"
	"unsafe"

	"github.com/po3rin/smartcircle"
)

type Cropper struct {
	inBuf                  []uint8
	outBuf                 bytes.Buffer
	onImgLoadCb, initMemCb js.Func
	sourceImg              image.Image

	console js.Value
	done    chan struct{}
}

func New() *Cropper {
	return &Cropper{
		console: js.Global().Get("console"),
		done:    make(chan struct{}),
	}
}

func (c *Cropper) Start() {
	// Setup functions
	c.setupInitMemCb()
	js.Global().Set("initMem", c.initMemCb)

	c.setupOnImgLoadCb()
	js.Global().Set("loadImage", c.onImgLoadCb)

	<-c.done
	c.log("Shutting down app")
	c.onImgLoadCb.Release()
}

func (c *Cropper) ConvertImage(argStartFlag string, argEndFlag string, argLoopFlag string) {
	// sourceImg is already decoded

	// Set image
	cropper, err := smartcircle.NewCropper(smartcircle.Params{Src: c.sourceImg})
	if err != nil {
		log.Fatal(err)
	}
	result, err := cropper.CropCircle()
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, result); err != nil {
		fmt.Println("error:png\n", err)
		return
	}

	c.outBuf = *buf
}

// updateImage writes the image to a byte buffer and then converts it to base64.
// Then it sets the value to the src attribute of the target image.
func (c *Cropper) updateImage(start time.Time) {
	c.console.Call("log", "updateImage:", start.String())
	c.ConvertImage("left", "right", "false")
	out := c.outBuf.Bytes()
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&out))
	ptr := uintptr(unsafe.Pointer(hdr.Data))
	// set pointer and length to JS function
	js.Global().Call("displayImage", ptr, len(out))
	c.console.Call("log", "time taken:", time.Now().Sub(start).String())
	c.outBuf.Reset()
}

// utility function to log a msg to the UI from inside a callback
func (s *Cropper) log(msg string) {
	js.Global().Get("document").
		Call("getElementById", "status").
		Set("innerText", msg)
}
