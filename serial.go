package dlt645

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/tarm/serial"
)

const (
	// Default timeout
	serialTimeout     = 5 * time.Second
	serialIdleTimeout = 60 * time.Second
)

type serialPort struct {
	serial.Config
	IdleTimeout time.Duration

	mu           sync.Mutex
	port         io.ReadWriteCloser
	lastActivity time.Time
	closeTimer   *time.Timer
	Logger       *log.Logger
}

func (mb *serialPort) Connect() (err error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	return mb.connect()
}

func (mb *serialPort) connect() error {
	if mb.port == nil {
		port, err := serial.OpenPort(&mb.Config)
		if err != nil {
			return err
		}
		mb.port = port
	}
	return nil
}

func (mb *serialPort) Close() (err error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	return mb.close()
}

func (mb *serialPort) close() (err error) {
	if mb.port != nil {
		err = mb.port.Close()
		mb.port = nil
	}
	return
}

func (mb *serialPort) logf(format string, v ...interface{}) {
	if mb.Logger != nil {
		mb.Logger.Printf(format, v...)
	}
}

func (mb *serialPort) startCloseTimer() {
	if mb.IdleTimeout <= 0 {
		return
	}
	if mb.closeTimer == nil {
		mb.closeTimer = time.AfterFunc(mb.IdleTimeout, mb.closeIdle)
	} else {
		mb.closeTimer.Reset(mb.IdleTimeout)
	}
}

func (mb *serialPort) closeIdle() {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if mb.IdleTimeout <= 0 {
		return
	}
	idle := time.Now().Sub(mb.lastActivity)
	if idle >= mb.IdleTimeout {
		mb.logf("dlt645: closing connection due to idle timeout: %v", idle)
		mb.close()
	}
}
