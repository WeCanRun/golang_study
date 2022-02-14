package etcd

import (
	"context"
	"encoding/json"
	"testing"
)

func TestInitConfig(t *testing.T) {
	defer Clear(context.Background())
	c := config{server{
		Addr:  "0.0.0.0",
		Port:  "8080",
		Name:  "test",
		Token: "test",
	}, database{
		Type:     "mysql",
		Url:      "mysql://127.0.0.1:3306?test",
		User:     "root",
		Password: "123456",
	}}
	marshal, err := json.Marshal(&c)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	appName := "test"
	if err := Put(context.Background(), configPath(appName), string(marshal)); err != nil {
		t.Fatalf("err: %v", err)
	}

	InitConfig(appName)

	if err := Put(context.Background(), configPath(appName), string(marshal)); err != nil {
		t.Fatalf("err: %v", err)
	}

	select {}
}
