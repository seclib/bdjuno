package local

import (
	"fmt"

	btclightclienttypes "github.com/babylonchain/babylon/x/btclightclient/types"
	"github.com/forbole/juno/v4/node/local"

	btclightclientsource "github.com/forbole/bdjuno/v4/modules/btclightclient/source"
)

var (
	_ btclightclientsource.Source = &Source{}
)

// Source implements govsource.Source by using a local node
type Source struct {
	*local.Source
	q btclightclienttypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, querier btclightclienttypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      querier,
	}
}

func (s Source) GetTip(height int64) (*btclightclienttypes.BTCHeaderInfo, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	tip, err := s.q.Tip(ctx, &btclightclienttypes.QueryTipRequest{})
	if err != nil {
		return nil, err
	}
	return tip.GetHeader(), nil
}
