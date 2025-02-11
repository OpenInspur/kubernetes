/*
Copyright 2017 The Kubernetes Authors.

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

package validation

import (
	"fmt"
	"time"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilvalidation "k8s.io/apimachinery/pkg/util/validation"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/features"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

// ValidateKubeletConfiguration validates `kc` and returns an error if it is invalid
func ValidateKubeletConfiguration(kc *kubeletconfig.KubeletConfiguration) error {
	allErrors := []error{}

	// Make a local copy of the global feature gates and combine it with the gates set by this configuration.
	// This allows us to validate the config against the set of gates it will actually run against.
	localFeatureGate := utilfeature.DefaultFeatureGate.DeepCopy()
	localFeatureGate.SetFromMap(kc.FeatureGates)

	if kc.NodeLeaseDurationSeconds <= 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: NodeLeaseDurationSeconds must be greater than 0"))
	}
	if !kc.CgroupsPerQOS && len(kc.EnforceNodeAllocatable) > 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: EnforceNodeAllocatable (--enforce-node-allocatable) is not supported unless CgroupsPerQOS (--cgroups-per-qos) feature is turned on"))
	}
	if kc.SystemCgroups != "" && kc.CgroupRoot == "" {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: SystemCgroups (--system-cgroups) was specified and CgroupRoot (--cgroup-root) was not specified"))
	}
	if kc.EventBurst < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: EventBurst (--event-burst) %v must not be a negative number", kc.EventBurst))
	}
	if kc.EventRecordQPS < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: EventRecordQPS (--event-qps) %v must not be a negative number", kc.EventRecordQPS))
	}
	if kc.HealthzPort != 0 && utilvalidation.IsValidPortNum(int(kc.HealthzPort)) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: HealthzPort (--healthz-port) %v must be between 1 and 65535, inclusive", kc.HealthzPort))
	}
	if localFeatureGate.Enabled(features.CPUCFSQuotaPeriod) && utilvalidation.IsInRange(int(kc.CPUCFSQuotaPeriod.Duration), int(1*time.Microsecond), int(time.Second)) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: CPUCFSQuotaPeriod (--cpu-cfs-quota-period) %v must be between 1usec and 1sec, inclusive", kc.CPUCFSQuotaPeriod))
	}
	if utilvalidation.IsInRange(int(kc.ImageGCHighThresholdPercent), 0, 100) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: ImageGCHighThresholdPercent (--image-gc-high-threshold) %v must be between 0 and 100, inclusive", kc.ImageGCHighThresholdPercent))
	}
	if utilvalidation.IsInRange(int(kc.ImageGCLowThresholdPercent), 0, 100) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: ImageGCLowThresholdPercent (--image-gc-low-threshold) %v must be between 0 and 100, inclusive", kc.ImageGCLowThresholdPercent))
	}
	if kc.ImageGCLowThresholdPercent >= kc.ImageGCHighThresholdPercent {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: ImageGCLowThresholdPercent (--image-gc-low-threshold) %v must be less than ImageGCHighThresholdPercent (--image-gc-high-threshold) %v", kc.ImageGCLowThresholdPercent, kc.ImageGCHighThresholdPercent))
	}
	if utilvalidation.IsInRange(int(kc.IPTablesDropBit), 0, 31) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: IPTablesDropBit (--iptables-drop-bit) %v must be between 0 and 31, inclusive", kc.IPTablesDropBit))
	}
	if utilvalidation.IsInRange(int(kc.IPTablesMasqueradeBit), 0, 31) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: IPTablesMasqueradeBit (--iptables-masquerade-bit) %v must be between 0 and 31, inclusive", kc.IPTablesMasqueradeBit))
	}
	if kc.KubeAPIBurst < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: KubeAPIBurst (--kube-api-burst) %v must not be a negative number", kc.KubeAPIBurst))
	}
	if kc.KubeAPIQPS < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: KubeAPIQPS (--kube-api-qps) %v must not be a negative number", kc.KubeAPIQPS))
	}
	if kc.MaxOpenFiles < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: MaxOpenFiles (--max-open-files) %v must not be a negative number", kc.MaxOpenFiles))
	}
	if kc.MaxPods < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: MaxPods (--max-pods) %v must not be a negative number", kc.MaxPods))
	}
	if utilvalidation.IsInRange(int(kc.OOMScoreAdj), -1000, 1000) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: OOMScoreAdj (--oom-score-adj) %v must be between -1000 and 1000, inclusive", kc.OOMScoreAdj))
	}
	if kc.PodsPerCore < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: PodsPerCore (--pods-per-core) %v must not be a negative number", kc.PodsPerCore))
	}
	if utilvalidation.IsValidPortNum(int(kc.Port)) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: Port (--port) %v must be between 1 and 65535, inclusive", kc.Port))
	}
	if kc.ReadOnlyPort != 0 && utilvalidation.IsValidPortNum(int(kc.ReadOnlyPort)) != nil {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: ReadOnlyPort (--read-only-port) %v must be between 0 and 65535, inclusive", kc.ReadOnlyPort))
	}
	if kc.RegistryBurst < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: RegistryBurst (--registry-burst) %v must not be a negative number", kc.RegistryBurst))
	}
	if kc.RegistryPullQPS < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: RegistryPullQPS (--registry-qps) %v must not be a negative number", kc.RegistryPullQPS))
	}
	if kc.RotateCertificates && !localFeatureGate.Enabled(features.RotateKubeletClientCertificate) {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: RotateCertificates %v requires feature gate RotateKubeletClientCertificate", kc.RotateCertificates))
	}
	if kc.ServerTLSBootstrap && !localFeatureGate.Enabled(features.RotateKubeletServerCertificate) {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: ServerTLSBootstrap %v requires feature gate RotateKubeletServerCertificate", kc.ServerTLSBootstrap))
	}
	if kc.ContainerBackOffResetSeconds < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: ContainerBackOffResetSeconds (--container-backoff-reset-seconds) %v must not be a negative number", kc.ContainerBackOffResetSeconds))
	}
	if kc.MaxContainerBackOffSeconds < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: MaxContainerBackOffSeconds (--max-container-backoff-seconds) %v must not be a negative number", kc.MaxContainerBackOffSeconds))
	}
	if kc.NodeLeaseRenewIntervalSeconds < 0 {
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: NodeLeaseRenewIntervalSeconds (--node-lease-renew-interval-seconds) %v must not be a negative number", kc.NodeLeaseRenewIntervalSeconds))
	}
	for _, val := range kc.EnforceNodeAllocatable {
		switch val {
		case kubetypes.NodeAllocatableEnforcementKey:
		case kubetypes.SystemReservedEnforcementKey:
		case kubetypes.KubeReservedEnforcementKey:
		case kubetypes.NodeAllocatableNoneKey:
			if len(kc.EnforceNodeAllocatable) > 1 {
				allErrors = append(allErrors, fmt.Errorf("invalid configuration: EnforceNodeAllocatable (--enforce-node-allocatable) may not contain additional enforcements when '%s' is specified", kubetypes.NodeAllocatableNoneKey))
			}
		default:
			allErrors = append(allErrors, fmt.Errorf("invalid configuration: option %q specified for EnforceNodeAllocatable (--enforce-node-allocatable). Valid options are %q, %q, %q, or %q",
				val, kubetypes.NodeAllocatableEnforcementKey, kubetypes.SystemReservedEnforcementKey, kubetypes.KubeReservedEnforcementKey, kubetypes.NodeAllocatableNoneKey))
		}
	}
	switch kc.HairpinMode {
	case kubeletconfig.HairpinNone:
	case kubeletconfig.HairpinVeth:
	case kubeletconfig.PromiscuousBridge:
	default:
		allErrors = append(allErrors, fmt.Errorf("invalid configuration: option %q specified for HairpinMode (--hairpin-mode). Valid options are %q, %q or %q",
			kc.HairpinMode, kubeletconfig.HairpinNone, kubeletconfig.HairpinVeth, kubeletconfig.PromiscuousBridge))
	}

	if err := validateKubeletOSConfiguration(kc); err != nil {
		allErrors = append(allErrors, err)
	}
	return utilerrors.NewAggregate(allErrors)
}
