package database

import (
	"fmt"

	btclightclienttypes "github.com/babylonchain/babylon/x/btclightclient/types"
	"github.com/rs/zerolog/log"
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

	_, err := db.SQL.Exec(stmt, header.Hash.MarshalHex(), header.Header.MarshalHex(), header.Height, header.Work.String())
	if err != nil {
		return fmt.Errorf("error while storing BTC header: %s", err)
	}

	log.Info().Msgf("Saved BTC header at height %d", header.Height)
	return nil
}
