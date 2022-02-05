package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/avast/retry-go"
	// "github.com/cosmos/relayer/relayer"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (i *Indexer) IndexIBCTransactions(height int64) error {
	block, err := i.Client.RPCClient.Block(context.Background(), &height)
	if err != nil {
		if err = retry.Do(func() error {
			block, err = i.Client.RPCClient.Block(context.Background(), &height)
			if err != nil {
				return err
			}

			return nil
		}, RtyAtt, RtyDel, RtyErr, retry.DelayType(retry.BackOffDelay), retry.OnRetry(func(n uint, err error) {
			i.LogRetryGetBlock(n, err, height)
		})); err != nil {
			return err
		}
	}
	if block != nil {
		i.ParseTxs(block)
	}
	return nil
}

func (i *Indexer) ParseTxs(block *coretypes.ResultBlock) {
	for index, tx := range block.Block.Data.Txs {
		sdkTx, err := i.Client.Codec.TxConfig.TxDecoder()(tx)
		if err != nil {
			// TODO application specific txs fail here (e.g. DEX swaps, Akash deployments, etc.)
			fmt.Printf("[Height %d] {%d/%d txs} - Failed to decode tx. Err: %s \n", block.Block.Height, index+1, len(block.Block.Data.Txs), err.Error())
			continue
		}

		txRes, err := i.Client.QueryTx(hex.EncodeToString(tx.Hash()))
		if err != nil {
			fmt.Printf("[Height %d] {%d/%d txs} - Failed to query tx results. Err: %s \n", block.Block.Height, index+1, len(block.Block.Data.Txs), err.Error())
			continue
		}

		fee := sdkTx.(sdk.FeeTx)
		var feeAmount, feeDenom string
		if len(fee.GetFee()) == 0 {
			feeAmount = "0"
			feeDenom = ""
		} else {
			feeAmount = fee.GetFee()[0].Amount.String()
			feeDenom = fee.GetFee()[0].Denom
		}

		if txRes.TxResult.Code > 0 {
			json := fmt.Sprintf("{\"error\":\"%s\"}", txRes.TxResult.Log)
			err = i.InsertTxRow(tx.Hash(), json, feeAmount, feeDenom, block.Block.Height, txRes.TxResult.GasUsed,
				txRes.TxResult.GasWanted, block.Block.Time, txRes.TxResult.Code)

			i.LogTxInsertion(err, index, len(sdkTx.GetMsgs()), len(block.Block.Data.Txs), block.Block.Height)
		} else {
			err = i.InsertTxRow(tx.Hash(), txRes.TxResult.Log, feeAmount, feeDenom, block.Block.Height, txRes.TxResult.GasUsed,
				txRes.TxResult.GasWanted, block.Block.Time, txRes.TxResult.Code)

			i.LogTxInsertion(err, index, len(sdkTx.GetMsgs()), len(block.Block.Data.Txs), block.Block.Height)
		}

		for msgIndex, msg := range sdkTx.GetMsgs() {
			i.HandleMsg(msg, msgIndex, block.Block.Height, tx.Hash())
		}
	}
}

func (i *Indexer) LogTxInsertion(err error, msgIndex, msgs, txs int, height int64) {
	if err != nil {
		i.logger.Info("[Height %d] {%d/%d txs} - Failed to write tx to db. Err: %s", height, msgIndex+1, txs, err.Error())
	} else {
		i.logger.Info("[Height %d] {%d/%d txs} - Successfuly wrote tx to db with %d msgs.", height, msgIndex+1, txs, msgs)
	}
}

func (i *Indexer) HandleMsg(msg sdk.Msg, msgIndex int, height int64, hash []byte) {
	switch m := msg.(type) {
	case *transfertypes.MsgTransfer:
		err := i.InsertMsgTransferRow(hash, m.Token.Denom, m.SourceChannel, m.Route(), m.Token.Amount.String(), m.Sender,
			i.Client.MustEncodeAccAddr(m.GetSigners()[0]), m.Receiver, m.SourcePort, msgIndex)
		if err != nil {
			i.logger.Info("Failed to insert MsgTransfer", "index", msgIndex, "height", height, "err", err.Error())
		}
	case *channeltypes.MsgRecvPacket:
		err := i.InsertMsgRecvPacketRow(hash, m.Signer, m.Packet.SourceChannel,
			m.Packet.DestinationChannel, m.Packet.SourcePort, m.Packet.DestinationPort, msgIndex)
		if err != nil {
			i.logger.Info("Failed to insert MsgRecvPacket", "index", msgIndex, "height", height, "err", err.Error())
		}
	case *channeltypes.MsgTimeout:
		err := i.InsertMsgTimeoutRow(hash, m.Signer, m.Packet.SourceChannel,
			m.Packet.DestinationChannel, m.Packet.SourcePort, m.Packet.DestinationPort, msgIndex)
		if err != nil {
			i.logger.Info("Failed to insert MsgTimeout", "index", msgIndex, "height", height, "err", err.Error())
		}
	case *channeltypes.MsgAcknowledgement:
		err := i.InsertMsgAckRow(hash, m.Signer, m.Packet.SourceChannel,
			m.Packet.DestinationChannel, m.Packet.SourcePort, m.Packet.DestinationPort, msgIndex)
		if err != nil {
			i.logger.Info("Failed to insert MsgAck", "index", msgIndex, "height", height, "err", err.Error())
		}
	default:
		// TODO: do we need to do anything here?
	}
}
