package tempest

import (
	"github.com/openstack-k8s-operators/lib-common/modules/storage"
	testv1beta1 "github.com/openstack-k8s-operators/test-operator/api/v1beta1"
	util "github.com/openstack-k8s-operators/test-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes - returns a list of volumes for the test pod
func GetVolumes(
	instance *testv1beta1.Tempest,
	customDataConfigMapName string,
	logsPVCName string,
	mountCerts bool,
	mountSSHKey bool,
	svc []storage.PropagationType,
) []corev1.Volume {

	volumes := []corev1.Volume{
		util.CreateConfigMapVolume("config-data", customDataConfigMapName, util.ScriptsVolumeDefaultMode),
		util.CreateOpenstackConfigMapVolume("openstack-config"),
		util.CreateOpenstackConfigSecretVolume(),
		util.CreateLogsPVCVolume(logsPVCName),
		util.CreateWorkdirVolume(),
		util.CreateTmpVolume(),
	}

	if mountCerts {
		volumes = util.AppendCACertsVolume(volumes)
	}

	if mountSSHKey {
		volumes = util.AppendSSHKeyVolumeWithPath(volumes, "ssh-key", instance.Spec.SSHKeySecretName, "ssh-privatekey", "ssh_key")
	}

	volumes = util.AppendExtraMountsVolumes(volumes, instance.Spec.ExtraMounts, svc)
	volumes = util.AppendExtraConfigmapsVolumes(volumes, instance.Spec.ExtraConfigmapsMounts, util.ScriptsVolumeDefaultMode)

	return volumes
}

// GetVolumeMounts - returns a list of volume mounts for the test container
func GetVolumeMounts(
	mountCerts bool,
	mountSSHKey bool,
	svc []storage.PropagationType,
	instance *testv1beta1.Tempest,
) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		util.CreateVolumeMount(util.TestOperatorEphemeralVolumeNameWorkdir, "/var/lib/tempest", false),
		util.CreateVolumeMount(util.TestOperatorEphemeralVolumeNameTmp, "/tmp", false),
		util.CreateVolumeMount("config-data", "/etc/test_operator", false),
		util.CreateVolumeMount(util.VolumeNameTestOperatorLogs, "/var/lib/tempest/external_files", false),
		util.CreateOpenstackConfigVolumeMount("/etc/openstack/clouds.yaml"),
		util.CreateOpenstackConfigVolumeMount("/var/lib/tempest/.config/openstack/clouds.yaml"),
		util.CreateOpenstackConfigSecretVolumeMount("/etc/openstack/secure.yaml"),
	}

	if mountCerts {
		volumeMounts = append(volumeMounts, util.CreateCACertVolumeMount("/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"))
	}

	if mountSSHKey {
		volumeMounts = append(volumeMounts,
			util.CreateVolumeMountWithSubPath("ssh-key", "/var/lib/tempest/id_ecdsa", "ssh_key", false),
		)
	}

	volumeMounts = util.AppendExtraMountsVolumeMounts(volumeMounts, instance.Spec.ExtraMounts, svc)
	volumeMounts = util.AppendExtraConfigmapsVolumeMounts(volumeMounts, instance.Spec.ExtraConfigmapsMounts)

	return volumeMounts
}
