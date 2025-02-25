package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/telepresenceio/telepresence/rpc/v2/manager"
	"github.com/telepresenceio/telepresence/v2/pkg/client/cli/cliutil"
)

// addPreviewFlags mutates 'flags', adding flags to it such that the flags set the appropriate
// fields in the given 'spec'.  If 'prefix' is given, long-flag names are prefixed with it.
func addPreviewFlags(prefix string, flags *pflag.FlagSet, spec *manager.PreviewSpec) {
	flags.BoolVarP(&spec.DisplayBanner, prefix+"banner", "b", true, "Display banner on preview page")
}

func previewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "preview",
		Args: OnlySubcommands,

		Short: "Create or remove preview domains for existing intercepts",
		RunE:  RunSubcommands,
	}

	var createSpec manager.PreviewSpec
	createCmd := &cobra.Command{
		Use:  "create [flags] <intercept_name>",
		Args: cobra.ExactArgs(1),

		Short: "Create a preview domain for an existing intercept",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := cliutil.EnsureLoggedIn(cmd.Context()); err != nil {
				return err
			}
			si := &sessionInfo{cmd: cmd}
			return si.withConnector(true, func(cs *connectorState) error {
				if createSpec.Ingress == nil {
					ingress, err := cs.selectIngress(cmd.Context(), cmd.InOrStdin(), cmd.OutOrStdout())
					if err != nil {
						return err
					}
					createSpec.Ingress = ingress
				}
				intercept, err := cs.managerClient.UpdateIntercept(cmd.Context(), &manager.UpdateInterceptRequest{
					Session: cs.info.SessionInfo,
					Name:    args[0],
					PreviewDomainAction: &manager.UpdateInterceptRequest_AddPreviewDomain{
						AddPreviewDomain: &createSpec,
					},
				})
				if err != nil {
					return err
				}
				fmt.Println(DescribeIntercept(intercept, nil, false))
				return nil
			})
		},
	}
	addPreviewFlags("", createCmd.Flags(), &createSpec)

	removeCmd := &cobra.Command{
		Use:  "remove <intercept_name>",
		Args: cobra.ExactArgs(1),

		Short: "Remove a preview domain from an intercept",
		RunE: func(cmd *cobra.Command, args []string) error {
			si := &sessionInfo{cmd: cmd}
			return si.withConnector(true, func(cs *connectorState) error {
				intercept, err := cs.managerClient.UpdateIntercept(cmd.Context(), &manager.UpdateInterceptRequest{
					Session: cs.info.SessionInfo,
					Name:    args[0],
					PreviewDomainAction: &manager.UpdateInterceptRequest_RemovePreviewDomain{
						RemovePreviewDomain: true,
					},
				})
				if err != nil {
					return err
				}
				fmt.Println(DescribeIntercept(intercept, nil, false))
				return nil
			})
		},
	}

	cmd.AddCommand(createCmd, removeCmd)

	return cmd
}
