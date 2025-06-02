package app

import (
  "fmt"

  "github.com/spf13/cobra"
)


func NewPdusmspVersionCommand() *cobra.Command {

        versionCmd := &cobra.Command{
                Use:   "version",
                Short: "Print the version number of PDU SMS",
                Long: `Print the version number of PDU SMS
                        Number before decimal point indicates smajor version.
                        Number after decimal point indicates minor version`,
                Run: func(cmd *cobra.Command, args []string) {
                        fmt.Println("PDU SMS v0.1 -- HEAD")
                },
        }

    return versionCmd
}

