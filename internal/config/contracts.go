package config

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type ContractConfiger interface {
	ContractInfo() ContractInfo
}

type ContractInfo struct {
	Address  common.Address `fig:"address"`
	GasLimit uint64         `fig:"gas_limit"`
}

func NewContractConfiger(getter kv.Getter) ContractConfiger {
	return &contract{
		getter: getter,
	}
}

type contract struct {
	getter kv.Getter
	once   comfig.Once
}

func (w *contract) ContractInfo() ContractInfo {
	return w.once.Do(func() interface{} {
		config := ContractInfo{}
		err := figure.
			Out(&config).
			With(figure.BaseHooks, figure.EthereumHooks).
			From(kv.MustGetStringMap(w.getter, "contract")).
			Please()
		if err != nil {
			panic(err)
		}

		return config
	}).(ContractInfo)
}
