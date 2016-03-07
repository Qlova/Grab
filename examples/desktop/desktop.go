package main

import "grab/lib"
import "image/color"
import "fmt"
import "time"

func main() {

	var window grab.Window
	var mouse grab.Mouse
	var cursor grab.Image
	
	if err := window.Connect(); err != nil {
		fmt.Println("Could not connect!", err)
		return
	}
	
	if err := cursor.Load("mouse.png"); err != nil {
		fmt.Println("Cursor did not load!!", err)
	}
	

	for {
		mouse.Update()
		window.Fill(color.RGBA{R:200, G:100, B:100, A: 255})
		window.Draw(cursor, mouse.X, mouse.Y)
		window.Render()
		
		time.Sleep(time.Millisecond*16)
	}
}
