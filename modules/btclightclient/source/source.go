package source

import btclightclienttypes "github.com/babylonchain/babylon/x/btclightclient/types"

type Source interface {
	GetTip(height int64) (*btclightclienttypes.BTCHeaderInfo, error)
}
