package resourceadapter

import (
	"context"
	"fmt"

	"github.com/shantanu-hashcash/go/amount"
	protocol "github.com/shantanu-hashcash/go/protocols/aurora"
	auroraContext "github.com/shantanu-hashcash/go/services/aurora/internal/context"
	"github.com/shantanu-hashcash/go/services/aurora/internal/db2/history"
	"github.com/shantanu-hashcash/go/support/errors"
	"github.com/shantanu-hashcash/go/support/render/hal"
	"github.com/shantanu-hashcash/go/xdr"
)

// PopulateLiquidityPool fills out the resource's fields
func PopulateLiquidityPool(
	ctx context.Context,
	dest *protocol.LiquidityPool,
	liquidityPool history.LiquidityPool,
	ledger *history.Ledger,
) error {
	dest.ID = liquidityPool.PoolID
	dest.FeeBP = liquidityPool.Fee
	typ, ok := xdr.LiquidityPoolTypeToString[liquidityPool.Type]
	if !ok {
		return errors.Errorf("unknown liquidity pool type: %d", liquidityPool.Type)
	}
	dest.Type = typ
	dest.TotalTrustlines = liquidityPool.TrustlineCount
	dest.TotalShares = amount.StringFromInt64(int64(liquidityPool.ShareCount))
	for _, reserve := range liquidityPool.AssetReserves {
		dest.Reserves = append(dest.Reserves, protocol.LiquidityPoolReserve{
			Asset:  reserve.Asset.StringCanonical(),
			Amount: amount.StringFromInt64(int64(reserve.Reserve)),
		})
	}

	dest.LastModifiedLedger = liquidityPool.LastModifiedLedger

	if ledger != nil {
		dest.LastModifiedTime = &ledger.ClosedAt
	}

	lb := hal.LinkBuilder{Base: auroraContext.BaseURL(ctx)}
	self := fmt.Sprintf("/liquidity_pools/%s", dest.ID)
	dest.Links.Self = lb.Link(self)
	dest.PT = dest.ID
	dest.Links.Transactions = lb.PagedLink(self, "transactions")
	dest.Links.Operations = lb.PagedLink(self, "operations")
	return nil
}
