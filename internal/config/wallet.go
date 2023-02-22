package config

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"reflect"
)

type WalletConfiger interface {
	WalletInfo() WalletInfo
}

type WalletInfo struct {
	Address    common.Address    `fig:"address,required"`
	PrivateKey *ecdsa.PrivateKey `fig:"private_key,required"`
}

func NewWalletConfiger(getter kv.Getter) WalletConfiger {
	return &walleter{
		getter: getter,
	}
}

type walleter struct {
	getter kv.Getter
	once   comfig.Once
}

func (w *walleter) WalletInfo() WalletInfo {
	return w.once.Do(func() interface{} {
		config := WalletInfo{}

		err := figure.
			Out(&config).
			With(figure.BaseHooks, WalletHook).
			From(kv.MustGetStringMap(w.getter, "wallet")).
			Please()
		if err != nil {
			panic(err)
		}

		return config
	}).(WalletInfo)
}

var WalletHook = figure.Hooks{
	"common.Address": func(value interface{}) (reflect.Value, error) {
		addressHex, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, err
		}

		address := common.HexToAddress(addressHex)
		return reflect.ValueOf(address), nil
	},
	"*ecdsa.PrivateKey": func(value interface{}) (reflect.Value, error) {
		key, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, err
		}

		privateKey, err := crypto.HexToECDSA(key)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to parse hex to private key eth wallet")
		}
		return reflect.ValueOf(privateKey), err
	},
}
