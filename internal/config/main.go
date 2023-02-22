package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	comfig.Logger
	NetworksConfiger
	WalletConfiger
	ContractConfiger
}

type config struct {
	comfig.Logger
	NetworksConfiger
	WalletConfiger
	ContractConfiger
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:           getter,
		ContractConfiger: NewContractConfiger(getter),
		WalletConfiger:   NewWalletConfiger(getter),
		NetworksConfiger: NewEthRPCConfiger(getter),
		Logger:           comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
