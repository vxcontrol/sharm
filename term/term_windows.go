// This PTY package using only on Windows system

// +build windows

package term

// #cgo CFLAGS: -I../include/
// #cgo windows,386 LDFLAGS: -static -L ${SRCDIR}/../lib/windows_386 -lpty -lmingwex -lmingw32 -lstdc++
// #cgo windows,amd64 LDFLAGS: -static -L ${SRCDIR}/../lib/windows_amd64 -lpty -lmingwex -lmingw32 -lstdc++
// #include <pty.h>
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

// DefBufferSize is max length buffer to read and to write from/to terminal
const DefBufferSize int = 10240

const pathToDefaultCMD string = "C:\\WINDOWS\\System32\\cmd.exe"

// Term is main terminal structure
type Term struct {
	pty unsafe.Pointer
}

// Start is function to prepare terminal for using
func (t *Term) Start(command string) error {
	if command == "" {
		command = pathToDefaultCMD
	}
	t.pty = C.PtyOpen(C.int(98), C.int(24), C.CString(command))
	if t.pty == nil {
		return errors.New("PtyOpen returned fail")
	}
	return nil
}

// Write is function to write data to opened terminal
func (t *Term) Write(b []byte) error {
	ln := len(b)
	var buf = make([]C.uchar, ln)
	for idx, bt := range b {
		buf[idx] = C.uchar(bt)
	}
	if n := C.PtyWrite(t.pty, &buf[0], C.ulong(ln)); C.int(n) == C.int(-1) {
		return errors.New("PtyWrite return fail: " + strconv.Itoa(int(C.int(n))))
	}
	return nil
}

// Read is function to read data from opened terminal
func (p *Term) Read(b []byte) (int, error) {
	var buf = make([]C.uchar, DefBufferSize)
	var size C.ulong = C.ulong(DefBufferSize)
	n := C.PtyRead(p.pty, &buf[0], size)
	if n > 0 && n <= size {
		s := C.GoBytes(unsafe.Pointer(&buf[0]), C.int(n))
		copy(b, s)
		return int(C.int(n)), nil
	} else if n == 0 {
		return 0, nil
	}
	return 0, errors.New("PtyRead return fail: " + strconv.Itoa(int(C.int(n))))
}

// Close is function to close opened terminal
func (p *Term) Close() {
	C.PtyKill(p.pty)
}

// Resize is function to resize opened terminal
func (p *Term) Resize(cols, rows int) error {
	if n := C.PtyResize(p.pty, C.int(cols), C.int(rows)); C.int(n) != C.int(1) {
		return errors.New("PtyResize return fail: " + strconv.Itoa(int(C.int(n))))
	}
	return nil
}
