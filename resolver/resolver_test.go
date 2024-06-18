package resolver_test

import (
	"testing"

	"github.com/xeitf/lamp"
	"github.com/xeitf/lamp-grpc/resolver"
	"google.golang.org/grpc"
)

func TestXxx(t *testing.T) {
	lc, err := lamp.NewClient("etcd://127.0.0.1:2379/services")
	if err != nil {
		t.Errorf("lamp.Init: %s", err.Error())
		return
	}
	defer lc.Close()

	resolver.Register(lc)

	cc, err := grpc.NewClient("lamp:///user-svr",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	if err != nil {
		t.Errorf("NewClient: %s", err.Error())
		return
	}
	defer cc.Close()
}
