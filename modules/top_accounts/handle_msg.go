package top_accounts

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/v3/modules/utils"
	juno "github.com/forbole/juno/v3/types"
	"github.com/gogo/protobuf/proto"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	// Refresh x/bank available account balances
	addresses, err := m.messageParser(m.cdc, msg)
	if err != nil {
		return fmt.Errorf("error while parsing account addresses of message type %s: %s", proto.MessageName(msg), err)
	}

	addresses = utils.FilterNonAccountAddresses(addresses)
	err = m.bankModule.UpdateBalances(addresses, tx.Height)
	if err != nil {
		return fmt.Errorf("error while updating account available balances: %s", err)
	}

	err = m.RefreshTopAccountsSum(addresses)
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while refreshing balance: %s", err)
	}

	// Handle x/staking delegations, redelegations, and unbondings
	switch cosmosMsg := msg.(type) {

	case *stakingtypes.MsgDelegate:
		return m.handleMsgDelegate(tx.Height, cosmosMsg.DelegatorAddress)

	case *stakingtypes.MsgBeginRedelegate:
		return m.handleMsgBeginRedelegate(tx, index, cosmosMsg.DelegatorAddress)

	case *stakingtypes.MsgUndelegate:
		return m.handleMsgUndelegate(tx, index, cosmosMsg.DelegatorAddress)

	// Handle x/distribution delegator rewards
	case *distritypes.MsgWithdrawDelegatorReward:
		return m.handleMsgWithdrawDelegatorReward(tx.Height, cosmosMsg.DelegatorAddress)

	}

	return nil
}

func (m *Module) handleMsgDelegate(height int64, delAddr string) error {
	err := m.stakingModule.RefreshDelegations(height, delAddr)
	if err != nil {
		return fmt.Errorf("error while refreshing delegations while handling MsgDelegate: %s", err)
	}

	err = m.RefreshTopAccountsSum([]string{delAddr})
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while handling MsgDelegate: %s", err)
	}

	return nil
}

func (m *Module) handleMsgBeginRedelegate(
	tx *juno.Tx, index int, delAddr string) error {

	err := m.stakingModule.RefreshRedelegations(tx.Height, delAddr)
	if err != nil {
		return fmt.Errorf("error while refreshing redelegations while handling MsgBeginRedelegate: %s", err)
	}

	err = m.RefreshTopAccountsSum([]string{delAddr})
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while handling MsgBeginRedelegate: %s", err)
	}

	event, err := tx.FindEventByType(index, stakingtypes.EventTypeRedelegate)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	// When the time expires, refresh the delegations & redelegations
	time.AfterFunc(time.Until(completionTime), m.refreshDelegations(tx.Height, delAddr))
	time.AfterFunc(time.Until(completionTime), m.refreshRedelegations(tx, delAddr))

	return nil
}

// handleMsgUndelegate handles a MsgUndelegate storing the data inside the database
func (m *Module) handleMsgUndelegate(tx *juno.Tx, index int, delAddr string) error {
	err := m.stakingModule.RefreshUnbondings(tx.Height, delAddr)
	if err != nil {
		return fmt.Errorf("error while refreshing undelegations while handling MsgUndelegate: %s", err)
	}

	err = m.RefreshTopAccountsSum([]string{delAddr})
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while handling MsgUndelegate: %s", err)
	}

	event, err := tx.FindEventByType(index, stakingtypes.EventTypeUnbond)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	// When the time expires, refresh the delegations & unbondings & available balance
	time.AfterFunc(time.Until(completionTime), m.refreshDelegations(tx.Height, delAddr))
	time.AfterFunc(time.Until(completionTime), m.refreshUnbondings(tx.Height, delAddr))
	time.AfterFunc(time.Until(completionTime), m.refreshBalance(tx.Height, delAddr))

	return nil
}

func (m *Module) handleMsgWithdrawDelegatorReward(height int64, delAddr string) error {
	err := m.distrModule.RefreshDelegatorRewards(height, []string{delAddr})
	if err != nil {
		return fmt.Errorf("error while refreshing delegator rewards: %s", err)
	}

	err = m.bankModule.UpdateBalances([]string{delAddr}, height)
	if err != nil {
		return fmt.Errorf("error while updating account available balances with MsgWithdrawDelegatorReward: %s", err)
	}

	err = m.RefreshTopAccountsSum([]string{delAddr})
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum while handling MsgWithdrawDelegatorReward: %s", err)
	}

	return nil
}
