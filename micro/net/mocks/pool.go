package mocks

import (
	"net"
	"time"
)

type Pool struct {
	idlesConns  chan *idleConn
	reqQueue    []connReq
	maxCnt      int
	cnt         int
	macIdleTime time.Duration
	initCat     int
	factory     func() (net.Conn, error)
}

func NewPool(initCnt int, macIdleCnt int, maxCnt int, maxIdleTime time.Duration, factory func() (net.Conn, error)) {
}
