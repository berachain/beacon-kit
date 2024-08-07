// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package config

import (
	"strings"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// NewPrefixedViper creates a new viper instance with the given environment
// prefix, and replaces all (.) and (-) with (_).
func NewPrefixedViper(prefix string) *viper.Viper {
	viper := viper.New()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	return viper
}

// BindFlags binds the given cobra command's flags to the given viper instance.
// It also sets the viper instance to automatically read from the environment.
func BindFlags(
	executableName string,
	cmd *cobra.Command,
	v *viper.Viper,
) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.Newf("bindFlags failed: %v", r)
		}
	}()

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flagName := strings.ToUpper(strings.ReplaceAll(flag.Name, "-", "_"))
		if err = v.BindEnv(
			flag.Name,
			strings.Join([]string{executableName, flagName}, "_"),
		); err != nil {
			panic(err)
		}
		if err = v.BindPFlag(flag.Name, flag); err != nil {
			panic(err)
		}

		if !flag.Changed && v.IsSet(flag.Name) {
			if err = cmd.Flags().Set(
				flag.Name,
				cast.ToString(v.Get(flag.Name)),
			); err != nil {
				panic(err)
			}
		}
	})

	return err
}
