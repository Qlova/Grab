package grab

import "github.com/gvalkov/golang-evdev"
import "fmt"

func init() {
	go processEvents()
}

func Updates() {

}

func processEvents() {
	devices, _ := evdev.ListInputDevices("/dev/input/*")
	for i := range devices {
		go func(i int) {
			for {
				events, err := devices[i].Read()
				if err == nil {
					for i := range events {
						processEvent(events[i])
					}
				}
			}
		}(i)
	}
}

var lastMouse Mouse
var touchingX, touchingY, touchedX, touchedY bool

func processEvent(event evdev.InputEvent) {

	if event.Type == evdev.EV_KEY {
		if event.Code == evdev.BTN_TOUCH {
			if event.Value == 0 {
				touchingX = false
				touchedX = false
				touchingY = false
				touchedY = false
			}
			if event.Value == 1 {
				touchedX = true
				touchedY = true
			}
		}
	}

	if event.Type == evdev.EV_ABS && processTouchPad(event) {
		return
	}
	if event.Type == evdev.EV_REL && processMouse(event) {
		return
	}

	var code_name string
	
	switch event.Type {
	case evdev.EV_KEY:
		val, haskey := evdev.KEY[int(event.Code)]
		if haskey {
			code_name = val
		} else {
			val, haskey := evdev.BTN[int(event.Code)]
			if haskey {
				code_name = val
			} else {
				code_name = "?"
			}
		}
	default:
		m, haskey := evdev.ByEventType[int(event.Type)]
		if haskey {
			code_name = m[int(event.Code)]
		} else {
			code_name = "?"
		}
	}
	
	//fmt.Println(code_name)
	fmt.Sprint(code_name)
}

func processMouse(event evdev.InputEvent) bool {
	if event.Code == evdev.REL_X {
		mouse.X += int(event.Value)
		return true
	}
	if event.Code == evdev.REL_Y {
		mouse.Y += int(event.Value)
		return true
	}
	return false
}

func processTouchPad(event evdev.InputEvent) bool {
	if event.Code == evdev.ABS_X {
		if touchingX {
			if lastMouse.X >= mouse.X+int(event.Value) {
				mouse.X = 0
			} else {
				mouse.X += (int(event.Value)-lastMouse.X)/2
			}
		}
		if touchedX {
			touchingX = true
		}
		lastMouse.X = int(event.Value)
		return true
	}
	if event.Code == evdev.ABS_Y {
		if touchingY {
			if lastMouse.Y >= mouse.Y+int(event.Value) {
				mouse.Y = 0
			} else {
				mouse.Y += (int(event.Value)-lastMouse.Y)/2
			}
		}
		if touchedY {
			touchingY = true
		}
		lastMouse.Y = int(event.Value)
		return true
	}
	return false
}
