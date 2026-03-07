package snet

import "github.com/Games55k/Sinx/siface"

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req siface.IRequest)  {}
func (br *BaseRouter) Handle(req siface.IRequest)     {}
func (br *BaseRouter) PostHandle(req siface.IRequest) {}