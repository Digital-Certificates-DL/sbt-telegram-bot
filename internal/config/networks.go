package config

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cast"
	figure "gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"reflect"
)

type NetworksConfiger interface {
	NetworksConfig() *NetworksConfig
}

type NetworksConfig struct {
	RPCEthEndpoint *ethclient.Client `fig:"rpc_eth"`
	IPFSEndpoint   string            `fig:"ipfs"`
	IpfsPrId       string            `fig:"ipfs_pr_id"`
	IpfsPrKey      string            `fig:"ipfs_pr_key"`
}

func NewEthRPCConfiger(getter kv.Getter) NetworksConfiger {
	return &ethRPCConfig{
		getter: getter,
	}
}

type ethRPCConfig struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *ethRPCConfig) NetworksConfig() *NetworksConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "networks")
		config := NetworksConfig{}
		err := figure.Out(&config).With(figure.BaseHooks, ClientHook).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*NetworksConfig)
}

var ClientHook = figure.Hooks{
	"*ethclient.Client": func(value interface{}) (reflect.Value, error) {
		endpoint, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, err
		}

		dial, err := ethclient.Dial(endpoint)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to dial ethclient")
		}
		return reflect.ValueOf(dial), nil
	},
}
