package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	MsgHandler ziface.IMsgHandle

	//当前Server的链接管理器
	ConnMgr ziface.IConnManger

	// 新增两个hook函数原型
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP:%s,Port:%d", s.IP, s.Port)
	go func() {

		// 开启线程池
		s.MsgHandler.StartWorkerPool()

		//获取一个tcp的addr
		tcpAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}

		// 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, tcpAddr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, err)
			return
		}

		// 打印监听信息
		fmt.Println("start Zinx server", s.Name, "succ at now...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accepc err:", err)
				continue
			}

			if s.ConnMgr.Len() >= utils.GLOBALOBJ.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConntion(s, conn, cid, s.MsgHandler)
			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Printf("[STOP] Zinx server,name:%s", s.Name)

	// 停止时清除所有连接
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	for {
		time.Sleep(time.Second * 10)
	}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)

	fmt.Println("Add Router succ! ")
}

func (s *Server) GetConnMgr() ziface.IConnManger {
	return s.ConnMgr
}

// 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 设置该Server的断开创建时Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// 调用断开OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func NewServer() ziface.IServer {
	//先初始化全局配置文件
	utils.GLOBALOBJ.Reload()

	return &Server{
		Name:       utils.GLOBALOBJ.Name,
		IPVersion:  "tcp4",
		IP:         utils.GLOBALOBJ.Host,
		Port:       utils.GLOBALOBJ.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
}
