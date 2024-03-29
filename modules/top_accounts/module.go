package top_accounts

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v3/database"

	"github.com/forbole/juno/v3/modules"
	junomessages "github.com/forbole/juno/v3/modules/messages"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
	// _ modules.PeriodicOperationsModule = &Module{}
)

// Module represent x/gov module
type Module struct {
	cdc           codec.Codec
	db            *database.Db
	messageParser junomessages.MessageAddressesParser
	bankModule    BankModule
	distrModule   DistrModule
	stakingModule StakingModule
}

// NewModule returns a new Module instance
func NewModule(
	bankModule BankModule,
	distrModule DistrModule,
	stakingModule StakingModule,
	messageParser junomessages.MessageAddressesParser,
	cdc codec.Codec,
	db *database.Db,
) *Module {
	return &Module{
		cdc:           cdc,
		bankModule:    bankModule,
		distrModule:   distrModule,
		messageParser: messageParser,
		stakingModule: stakingModule,
		db:            db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "top_accounts"
}
