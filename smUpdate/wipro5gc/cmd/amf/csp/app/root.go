package app

import "github.com/spf13/cobra"

func NewCspRootCommand() *cobra.Command {

	var cfgFile, etcdServer, etcdConfigKey string //, requestType string
	var reset bool

	rootCmd := &cobra.Command{
		Use:   "csp",
		Short: "csp is an AMF microservice in 5GC to handle pdu sms signalling",
		Long: `csp provides the session management services in the amf
					Complete documentation is available at http://wipro.com/5gc`,
		Run: func(cmd *cobra.Command, args []string) {
			RunCsp(cfgFile, etcdServer, etcdConfigKey, reset)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is w5gc.io/wipro5gcore/configs/csp)")
	rootCmd.PersistentFlags().StringVar(&etcdServer, "etcd-server", "", "etcd server to read config file (default is /w5gc/config/csp.json)")
	rootCmd.PersistentFlags().StringVar(&etcdConfigKey, "etcd-config-key", "", "etcd server key for config file (default is /w5gc/config/csp.json)")
	rootCmd.PersistentFlags().BoolVarP(&reset, "reset", "r", false, "reset flag")
	rootCmd.PersistentFlags().StringP("author", "a", "Wipro", "author name for copyright attribution")
	//to be removed later
	// rootCmd.PersistentFlags().StringVar(&requestType, "request-type", "", "type of SM context request : create, update, release, retrieve")
	versionCmd := NewCspVersionCommand()

	rootCmd.AddCommand(versionCmd)

	return rootCmd
}
