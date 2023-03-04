package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	
	RootCmd = cobra.Command{
		Use: "im-chatgpt",
	}
)
