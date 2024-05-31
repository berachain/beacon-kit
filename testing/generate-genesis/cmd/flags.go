package cmd

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/spf13/cobra"
)

const (
	outputFileFlag         = "output"
	predeployAddressesFlag = "predeployAddresses"
	predeployCodesFlag     = "predeployCodes"
	predeployBalancesFlag  = "predeployBalances"
	predeployNoncesFlag    = "predeployNonces"
	accountAddressesFlag   = "accountAddresses"
	accountBalancesFlag    = "accountBalances"
)

// Returns slices of strings for each predeploy flag in order of:
// accountAddresses, accountBalances, predeployAddresses, predeployCodes, predeployBalances, predeployNonces, outputFile
// TODO: maybe unhood this idk does it really matter?
func sanitizeFlags(cmd *cobra.Command) (
	[]string, []string, []string, []string, []string, []string, string, error,
) {
	// Check if all predeploy flags have the same length
	predeployAddresses, err := cmd.Flags().GetStringSlice(predeployAddressesFlag)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, "", errors.Wrap(err, "failed to get predeployAddresses flag")
	}
	predeployCodes, err := cmd.Flags().GetStringSlice(predeployCodesFlag)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, "", errors.Wrap(err, "failed to get predeployCodes flag")
	}
	predeployBalances, err := cmd.Flags().GetStringSlice(predeployBalancesFlag)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, "", errors.Wrap(err, "failed to get predeployBalances flag")
	}
	predeployNonces, err := cmd.Flags().GetStringSlice(predeployNoncesFlag)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, "", errors.Wrap(err, "failed to get predeployNonces flag")
	}
	if len(predeployAddresses) != len(predeployCodes) ||
		len(predeployAddresses) != len(predeployBalances) ||
		len(predeployAddresses) != len(predeployNonces) ||
		len(predeployCodes) != len(predeployNonces) {
		return nil, nil, nil, nil, nil, nil, "", errPredeployFlagsLength
	}

	// Check if accountAddresses and accountBalances have the same length
	accountAddresses, err := cmd.Flags().GetStringSlice(accountAddressesFlag)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, "", errors.Wrap(err, "failed to get accountAddresses flag")
	}
	accountBalances, err := cmd.Flags().GetStringSlice(accountBalancesFlag)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, "", errors.Wrap(err, "failed to get accountBalances flag")
	}
	if len(accountAddresses) != len(accountBalances) {
		return nil, nil, nil, nil, nil, nil, "", errAccountFlagsLength
	}

	// Get the output file name
	outputFile, err := cmd.Flags().GetString(outputFileFlag)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, "", errors.Wrap(err, "failed to get output flag")
	}

	return accountAddresses,
		accountBalances,
		predeployAddresses,
		predeployCodes,
		predeployBalances,
		predeployNonces,
		outputFile,
		nil
}
