package cmd

import (
	"encoding/base64"
	"net"
	"net/url"
	"time"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure chkit default values",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := chlib.GetUserInfo()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		httpApi, err := chlib.GetHttpApiCfg()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		tcpApi, err := chlib.GetTcpApiConfig()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}

		if cmd.Flags().NFlag() == 0 {
			jww.FEEDBACK.Println("Token: ", info.Token)
			jww.FEEDBACK.Println("Namespace: ", info.Namespace)
			jww.FEEDBACK.Println("HTTP API")
			jww.FEEDBACK.Println("\tServer: ", httpApi.Server)
			jww.FEEDBACK.Println("\tTimeout: ", httpApi.Timeout)
			jww.FEEDBACK.Println("TCP API")
			jww.FEEDBACK.Printf("\tServer: %s:%d", tcpApi.Address, tcpApi.Port)
			jww.FEEDBACK.Println("\tBuffer size: ", tcpApi.BufferSize)
			return
		}

		if cmd.Flag("set-default-namespace").Changed {
			info.Namespace = cmd.Flag("set-default-namespace").Value.String()
			jww.FEEDBACK.Printf("Namespace changed to: %s\n", info.Namespace)
		}
		if cmd.Flag("set-token").Changed {
			enteredToken := cmd.Flag("set-token").Value.String()
			if _, err := base64.StdEncoding.DecodeString(enteredToken); err != nil {
				jww.FEEDBACK.Println("Invalid token given")
				return
			}
			info.Token = enteredToken
			jww.FEEDBACK.Printf("Token changed to: %s\n", info.Token)
		}
		if cmd.Flag("set-http-server-address").Changed {
			address := cmd.Flag("set-http-server-address").Value.String()
			if _, err := url.ParseRequestURI(address); err != nil {
				jww.FEEDBACK.Printf("Invalid HTTP API server address given")
				return
			}
			httpApi.Server = address
			jww.FEEDBACK.Printf("HTTP API server address changed to: %s", address)
		}
		if cmd.Flag("set-http-server-timeout").Changed {
			tm, err := cmd.Flags().GetDuration("set-http-server-timeout")
			if err != nil {
				jww.FEEDBACK.Printf("Invalid HTTP API timeout given")
				return
			}
			httpApi.Timeout = tm
			jww.FEEDBACK.Printf("HTTP API timeout changed to: %s", tm)
		}
		if cmd.Flag("set-tcp-server-address").Changed {
			ip, err := cmd.Flags().GetIP("set-tcp-server-address")
			if err != nil {
				jww.FEEDBACK.Println("Invalid IP address given")
				return
			}
			tcpApi.Address = ip
			jww.FEEDBACK.Printf("TCP API server address changed to: %s", ip)
		}
		if cmd.Flag("set-tcp-server-port").Changed {
			port, err := cmd.Flags().GetInt("set-tcp-server-port")
			if err != nil || port < 0 || port > 65535 {
				jww.FEEDBACK.Println("Invalid port number given")
				return
			}
			tcpApi.Port = port
			jww.FEEDBACK.Printf("TCP API server port changed to: %d", port)
		}
		if cmd.Flag("set-tcp-buffer-size").Changed {
			bufsz, err := cmd.Flags().GetInt("set-tcp-buffer-size")
			if err != nil || bufsz < 0 {
				jww.FEEDBACK.Println("Invalid buffer size given")
				return
			}
			tcpApi.BufferSize = bufsz
			jww.FEEDBACK.Println("TCP API buffer size changed to: %d", bufsz)
		}

		err = chlib.UpdateUserInfo(info)
		if err != nil {
			jww.ERROR.Println(err)
		}
		err = chlib.UpdateHttpApiCfg(httpApi)
		if err != nil {
			jww.ERROR.Println(err)
		}
		err = chlib.UpdateTcpApiConfig(tcpApi)
		if err != nil {
			jww.ERROR.Println(err)
		}
	},
}

func init() {
	configCmd.PersistentFlags().StringP("set-token", "t", "", "Set user token")
	configCmd.PersistentFlags().StringP("set-default-namespace", "n", chlib.DefaultNameSpace, "Default namespace")
	configCmd.PersistentFlags().String("set-http-server-address", "http://0.0.0.0:3333", "HTTP API server address")
	configCmd.PersistentFlags().Duration("set-http-server-timeout", 10*time.Second, "HTTP API calls timeout")
	configCmd.PersistentFlags().IP("set-tcp-server-address", net.IPv6zero, "TCP API server address")
	configCmd.PersistentFlags().Int("set-tcp-server-port", 3000, "TCP API server port")
	configCmd.PersistentFlags().Int("set-tcp-buffer-size", 1024, "TCP API buffer size")
	RootCmd.AddCommand(configCmd)
}
