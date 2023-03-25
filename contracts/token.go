package contracts

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xueqianLu/routegen/contracts/erc20"
	"github.com/xueqianLu/routegen/log"
)

var (
	callOpt = &bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}
)

func GetTokenName(client *ethclient.Client, address string) string {
	addr := common.HexToAddress(address)
	contract, _ := erc20.NewErc20(addr, client)
	name, err := contract.Name(callOpt)
	if err != nil {
		log.WithField("err", err).WithField("addr", address).Error("get token name failed")
	}
	return name
}
