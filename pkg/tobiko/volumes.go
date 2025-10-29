package tobiko

import (
	"github.com/openstack-k8s-operators/lib-common/modules/storage"
	testv1beta1 "github.com/openstack-k8s-operators/test-operator/api/v1beta1"
	"github.com/openstack-k8s-operators/test-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes - returns a list of volumes for the test pod
func GetVolumes(
	instance *testv1beta1.Tobiko,
	logsPVCName string,
	mountCerts bool,
	mountKeys bool,
	mountKubeconfig bool,
	svc []storage.PropagationType,
) []corev1.Volume {

	volumes := []corev1.Volume{
		util.CreateConfigMapVolume("tobiko-config", instance.Name+"tobiko-config", util.ScriptsVolumeDefaultMode),
		util.CreateOpenstackConfigMapVolume(util.TestOperatorCloudsConfigMapName),
		util.CreateOpenstackConfigSecretVolume(),
		util.CreateLogsPVCVolume(logsPVCName),
		util.CreateWorkdirVolume(),
		util.CreateTmpVolume(),
	}

	if mountCerts {
		volumes = util.AppendCACertsVolume(volumes)
	}

	if mountKeys {
		volumes = append(volumes,
			util.CreateConfigMapVolume("tobiko-private-key", instance.Name+"tobiko-private-key", util.PrivateKeyMode),
			util.CreateConfigMapVolume("tobiko-public-key", instance.Name+"tobiko-public-key", util.PublicKeyMode),
		)
	}

	if mountKubeconfig {
		volumes = util.AppendKubeconfigVolume(volumes, instance.Spec.KubeconfigSecretName)
	}

	volumes = util.AppendExtraMountsVolumes(volumes, instance.Spec.ExtraMounts, svc)
	volumes = util.AppendExtraConfigmapsVolumes(volumes, instance.Spec.ExtraConfigmapsMounts, util.PublicInfoMode)

	return volumes
}

// GetVolumeMounts - returns a list of volume mounts for the test container
func GetVolumeMounts(
	mountCerts bool,
	mountKeys bool,
	mountKubeconfig bool,
	svc []storage.PropagationType,
	instance *testv1beta1.Tobiko,
) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		util.CreateVolumeMount(util.TestOperatorEphemeralVolumeNameWorkdir, "/var/lib/tobiko", false),
		util.CreateVolumeMount(util.TestOperatorEphemeralVolumeNameTmp, "/tmp", false),
		util.CreateVolumeMount(util.VolumeNameTestOperatorLogs, "/var/lib/tobiko/external_files", false),
		util.CreateTestOperatorCloudsConfigVolumeMount("/var/lib/tobiko/.config/openstack/clouds.yaml"),
		util.CreateTestOperatorCloudsConfigVolumeMount("/etc/openstack/clouds.yaml"),
		util.CreateOpenstackConfigSecretVolumeMount("/etc/openstack/secure.yaml"),
		util.CreateVolumeMountWithSubPath("tobiko-config", "/etc/tobiko/tobiko.conf", "tobiko.conf", false),
	}

	if mountCerts {
		volumeMounts = append(volumeMounts,
			util.CreateCACertVolumeMount("/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"),
			util.CreateCACertVolumeMount("/etc/pki/tls/certs/ca-bundle.trust.crt"),
		)
	}

	if mountKeys {
		volumeMounts = append(volumeMounts,
			util.CreateVolumeMountWithSubPath("tobiko-private-key", "/etc/test_operator/id_ecdsa", "id_ecdsa", true),
			util.CreateVolumeMountWithSubPath("tobiko-public-key", "/etc/test_operator/id_ecdsa.pub", "id_ecdsa.pub", true),
		)
	}

	if mountKubeconfig {
		volumeMounts = append(volumeMounts,
			util.CreateVolumeMountWithSubPath(util.VolumeNameKubeconfig, "/var/lib/tobiko/.kube/config", "config", true),
		)
	}

	volumeMounts = util.AppendExtraMountsVolumeMounts(volumeMounts, instance.Spec.ExtraMounts, svc)
	volumeMounts = util.AppendExtraConfigmapsVolumeMounts(volumeMounts, instance.Spec.ExtraConfigmapsMounts)

	return volumeMounts
}
