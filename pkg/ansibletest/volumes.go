package ansibletest

import (
	"github.com/openstack-k8s-operators/lib-common/modules/storage"
	testv1beta1 "github.com/openstack-k8s-operators/test-operator/api/v1beta1"
	util "github.com/openstack-k8s-operators/test-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes - returns a list of volumes for the test pod
func GetVolumes(
	instance *testv1beta1.AnsibleTest,
	logsPVCName string,
	mountCerts bool,
	svc []storage.PropagationType,
	externalWorkflowCounter int,
) []corev1.Volume {

	volumes := []corev1.Volume{
		util.CreateOpenstackConfigMapVolume("openstack-config"),
		util.CreateOpenstackConfigSecretVolume(),
		util.CreateLogsPVCVolume(logsPVCName),
		util.CreateWorkdirVolume(),
		util.CreateTmpVolume(),
	}

	if mountCerts {
		volumes = util.AppendCACertsVolume(volumes)
	}

	volumes = util.AppendSSHKeyVolume(volumes, "compute-ssh-secret", instance.Spec.ComputeSSHKeySecretName)
	volumes = util.AppendSSHKeyVolume(volumes, "workload-ssh-secret", instance.Spec.WorkloadSSHKeySecretName)

	volumes = util.AppendExtraMountsVolumes(volumes, instance.Spec.ExtraMounts, svc)
	volumes = util.AppendExtraConfigmapsVolumes(volumes, instance.Spec.ExtraConfigmapsMounts, util.ScriptsVolumeDefaultMode)

	if len(instance.Spec.Workflow) > 0 && instance.Spec.Workflow[externalWorkflowCounter].ExtraConfigmapsMounts != nil {
		volumes = util.AppendExtraConfigmapsVolumes(volumes, *instance.Spec.Workflow[externalWorkflowCounter].ExtraConfigmapsMounts, util.ScriptsVolumeDefaultMode)
	}

	return volumes
}

// GetVolumeMounts - returns a list of volume mounts for the test container
func GetVolumeMounts(
	mountCerts bool,
	svc []storage.PropagationType,
	instance *testv1beta1.AnsibleTest,
	externalWorkflowCounter int,
) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		util.CreateVolumeMount(util.TestOperatorEphemeralVolumeNameWorkdir, "/var/lib/ansible", false),
		util.CreateVolumeMount(util.TestOperatorEphemeralVolumeNameTmp, "/tmp", false),
		util.CreateVolumeMount(util.VolumeNameTestOperatorLogs, "/var/lib/AnsibleTests/external_files", false),
		util.CreateOpenstackConfigVolumeMount("/etc/openstack/clouds.yaml"),
		util.CreateOpenstackConfigVolumeMount("/var/lib/ansible/.config/openstack/clouds.yaml"),
		util.CreateOpenstackConfigSecretVolumeMount("/var/lib/ansible/.config/openstack/secure.yaml"),
	}

	if mountCerts {
		volumeMounts = append(volumeMounts,
			util.CreateCACertVolumeMount("/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"),
			util.CreateCACertVolumeMount("/etc/pki/tls/certs/ca-bundle.trust.crt"),
		)
	}

	if instance.Spec.WorkloadSSHKeySecretName != "" {
		volumeMounts = append(volumeMounts,
			util.CreateVolumeMountWithSubPath("workload-ssh-secret", "/var/lib/ansible/test_keypair.key", "ssh-privatekey", true),
		)
	}

	volumeMounts = append(volumeMounts,
		util.CreateVolumeMountWithSubPath("compute-ssh-secret", "/var/lib/ansible/.ssh/compute_id", "ssh-privatekey", true),
	)

	volumeMounts = util.AppendExtraMountsVolumeMounts(volumeMounts, instance.Spec.ExtraMounts, svc)
	volumeMounts = util.AppendExtraConfigmapsVolumeMounts(volumeMounts, instance.Spec.ExtraConfigmapsMounts)

	if len(instance.Spec.Workflow) > 0 && instance.Spec.Workflow[externalWorkflowCounter].ExtraConfigmapsMounts != nil {
		volumeMounts = util.AppendExtraConfigmapsVolumeMounts(volumeMounts, *instance.Spec.Workflow[externalWorkflowCounter].ExtraConfigmapsMounts)
	}

	return volumeMounts
}
