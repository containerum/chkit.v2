package cliserv

import (
	"fmt"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/strset"
	"github.com/spf13/cobra"
)

var aliases = []string{"srv", "services", "svc", "serv"}

func Get(ctx *context.Context) *cobra.Command {
	var getServiceConfig = struct {
		configuration.ExportConfig
	}{}
	command := &cobra.Command{
		Use:     "service",
		Aliases: aliases,
		Short:   "show service info",
		Long:    "Show service info.",
		Example: "chkit get service service_label [-o yaml/json] [-f output_file]",
		Run: func(cmd *cobra.Command, args []string) {
			serviceData, err := func() (model.Renderer, error) {
				switch len(args) {
				case 0:
					var list service.ServiceList
					var err error
					if solutionName, _ := cmd.Flags().GetString("solution_name"); solutionName != "" {
						list, err = ctx.Client.GetSolutionServices(ctx.Namespace.ID, solutionName)
					} else {
						list, err = ctx.Client.GetServiceList(ctx.Namespace.ID)
					}
					return list, err
				case 1:
					svc, err := ctx.Client.GetService(ctx.Namespace.ID, args[0])
					return svc, err
				default:
					list, err := ctx.Client.GetServiceList(ctx.Namespace.ID)
					var filteredList service.ServiceList
					names := strset.NewSet(args)
					for _, svc := range list {
						if names.Have(svc.Name) {
							filteredList = append(filteredList, svc)
						}
					}
					return filteredList, err
				}
			}()
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := configuration.ExportData(serviceData, getServiceConfig.ExportConfig); err != nil {
				fmt.Println(err)
				return
			}

		},
	}
	command.PersistentFlags().
		StringVarP((*string)(&getServiceConfig.Format), "output", "o", "", "output format [yaml/json]")
	command.PersistentFlags().
		StringVarP(&getServiceConfig.Filename, "file", "f", "-", "output file")
	command.PersistentFlags().
		StringP("solution_name", "s", "", "solution name")

	return command
}
