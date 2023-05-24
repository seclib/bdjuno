package btclightclient

import (
	"fmt"

	btclightclienttypes "github.com/babylonchain/babylon/x/btclightclient/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	junotypes "github.com/forbole/juno/v4/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *junotypes.Tx) error {
	if _, ok := msg.(*btclightclienttypes.MsgInsertHeader); ok {
		return m.saveBtcHeader()
	}

	return nil
}

func (m *Module) saveBtcHeader() error {
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	tip, err := m.source.GetTip(height)
	if err != nil {
		return fmt.Errorf("error while getting latest tip: %s", err)
	}

	return m.db.SaveBtcHeader(tip)
}
