// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package config_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sgconfig "github.com/itsdevbear/bolaris/config"
)

func TestConfiguration(t *testing.T) {
	t.Run("should set CoinType", func(t *testing.T) {
		config := sdk.GetConfig()

		if int(config.GetCoinType()) != sdk.CoinType {
			t.Errorf("expected CoinType %d, got %d", sdk.CoinType, config.GetCoinType())
		}
		if config.GetFullBIP44Path() != sdk.FullFundraiserPath {
			t.Errorf("expected FullBIP44Path %s, got %s",
				sdk.FullFundraiserPath, config.GetFullBIP44Path())
		}

		sgconfig.SetupCosmosConfig()

		if int(config.GetCoinType()) != 60 {
			t.Errorf("expected CoinType %d, got %d", 60, config.GetCoinType())
		}
		if config.GetCoinType() != sdk.GetConfig().GetCoinType() {
			t.Errorf("expected CoinType %d, got %d",
				sdk.GetConfig().GetCoinType(), config.GetCoinType())
		}
		if config.GetFullBIP44Path() != sdk.GetConfig().GetFullBIP44Path() {
			t.Errorf("expected FullBIP44Path %s, got %s",
				sdk.GetConfig().GetFullBIP44Path(), config.GetFullBIP44Path())
		}
	})

	t.Run("should generate HD path", func(t *testing.T) {
		params := *hd.NewFundraiserParams(0, 60, 0)
		hdPath := params.String()

		if hdPath != "m/44'/60'/0'/0/0" {
			t.Errorf("expected HD path %s, got %s", "m/44'/60'/0'/0/0", hdPath)
		}
	})
}
