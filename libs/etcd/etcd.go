package etcd

import (
	"context"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
)

const (
	DefaultEndpoints = "127.0.0.1:2379"
	DefaultNamespace = "global-prefix/"
	DialTimeout      = 5
	RequestTimeout   = 3
)

type EtcdOptions struct {
	Endpoints []string `yaml:"endpoints" mapstructure:"endpoints"`
	Namespace string   `yaml:"namespace" mapstructure:"namespace"`
	Username  string   `yaml:"username" mapstructure:"username"`
	Password  string   `yaml:"password" mapstructure:"password"`
}

type EtcdConfig struct {
	client  *clientv3.Client
	configs sync.Map
}

func (opts *EtcdOptions) loadDefaults() {
	if len(opts.Endpoints) == 0 {
		opts.Endpoints = []string{DefaultEndpoints}
	}
	if opts.Namespace == "" {
		opts.Namespace = DefaultNamespace
	}
}

func NewEtcdConfig(opts EtcdOptions) (*EtcdConfig, error) {
	opts.loadDefaults()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   opts.Endpoints,
		DialTimeout: DialTimeout * time.Second,
		Username:    opts.Username,
		Password:    opts.Password,
	})
	if err != nil {
		return nil, err
	}
	cli.KV = namespace.NewKV(cli.KV, opts.Namespace)
	cli.Watcher = namespace.NewWatcher(cli.Watcher, opts.Namespace)
	cli.Lease = namespace.NewLease(cli.Lease, opts.Namespace)

	ec := &EtcdConfig{
		client: cli,
	}
	// Watch whole namespace key when watch's key is empty
	go ec.watch("")

	return ec, nil
}

func (ec *EtcdConfig) watch(key string) {
	rch := ec.client.Watch(context.Background(), key, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type.String() == "PUT" {
				ec.configs.Store(string(ev.Kv.Key), string(ev.Kv.Value))
			} else if ev.Type.String() == "DELETE" {
				ec.configs.Delete(string(ev.Kv.Key))
			}
		}
	}
}

func (ec *EtcdConfig) Get(key string) (s string, err error) {
	if v, ok := ec.configs.Load(key); ok {
		s = v.(string)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout*time.Second)
		resp, err := ec.client.KV.Get(ctx, key)
		defer cancel()
		if err != nil {
			return "", err
		}
		for _, ev := range resp.Kvs {
			if string(ev.Key) == key {
				ec.configs.Store(key, string(ev.Value))
				s = string(ev.Value)
			}
		}
	}

	return
}

func (ec *EtcdConfig) Set(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout*time.Second)
	_, err := ec.client.Put(ctx, key, value)
	defer cancel()
	if err != nil {
		return err
	}
	ec.configs.Store(key, value)
	return nil
}

func (ec *EtcdConfig) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout*time.Second)
	_, err := ec.client.Delete(ctx, key)
	defer cancel()
	if err != nil {
		return err
	}
	ec.configs.Delete(key)
	return nil
}
