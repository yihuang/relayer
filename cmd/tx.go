/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/relayer/relayer"
	"github.com/spf13/cobra"
)

// transactionCmd represents the tx command
func transactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transact",
		Aliases: []string{"tx"},
		Short:   "IBC Transaction Commands",
		Long: strings.TrimSpace(`Commands to create IBC transactions on configured chains. 
		Most of these commands take a '[path]' argument. Make sure:
	1. Chains are properly configured to relay over by using the 'rly chains list' command
	2. Path is properly configured to relay over by using the 'rly paths list' command`),
	}

	cmd.AddCommand(
		linkCmd(),
		linkThenStartCmd(),
		relayMsgsCmd(),
		relayAcksCmd(),
		xfersend(),
		flags.LineBreak,
		createClientsCmd(),
		updateClientsCmd(),
		createConnectionCmd(),
		createChannelCmd(),
		closeChannelCmd(),
		flags.LineBreak,
		rawTransactionCmd(),
	)

	return cmd
}

func createClientsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clients [path-name]",
		Aliases: []string{"clnts"},
		Short:   "create a clients between two configured chains with a configured path",
		Long: "Creates a working ibc client for chain configured on each end of the" +
			" path by querying headers from each chain and then sending the corresponding create-client messages",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			// ensure that keys exist
			if _, err = c[src].GetAddress(); err != nil {
				return err
			}
			if _, err = c[dst].GetAddress(); err != nil {
				return err
			}

			return c[src].CreateClients(c[dst])
		},
	}
	return cmd
}

func updateClientsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-clients [path-name]",
		Aliases: []string{"update", "uc"},
		Short:   "update a clients between two configured chains with a configured path",
		Long: "Updates a working ibc client for chain configured on each end of the " +
			"path by querying headers from each chain and then sending the corresponding update-client messages",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			// ensure that keys exist
			if _, err = c[src].GetAddress(); err != nil {
				return err
			}
			if _, err = c[dst].GetAddress(); err != nil {
				return err
			}

			return c[src].UpdateClients(c[dst])
		},
	}
	return cmd
}

func createConnectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connection [path-name]",
		Aliases: []string{"conn", "con"},
		Short:   "create a connection between two configured chains with a configured path",
		Long: strings.TrimSpace(`This command is meant to be used to repair or create 
		a connection between two chains with a configured path in the config file`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			to, err := getTimeout(cmd)
			if err != nil {
				return err
			}

			// ensure that keys exist
			if _, err = c[src].GetAddress(); err != nil {
				return err
			}
			if _, err = c[dst].GetAddress(); err != nil {
				return err
			}

			return c[src].CreateConnection(c[dst], to)
		},
	}

	return timeoutFlag(cmd)
}

func createChannelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "channel [path-name]",
		Aliases: []string{"chan", "ch"},
		Short:   "create a channel between two configured chains with a configured path",
		Long: strings.TrimSpace(`This command is meant to be used to repair or 
		create a channel between two chains with a configured path in the config file`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			to, err := getTimeout(cmd)
			if err != nil {
				return err
			}

			// ensure that keys exist
			if _, err = c[src].GetAddress(); err != nil {
				return err
			}
			if _, err = c[dst].GetAddress(); err != nil {
				return err
			}

			return c[src].CreateChannel(c[dst], false, to)
		},
	}

	return timeoutFlag(cmd)
}

func closeChannelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "channel-close [path-name]",
		Aliases: []string{"chan-cl", "close", "cl"},
		Short:   "close a channel between two configured chains with a configured path",
		Long:    "This command is meant to close a channel",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			to, err := getTimeout(cmd)
			if err != nil {
				return err
			}

			// ensure that keys exist
			if _, err = c[src].GetAddress(); err != nil {
				return err
			}
			if _, err = c[dst].GetAddress(); err != nil {
				return err
			}

			return c[src].CloseChannel(c[dst], to)
		},
	}

	return timeoutFlag(cmd)
}

func linkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "link [path-name]",
		Aliases: []string{"full-path", "connect", "path", "pth"},
		Short:   "create clients, connection, and channel between two configured chains with a configured path",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			to, err := getTimeout(cmd)
			if err != nil {
				return err
			}

			// ensure that keys exist
			if _, err = c[src].GetAddress(); err != nil {
				return err
			}
			if _, err = c[dst].GetAddress(); err != nil {
				return err
			}

			// create clients if they aren't already created
			if err = c[src].CreateClients(c[dst]); err != nil {
				return err
			}

			// create connection if it isn't already created
			if err = c[src].CreateConnection(c[dst], to); err != nil {
				return err
			}

			// create channel if it isn't already created
			return c[src].CreateChannel(c[dst], true, to)
		},
	}

	return timeoutFlag(cmd)
}

func linkThenStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link-then-start [path-name]",
		Short: "wait for a link to come up, then start relaying packets",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			lCmd := linkCmd()
			for err := lCmd.RunE(cmd, args); err != nil; err = lCmd.RunE(cmd, args) {
				fmt.Printf("retrying link: %s\n", err)
				time.Sleep(1 * time.Second)
			}
			sCmd := startCmd()
			return sCmd.RunE(cmd, args)
		},
	}
	return strategyFlag(timeoutFlag(cmd))
}

func relayMsgsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "relay-packets [path-name]",
		Aliases: []string{"rly", "pkts", "relay"},
		Short:   "relay any packets that remain to be relayed on a given path, in both directions",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			if err = ensureKeysExist(c); err != nil {
				return err
			}

			sh, err := relayer.NewSyncHeaders(c[src], c[dst])
			if err != nil {
				return err
			}

			strategy, err := GetStrategyWithOptions(cmd, config.Paths.MustGet(args[0]).MustGetStrategy())
			if err != nil {
				return err
			}

			sp, err := strategy.UnrelayedSequences(c[src], c[dst], sh)
			if err != nil {
				return err
			}

			if err = strategy.RelayPackets(c[src], c[dst], sp, sh); err != nil {
				return err
			}

			return nil
		},
	}

	return strategyFlag(cmd)
}

func relayAcksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "relay-acknowledgements [path-name]",
		Aliases: []string{"acks"},
		Short:   "relay any acknowledgements that remain to be relayed on a given path, in both directions",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			if err = ensureKeysExist(c); err != nil {
				return err
			}

			sh, err := relayer.NewSyncHeaders(c[src], c[dst])
			if err != nil {
				return err
			}

			strategy, err := GetStrategyWithOptions(cmd, config.Paths.MustGet(args[0]).MustGetStrategy())
			if err != nil {
				return err
			}

			// sp.Src contains all sequences acked on SRC but acknowledgement not processed on DST
			// sp.Dst contains all sequences acked on DST but acknowledgement not processed on SRC
			sp, err := strategy.UnrelayedAcknowledgements(c[src], c[dst], sh)
			if err != nil {
				return err
			}

			if err = strategy.RelayAcknowledgements(c[src], c[dst], sp, sh); err != nil {
				return err
			}

			return nil
		},
	}

	return strategyFlag(cmd)
}

// Returns an error if a configured key for a given chain doesn't exist
func ensureKeysExist(chains map[string]*relayer.Chain) (err error) {
	for _, v := range chains {
		if _, err = v.GetAddress(); err != nil {
			return
		}
	}
	return nil
}
