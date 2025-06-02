package app

import "github.com/spf13/cobra"

func NewNgapRootCommand() *cobra.Command {

	var cfgFile, etcdServer, etcdConfigKey, requestType string
	var reset bool

	rootCmd := &cobra.Command{
		Use:   "ngap",
		Short: "ngap is an AMF microservice in 5GC to handle connections with gNB",
		Long: `ngap provides connection services to gNB in the amf
					Complete documentation is available at http://wipro.com/5gc`,
		Run: func(cmd *cobra.Command, args []string) {
			RunNgap(cfgFile, etcdServer, etcdConfigKey, reset, requestType)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is w5gc.io/wipro5gcore/configs/ngap)")
	rootCmd.PersistentFlags().StringVar(&etcdServer, "etcd-server", "", "etcd server to read config file (default is /w5gc/config/ngap.json)")
	rootCmd.PersistentFlags().StringVar(&etcdConfigKey, "etcd-config-key", "", "etcd server key for config file (default is /w5gc/config/ngap.json)")
	rootCmd.PersistentFlags().BoolVarP(&reset, "reset", "r", false, "reset flag")
	rootCmd.PersistentFlags().StringP("author", "a", "Wipro", "author name for copyright attribution")
	versionCmd := NewNgapVersionCommand()

	rootCmd.AddCommand(versionCmd)

	return rootCmd
}
