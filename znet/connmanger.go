package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection, 0),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 添加写锁，保护资源
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", len(connMgr.connections))
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 添加写锁保护资源
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", len(connMgr.connections))
}

func (connMgr *ConnManager) GetConn(connID uint32) (ziface.IConnection, error) {
	// 添加读锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 获取当前连接个数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn num = ", len(connMgr.connections))
}
