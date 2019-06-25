package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/sharm/term"
)

const defSharmHost = "sharm.io"
const defSharmLogFile = "sharm.log"

var log *logrus.Logger

var cstDialer = websocket.Dialer{
	Subprotocols:    []string{"sharm-stream"},
	ReadBufferSize:  1024 * 100,
	WriteBufferSize: 1024 * 100,
}

func exitHandler() {
	log.Info("Sharm client was exited")
}

func initLogging(logPath, logLevel string) {
	logrus.RegisterExitHandler(exitHandler)
	log = logrus.New()
	log.Out = os.Stderr

	if logLevel != "" {
		logPathFile := path.Join(path.Clean(logPath), defSharmLogFile)
		logFile, err := os.OpenFile(logPathFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = logFile
		} else {
			log.Error("Failed to log to file, using default stderr")
		}
	}

	switch logLevel {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warning":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.FatalLevel)
	}
}

func main() {
	token := os.Getenv("TOKEN")
	if len(os.Args) < 2 && token == "" {
		fmt.Printf("Usage: %s <token> [command]\n", os.Args[0])
		os.Exit(0)
	}

	host := defSharmHost
	logLevel := ""
	logPath := "."
	command := ""

	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	if os.Getenv("LOG_PATH") != "" {
		logPath = os.Getenv("LOG_PATH")
	}
	if os.Getenv("COMMAND") != "" {
		command = os.Getenv("COMMAND")
	}
	if token == "" {
		token = os.Args[1]
	}
	if command == "" && len(os.Args) >= 3 {
		command = strings.Join(os.Args[2:], " ")
	}

	initLogging(logPath, logLevel)
	log.WithFields(logrus.Fields{
		"HOST":      host,
		"LOG_LEVEL": logLevel,
		"LOG_PATH":  logPath,
		"COMMAND":   command,
		"TOKEN":     token,
	}).Debug("Sharm client has initialized by next vars")

	dialer := cstDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	conn, _, err := dialer.Dial(fmt.Sprintf("wss://%s/stream/%s", host, token), nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Can't create connection to Sharm server")
	}
	defer conn.Close()

	pty := &term.Term{}
	if err = pty.Start(command); err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Can't create new terminal")
	}
	defer pty.Close()

	log.Info("Sharm client has started")
	context := &clientContext{
		connection: conn,
		pty:        pty,
		writeMutex: &sync.Mutex{},
	}
	// TODO: here may use reconnect logic
	context.goHandleClient()
}
