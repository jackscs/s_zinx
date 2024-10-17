package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer ziface.IServer

	// 当前socket tcp连接套接字
	Conn *net.TCPConn

	// 当前连接的ID
	ConnID uint32

	// 当前连接的状态
	isClosed bool

	// 该连接处理方法的API
	MsgHandler ziface.IMsgHandle

	// 告知该连接已经停止/退出的 channel
	ExitBufChan chan bool

	Router ziface.IRouter

	// 创建读写chan
	msgChan chan []byte

	// 创建有缓冲的读写chan
	msgBuffChan chan []byte

	// 连接属性
	property map[string]interface{}
	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

func NewConntion(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandle,
		ExitBufChan: make(chan bool, 1),
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GLOBALOBJ.MaxMsgChanLen),
		property:    make(map[string]interface{}),
	}

	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// 创建处理conn的函数
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {

		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte, dp.HeadDataLen())

		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBufChan <- true
			continue
		}

		// 拆包的到dataLen和data,放在msg.Data
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBufChan <- true
			continue
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBufChan <- true
				continue
			}

		}
		msg.SetData(data)

		req := Request{
			Conn:    c,
			Message: msg,
		}

		if utils.GLOBALOBJ.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsghandle(&req)
		}

	}

}

/*
写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			// 有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Data error:, ", err, " Conn Writer exit")
					return
				} else {
					fmt.Println("msgBuffChan is Closed")
					break
				}
			}
		case <-c.ExitBufChan:
			// conn连接关闭
			return
		}
	}
}

// 启动，让当前链接工作
func (c *Connection) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	go c.StartWriter()

	// 执行用户自定义钩子函数
	c.TcpServer.CallOnConnStart(c)

	for {
		select {
		//得到退出消息，不再阻塞
		case <-c.ExitBufChan:
			return
		}
	}

}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}

	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 通知从缓冲队列中读取数据的的业务,通道已经关闭
	c.ExitBufChan <- true

	// 将连接从管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	close(c.ExitBufChan)
	close(c.msgChan)
}

// 从当前链接中获取原始的tcp链接
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的地址
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection close when send msg")
	}

	// 将data数据打包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg ")
	}

	// 写回客户端
	c.msgChan <- msg

	return nil
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection close when send msg")
	}

	// 将data数据打包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg ")
	}

	// 写回客户端
	c.msgBuffChan <- msg

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) error {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if _, ok := c.property[key]; ok {
		return errors.New("key already existed")
	} else {
		c.property[key] = value
		return nil
	}
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("key no foued")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
