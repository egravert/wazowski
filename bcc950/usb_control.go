package bcc950

/*
#include <libusb.h>

#cgo pkg-config: libusb-1.0
*/
import "C"
import (
	"errors"
	"time"
)

// Camera represents the BC950 webcam.
// It can be used to pan and zoom the camera.
type Camera struct {
	handle   *[0]byte
	OnTimeMs time.Duration // the duration to run the controlling motors (relevant for pan & tilt)
}

const (
	logitech_vendid = 0x046d
	bcc950_prodid   = 0x0837
	defaultOnTimeMs = 20 * time.Millisecond
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

// NewCamera returns a Camera that can be used to control the BCC950 webcam.
func NewCamera() (*Camera, error) {
	handle := C.libusb_open_device_with_vid_pid(nil, logitech_vendid, bcc950_prodid)
	if handle == nil {
		return nil, errors.New("Cannot open device.")
	}

	camera := &Camera{
		handle,
		defaultOnTimeMs,
	}
	return camera, nil
}

// rotates the camera right(clockwise).
// The camera rotation is relative and controlled by turnning on the motor
// with a direction and stopping it after a period of time.
func (camera *Camera) PanRight() {
	camera.moveCamera(panRightCommand)
}

// rotates the camera left(counter-clockwise).
// The camera rotation is relative and controlled by turnning on the motor
// with a direction and stopping it after a period of time.
func (camera *Camera) PanLeft() {
	camera.moveCamera(panLeftCommand)
}

// Tilts the camera upward.
// The camera rotation is relative and controlled by turnning on the motor
// with a direction and stopping it after a period of time.
func (camera *Camera) TiltUp() {
	camera.moveCamera(tiltUpCommand)
}

// Tilts the camera upward.
// The camera rotation is relative and controlled by turnning on the motor
// with a direction and stopping it after a period of time.
func (camera *Camera) TiltDown() {
	camera.moveCamera(tiltDownCommand)
}

func (camera *Camera) moveCamera(command [4]C.uchar) {
	camera.controlTransfer(command)
	time.Sleep(camera.OnTimeMs)
	camera.controlTransfer(stopCommand)
}

func (camera *Camera) controlTransfer(command [4]C.uchar) {
	C.libusb_control_transfer(camera.handle,
		0x21,
		0x01,
		0x0E00,
		0x0100,
		&command[0],
		0x0004,
		1000)
}
