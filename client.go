package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/sharm/term"
)

const (
	// UnknownInput is unknown message type, maybe sent by a bug
	UnknownInput = '0'
	// Input is user input typically from a keyboard
	Input = '1'
	// Ping to the server
	Ping = '2'
	// ResizeTerminal is a notification that the browser size has been changed
	ResizeTerminal = '3'
	// Quit is a notification that client must close current connection and exit
	Quit = '4'
)

const (
	// UnknownOutput is unknown message type, maybe set by a bug
	UnknownOutput = '0'
	// Output is normal output to the terminal
	Output = '1'
	// Pong to the browser
	Pong = '2'
	// SetWindowTitle is set window title of the terminal
	SetWindowTitle = '3'
	// SetPreferences is set terminal preference
	SetPreferences = '4'
	// SetReconnect is make terminal to reconnect
	SetReconnect = '5'
)

type clientContext struct {
	request    *http.Request
	connection *websocket.Conn
	pty        *term.Term
	closed     bool
	writeMutex *sync.Mutex
}

type argResizeTerminal struct {
	Columns int
	Rows    int
}

func (context *clientContext) goHandleClient() error {
	var wg sync.WaitGroup
	var errSend, errRecv error

	wg.Add(1)
	go func() {
		defer wg.Done()
		errSend = context.processSend()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errRecv = context.processReceive()
	}()

	log.Info("Sharm client handler is waiting")
	wg.Wait()
	log.Info("Sharm client handler was released")

	if errSend != nil {
		return errSend
	} else if errRecv != nil {
		return errRecv
	}
	return nil
}

func (context *clientContext) processSend() (errResult error) {
	log.Info("Sharm sender has running")
	defer log.Info("Sharm sender was stopped")
	buf := make([]byte, term.DefBufferSize)
	var size int
	var err error
	defer func() {
		context.closed = true
		if errResult != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error(errResult.Error())
		}
	}()

	for !context.closed {
		if size, err = context.pty.Read(buf); err != nil {
			return errors.New("Failed to read data from termital")
		}
		if size > 0 {
			if err = context.write(append([]byte("1"), buf[:size]...)); err != nil {
				return errors.New("Failed to send data to server")
			}
		}
	}
	return nil
}

func (context *clientContext) write(data []byte) error {
	context.writeMutex.Lock()
	defer context.writeMutex.Unlock()
	context.connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return context.connection.WriteMessage(websocket.TextMessage, data)
}

func (context *clientContext) processReceive() (errResult error) {
	log.Info("Sharm receiver has running")
	defer log.Info("Sharm receiver was stopped")
	var data []byte
	var err error
	defer func() {
		context.closed = true
		if errResult != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error(errResult.Error())
		}
	}()

	for !context.closed {
		_, data, err = context.connection.ReadMessage()
		if err != nil {
			return errors.New("Failed to get message from server")
		}
		if len(data) == 0 {
			return errors.New("Received empty data from server")
		}

		log.WithFields(logrus.Fields{
			"data": data,
		}).Trace("Received new data from server")
		payload := data[1:]
		switch data[0] {
		case Input:
			if err = context.pty.Write(payload); err != nil {
				return errors.New("Failed to write data to termital")
			}

		case Ping:
			log.Debug("Received healthcheck packet")
			if err = context.write([]byte{Pong}); err != nil {
				return errors.New("Failed to send healthcheck packet")
			}

		case ResizeTerminal:
			log.Debug("Received resize packet")
			var args argResizeTerminal
			if err = json.Unmarshal(payload, &args); err != nil {
				return errors.New("Failed to parse resize message")
			}
			if err = context.pty.Resize(args.Columns, args.Rows); err != nil {
				return errors.New("Failed to resize terminal")
			}

		case Quit:
			log.Debug("Received quit packet")
			return nil

		default:
			return errors.New("Received unknown message type")
		}
	}
	return nil
}
