package siface

type IServer interface {
	Start()
	Stop()
	Serve()
}