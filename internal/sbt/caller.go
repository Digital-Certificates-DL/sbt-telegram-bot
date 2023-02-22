package sbt

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/config"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/contracts"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/helpers"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/ipfs"
	"math/big"
	"strconv"
)

type Caller struct {
	Address  common.Address
	Client   *ethclient.Client
	ChainID  *big.Int
	Nonce    uint64
	Cfg      config.Config
	Contract *contracts.Factory
	Log      *logan.Entry
}

func NewCaller(cfg config.Config) (*Caller, error) {
	client := cfg.NetworksConfig().RPCEthEndpoint
	address := cfg.WalletInfo().Address
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nonce")
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nonce")
	}
	contractAddr := cfg.ContractInfo().Address
	contract, err := contracts.NewFactory(contractAddr, cfg.NetworksConfig().RPCEthEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new contract")
	}
	return &Caller{
		Address:  address,
		Client:   client,
		ChainID:  chainID,
		Nonce:    nonce,
		Contract: contract,
		Cfg:      cfg,
		Log:      cfg.Log(),
	}, nil
}

func (c Caller) TxOpts(cfg config.Config) (*bind.TransactOpts, error) {

	gasPrice, err := c.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get gas price")
	}
	gasLimit := cfg.ContractInfo().GasLimit
	return &bind.TransactOpts{
		From:  c.Address,
		Nonce: big.NewInt(int64(c.Nonce)),
		Signer: func(address common.Address, tx *ethtypes.Transaction) (*ethtypes.Transaction, error) {
			signer := ethtypes.NewEIP155Signer(c.ChainID)
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), cfg.WalletInfo().PrivateKey)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
		Value:     big.NewInt(0),
		GasPrice:  gasPrice,
		GasFeeCap: nil,
		GasTipCap: nil,
		GasLimit:  gasLimit,
		Context:   nil,
		NoSend:    false,
	}, nil
}

func (c Caller) Mint() error {

	fmt.Println("Enter the img path")

	imgPath, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get img")
	}

	fmt.Println("Enter the name of token")

	name, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get name")
	}

	fmt.Println("Enter the definition of token")
	definition, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get token's definition")
	}
	fmt.Println("Enter the address for mint")
	address, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get address for mint")
	}

	recipientAdd := common.HexToAddress(address)
	connector := ipfs.NewConnector(c.Cfg)
	img, err := connector.PrepareImage(imgPath)
	if err != nil {
		return errors.Wrap(err, "failed to send image to ipfs")
	}
	imgHash, err := connector.Upload(img)
	jsonHash, err := connector.PrepareJSON(name, definition, imgHash)
	if err != nil {
		return errors.Wrap(err, "failed to  prepare json")
	}

	preparedURI, err := connector.Upload(jsonHash)
	if err != nil {
		return errors.Wrap(err, "failed to upload ")
	}
	c.Log.Info(preparedURI)
	transactOpts, err := c.TxOpts(c.Cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create tx Opt")
	}

	mintTx, err := c.Contract.SafeMint(transactOpts, recipientAdd, preparedURI)
	if err != nil {
		return errors.Wrap(err, "failed to send tx to chain")
	}
	c.Log.Info("Token was minted")
	c.Log.Info("tx hash: ", mintTx.Hash())

	return nil
}

func (c Caller) Burn() error {

	fmt.Println("Enter the tokenID for  burning")

	tokenIDString, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get token id")
	}

	tokenID, err := strconv.Atoi(tokenIDString)
	if err != nil {
		return errors.Wrap(err, "failed to convert string to int")
	}
	transactOpts, err := c.TxOpts(c.Cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create tx Opt")
	}

	mintTx, err := c.Contract.Burn(transactOpts, big.NewInt(int64(tokenID)))
	if err != nil {
		return errors.Wrap(err, "failed to send tx to chain")
	}
	c.Log.Info("Token was burned")
	c.Log.Info("tx hash: ", mintTx.Hash())

	return nil
}

func (c Caller) Transfer() error {

	fmt.Println("Enter the address of sender")
	address, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get address of sender")
	}
	from := common.HexToAddress(address)

	fmt.Println("Enter the address of sender")
	address, err = helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get address of recipient")
	}
	to := common.HexToAddress(address)

	fmt.Println("Enter the tokenID for transfer")

	tokenIDString, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get token id")
	}
	tokenID, err := strconv.Atoi(tokenIDString)
	if err != nil {
		return errors.Wrap(err, "failed to convert string to int")
	}

	transactOpts, err := c.TxOpts(c.Cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create tx Opt")
	}

	mintTx, err := c.Contract.TransferToken(transactOpts, from, to, big.NewInt(int64(tokenID)))
	if err != nil {
		return errors.Wrap(err, "failed to send tx to chain")
	}
	c.Log.Info("Token was sent")
	c.Log.Info("tx hash: ", mintTx.Hash())

	return nil
}

func (c Caller) NewAdmin() error {
	fmt.Println("Enter the address of new admin")

	adminString, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get admin")
	}
	newAdmin := common.HexToAddress(adminString)

	transactOpts, err := c.TxOpts(c.Cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create tx Opt")
	}

	mintTx, err := c.Contract.SetNewAdmin(transactOpts, newAdmin)
	if err != nil {
		return errors.Wrap(err, "failed to send tx to chain")
	}
	c.Log.Info("New admin was added")
	c.Log.Info("tx hash: ", mintTx.Hash())

	return nil
}

func (c Caller) DeleteAdmin() error {

	fmt.Println("Enter the address that you want to delete")
	adminString, err := helpers.Read()
	if err != nil {
		return errors.Wrap(err, "failed to get admin")
	}
	admin := common.HexToAddress(adminString)

	transactOpts, err := c.TxOpts(c.Cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create tx Opt")
	}
	mintTx, err := c.Contract.DeleteAdmin(transactOpts, admin)
	if err != nil {
		return errors.Wrap(err, "failed to send tx to chain")
	}
	c.Log.Info("New admin was deleted")
	c.Log.Info("tx hash: ", mintTx.Hash())

	return nil
}

func (c Caller) OwnerOf() error {
	address, err := c.Contract.OwnerOf(&bind.CallOpts{}, big.NewInt(1))
	if err != nil {
		return errors.Wrap(err, "failed to get owner of token")
	}
	c.Log.Info(address)
	return nil
}

func (c Caller) Name() error {
	name, err := c.Contract.Name(&bind.CallOpts{})
	if err != nil {
		return errors.Wrap(err, "failed to get name")
	}
	c.Log.Info(name)
	return nil
}
