/*
Copyright 2018 The Kubernetes Authors.

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

package managed

import (
	"time"

	"github.com/kubernetes-sigs/federation-v2/pkg/apis/core/typeconfig"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/dnsendpoint"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/federatedcluster"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/ingressdns"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/schedulingmanager"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/servicedns"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/status"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/sync"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/util"
	"github.com/kubernetes-sigs/federation-v2/test/common"
)

// ControllerFixture manages a federation controller for testing.
type ControllerFixture struct {
	stopChan chan struct{}
}

// NewSyncControllerFixture initializes a new sync controller fixture.
func NewSyncControllerFixture(tl common.TestLogger, controllerConfig *util.ControllerConfig, typeConfig typeconfig.Interface) *ControllerFixture {
	f := &ControllerFixture{
		stopChan: make(chan struct{}),
	}
	err := sync.StartFederationSyncController(controllerConfig, f.stopChan, typeConfig)
	if err != nil {
		tl.Fatalf("Error starting sync controller: %v", err)
	}
	if typeConfig.GetEnableStatus() {
		err := status.StartFederationStatusController(controllerConfig, f.stopChan, typeConfig)
		if err != nil {
			tl.Fatalf("Error starting status controller: %v", err)
		}
	}
	return f
}

// NewServiceDNSControllerFixture initializes a new service-dns controller fixture.
func NewServiceDNSControllerFixture(tl common.TestLogger, config *util.ControllerConfig) *ControllerFixture {
	f := &ControllerFixture{
		stopChan: make(chan struct{}),
	}
	err := servicedns.StartController(config, f.stopChan)
	if err != nil {
		tl.Fatalf("Error starting service dns controller: %v", err)
	}
	err = dnsendpoint.StartServiceDNSEndpointController(config, f.stopChan)
	if err != nil {
		tl.Fatalf("Error starting service dns endpoint controller: %v", err)
	}
	return f
}

// NewIngressDNSControllerFixture initializes a new ingress-dns controller fixture.
func NewIngressDNSControllerFixture(tl common.TestLogger, config *util.ControllerConfig) *ControllerFixture {
	f := &ControllerFixture{
		stopChan: make(chan struct{}),
	}
	err := ingressdns.StartController(config, f.stopChan)
	if err != nil {
		tl.Fatalf("Error starting ingress dns controller: %v", err)
	}
	err = dnsendpoint.StartIngressDNSEndpointController(config, f.stopChan)
	if err != nil {
		tl.Fatalf("Error starting ingress dns endpoint controller: %v", err)
	}
	return f
}

// NewClusterControllerFixture initializes a new cluster controller fixture.
func NewClusterControllerFixture(config *util.ControllerConfig) *ControllerFixture {
	f := &ControllerFixture{
		stopChan: make(chan struct{}),
	}
	monitorPeriod := 1 * time.Second
	federatedcluster.StartClusterController(config, f.stopChan, monitorPeriod)
	return f
}

func NewSchedulerControllerFixture(tl common.TestLogger, config *util.ControllerConfig) *ControllerFixture {
	f := &ControllerFixture{
		stopChan: make(chan struct{}),
	}
	schedulingmanager.StartSchedulerController(config, f.stopChan)
	return f
}

func (f *ControllerFixture) TearDown(tl common.TestLogger) {
	if f.stopChan != nil {
		close(f.stopChan)
		f.stopChan = nil
	}
}
