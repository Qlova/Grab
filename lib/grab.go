package grab

import (
	"image/color"
	"net"
	"os"
	"fmt"
	"bytes"
	"encoding/binary"
	"errors"
)

var Framebuffer *os.File

var ScreenWidth, ScreenHeight, ScreenDepth int

func XYToI(x, y, width int) int64 {
	return int64(x)*int64(ScreenDepth)+int64(y)*int64(width)*int64(ScreenDepth)
}

func IToXY(i, width int) (int, int) {
	return i % width, i / width
}


var initialised bool
func Init() {
	
	if initialised {
		return
	}
    
    var err error
   	Framebuffer, err = os.OpenFile("/dev/fb0", os.O_WRONLY, 662)
	if err != nil {
		fmt.Println("Cannot open the framebuffer, please run as root!")
		os.Exit(1)
	}
	
	 
	ScreenWidth, ScreenHeight, ScreenDepth = Geometry(Framebuffer.Fd())
	//fmt.Println("Geometry", ScreenWidth, ScreenHeight, ScreenDepth)
	initialised = true
}

type Window struct {
	X, Y, Width, Height uint16
	Connection net.Conn
	Buffer []byte
}

type Panel struct {
	Window
	Direction, Size int
}

func (window Window) Draw(image Image, X int, Y int) {
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
			int64(XYToI(X, Y+int(y), int(window.Width))))
	}
}

func (window Window) WriteBufferAt(bytes []byte, offset64 int64) {
	var offset = int(offset64)
	for i := 0; i < len(bytes); i+=int(ScreenDepth) {
		if bytes[i+3] != 0 && len(window.Buffer) > offset+i+3 {
			window.Buffer[offset+i]   = bytes[i]
			window.Buffer[offset+i+1] = bytes[i+1]
			window.Buffer[offset+i+2] = bytes[i+2]
			window.Buffer[offset+i+3] = bytes[i+3]
		} 
	}
}

func (window Window) Render() {
	for y := int(0); y < int(window.Height); y++ {
		//fmt.Println(window.Width*y+window.Width)
		//if len(window.Buffer) > window.Width*y+window.Width { 
			Framebuffer.WriteAt(
				window.Buffer[int(ScreenDepth)*int(window.Width)*int(y):int(ScreenDepth)*int(window.Width)*int(y)+int(ScreenDepth)*int(window.Width)], 
				int64(XYToI(int(window.X), int(window.Y)+y, ScreenWidth+42)))
		//}
	}
}

func (window Window) Fill(c color.Color) {
	if window.Hidden() {
		return
	}
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	
	var row []byte = make([]byte, int(window.Width)*ScreenDepth)
	for i := 0; i < int(window.Width)*int(ScreenDepth); i+=int(ScreenDepth) {
		if ScreenDepth == 4 {
			row[i]	 = rgba.R
			row[i+1] = rgba.G
			row[i+2] = rgba.B
			row[i+3] = rgba.A
		}
	}
	
	for y := int(0); y < int(window.Height); y++ {
		window.WriteBufferAt( row, int64(XYToI(int(window.X), int(window.Y)+y, int(window.Width))))
	}
		
}

func (window Window) Hidden() bool {
	if window.Height == 0 || window.Width == 0 {
		return true
	}
	fmt.Print("Window Hidden!\r")
	return false
	
}

const (
	Top = iota
	Bottom
	Left
	Right
)

func (panel *Panel) Connect() error {
	Init()
	
	var x, y, width, height int
	if panel.Direction == Top {
		x, y, width, height  = 0, 0, ScreenWidth, panel.Size
	}
	if panel.Direction == Bottom {
		x, y, width, height  = 0, ScreenHeight-panel.Size, ScreenWidth, panel.Size
	}
	if panel.Direction == Left {
		x, y, width, height  = 0, 0, panel.Size, ScreenHeight
	}
	if panel.Direction == Right {
		x, y, width, height  = ScreenWidth-panel.Size, 0, panel.Size, ScreenHeight
	}
	
	Connection, err := net.Dial("unix", "/var/run/grab.sock")
    if err != nil {
       return errors.New("Cannot grab :S "+err.Error())
    }
	
	var buff bytes.Buffer
	binary.Write(&buff, binary.BigEndian, uint16(width))
	binary.Write(&buff, binary.BigEndian, uint16(height))
	binary.Write(&buff, binary.BigEndian, uint16(x))
	binary.Write(&buff, binary.BigEndian, uint16(y))
	
	if _, err := Connection.Write(buff.Bytes()); err != nil {
		fmt.Println("Error reading from server! Closed connection :(")
		Connection.Close()
		return errors.New("Cannot grab :S "+err.Error())
	}
	
	var buffer = make([]byte, 8)
	
	_, err = Connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server! Closed connection :(")
		Connection.Close()
		return errors.New("Cannot grab :S "+err.Error())	
	}
	{
	
		x := bytes.NewReader(buffer[:2])
		binary.Read(x, binary.BigEndian, &panel.X)
		y := bytes.NewReader(buffer[2:4])
		binary.Read(y, binary.BigEndian, &panel.Y)
		width := bytes.NewReader(buffer[4:6])
		binary.Read(width, binary.BigEndian, &panel.Width)
		height := bytes.NewReader(buffer[6:8])
		binary.Read(height, binary.BigEndian, &panel.Height)
	
		panel.Connection = Connection
		panel.Buffer = make([]byte, int(panel.Height)*int(panel.Width)*int(ScreenDepth))
	}
	
	return nil
}

func (window *Window) Connect() (error) {
	Init()
	
	Connection, err := net.Dial("unix", "/var/run/grab.sock")
    if err != nil {
       return errors.New("Cannot grab :S "+err.Error())
    }
	
	var buff bytes.Buffer
	binary.Write(&buff, binary.BigEndian, uint16(ScreenWidth))
	binary.Write(&buff, binary.BigEndian, uint16(ScreenHeight))
	binary.Write(&buff, binary.BigEndian, uint16(0))
	binary.Write(&buff, binary.BigEndian, uint16(0))
	
	if _, err := Connection.Write(buff.Bytes()); err != nil {
		fmt.Println("Error reading from server! Closed connection :(")
		Connection.Close()
		return errors.New("Cannot grab :S "+err.Error())
	}
	
	var buffer = make([]byte, 8)
	
	_, err = Connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server! Closed connection :(")
		Connection.Close()
		return errors.New("Cannot grab :S "+err.Error())	
	}
	
	x := bytes.NewReader(buffer[:2])
	binary.Read(x, binary.BigEndian, &window.X)
	y := bytes.NewReader(buffer[2:4])
	binary.Read(y, binary.BigEndian, &window.Y)
	width := bytes.NewReader(buffer[4:6])
	binary.Read(width, binary.BigEndian, &window.Width)
	height := bytes.NewReader(buffer[6:8])
	binary.Read(height, binary.BigEndian, &window.Height)
	
	window.Connection = Connection
	window.Buffer = make([]byte, int(window.Height)*int(window.Width)*int(ScreenDepth))
	
	return nil
}
