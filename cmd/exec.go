/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/meln5674/kink/pkg/command"
	"github.com/meln5674/kink/pkg/kubectl"
	"github.com/spf13/cobra"
)

var (
	execCommand []string
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a new process with access to a cluster",
	Long: `Because kink clusters are contained within another cluster, their controlplane may not
be accessible from where you are running kink, unless you have made extra provisions such as Ingress
or a LoadBalancer Service.

To work around this, kink can use Kubernetes port-forwarding to provide
access to that controlplane. This command sets up that port forwarding, sets the KUBECONFIG variable
to a temporary file that will connect to it, and executes your shell command, then stop forwarding and
clean up the temporary kubeconfig once it has exited.

This command does not perform variable replacements or glob expansions. To do this, use 'kink sh'`,
	Run: func(cmd *cobra.Command, args []string) {
		ec, err := func() (*int, error) {
			ctx := context.TODO()

			kubeconfig, err := os.CreateTemp("", "kink-kubeconfig-*")
			defer kubeconfig.Close()
			defer os.Remove(kubeconfig.Name())
			if err != nil {
				return nil, err
			}
			kubectlCp := kubectl.Cp(&kubectlFlags, &kubeFlags, releaseFlags.Namespace, "kink-kink-controlplane-0", "/etc/rancher/k3s/k3s.yaml", kubeconfig.Name())
			err = command.Command(ctx, kubectlCp...).ForwardOutErr().Run()
			if err != nil {
				return nil, err
			}

			// TODO: Get service name/remote port from chart (helm get manifest)
			// TODO: Make local port configurable with flag
			kubectlPortForward := kubectl.PortForward(&kubectlFlags, &kubeFlags, releaseFlags.Namespace, "svc/kink-kink-controlplane", map[string]string{"6443": "6443"})
			kubectlPortForwardCmd := command.
				Command(ctx, kubectlPortForward...).
				ForwardOutErr()

			err = kubectlPortForwardCmd.Start()

			if err != nil {
				return nil, err
			}
			defer func() {
				// Deliberately ignoing the errors here
				kubectlPortForwardCmd.Kill()
				kubectlPortForwardCmd.Wait()
			}()

			log.Println("Waiting for cluster to be accessible on localhost...")
			kubectlVersion := kubectl.Version(&kubectlFlags, &kubeFlags)
			for err = errors.New("dummy"); err != nil; err = command.Command(ctx, kubectlVersion...).ForwardOutErr().Run() {
				time.Sleep(5 * time.Second)
			}

			// TODO: Make shell configurable with flag
			err = command.
				Command(ctx, args...).
				WithParentEnv().
				WithEnv(map[string]string{
					"KUBECONFIG": kubeconfig.Name(),
				}).
				ForwardAll().
				Run()
			var exitError *exec.ExitError
			if errors.As(err, &exitError) {
				ec := exitError.ProcessState.ExitCode()
				return &ec, nil
			}
			if err != nil {
				return nil, err
			}
			return nil, nil
		}()
		if err != nil {
			log.Fatal(err)
		}
		if ec != nil {
			os.Exit(*ec)
		}

	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
