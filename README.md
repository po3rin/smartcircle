# smartcircle

<img src="https://img.shields.io/badge/go-v1.11-blue.svg"/>

Package smartcircle lets you crop a circular image out of a rectangular smartly. It will automatically determine where to crop.

<img src="./src/what.png">

## Installation

```bash
$ go get github.com/po3rin/smartcircle/cmd/smartcircle
```

## Try this on Web

you enabled to try here !! (developed by Go + Wasm)

https://po3rin.github.io/smartcircle/web/

## Usage

as CLI tool.

```bash
$ smartcircle -f testdata/gopher.jpeg -o cropped.png
```

as Code.

```go
package main

import (
    _ "image/jpeg"
    "image/png"
    "os"

    "github.com/po3rin/smartcircle"
)

func main(){
    img, _ := os.Open(*imgPath)
    defer img.Close()
    src, _, _ := image.Decode(img)

    // use smartcircle packege.
    c, _ := smartcircle.NewCropper(smartcircle.Params{Src: src})
    result, _ := c.CropCircle()

    file, _ := os.Create("cropped.png")
    defer file.Close()
    _ = png.Encode(file, result)
}
```
