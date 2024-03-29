package profiles

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/juno/v3/modules"

	"github.com/forbole/bdjuno/v3/database"
	profilessource "github.com/forbole/bdjuno/v3/modules/profiles/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
)

// Module represent x/profiles module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source profilessource.Source
}

// NewModule returns a new Module instance
func NewModule(source profilessource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles"
}
