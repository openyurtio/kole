/*
Copyright 2022 The OpenYurt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/kole-controller/app/options"
	apiv1alpha1 "github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/controller"
	"github.com/openyurtio/kole/pkg/projectinfo"
)

const (
	// component name
	componentController = "kole-controller"
)

var (
	scheme = apiruntime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = apiv1alpha1.AddToScheme(clientgoscheme.Scheme)
	_ = apiv1alpha1.AddToScheme(scheme)
}

// NewKoleControllerCommand creates a *cobra.Command object with default parameters
func NewKoleControllerCommand() *cobra.Command {
	cleanFlagSet := pflag.NewFlagSet(componentController, pflag.ContinueOnError)
	koleControllerFlags := options.NewKoleControllerFlags()

	cmd := &cobra.Command{

		Use: componentController,
		Long: `The kole-controller is the cloud controller.  
`,
		// The kole-controller has special flag parsing requirements to enforce flag precedence rules,
		// so we do all our parsing manually in Run, below.
		// DisableFlagParsing=true provides the full set of flags passed to the lite kubelet in the
		// `args` arg to Run, without Cobra's interference.
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			// initial flag parse, since we disable cobra's flag parsing
			if err := cleanFlagSet.Parse(args); err != nil {
				cmd.Usage()
				klog.Fatal(err)
			}

			// check if there are non-flag arguments in the command line
			cmds := cleanFlagSet.Args()
			if len(cmds) > 0 {
				cmd.Usage()
				klog.Fatalf("unknown command: %s", cmds[0])
			}

			// short-circuit on help
			help, err := cleanFlagSet.GetBool("help")
			if err != nil {
				klog.Fatal(`"help" flag is non-bool, programmer error, please correct`)
			}
			if help {
				cmd.Help()
				return
			}

			if err := options.SetKoleControllerFlagsByEnv(koleControllerFlags); err != nil {
				klog.Fatal(err)
			}

			// validate the initial KubeletFlags
			if err := options.ValidateKoleControllerFlags(koleControllerFlags); err != nil {
				klog.Fatal(err)
			}

			shutdownHandler := make(chan os.Signal, 2)
			var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
			ctx, cancel := context.WithCancel(context.Background())
			signal.Notify(shutdownHandler, shutdownSignals...)
			go func() {
				<-shutdownHandler
				cancel()
				<-shutdownHandler
				os.Exit(1) // second signal. Exit directly.
			}()

			// run the kubelet
			klog.Infof("kole-controller Config: %#v", *koleControllerFlags)
			klog.Infof("Version:%#v", projectinfo.Get())

			if err := Run(ctx, koleControllerFlags); err != nil {
				klog.Fatal(err)
			}
		},
	}

	// keep cleanFlagSet separate, so Cobra doesn't pollute it with the global flags
	koleControllerFlags.AddFlags(cleanFlagSet)
	// DELETE BY zhangjie
	options.AddGlobalFlags(cleanFlagSet)
	cleanFlagSet.BoolP("help", "h", false, fmt.Sprintf("help for %s", cmd.Name()))

	// ugly, but necessary, because Cobra's default UsageFunc and HelpFunc pollute the flagset with global flags
	const usageFmt = "Usage:\n  %s\n\nFlags:\n%s"
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine(), cleanFlagSet.FlagUsagesWrapped(2))
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine(), cleanFlagSet.FlagUsagesWrapped(2))
	})

	return cmd
}

func Run(ctx context.Context, config *options.KoleControllerFlags) error {
	defer runtime.HandleCrash()

	stop := make(chan struct{})
	defer close(stop)

	if err := RunKoleController(stop, config); err != nil {
		return err
	}
	// must exec after SubscribeTopics
	select {
	case <-stop:
		break
	case <-ctx.Done():
		break
	}
	return nil
}

func RunKoleController(stop chan struct{}, config *options.KoleControllerFlags) error {

	lite, err := controller.NewMainKoleController(stop,
		config,
		&controller.Gzip{})
	if err != nil {
		klog.Errorf("NewMainKoleController error %v", err)
		return err
	}
	return lite.Run()
}
