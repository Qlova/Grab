#Grab
=======
The Display Server for the future.
Enough is the time of the old, the time of the so called "windows".
It is time for a common sense display server which meets the needs of the few who are the many.
As they said in a quote "In a world without walls, who needs windows and gates?".
That is the philosophy of Grab.

Grab is based on the idea of sharing.
You share with the computer, the functions you require.
For example, instead of drawing your own buttons and dealing with the havoc of getting them in the right spot. Grab a button.
Instead of designing a great user interface for your program, let the display server design you.

There are 4 different things you can Grab: (More may be sent back from the future)

* Images.
* Buttons.
* Text Boxes.
* Menus.

Each of these items can be grabbed by your very own GUI application.
Now this application can be written in ANY language.
That's right! You can write it in ANY language.

Grab operates on a socket. You send your gui to it. OpenGL and video can only be used fullscreen.

If you want to grab a button, send a button request to the Grab Server.
eg.
```
	"GRAB button; press me!"
```
The server will respond with an identification number:
```
	"ID 232"
```
Now whenever the button is pressed, you will recieve:
```
	"PRESS 232"
```
Fantastic!!!! A GUI for everyone =)

##Theme
=========
Ok cool so what does grab actually look like??
The protocol does not specify this, this is what a grab theme would be responable for.

Grab does have a default theme though.

###SteamPunk
Grab has a steampunk theme, bringing the feel from the past's future to our present.
The way this works is a series of buttons on the screen, the posistion of these buttons maps to the keys on your keyboard, your keyboard becomes the touch-screen for your monitor.
This is very intuitive and allows for quick and easy navigation, typing mode can be entered and exited.
