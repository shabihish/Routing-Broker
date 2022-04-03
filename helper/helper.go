package helper

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// From https://golangbyexample.com/print-output-text-color-console/

const ColorReset = "\033[0m"

const ColorRed = "\033[31m"
const ColorGreen = "\033[32m"
const ColorYellow = "\033[33m"
const ColorBlue = "\033[34m"
const ColorCyan = "\033[36m"
const ColorWhite = "\033[0m"

type Logger struct {
	colorMu sync.Mutex
}

func (logger *Logger) PrintLogInColor(color string, format string, args ...interface{}) {
	logger.colorMu.Lock()
	fmt.Print(color)
	fmt.Printf(strconv.Itoa(int(time.Now().UnixNano()/int64(time.Millisecond))) /*.Format("15:04:05")*/ +" - "+format, args...)
	fmt.Print(ColorReset)
	logger.colorMu.Unlock()
}

type Msg struct {
	valid bool
	Id    int
	data  string
}

func NewMsg(valid bool, msgId int, data string) *Msg {
	return &Msg{valid, msgId, data}
}

func (msg *Msg) Invalidate() {
	msg.valid = false
}

func (msg *Msg) IsValid() bool {
	return msg.valid
}

type ClientInterface interface {
	PutAcknowledgement(msgId int)
	PutNewServerResponse(msg Msg, responseToMessageId int)
	GetClientId() int
}

type ServerInterface interface {
	IsRunning() bool
	PutMessage(msg Msg, clientId int) bool
	PutAcknowledgement(responseId int, clientId int)
}
