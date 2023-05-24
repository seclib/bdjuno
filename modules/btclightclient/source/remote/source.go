package remote

import (
	btclightclienttypes "github.com/babylonchain/babylon/x/btclightclient/types"

	btclightclientsource "github.com/forbole/bdjuno/v4/modules/btclightclient/source"

	"github.com/forbole/juno/v4/node/remote"
)

var (
	_ btclightclientsource.Source = &Source{}
)

// Source implements marketsource.Source using a remote node
type Source struct {
	*remote.Source
	queryClient btclightclienttypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, queryClient btclightclienttypes.QueryClient) *Source {
	return &Source{
		Source:      source,
		queryClient: queryClient,
	}
}

func (s Source) GetTip(height int64) (*btclightclienttypes.BTCHeaderInfo, error) {
	tip, err := s.queryClient.Tip(remote.GetHeightRequestContext(s.Ctx, height), &btclightclienttypes.QueryTipRequest{})
	if err != nil {
		return nil, err
	}
	return tip.GetHeader(), nil
}
