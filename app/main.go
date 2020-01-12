package main

import (
	"fmt"
	"os"
	"time"

	"../mouse"
	"github.com/tajtiattila/xinput"
)

func main() {

	if err := xinput.Load(); err != nil {
		panic(err)
	}

	state := &xinput.State{}
	vibration := &xinput.Vibration{}
	point := &mouse.POINT{}
	speed := int32(10)
	mouseDown := false
	for {
		if err := xinput.GetState(0, state); err != nil {
			panic(fmt.Errorf("Controller 0 error: %w", err))
		}

		vibration.LeftMotorSpeed = uint16(state.Gamepad.LeftTrigger) * 257
		vibration.RightMotorSpeed = uint16(state.Gamepad.RightTrigger) * 257

		xinput.SetState(0, vibration)

		mouse.GetCursorPos(point)

		if xinput.DPAD_UP&xinput.Button(state.Gamepad.Buttons) == xinput.DPAD_UP {
			point.Y -= speed
		}
		if xinput.DPAD_DOWN&xinput.Button(state.Gamepad.Buttons) == xinput.DPAD_DOWN {
			point.Y += speed
		}
		if xinput.DPAD_LEFT&xinput.Button(state.Gamepad.Buttons) == xinput.DPAD_LEFT {
			point.X -= speed
		}
		if xinput.DPAD_RIGHT&xinput.Button(state.Gamepad.Buttons) == xinput.DPAD_RIGHT {
			point.X += speed
		}

		if mouseDown == false && (xinput.BUTTON_X&xinput.Button(state.Gamepad.Buttons) == xinput.BUTTON_X) {
			mouse.MouseEvent(mouse.MOUSEEVENTF_LEFTDOWN, point)
			mouseDown = true
			fmt.Println("Mousedown", point)
		}
		if mouseDown == true && (xinput.BUTTON_X&xinput.Button(state.Gamepad.Buttons) != xinput.BUTTON_X) {
			mouse.MouseEvent(mouse.MOUSEEVENTF_LEFTUP, point)
			mouseDown = false
			fmt.Println("Mouseup", point)
		}

		if xinput.BACK&xinput.Button(state.Gamepad.Buttons) == xinput.BACK {
			fmt.Println("Exiting")
			os.Exit(0)
		}
		mouse.SetCursorPos(point.X, point.Y)
		time.Sleep(time.Second / 60)
	}

}
