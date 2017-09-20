package main

import "fmt"
import "os"
import "image"
import "image/color"
import "github.com/fogleman/gg"

import "./widget"

func XYToI(x, y, width int, window *Framebuffer) int64 {
	return int64(x)*int64(window.Depth)+int64(y)*int64(width)*int64(window.Depth)
}

func IToXY(i, width int) (int, int) {
	return i % width, i / width
}

type Framebuffer struct {
	Width, Height, Depth int
	Buffer []byte
	File *os.File
}

func OpenFrameBuffer() *Framebuffer {
	framebuffer := new(Framebuffer)

	var file, err = os.OpenFile("/dev/fb0", os.O_WRONLY, 662)
	if err != nil {
		fmt.Println("Cannot open the framebuffer, please run as root!")
		os.Exit(1)
	}
	
	framebuffer.File = file
	 
	framebuffer.Width, framebuffer.Height, framebuffer.Depth = Geometry(file.Fd())
	
	framebuffer.Buffer = make([]byte, framebuffer.Width*framebuffer.Height*framebuffer.Depth)
	
	return framebuffer
}

func (window *Framebuffer) RoundedSquare() {
	const S = 256
	dc := gg.NewContext(window.Width, window.Height)
	dc.SetRGBA(0, 0, 0, 1)
	dc.Clear()
	
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/ubuntu-font-family/Ubuntu-M.ttf", 24); err != nil {
		panic(err)
	}
	
	dc.SetRGBA(1, 1, 1, 1)
	
	var x float64 = 0
	var y float64 = 0
	for i:=0; i < 11; i++ {
		x += S/2
		dc.DrawRoundedRectangle(20+x, 20+y, S/2, S/2, 25)
		dc.SetLineWidth(5)
		dc.Stroke()
		dc.Fill()
		
		dc.DrawStringAnchored("Empty", 20+x+S/4, 20+y+S/4, 0.5, 0.5)
		
		if i == 5 {
			y += S/2
			x  = 0
		}
	}
	
	//Need to draw string wrapped :) DONNE
	dc.DrawRoundedRectangle(100, 300, float64(window.Width)-200, float64(window.Height-300), 25)
	dc.SetLineWidth(5)
	dc.Stroke()
	dc.Fill()
	
	//THIS IS BUGGY, TODO notify foggleman. Should wrap on width.
	dc.DrawStringWrapped(textbox.GetFormattedText(), 150, 320, 0, 0, float64(window.Width)-220, 1.5, gg.AlignLeft)
	
	//dc.DrawString(textbox.GetText(), 100, 300)
	
	img := dc.Image().(*image.RGBA)
	window.Draw(Image{Data:img.Pix, Stride:img.Stride, Width:window.Width, Height:window.Height}, 0, 0)
}

func init() {
	
}

func (window *Framebuffer) Print() {
	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/ubuntu-font-family/Ubuntu-M.ttf", 96); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)
	img := dc.Image().(*image.RGBA)
	window.Draw(Image{Data:img.Pix, Stride:img.Stride, Width:1024, Height:1024}, 0, 0)
}

func (window *Framebuffer) Draw(image Image, X int, Y int) {
	var stride = image.Stride
	if X+image.Width > int(window.Width) {
		stride = image.Stride-(4*(int(X)+int(image.Width)-int(window.Width)))
		if stride <= 1 {
			return
		}
		if stride > image.Stride {
			stride = image.Stride
		}
	}
	for y := 0; y < int(image.Height); y++ {
		window.WriteBufferAt( 
			image.Data[image.Stride*y:image.Stride*y+stride], 
			int64(XYToI(X, Y+int(y), int(window.Width), window)))
	}
}

func (window *Framebuffer) WriteBufferAt(bytes []byte, offset64 int64) {
	var offset = int(offset64)
	for i := 0; i < len(bytes); i+=int(window.Depth) {
		if bytes[i+3] != 0 && len(window.Buffer) > offset+i+3 {
			window.Buffer[offset+i]   = bytes[i]
			window.Buffer[offset+i+1] = bytes[i+1]
			window.Buffer[offset+i+2] = bytes[i+2]
			window.Buffer[offset+i+3] = bytes[i+3]
		} 
	}
}

func (window *Framebuffer) Render() {
	for y := int(0); y < int(window.Height); y++ {
		//fmt.Println(window.Width*y+window.Width)
		//if len(window.Buffer) > window.Width*y+window.Width { 
			window.File.WriteAt(
				window.Buffer[int(window.Depth)*int(window.Width)*int(y):int(window.Depth)*int(window.Width)*int(y)+int(window.Depth)*int(window.Width)], 
				int64(XYToI(int(0), int(0)+y, window.Width+42, window)))
		//}
	}
}

func (window *Framebuffer) Fill(c color.Color) {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	
	var row []byte = make([]byte, int(window.Width)*window.Depth)
	for i := 0; i < int(window.Width)*int(window.Depth); i+=int(window.Depth) {
		if window.Depth == 4 {
			row[i]	 = rgba.R
			row[i+1] = rgba.G
			row[i+2] = rgba.B
			row[i+3] = rgba.A
		}
	}
	
	for y := int(0); y < int(window.Height); y++ {
		window.WriteBufferAt( row, int64(XYToI(int(0), int(0)+y, int(window.Width), window)))
	}
		
}

var textbox widget.TextBox

func Renderer() {
   	var framebuffer = OpenFrameBuffer()
	
	img := Image{}
	err := img.Load("screenshot.png")
	if err != nil {
		println(err.Error())
	}
	
	framebuffer.Draw(img, 0, 0)
	
	framebuffer.RoundedSquare()
	
	framebuffer.Render()
	
	for {
		<- widget.Changed

		framebuffer.Draw(img, 0, 0)
		framebuffer.Print()
	
		framebuffer.RoundedSquare()
	
		framebuffer.Render()
	}
	
}
