package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/public-awesome/stargaze/v12/testutil/keeper"
	"github.com/public-awesome/stargaze/v12/testutil/sample"
	"github.com/public-awesome/stargaze/v12/x/cron/keeper"
	"github.com/public-awesome/stargaze/v12/x/cron/types"
	"github.com/stretchr/testify/require"
)

func TestPromoteToPrivilegedContract(t *testing.T) {
	testCases := []struct {
		testCase    string
		prepare     func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgPromoteToPrivilegedContract
		expectError bool
	}{
		{
			"invalid sender address",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgPromoteToPrivilegedContract {
				msg := types.MsgPromoteToPrivilegedContract{
					Authority: "👻",
					Contract:  sample.AccAddress().String(),
				}
				return &msg
			},
			true,
		},
		{
			"sender not gov module",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgPromoteToPrivilegedContract {
				sender := sample.AccAddress()
				msg := types.MsgPromoteToPrivilegedContract{
					Authority: sender.String(),
					Contract:  sample.AccAddress().String(),
				}
				return &msg
			},
			true,
		},
		{
			"contract does not exist",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgPromoteToPrivilegedContract {
				govModuleAddr := keeper.GetAuthority()
				msg := types.MsgPromoteToPrivilegedContract{
					Authority: govModuleAddr,
					Contract:  sample.AccAddress().String(),
				}
				return &msg
			},
			true,
		},
		{
			"valid",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgPromoteToPrivilegedContract {
				govModuleAddr := keeper.GetAuthority()
				msg := types.MsgPromoteToPrivilegedContract{
					Authority: govModuleAddr,
					Contract:  "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
				}
				return &msg
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testCase, func(t *testing.T) {
			k, c := keepertest.CronKeeper(t)
			msgSrvr, ctx := keeper.NewMsgServerImpl(k), sdk.WrapSDKContext(c)

			msg := tc.prepare(c, k)

			_, err := msgSrvr.PromoteToPrivilegedContract(ctx, msg)

			if tc.expectError {
				require.Error(t, err, tc)
				isPrivileged := k.IsPrivileged(c, sdk.MustAccAddressFromBech32(msg.Contract))
				require.False(t, isPrivileged)
			} else {
				require.NoError(t, err, tc)
				isPrivileged := k.IsPrivileged(c, sdk.MustAccAddressFromBech32(msg.Contract))
				require.True(t, isPrivileged)
			}
		})
	}
}

func TestDemoteFromPrivilegedContract(t *testing.T) {
	testCases := []struct {
		testCase    string
		prepare     func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgDemoteFromPrivilegedContract
		expectError bool
	}{
		{
			"invalid sender address",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgDemoteFromPrivilegedContract {
				msg := types.MsgDemoteFromPrivilegedContract{
					Authority: "👻",
					Contract:  sample.AccAddress().String(),
				}
				return &msg
			},
			true,
		},
		{
			"sender not gov module",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgDemoteFromPrivilegedContract {
				sender := sample.AccAddress()
				msg := types.MsgDemoteFromPrivilegedContract{
					Authority: sender.String(),
					Contract:  sample.AccAddress().String(),
				}
				return &msg
			},
			true,
		},
		{
			"contract does not exist",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgDemoteFromPrivilegedContract {
				govModuleAddr := keeper.GetAuthority()
				msg := types.MsgDemoteFromPrivilegedContract{
					Authority: govModuleAddr,
					Contract:  sample.AccAddress().String(),
				}
				return &msg
			},
			true,
		},
		{
			"contract curretly does not have privilege to demote it",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgDemoteFromPrivilegedContract {
				govModuleAddr := keeper.GetAuthority()
				msg := types.MsgDemoteFromPrivilegedContract{
					Authority: govModuleAddr,
					Contract:  "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
				}

				return &msg
			},
			true,
		},
		{
			"valid",
			func(ctx sdk.Context, keeper keeper.Keeper) *types.MsgDemoteFromPrivilegedContract {
				contractAddr := "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du"
				err := keeper.SetPrivileged(ctx, sdk.MustAccAddressFromBech32(contractAddr))
				require.NoError(t, err)

				govModuleAddr := keeper.GetAuthority()
				msg := types.MsgDemoteFromPrivilegedContract{
					Authority: govModuleAddr,
					Contract:  contractAddr,
				}
				return &msg
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testCase, func(t *testing.T) {
			k, c := keepertest.CronKeeper(t)
			msgSrvr, ctx := keeper.NewMsgServerImpl(k), sdk.WrapSDKContext(c)

			msg := tc.prepare(c, k)

			_, err := msgSrvr.DemoteFromPrivilegedContract(ctx, msg)

			if tc.expectError {
				require.Error(t, err, tc)
			} else {
				require.NoError(t, err, tc)
				isPrivileged := k.IsPrivileged(c, sdk.MustAccAddressFromBech32(msg.Contract))
				require.False(t, isPrivileged)
			}
		})
	}
}
