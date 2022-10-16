/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"

	"github.com/meln5674/kink/pkg/command"
	"github.com/meln5674/kink/pkg/helm"

	"github.com/spf13/cobra"
)

// deleteClusterCmd represents the deleteCluster command
var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := func() error {

			ctx := context.TODO()

			log.Println("Deleting release...")
			// TODO: Add flag to also delete PVCs

			helmDelete := helm.Delete(&helmFlags, &chartFlags, &releaseFlags, &kubeFlags)
			err := command.
				Command(ctx, helmDelete...).
				ForwardOutErr().
				Run()
			if err != nil {
				return err
			}
			log.Println("Cluster deleted")
			return nil
		}()
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteClusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteClusterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
