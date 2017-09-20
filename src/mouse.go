package main


var mouse Mouse

//This has all the information for a mouse.
type Mouse struct {
	X, Y int
}

func (ms *Mouse) Update() {
	/*if mouse.X < 0 {
		mouse.X = 0
	}
	if mouse.Y < 0 {
		mouse.Y = 0
	}
	if mouse.X > int(ScreenWidth) {
		mouse.X = int(ScreenWidth)
	}
	if mouse.Y > int(ScreenHeight) {
		mouse.Y = int(ScreenHeight)
	}
	ms.X, ms.Y = mouse.X, mouse.Y*/
}
