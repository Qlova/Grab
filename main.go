package main

import "os"
//import "io"
import "fmt"
import "time"

import (
	"net"
	"encoding/binary"
	"bytes"
	"grab/lib"
	"syscall"
	"io/ioutil"
)


//4325376 bytes
//4196352 bytes?
var image []byte = make([]byte, (1366+42)*768*4)

var loop chan bool = make(chan bool)

func draw() {
	//Open the framebuffer for writing.
	fb0, err := os.OpenFile("/dev/fb0", os.O_WRONLY, 660)
	if err != nil {
		fmt.Println("Cannot open the framebuffer, please run as root!")
		os.Exit(1)
	}
	defer fb0.Close()
	
	for i := range image {
		image[i] = 200
	}
	
	fb0.WriteAt(image, 0)

	time.Sleep(5000 * time.Millisecond)
}

var pid int

type Grab struct {
	X, Y, Width, Height uint16
}

type Job struct {
	Grab
	Pid int
	Response chan Grab
}

type Client struct {
	Name string
	Pid int
	
	Connection net.Conn
	Grab
}

func CreateClient(connection net.Conn) *Client {
	var client Client
	client.Connection = connection
	client.Pid = pid
	pid++
	
	return &client
}

func (client *Client) Do(request Grab) error {
	var channel = make(chan Grab)
	Jobs <- Job{Grab:request, Pid: client.Pid, Response:channel}
	response := <- channel
	
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, response.X)
	binary.Write(&buffer, binary.BigEndian, response.Y)
	binary.Write(&buffer, binary.BigEndian, response.Width)
	binary.Write(&buffer, binary.BigEndian, response.Height)
	
	if _, err := client.Connection.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

func (client *Client) Serve() {
	var buffer = make([]byte, 8)
	for {
		_, err := client.Connection.Read(buffer)
		if err != nil {
			switch err := err.(type) {
				case net.Error:
					if !err.Timeout() {
						continue
					}
				default:
					fmt.Println("Error reading from client! Closed client >:)")
					Jobs <- Job{Grab:Grab{}, Pid: client.Pid, Response:nil}
					client.Connection.Close()
					return		
			}
		}
		
		var request Grab
		x := bytes.NewReader(buffer[4:6])
		binary.Read(x, binary.BigEndian, &request.X)
		y := bytes.NewReader(buffer[6:8])
		binary.Read(y, binary.BigEndian, &request.Y)
		width := bytes.NewReader(buffer[:2])
		binary.Read(width, binary.BigEndian, &request.Width)
		height := bytes.NewReader(buffer[2:4])
		binary.Read(height, binary.BigEndian, &request.Height)

		if err := client.Do(request); err != nil {
			fmt.Println("Error writing to client! Closed client >:)")
			Jobs <- Job{Grab:Grab{}, Pid: client.Pid, Response:nil}
			client.Connection.Close()
			return
		}
	}
}

var Jobs = make(chan Job)

var Clients = make(map[int]Grab)

//We want to manage grabs here, saying who can render and where.
func work() {
	for {
		job := <- Jobs
		if len(Clients) == 0 {
			Clients[job.Pid] = job.Grab
			fmt.Println(job.Grab)
		} else {
			if job.Response == nil {
				fmt.Println("deleting client")
				delete(Clients, job.Pid)
				fmt.Println(len(Clients), " clients")
				continue
			}
			for i, v := range Clients {
				if i != job.Pid {
					if v.Width <= v.Height && job.X < v.X+v.Width {
						job.X = v.X+v.Width+1
					}
					if v.Width > v.Height && job.Y < v.Y+v.Height {
						job.Y = v.Y+v.Height+1
					}
				}
			}
			if int(job.X+job.Width) > grab.ScreenWidth {
				job.Width = job.Width-(job.X+job.Width-uint16(grab.ScreenWidth))
			}
			if int(job.Y+job.Height) > grab.ScreenHeight {
				job.Height = job.Height-(job.Y+job.Height-uint16(grab.ScreenWidth))
			}
			Clients[job.Pid] = job.Grab
		}		
		job.Response <- Clients[job.Pid]
	}
}

func client() {
	grab.Init()
}


func main() {

	os.Remove("/var/run/grab.sock")
	syscall.Umask(0700)
	listener, err := net.Listen("unix", "/var/run/grab.sock")
	if err != nil {
		fmt.Println("Cannot open socket, please run as root!", err)
		os.Exit(1)
	}
	os.Chmod("/dev/fb0", 0662)
	
	files, _ := ioutil.ReadDir("/dev/input")
    for _, f := range files {
           os.Chmod("/dev/input/"+f.Name(), 0664)
    }
	
	go work()
	go client()
	
	fmt.Println("You can now Grab!")
	for {
        connection, err := listener.Accept()
        if err != nil {
        	fmt.Println("Error accepting client: "+err.Error())
        }

        go CreateClient(connection).Serve()
    }
	
}
