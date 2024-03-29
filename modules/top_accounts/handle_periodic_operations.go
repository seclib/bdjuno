package top_accounts

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v3/modules/utils"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "top accounts").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.RefreshRewards)
	}); err != nil {
		return fmt.Errorf("error while setting up top accounts periodic operation: %s", err)
	}

	return nil
}

// RefreshRewards refreshes the rewards for all delegators
func (m *Module) RefreshRewards() error {
	log.Trace().Str("module", "top accounts").Str("operation", "refresh rewards").
		Msg("refreshing delegators rewards")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	// Get the delegators
	delegators, err := m.db.GetDelegators()
	if err != nil {
		return fmt.Errorf("error while getting delegators: %s", err)
	}

	if len(delegators) == 0 {
		return nil
	}

	// Refresh rewards
	err = m.distrModule.RefreshDelegatorRewards(height, delegators)
	if err != nil {
		return fmt.Errorf("error while refreshing delegators rewards: %s", err)
	}

	err = m.RefreshTopAccountsSum(delegators)
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum value: %s", err)
	}

	return nil
}
