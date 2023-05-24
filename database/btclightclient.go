package database

import (
	"fmt"

	btclightclienttypes "github.com/babylonchain/babylon/x/btclightclient/types"
)

func (db *Db) SaveBtcHeader(header *btclightclienttypes.BTCHeaderInfo) error {
	stmt := `
INSERT INTO btc_header_info (hash, header, height, work) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (height) DO UPDATE 
    SET hash = excluded.hash,
        header = excluded.header,
        work = excluded.work
WHERE btc_header_info.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, header.Hash.MarshalHex(), header.Header.MarshalHex(), header.Height, header.Work.Uint64())
	if err != nil {
		return fmt.Errorf("error while storing BTC header: %s", err)
	}

	return nil
}
