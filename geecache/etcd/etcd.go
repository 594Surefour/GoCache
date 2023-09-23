package etcd

import (
	"fmt"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/clientv3"
)

func Register(service string, addr string, stop chan error) error {
	config := clientv3.Config{}

	cli, err := clientv3.New(config)
	if err != nil {
		return fmt.Errorf("create etcd client failed: %v", err)
	}
	defer cli.Close()
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		return fmt.Errorf("create lease failed: %v", err)
	}
	leaseId = resp.ID
	//注册服务
	err = etcdAdd(cli)
	//设置服务租约检测
}
