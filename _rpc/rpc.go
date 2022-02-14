package _rpc

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const HelloServiceName = "study/_rpc/helloService"

type HelloServiceInterface interface {
	Hello(req string, resp *string) error
}

type helloService struct {
}

func NewHelloService() HelloServiceInterface {
	return &helloService{}
}

func (p *helloService) Hello(req string, resp *string) error {
	*resp = HelloServiceName + req
	return nil
}

func Register(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

type HelloServiceClient struct {
	*rpc.Client
}

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Printf("init client err, err: %v\n", err)
		return nil, err
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return &HelloServiceClient{Client: client}, nil
}

func (p *HelloServiceClient) Hello(req string, resp *string) error {
	return p.Client.Call(HelloServiceName+".Hello", req, resp)
}
