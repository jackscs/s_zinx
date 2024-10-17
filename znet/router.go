package znet

import "zinx/ziface"

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req *ziface.IRouter)  {}
func (br *BaseRouter) Handle(req *ziface.IRouter)     {}
func (br *BaseRouter) PostHandle(req *ziface.IRouter) {}
