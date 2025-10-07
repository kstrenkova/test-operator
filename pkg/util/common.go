// Package util provides common utility functions and constants for test operations
package util //nolint:revive // util is a legitimate package name for utility functions

import (
	"github.com/openstack-k8s-operators/lib-common/modules/storage"
	corev1 "k8s.io/api/core/v1"
)

const (
	// TestOperatorCloudsConfigMapName is the name of the ConfigMap which contains
	// modified clouds.yaml obtained from openstack-config ConfigMap. The modified
	// CM is needed by some test frameworks (e.g., HorizonTest and Tobiko)
	TestOperatorCloudsConfigMapName = "test-operator-clouds-config"

	// TestOperatorEphemeralVolumeNameWorkdir is the name of the ephemeral workdir volume
	TestOperatorEphemeralVolumeNameWorkdir = "test-operator-ephemeral-workdir"

	// TestOperatorEphemeralVolumeNameTmp is the name of the ephemeral temporary volume
	TestOperatorEphemeralVolumeNameTmp = "test-operator-ephemeral-temporary"

	// ExtraVolTypeUndefined can be used to label an extraMount which is
	// not associated to anything in particular
	ExtraVolTypeUndefined storage.ExtraVolType = "Undefined"
)

// GetSecurityContext returns a security context with the specified configuration
func GetSecurityContext(
	runAsUser int64,
	addCapabilities []corev1.Capability,
	privileged bool,
) corev1.SecurityContext {
	falseVar := false
	trueVar := true

	securityContext := corev1.SecurityContext{
		RunAsUser:                &runAsUser,
		RunAsGroup:               &runAsUser,
		ReadOnlyRootFilesystem:   &trueVar,
		RunAsNonRoot:             &trueVar,
		AllowPrivilegeEscalation: &falseVar,
		Capabilities:             &corev1.Capabilities{},
		SeccompProfile: &corev1.SeccompProfile{
			Type: corev1.SeccompProfileTypeRuntimeDefault,
		},
	}

	if privileged {
		// Sometimes we require the test pods run sudo to be able to install
		// additional packages or run commands with elevated privileges (e.g.,
		// tcpdump in case of Tobiko)
		securityContext.RunAsNonRoot = &falseVar

		// We need to run pods with AllowPrivilegedEscalation: true to remove
		// nosuid from the pod (in order to be able to run sudo)
		securityContext.AllowPrivilegeEscalation = &trueVar

		// We need to run pods with ReadOnlyRootFileSystem: false when installing
		// additional tests using extraRPMs parameter in Tempest CR
		securityContext.ReadOnlyRootFilesystem = &falseVar
		securityContext.Capabilities.Add = addCapabilities
	}

	if !privileged {
		// We need to keep default capabilities in order to be able to use sudo
		securityContext.Capabilities.Drop = []corev1.Capability{"ALL"}
	}

	return securityContext
}
