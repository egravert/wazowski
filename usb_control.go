package main

/*
#include <libusb.h>

#cgo pkg-config: libusb-1.0
*/
import "C"
import (
	"errors"
	"time"
)

type Camera struct {
	handle *[0]byte
}

const (
	logitech_vendid = 0x046d
	bcc950_prodid   = 0x0837
)

var (
	panRightCommand = [4]C.uchar{0x01, 0x01, 0x00, 0x01}
	panLeftCommand  = [4]C.uchar{0xFF, 0x01, 0x00, 0x01}
	tiltUpCommand   = [4]C.uchar{0x00, 0x01, 0x01, 0x01}
	tiltDownCommand = [4]C.uchar{0x00, 0x01, 0xFF, 0x01}
	stopCommand     = [4]C.uchar{0x00, 0x01, 0x00, 0x01}
)

func init() {
	C.libusb_init(nil)
	C.libusb_set_debug(nil, 3)
}

func NewCamera() (*Camera, error) {
	handle := C.libusb_open_device_with_vid_pid(nil, logitech_vendid, bcc950_prodid)
	if handle == nil {
		return nil, errors.New("Cannot open device.")
	}

	camera := &Camera{
		handle,
	}
	return camera, nil
}

// rotates the camera clockwise for x milliseconds.
// The camera rotation is relative and controlled by turnning on the motor
// with a direction and stopping it after a period of time.
func (camera *Camera) PanRight() {
	camera.moveCamera(panRightCommand)
	time.Sleep(20 * time.Millisecond)
	camera.moveCamera(stopCommand)
}
func (camera *Camera) PanLeft() {
	camera.moveCamera(panLeftCommand)
	time.Sleep(20 * time.Millisecond)
	camera.moveCamera(stopCommand)
}

func (camera *Camera) TiltUp() {
	camera.moveCamera(tiltUpCommand)
	time.Sleep(20 * time.Millisecond)
	camera.moveCamera(stopCommand)
}

func (camera *Camera) TiltDown() {
	camera.moveCamera(tiltDownCommand)
	time.Sleep(20 * time.Millisecond)
	camera.moveCamera(stopCommand)
}

func (camera *Camera) moveCamera(command [4]C.uchar) {
	C.libusb_control_transfer(camera.handle,
		0x21,
		0x01,
		0x0E00,
		0x0100,
		&command[0],
		0x0004,
		1000)
}
