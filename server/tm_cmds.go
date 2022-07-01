package server

// DONTCOVER

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	tcmd "github.com/Fincor-Blockchain/pou/cmd/pou/commands"
	"github.com/Fincor-Blockchain/pou/libs/cli"
	"github.com/Fincor-Blockchain/pou/p2p"
	pvm "github.com/Fincor-Blockchain/pou/privval"
	tversion "github.com/Fincor-Blockchain/pou/version"

	"github.com/Fincor-Blockchain/cusp-sdk/codec"
	sdk "github.com/Fincor-Blockchain/cusp-sdk/types"
)

// ShowNodeIDCmd - ported from pou, dump node ID to stdout
func ShowNodeIDCmd(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			fmt.Println(nodeKey.ID())
			return nil
		},
	}
}

// ShowValidator - ported from pou, show this node's validator info
func ShowValidatorCmd(ctx *Context) *cobra.Command {
	cmd := cobra.Command{
		Use:   "show-validator",
		Short: "Show this node's pou validator info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg := ctx.Config
			UpgradeOldPrivValFile(cfg)
			privValidator := pvm.LoadOrGenFilePV(
				cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
			valPubKey := privValidator.GetPubKey()

			if viper.GetString(cli.OutputFlag) == "json" {
				return printlnJSON(valPubKey)
			}

			pubkey, err := sdk.Bech32ifyConsPub(valPubKey)
			if err != nil {
				return err
			}

			fmt.Println(pubkey)
			return nil
		},
	}

	cmd.Flags().StringP(cli.OutputFlag, "o", "text", "Output format (text|json)")
	return &cmd
}

// ShowAddressCmd - show this node's validator address
func ShowAddressCmd(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address",
		Short: "Shows this node's pou validator consensus address",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg := ctx.Config
			UpgradeOldPrivValFile(cfg)
			privValidator := pvm.LoadOrGenFilePV(
				cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
			valConsAddr := (sdk.ConsAddress)(privValidator.GetAddress())

			if viper.GetString(cli.OutputFlag) == "json" {
				return printlnJSON(valConsAddr)
			}

			fmt.Println(valConsAddr.String())
			return nil
		},
	}

	cmd.Flags().StringP(cli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

// VersionCmd prints pou and ABCI version numbers.
func VersionCmd(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print pou libraries' version",
		Long: `Print protocols' and libraries' version numbers
against which this app has been compiled.
`,
		RunE: func(cmd *cobra.Command, args []string) error {

			bs, err := yaml.Marshal(&struct {
				pou    string
				ABCI          string
				BlockProtocol uint64
				P2PProtocol   uint64
			}{
				pou:    tversion.Version,
				ABCI:          tversion.ABCIVersion,
				BlockProtocol: tversion.BlockProtocol.Uint64(),
				P2PProtocol:   tversion.P2PProtocol.Uint64(),
			})
			if err != nil {
				return err
			}

			fmt.Println(string(bs))
			return nil
		},
	}
	return cmd
}

func printlnJSON(v interface{}) error {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	marshalled, err := cdc.MarshalJSON(v)
	if err != nil {
		return err
	}
	fmt.Println(string(marshalled))
	return nil
}

// UnsafeResetAllCmd - extension of the pou command, resets initialization
func UnsafeResetAllCmd(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-reset-all",
		Short: "Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			tcmd.ResetAll(cfg.DBDir(), cfg.P2P.AddrBookFile(), cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile(), ctx.Logger)
			return nil
		},
	}
}
