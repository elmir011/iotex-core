// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimesd. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/iotexproject/iotex-core/pkg/version"
	"github.com/iotexproject/iotex-core/proto"
)

const (
	// SettleDepositIntrinsicGas represents the intrinsic gas for the deposit action
	SettleDepositIntrinsicGas = uint64(10000)
)

// SettleDeposit represents the action to settle a deposit on the sub-chain
type SettleDeposit struct {
	AbstractAction
	amount *big.Int
	index  uint64
}

// NewSettleDeposit instantiates a deposit settlement to sub-chain action struct
func NewSettleDeposit(
	nonce uint64,
	amount *big.Int,
	index uint64,
	sender string,
	recipient string,
	gasLimit uint64,
	gasPrice *big.Int,
) *SettleDeposit {
	return &SettleDeposit{
		AbstractAction: AbstractAction{
			version:  version.ProtocolVersion,
			nonce:    nonce,
			srcAddr:  sender,
			dstAddr:  recipient,
			gasLimit: gasLimit,
			gasPrice: gasPrice,
		},
		amount: amount,
		index:  index,
	}
}

// Amount returns the amount
func (sd *SettleDeposit) Amount() *big.Int { return sd.amount }

// Index returns the index of the deposit on main-chain's sub-chain account
func (sd *SettleDeposit) Index() uint64 { return sd.index }

// Sender returns the sender address. It's the wrapper of Action.SrcAddr
func (sd *SettleDeposit) Sender() string { return sd.SrcAddr() }

// SenderPublicKey returns the sender public key. It's the wrapper of Action.SrcPubkey
func (sd *SettleDeposit) SenderPublicKey() keypair.PublicKey { return sd.SrcPubkey() }

// Recipient returns the recipient address. It's the wrapper of Action.DstAddr. The recipient should be an address on
// the sub-chain
func (sd *SettleDeposit) Recipient() string { return sd.DstAddr() }

// ByteStream returns a raw byte stream of the settle deposit action
func (sd *SettleDeposit) ByteStream() []byte {
	return byteutil.Must(proto.Marshal(sd.Proto()))
}

// Proto converts SettleDeposit to protobuf's ActionPb
func (sd *SettleDeposit) Proto() *iproto.SettleDepositPb {
	act := &iproto.SettleDepositPb{
		Recipient: sd.dstAddr,
		Index:     sd.index,
	}
	if sd.amount != nil && len(sd.amount.Bytes()) > 0 {
		act.Amount = sd.amount.Bytes()
	}
	return act
}

// LoadProto converts a protobuf's ActionPb to SettleDeposit
func (sd *SettleDeposit) LoadProto(pbDpst *iproto.SettleDepositPb) error {
	if pbDpst == nil {
		return errors.New("empty action proto to load")
	}
	if sd == nil {
		return errors.New("nil action to load proto")
	}
	*sd = SettleDeposit{}

	sd.amount = big.NewInt(0)
	sd.amount.SetBytes(pbDpst.GetAmount())
	sd.index = pbDpst.GetIndex()
	return nil
}

// IntrinsicGas returns the intrinsic gas of a settle deposit
func (sd *SettleDeposit) IntrinsicGas() (uint64, error) { return SettleDepositIntrinsicGas, nil }

// Cost returns the total cost of a settle deposit
func (sd *SettleDeposit) Cost() (*big.Int, error) {
	intrinsicGas, err := sd.IntrinsicGas()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get intrinsic gas for the settle deposit")
	}
	return big.NewInt(0).Mul(sd.GasPrice(), big.NewInt(0).SetUint64(intrinsicGas)), nil
}
