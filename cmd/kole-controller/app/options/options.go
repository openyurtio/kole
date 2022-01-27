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

package options

import (
	"fmt"
	"os"
	"strconv"

	"k8s.io/klog/v2"

	"github.com/spf13/pflag"
)

type KoleControllerFlags struct {
	Mqtt3Flags *Mqtt3Flags
	Mqtt5Flags *Mqtt5Flags
	IsMqtt5    bool

	KubeConfig       string
	NameSpace        string
	SnapshotInterval int
	// s
	HBTimeOut int
}

type Mqtt3Flags struct {
	MqttBroker     string
	MqttBrokerPort int
	MqttGroup      string
	MqttInstance   string
}

type Mqtt5Flags struct {
	MqttServer string
}

func NewKoleControllerFlags() *KoleControllerFlags {
	ns := os.Getenv("NAME_SPACE")
	if len(ns) == 0 {
		ns = "kole"
	}
	return &KoleControllerFlags{
		SnapshotInterval: 60,     // second
		HBTimeOut:        60 * 5, // second
		NameSpace:        ns,
		Mqtt3Flags:       &Mqtt3Flags{},
		Mqtt5Flags:       &Mqtt5Flags{},
	}
}

// AddFlags adds flags for a specific KoleControllerFlags to the specified FlagSet
func (f *KoleControllerFlags) AddFlags(mainfs *pflag.FlagSet) {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	defer func() {
		// Unhide deprecated flags. We want deprecated flags to show in Kubelet help.
		// We have some hidden flags, but we might as well unhide these when they are deprecated,
		// as silently deprecating and removing (even hidden) things is unkind to people who use them.
		fs.VisitAll(func(f *pflag.Flag) {
			if len(f.Deprecated) > 0 {
				f.Hidden = false
			}
		})
		mainfs.AddFlagSet(fs)
	}()
	if f.Mqtt3Flags == nil {
		f.Mqtt3Flags = &Mqtt3Flags{}
	}

	if f.Mqtt5Flags == nil {
		f.Mqtt5Flags = &Mqtt5Flags{}
	}

	fs.StringVar(&f.Mqtt3Flags.MqttBroker, "mqtt3-broker", f.Mqtt3Flags.MqttBroker, "the address of mqtt broker")
	fs.IntVar(&f.Mqtt3Flags.MqttBrokerPort, "mqtt3-broker-port", f.Mqtt3Flags.MqttBrokerPort, "the port of mqtt broker")
	fs.StringVar(&f.Mqtt3Flags.MqttGroup, "mqtt3-group", f.Mqtt3Flags.MqttGroup, "the mqtt group")
	fs.StringVar(&f.Mqtt3Flags.MqttInstance, "mqtt3-instance", f.Mqtt3Flags.MqttInstance, "mqtt instance name")

	fs.StringVar(&f.Mqtt5Flags.MqttServer, "mqtt5-server", f.Mqtt5Flags.MqttServer, "mqtt5 server")

	fs.StringVar(&f.KubeConfig, "kubeconfig", f.KubeConfig, "Path to a kubeconfig file, specifying how to connect to the API server.")
	fs.IntVar(&f.SnapshotInterval, "snapshot-interval", f.SnapshotInterval, "snapshot interval (second)")
	fs.IntVar(&f.HBTimeOut, "hb-timeout", f.HBTimeOut, "hb time out(second)")
}

// ValidateKoleControllerFlags validates litekubelet's configuration flags and returns an error if they are invalid.
func ValidateKoleControllerFlags(f *KoleControllerFlags) error {
	// ensure that nobody sets DynamicConfigDir if the dynamic config feature gate is turned off

	// mqtt3
	if len(f.Mqtt5Flags.MqttServer) == 0 {
		switch {
		case len(f.Mqtt3Flags.MqttInstance) == 0:
			return fmt.Errorf("need set mqtt3-instance")
		case len(f.Mqtt3Flags.MqttBroker) == 0:
			return fmt.Errorf("need set mqtt3-broker")
		case f.Mqtt3Flags.MqttBrokerPort == 0:
			return fmt.Errorf("need set mqtt3-broker-port")
		case len(f.Mqtt3Flags.MqttGroup) == 0:
			return fmt.Errorf("need set mqtt3-group")
		}

	} else {
		f.IsMqtt5 = true
	}

	return nil
}

func SetKoleControllerFlagsByEnv(f *KoleControllerFlags) error {
	if numStr := os.Getenv("HB_TIMEOUT"); len(numStr) != 0 {
		if num, err := strconv.Atoi(numStr); err != nil {
			klog.Errorf("Can not atoi %s, error %v", numStr, err)
			return err
		} else {
			f.HBTimeOut = num
			klog.Infof("Set --hb-timeout value to %d by env", f.HBTimeOut)
		}
	}

	if numStr := os.Getenv("MQTT5_SERVER"); len(numStr) != 0 {
		f.Mqtt5Flags.MqttServer = numStr
		klog.Infof("Set --mqtt5-server value to %s by env", f.Mqtt5Flags.MqttServer)
	}

	return nil
}
