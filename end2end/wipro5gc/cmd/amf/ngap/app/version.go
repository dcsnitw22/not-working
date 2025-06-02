package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewNgapVersionCommand() *cobra.Command {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of NGAP",
		Long: `Print the version number of NGAP
                        Number before decimal point indicates smajor version.
                        Number after decimal point indicates minor version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("NGAP v0.1 -- HEAD")
		},
	}

	return versionCmd
}
