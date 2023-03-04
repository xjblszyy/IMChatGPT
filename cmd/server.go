package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xjblszyy/im-chatgpt/config"
	"github.com/xjblszyy/im-chatgpt/ims"
)

func serverCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			setups := []func() error{
				setupLogger,
				setupGPT,
			}
			
			for _, setup := range setups {
				if err := setup(); err != nil {
					panic(err)
				}
			}
			imBots := ims.NewBots(*config.C)
			imBots.Start()
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			return nil
		},
	}
	
	c.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	return c
}

func init() {
	RootCmd.AddCommand(serverCmd())
}
