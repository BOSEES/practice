package feeshare_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/suite"

	"github.com/BOSEES/practice/v12/app"
	"github.com/BOSEES/practice/v12/x/feeshare"
	"github.com/BOSEES/practice/v12/x/feeshare/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type GenesisTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app     *app.App
	genesis types.GenesisState
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	app := Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.ctx = ctx

	suite.genesis = *types.DefaultGenesisState()
}

func (suite *GenesisTestSuite) TestFeeShareInitGenesis() {
	testCases := []struct {
		name     string
		genesis  types.GenesisState
		expPanic bool
	}{
		{
			"default genesis",
			suite.genesis,
			false,
		},
		{
			"custom genesis - feeshare disabled",
			types.GenesisState{
				Params: types.Params{
					EnableFeeShare:  false,
					DeveloperShares: types.DefaultDeveloperShares,
					AllowedDenoms:   []string{"upractice"},
				},
			},
			false,
		},
		{
			"custom genesis - feeshare enabled, 0% developer shares",
			types.GenesisState{
				Params: types.Params{
					EnableFeeShare:  true,
					DeveloperShares: sdk.NewDecWithPrec(0, 2),
					AllowedDenoms:   []string{"upractice"},
				},
			},
			false,
		},
		{
			"custom genesis - feeshare enabled, 100% developer shares",
			types.GenesisState{
				Params: types.Params{
					EnableFeeShare:  true,
					DeveloperShares: sdk.NewDecWithPrec(100, 2),
					AllowedDenoms:   []string{"upractice"},
				},
			},
			false,
		},
		{
			"custom genesis - feeshare enabled, all denoms allowed",
			types.GenesisState{
				Params: types.Params{
					EnableFeeShare:  true,
					DeveloperShares: sdk.NewDecWithPrec(10, 2),
					AllowedDenoms:   []string(nil),
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			if tc.expPanic {
				suite.Require().Panics(func() {
					feeshare.InitGenesis(suite.ctx, suite.app.FeeShareKeeper, tc.genesis)
				})
			} else {
				suite.Require().NotPanics(func() {
					feeshare.InitGenesis(suite.ctx, suite.app.FeeShareKeeper, tc.genesis)
				})

				params := suite.app.FeeShareKeeper.GetParams(suite.ctx)
				suite.Require().Equal(tc.genesis.Params, params)
			}
		})
	}
}
