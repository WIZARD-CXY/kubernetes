/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

// Package options provides the scheduler flags
package options

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/componentconfig"
	"k8s.io/kubernetes/pkg/client/leaderelection"
	"k8s.io/kubernetes/pkg/master/ports"
	"k8s.io/kubernetes/plugin/pkg/scheduler/factory"

	"github.com/spf13/pflag"
)

// SchedulerServer has all the context and params needed to run a Scheduler
type SchedulerServer struct {
	componentconfig.KubeSchedulerConfiguration
	// Master is the address of the Kubernetes API server (overrides any
	// value in kubeconfig).
	Master string
	// Kubeconfig is Path to kubeconfig file with authorization and master
	// location information.
	Kubeconfig string
}

// NewSchedulerServer creates a new SchedulerServer with default parameters
func NewSchedulerServer() *SchedulerServer {
	s := SchedulerServer{
		KubeSchedulerConfiguration: componentconfig.KubeSchedulerConfiguration{
			Port:              ports.SchedulerPort,
			Address:           "0.0.0.0",
			AlgorithmProvider: factory.DefaultProvider,
			KubeAPIQPS:        50.0,
			KubeAPIBurst:      100,
			SchedulerName:     api.DefaultSchedulerName,
			LeaderElection:    leaderelection.DefaultLeaderElectionConfiguration(),
		},
	}
	return &s
}

// AddFlags adds flags for a specific SchedulerServer to the specified FlagSet
func (s *SchedulerServer) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&s.Port, "port", s.Port, "The port that the scheduler's http service runs on")
	fs.StringVar(&s.Address, "address", s.Address, "The IP address to serve on (set to 0.0.0.0 for all interfaces)")
	fs.StringVar(&s.AlgorithmProvider, "algorithm-provider", s.AlgorithmProvider, "The scheduling algorithm provider to use, one of: "+factory.ListAlgorithmProviders())
	fs.StringVar(&s.PolicyConfigFile, "policy-config-file", s.PolicyConfigFile, "File with scheduler policy configuration")
	fs.BoolVar(&s.EnableProfiling, "profiling", true, "Enable profiling via web interface host:port/debug/pprof/")
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	var unusedBindPodsQPS float32
	fs.Float32Var(&unusedBindPodsQPS, "bind-pods-qps", 0, "unused, use --kube-api-qps")
	fs.MarkDeprecated("bind-pods-qps", "flag is unused and will be removed. Use kube-api-qps instead.")
	var unusedBindPodsBurst int
	fs.IntVar(&unusedBindPodsBurst, "bind-pods-burst", 0, "unused, use --kube-api-burst")
	fs.MarkDeprecated("bind-pods-burst", "flag is unused and will be removed. Use kube-api-burst instead.")
	fs.Float32Var(&s.KubeAPIQPS, "kube-api-qps", s.KubeAPIQPS, "QPS to use while talking with kubernetes apiserver")
	fs.IntVar(&s.KubeAPIBurst, "kube-api-burst", s.KubeAPIBurst, "Burst to use while talking with kubernetes apiserver")
	fs.StringVar(&s.SchedulerName, "scheduler-name", s.SchedulerName, "Name of the scheduler, used to select which pods will be processed by this scheduler, based on pod's annotation with key 'scheduler.alpha.kubernetes.io/name'")
	leaderelection.BindFlags(&s.LeaderElection, fs)
}
