package horizontest

import (
	"github.com/openstack-k8s-operators/lib-common/modules/storage"
	testv1beta1 "github.com/openstack-k8s-operators/test-operator/api/v1beta1"
	util "github.com/openstack-k8s-operators/test-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumes - returns a list of volumes for the test pod
func GetVolumes(
	instance *testv1beta1.HorizonTest,
	logsPVCName string,
	mountCerts bool,
	mountKubeconfig bool,
	svc []storage.PropagationType,
) []corev1.Volume {

	volumes := []corev1.Volume{
		util.CreateConfigMapVolume("horizontest-config", instance.Name+"horizontest-config", util.ScriptsVolumeDefaultMode),
		util.CreateOpenstackConfigMapVolume(util.TestOperatorCloudsConfigMapName),
		util.CreateOpenstackConfigSecretVolume(),
		util.CreateLogsPVCVolume(logsPVCName),
		util.CreateWorkdirVolume(),
		util.CreateTmpVolume(),
	}

	if mountCerts {
		volumes = util.AppendCACertsVolume(volumes)
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
	instance *testv1beta1.HorizonTest,
) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		util.CreateVolumeMount(util.TestOperatorEphemeralVolumeNameWorkdir, "/var/lib/horizontest", false),
		util.CreateVolumeMount(util.TestOperatorEphemeralVolumeNameTmp, "/tmp", false),
		util.CreateVolumeMount(util.VolumeNameTestOperatorLogs, "/var/lib/horizontest/external_files", false),
		util.CreateTestOperatorCloudsConfigVolumeMount("/var/lib/horizontest/.config/openstack/clouds.yaml"),
		util.CreateTestOperatorCloudsConfigVolumeMount("/etc/openstack/clouds.yaml"),
		util.CreateOpenstackConfigSecretVolumeMount("/etc/openstack/secure.yaml"),
	}

	if mountCerts {
		volumeMounts = append(volumeMounts,
			util.CreateCACertVolumeMount("/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"),
			util.CreateCACertVolumeMount("/etc/pki/tls/certs/ca-bundle.trust.crt"),
		)
	}

	if mountKeys {
		volumeMounts = append(volumeMounts,
			util.CreateVolumeMountWithSubPath("horizontest-private-key", "/etc/test_operator/id_ecdsa", "id_ecdsa", true),
			util.CreateVolumeMountWithSubPath("horizontest-public-key", "/etc/test_operator/id_ecdsa.pub", "id_ecdsa.pub", true),
		)
	}

	if mountKubeconfig {
		volumeMounts = append(volumeMounts,
			util.CreateVolumeMountWithSubPath(util.VolumeNameKubeconfig, "/var/lib/horizontest/.kube/config", "config", true),
		)
	}

	volumeMounts = util.AppendExtraMountsVolumeMounts(volumeMounts, instance.Spec.ExtraMounts, svc)
	volumeMounts = util.AppendExtraConfigmapsVolumeMounts(volumeMounts, instance.Spec.ExtraConfigmapsMounts)

	return volumeMounts
}
