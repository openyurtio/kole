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
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/lite-kubelet/app/options"
	"github.com/openyurtio/kole/pkg/litekubelet"
	"github.com/openyurtio/kole/pkg/projectinfo"
)

const (
	// component name
	componentLiteKubelet = "lite-kubelet"
)

// NewLiteKubeletCommand creates a *cobra.Command object with default parameters
func NewLiteKubeletCommand() *cobra.Command {
	cleanFlagSet := pflag.NewFlagSet(componentLiteKubelet, pflag.ContinueOnError)
	liteKubeletFlags := options.NewLiteKubeletFlags()

	cmd := &cobra.Command{

		Use: componentLiteKubelet,
		Long: `The lite-kubelet is the lite kubelet which communicate with apiserver using MQTT.  
`,
		// The lite-kubelet has special flag parsing requirements to enforce flag precedence rules,
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

			if err := options.SetLiteKubeletFlagsByEnv(liteKubeletFlags); err != nil {
				klog.Fatal(err)
			}

			// validate the initial KubeletFlags
			if err := options.ValidateLiteKubeletFlags(liteKubeletFlags); err != nil {
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
			klog.Infof("LiteKubelet Config: %#v", *liteKubeletFlags)
			klog.Infof("Version:%#v", projectinfo.Get())

			if err := Run(ctx, liteKubeletFlags); err != nil {
				klog.Fatal(err)
			}
		},
	}

	// keep cleanFlagSet separate, so Cobra doesn't pollute it with the global flags
	liteKubeletFlags.AddFlags(cleanFlagSet)
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

func Run(ctx context.Context, config *options.LiteKubeletFlags) error {

	if err := RunLiteKubelet(ctx, config); err != nil {
		return err
	}
	return nil
}

func CreateSignal(start chan<- struct{}, receiveStartTime *time.Time, deps *options.LiteKubeletFlags) error {
	stop := make(chan struct{}, 1)
	once := sync.Once{}
	c, err := clientcmd.BuildConfigFromFlags("", deps.KubeConfig)
	if err != nil {
		klog.Errorf("BuildConfigFromFlags error %v", err)
		return err
	}

	// set rate limit
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(2000, 3000)

	// 实例化clientset对象
	client, err := kubernetes.NewForConfig(c)
	if err != nil {
		klog.Errorf("NewForConfig error %v", err)
		return err
	}

	factory := informers.NewSharedInformerFactory(client, 0)
	cmInfor := factory.Core().V1().ConfigMaps().Informer()

	go factory.Start(stop)

	if !cache.WaitForCacheSync(wait.NeverStop,
		cmInfor.HasSynced,
	) {
		return fmt.Errorf("time out")
	}

	// 使用自定义 handler
	cmInfor.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			cm := obj.(*v1.ConfigMap)
			if cm.Namespace == deps.NameSpace && cm.Name == deps.SignalConfigMapName {
				return true
			}
			return false
		},
		Handler: cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				oc := oldObj.(*v1.ConfigMap)
				nc := newObj.(*v1.ConfigMap)
				if !reflect.DeepEqual(oc.Data, nc.Data) {
					*receiveStartTime = time.Now()
					klog.Infof("Find configmap[%s][%s] test value changed from %#v to %#v, close start channel",
						deps.NameSpace, deps.SignalConfigMapName, oc.Data, oc.Data)
					once.Do(func() {
						close(start)
					})
				}

			},
		},
	})
	return nil
}

func RunLiteKubelet(ctx context.Context, deps *options.LiteKubeletFlags) error {
	defer runtime.HandleCrash()

	mqtt5 := true
	start := make(chan struct{}, 1)
	receiveStartTime := &time.Time{}
	allLites := make([]*litekubelet.LiteKubelet, 0, deps.SimulationsNums)

	go func() {
		CreateSignal(start, receiveStartTime, deps)
	}()

	select {
	case <-start:
		klog.Infof("Receive start signal , so start to register")
		break
	}

	for i := 0; i < deps.SimulationsNums; i++ {
		func(index int) {
			lite, err := litekubelet.NewMainLiteKubelet(deps, index, mqtt5)
			if err != nil {
				klog.Errorf("NewMainLiteKublet error %v", err)
				return
			}
			go lite.Run()
			allLites = append(allLites, lite)
		}(i)
		time.Sleep(time.Millisecond * time.Duration(deps.CreateClientInterval))
	}
	klog.Infof("All simulate lite-kubelet connect success")

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sussLen := 0
				for _, l := range allLites {
					if l.Registerd {
						sussLen++
					}
				}
				klog.V(4).Infof("Registering successful num is %d, not registering num is %d", sussLen, deps.SimulationsNums-sussLen)
				if sussLen == deps.SimulationsNums {
					klog.V(4).Infof("All lite-kubelet registering successfull ##%d", time.Now().Unix())
					klog.Infof("All lite-kubelet registering successfull need ms##%d", time.Now().Sub(*receiveStartTime).Milliseconds())
					return
				}
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				allNums := 0
				for _, l := range allLites {
					if l.ReceivePodDataNum == 1 {
						allNums++
					}
				}
				if allNums == deps.SimulationsNums {
					klog.Infof("All lite-kubelet receive pod data##%d", time.Now().Unix())
					return
				}
			}
		}
	}()

	// must exec after SubscribeTopics
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}
