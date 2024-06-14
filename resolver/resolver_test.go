package resolver_test

import (
	"testing"
	"time"

	"github.com/xeitf/lamp"
	"github.com/xeitf/lamp-grpc/resolver"
	"google.golang.org/grpc"
)

func TestXxx(t *testing.T) {
	close, err := resolver.Init("etcd://127.0.0.1:2379/services")
	if err != nil {
		t.Errorf("Init: %s", err.Error())
		return
	}
	defer close()

	cancel, err := lamp.Register("user-svr",
		lamp.WithTTL(5),
		lamp.WithPublic("127.0.0.1:8999"),
		lamp.WithPublic("127.0.0.1:80", "http"),
	)
	if err != nil {
		t.Errorf("Register: %s", err.Error())
		return
	}
	defer cancel()

	cancel2, err := lamp.Register("user-svr",
		lamp.WithPublic("127.0.0.1:8990"),
	)
	if err != nil {
		t.Errorf("Register: %s", err.Error())
		return
	}
	defer cancel2()

	conn, err := grpc.NewClient("lamp:///user-svr",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	if err != nil {
		t.Errorf("NewClient: %s", err.Error())
		return
	}
	defer conn.Close()

	time.Sleep(10 * time.Second)
}
