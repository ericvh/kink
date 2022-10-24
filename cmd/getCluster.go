/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"k8s.io/klog/v2"

	"github.com/meln5674/gosh"
	"github.com/meln5674/kink/pkg/helm"

	"github.com/spf13/cobra"
)

// getClusterCmd represents the getCluster command
var getClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Lists existing kink clusters by their name",
	Run: func(cmd *cobra.Command, args []string) {
		err := func() error {

			ctx := context.TODO()

			releases := make([]map[string]interface{}, 0)

			helmList := helm.List(&helmFlags, &chartFlags, &releaseFlags, &kubeFlags)
			err := gosh.
				Command(helmList...).
				WithContext(ctx).
				WithStreams(
					gosh.ForwardErr,
					gosh.FuncOut(gosh.SaveJSON(&releases)),
				).
				Run()
			if err != nil {
				return err
			}

			for _, release := range releases {
				if !helm.IsKinkRelease(release["name"].(string)) {
					continue
				}
				if releaseFlags.Namespace != "" {
					fmt.Printf("%s %s\n", release["namespace"], release["name"])
				} else {
					fmt.Println(release["name"])
				}
			}

			return nil
		}()
		if err != nil {
			klog.Fatal(err)
		}

	},
}

func init() {
	getCmd.AddCommand(getClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getClusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getClusterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
