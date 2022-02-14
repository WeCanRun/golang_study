package _rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

func TestServer(t *testing.T) {
	service := NewHelloService()
	if err := Register(service); err != nil {
		t.Fatal(err)
	}

	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("init server successfully")
	for {
		conn, err := listen.Accept()
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			t.Logf("handle request, %v", conn)
			rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		}()
	}
}

func TestClient(t *testing.T) {
	client, err := DialHelloService("tcp", ":1234")
	if err != nil {
		t.Fatal(err)
	}
	var res string
	for i := 1; i <= 10; i++ {
		req := fmt.Sprintf("test_client: %d", i)
		client.Hello(req, &res)
		t.Log(res)
	}
}
