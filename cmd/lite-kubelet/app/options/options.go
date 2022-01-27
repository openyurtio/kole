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

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

type LiteKubeletFlags struct {
	Mqtt5Flags *Mqtt5Flags
	Mqtt3Flags *Mqtt3Flags
	IsMqtt5    bool

	NameSpace           string
	SignalConfigMapName string

	SimulationsNums int
	// m
	HeartBeatInterval int
	// ms
	CreateClientInterval int
	PersistentDir        string
	KubeConfig           string
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

func NewLiteKubeletFlags() *LiteKubeletFlags {

	ns := os.Getenv("NAME_SPACE")
	if len(ns) == 0 {
		ns = "kole"
	}

	return &LiteKubeletFlags{
		SimulationsNums:      1,
		HeartBeatInterval:    120,
		PersistentDir:        "/etc/lite-kubelet/",
		CreateClientInterval: 500,
		Mqtt3Flags:           &Mqtt3Flags{},
		Mqtt5Flags:           &Mqtt5Flags{},
		NameSpace:            ns,
		SignalConfigMapName:  "lite-kubelet-start-signal",
	}
}

// AddFlags adds flags for a specific LiteKubeletFlags to the specified FlagSet
func (f *LiteKubeletFlags) AddFlags(mainfs *pflag.FlagSet) {
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

	if f.Mqtt5Flags == nil {
		f.Mqtt5Flags = &Mqtt5Flags{}
	}
	if f.Mqtt3Flags == nil {
		f.Mqtt3Flags = &Mqtt3Flags{}
	}

	fs.StringVar(&f.Mqtt3Flags.MqttBroker, "mqtt3-broker", f.Mqtt3Flags.MqttBroker, "the address of mqtt broker")
	fs.IntVar(&f.Mqtt3Flags.MqttBrokerPort, "mqtt3-broker-port", f.Mqtt3Flags.MqttBrokerPort, "the port of mqtt broker")
	fs.StringVar(&f.Mqtt3Flags.MqttGroup, "mqtt3-group", f.Mqtt3Flags.MqttGroup, "the mqtt group")
	fs.StringVar(&f.Mqtt3Flags.MqttInstance, "mqtt3-instance", f.Mqtt3Flags.MqttInstance, "mqtt instance name")

	fs.StringVar(&f.Mqtt5Flags.MqttServer, "mqtt5-server", f.Mqtt5Flags.MqttServer, "mqtt5 server name")

	fs.IntVar(&f.HeartBeatInterval, "heartbeat-interval", f.HeartBeatInterval, "heartbeat-interval (s)")
	fs.IntVar(&f.CreateClientInterval, "create-client-interval", f.CreateClientInterval, "create mqtt client interval (ms)")
	fs.IntVar(&f.SimulationsNums, "simulations-nums", f.SimulationsNums, "the number of simulations")
	fs.StringVar(&f.PersistentDir, "persistent-dir", f.PersistentDir, "persistent dir")
	fs.StringVar(&f.KubeConfig, "kubeconfig", f.KubeConfig, "Path to a kubeconfig file, specifying how to connect to the API server.")
	fs.StringVar(&f.SignalConfigMapName, "signal-cm-name", f.SignalConfigMapName, "the name of configmap name which trigger lite-kubelet start.")
}

func SetLiteKubeletFlagsByEnv(f *LiteKubeletFlags) error {
	numStr := os.Getenv("SIMULATIONS_NUMS")
	if len(numStr) != 0 {
		if num, err := strconv.Atoi(numStr); err != nil {
			klog.Errorf("Can not atoi %s, error %v", numStr, err)
			return err
		} else {
			f.SimulationsNums = num
			klog.Infof("Set --simulations-nums value to %d by env", f.SimulationsNums)
		}
	}

	createIntervalStr := os.Getenv("CREATE_CLIENT_INTERVAL")
	if len(createIntervalStr) != 0 {
		if interval, err := strconv.Atoi(createIntervalStr); err != nil {
			klog.Errorf("Can not atoi %s, error %v", createIntervalStr, err)
			return err
		} else {
			f.CreateClientInterval = interval
			klog.Infof("Set --create-client-interval value to %d by env", f.CreateClientInterval)
		}
	}

	heatBeatIntervalStr := os.Getenv("HEAT_BEAT_INTERVAL")
	if len(heatBeatIntervalStr) != 0 {
		if interval, err := strconv.Atoi(heatBeatIntervalStr); err != nil {
			klog.Errorf("Can not atoi %s, error %v", heatBeatIntervalStr, err)
			return err
		} else {
			f.HeartBeatInterval = interval
			klog.Infof("Set --heartbeat-interval value to %d by env", f.CreateClientInterval)
		}
	}

	if numStr := os.Getenv("MQTT5_SERVER"); len(numStr) != 0 {
		f.Mqtt5Flags.MqttServer = numStr
		klog.Infof("Set --mqtt5-server value to %s by env", f.Mqtt5Flags.MqttServer)
	}

	if numStr := os.Getenv("SIGNAL_CM_NAME"); len(numStr) != 0 {
		f.SignalConfigMapName = numStr
		klog.Infof("Set --signal-cm-name value to %s by env", f.SignalConfigMapName)
	}

	return nil
}

// ValidateLiteKubeletFlags validates litekubelet's configuration flags and returns an error if they are invalid.
func ValidateLiteKubeletFlags(f *LiteKubeletFlags) error {
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
