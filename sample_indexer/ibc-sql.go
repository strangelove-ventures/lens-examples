package main

import (
	"fmt"
	"time"
)

func (i *Indexer) CreateTables() error {
	txs := "CREATE TABLE IF NOT EXISTS txs ( " +
		"hash bytea PRIMARY KEY, " +
		"block_time TIMESTAMP NOT NULL, " +
		"chainid TEXT NOT NULL, " +
		"block_height BIGINT NOT NULL, " +
		"raw_log JSONB NOT NULL," +
		"code INT NOT NULL, " +
		"fee_amount TEXT, " +
		"fee_denom TEXT, " +
		"gas_used BIGINT NOT NULL," +
		"gas_wanted BIGINT NOT NULL" +
		")"

	transfer := "CREATE TABLE IF NOT EXISTS msg_transfer (" +
		"tx_hash bytea," +
		"msg_index INT," +
		"signer TEXT NOT NULL," +
		"sender TEXT NOT NULL," +
		"receiver TEXT NOT NULL," +
		"amount TEXT NOT NULL," +
		"denom TEXT NOT NULL," +
		"src_chan TEXT NOT NULL," +
		"src_port TEXT NOT NULL," +
		"route TEXT NOT NULL," +
		"PRIMARY KEY (tx_hash, msg_index)," +
		"FOREIGN KEY (tx_hash) REFERENCES txs(hash) ON DELETE CASCADE" +
		")"

	recvpacket := "CREATE TABLE IF NOT EXISTS msg_recvpacket ( " +
		"tx_hash bytea," +
		"msg_index INT," +
		"signer TEXT NOT NULL," +
		"src_chan TEXT NOT NULL," +
		"dst_chan TEXT NOT NULL," +
		"src_port TEXT NOT NULL," +
		"dst_port TEXT NOT NULL," +
		"PRIMARY KEY (tx_hash, msg_index)," +
		"FOREIGN KEY (tx_hash) REFERENCES txs(hash) ON DELETE CASCADE" +
		")"

	timeout := "CREATE TABLE IF NOT EXISTS msg_timeout (" +
		"tx_hash bytea," +
		"msg_index INT," +
		"signer TEXT NOT NULL," +
		"src_chan TEXT NOT NULL," +
		"dst_chan TEXT NOT NULL," +
		"src_port TEXT NOT NULL," +
		"dst_port TEXT NOT NULL," +
		"PRIMARY KEY (tx_hash, msg_index)," +
		"FOREIGN KEY (tx_hash) REFERENCES txs(hash) ON DELETE CASCADE" +
		")"

	acks := "CREATE TABLE IF NOT EXISTS msg_ack (" +
		"tx_hash bytea," +
		"msg_index INT," +
		"signer TEXT NOT NULL," +
		"src_chan TEXT NOT NULL," +
		"dst_chan TEXT NOT NULL," +
		"src_port TEXT NOT NULL," +
		"dst_port TEXT NOT NULL," +
		"PRIMARY KEY (tx_hash, msg_index)," +
		"FOREIGN KEY (tx_hash) REFERENCES txs(hash) ON DELETE CASCADE" +
		")"

	tables := []string{txs, transfer, recvpacket, timeout, acks}
	for _, table := range tables {
		if _, err := i.DB.Exec(table); err != nil {
			return err
		}
	}
	return nil
}

func (i *Indexer) InsertTxRow(hash []byte, log, feeAmount, feeDenom string, height, gasUsed, gasWanted int64, timestamp time.Time, code uint32) error {
	query := "INSERT INTO txs(hash, block_time, chainid, block_height, raw_log, code, gas_used, gas_wanted, fee_amount, fee_denom) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	stmt, err := i.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("Fail to create query for new tx. Err: %s \n", err.Error())
	}

	_, err = stmt.Exec(hash, timestamp, i.Client.Config.ChainID, height, log, code, gasUsed, gasWanted, feeAmount, feeDenom)
	if err != nil {
		return fmt.Errorf("Fail to execute query for new tx. Err: %s \n", err.Error())
	}

	return nil
}

func (i *Indexer) InsertMsgTransferRow(hash []byte, denom, srcChan, route, amount, sender, signer, receiver, port string, msgIndex int) error {
	query := "INSERT INTO msg_transfer(tx_hash, msg_index, amount, denom, src_chan, route, signer, sender, receiver, src_port) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	stmt, err := i.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("Fail to create query for MsgTransfer. Err: %s \n", err.Error())
	}

	_, err = stmt.Exec(hash, msgIndex, amount, denom, srcChan, route, signer, sender, receiver, port)
	if err != nil {
		return fmt.Errorf("Fail to execute query for MsgTransfer. Err: %s \n", err.Error())
	}

	return nil
}

func (i *Indexer) InsertMsgTimeoutRow(hash []byte, signer, srcChan, dstChan, srcPort, dstPort string, msgIndex int) error {
	query := "INSERT INTO msg_timeout(tx_hash, msg_index, signer, src_chan, dst_chan, src_port, dst_port) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7)"
	stmt, err := i.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("Fail to create query for MsgTimeout. Err: %s \n", err.Error())
	}

	_, err = stmt.Exec(hash, msgIndex, signer, srcChan, dstChan, srcPort, dstPort)
	if err != nil {
		return fmt.Errorf("Fail to execute query for MsgTimeout. Err: %s \n", err.Error())
	}

	return nil
}

func (i *Indexer) InsertMsgRecvPacketRow(hash []byte, signer, srcChan, dstChan, srcPort, dstPort string, msgIndex int) error {
	query := "INSERT INTO msg_recvpacket(tx_hash, msg_index, signer, src_chan, dst_chan, src_port, dst_port) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7)"
	stmt, err := i.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("Fail to create query for MsgRecvPacket. Err: %s \n", err.Error())
	}

	_, err = stmt.Exec(hash, msgIndex, signer, srcChan, dstChan, srcPort, dstPort)
	if err != nil {
		return fmt.Errorf("Fail to execute query for MsgRecvPacket. Err: %s \n", err.Error())
	}

	return nil
}

func (i *Indexer) InsertMsgAckRow(hash []byte, signer, srcChan, dstChan, srcPort, dstPort string, msgIndex int) error {
	query := "INSERT INTO msg_ack(tx_hash, msg_index, signer, src_chan, dst_chan, src_port, dst_port) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7)"
	stmt, err := i.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("Fail to create query for MsgAck. Err: %s \n", err.Error())
	}

	_, err = stmt.Exec(hash, msgIndex, signer, srcChan, dstChan, srcPort, dstPort)
	if err != nil {
		return fmt.Errorf("Fail to execute query for MsgAck. Err: %s \n", err.Error())
	}

	return nil
}

func (i *Indexer) GetLastStoredBlock(chainId string) (int64, error) {
	var height int64
	if err := i.DB.QueryRow("SELECT MAX(block_height) FROM txs WHERE chainid=$1", chainId).Scan(&height); err != nil {
		return 1, err
	}
	return height, nil
}
