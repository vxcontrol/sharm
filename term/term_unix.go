// This PTY package using only on UNIX system

// +build darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package term

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/kr/pty"
)

// DefBufferSize is max length buffer to read and to write from/to terminal
const DefBufferSize int = 10240

const pathToDefaultCMD string = "/bin/bash"

// Term is main terminal structure
type Term struct {
	cmd *exec.Cmd
	pty *os.File
}

// Start is function to prepare terminal for using
func (t *Term) Start(command string) (err error) {
	if command == "" {
		command = pathToDefaultCMD
	}
	cmds := strings.Split(command, " ")
	t.cmd = exec.Command(cmds[0], cmds[1:]...)
	t.cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	t.pty, err = pty.Start(t.cmd)
	return t.Resize(98, 24)
}

// Write is function to write data to opened terminal
func (t *Term) Write(b []byte) (err error) {
	_, err = t.pty.Write(b)
	return
}

// Read is function to read data from opened terminal
func (t *Term) Read(b []byte) (size int, err error) {
	size, err = t.pty.Read(b)
	return
}

// Close is function to close opened terminal
func (t *Term) Close() {
	t.pty.Close()
	t.cmd.Process.Kill()
	t.cmd.Wait()
}

// Resize is function to resize opened terminal
func (t *Term) Resize(cols, rows int) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(rows),
		uint16(cols),
		0,
		0,
	}
	syscall.Syscall(
		syscall.SYS_IOCTL,
		t.pty.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	return nil
}
