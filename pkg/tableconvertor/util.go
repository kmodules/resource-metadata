/*
Copyright AppsCode Inc. and Contributors

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

package tableconvertor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	ValueNone           = "<none>"
	ResourceKindMongoDB = "MongoDB"
)

// ref: https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubectl/pkg/describe/describe.go
func describeVolume(volume core.Volume) string {
	switch {
	case volume.VolumeSource.HostPath != nil:
		hostPath := volume.VolumeSource.HostPath
		hostPathType := ValueNone
		if hostPath.Type != nil {
			hostPathType = string(*hostPath.Type)
		}
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"HostPath\","+
			"\"HostPath Type\": %q,"+
			"\"Path\": \"%v\""+
			"}", volume.Name, hostPathType, hostPath.Path)
	case volume.VolumeSource.EmptyDir != nil:
		var sizeLimit string
		emptyDir := volume.VolumeSource.EmptyDir
		if emptyDir.SizeLimit != nil && emptyDir.SizeLimit.Cmp(resource.Quantity{}) > 0 {
			sizeLimit = fmt.Sprintf("%v", emptyDir.SizeLimit)
		} else {
			sizeLimit = "<unset>"
		}
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"EmptyDir\","+
			"\"Medium\": \"%v\","+
			"\"SizeLimit\": %q"+
			"}", volume.Name, emptyDir.Medium, sizeLimit)
	case volume.VolumeSource.GCEPersistentDisk != nil:
		gce := volume.VolumeSource.GCEPersistentDisk
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"GCEPersistentDisk\","+
			"\"PDName\": %q,"+
			"\"FSType\": %q,"+
			"\"Partition\": %q,"+
			"\"ReadOnly\": \"%v\","+
			"}", volume.Name, gce.PDName, gce.FSType, gce.Partition, gce.ReadOnly)
	case volume.VolumeSource.AWSElasticBlockStore != nil:
		aws := volume.VolumeSource.AWSElasticBlockStore
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"AWSElasticBlockStore\","+
			"\"VolumeID\": %q,"+
			"\"FSType\": %q,"+
			"\"Partition\": %q,"+
			"\"ReadOnly\": \"%v\","+
			"}", volume.Name, aws.VolumeID, aws.FSType, aws.Partition, aws.ReadOnly)
	case volume.VolumeSource.Secret != nil:
		secret := volume.VolumeSource.Secret
		optional := secret.Optional != nil && *secret.Optional
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"Secret\","+
			"\"SecretName\": %q,"+
			"\"Optional\": \"%v\""+
			"}", volume.Name, secret.SecretName, optional)
	case volume.VolumeSource.ConfigMap != nil:
		configMap := volume.VolumeSource.ConfigMap
		optional := configMap.Optional != nil && *configMap.Optional
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"ConfigMap\","+
			"\"ConfigMap Name\": %q,"+
			"\"Optional\": \"%v\""+
			"}", volume.Name, configMap.Name, optional)
	case volume.VolumeSource.NFS != nil:
		nfs := volume.VolumeSource.NFS
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"NFS\","+
			"\"Server\": %q,"+
			"\"Path\": \"%v\","+
			"\"ReadOnly\": \"%v\""+
			"}", volume.Name, nfs.Server, nfs.Path, nfs.ReadOnly)
	case volume.VolumeSource.ISCSI != nil:
		iscsi := volume.VolumeSource.ISCSI
		initiator := ValueNone
		if iscsi.InitiatorName != nil {
			initiator = *iscsi.InitiatorName
		}
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"ISCSI\","+
			"\"TargetPortal\": %q,"+
			"\"IQN\": %q,"+
			"\"Lun\": %q,"+
			"\"ISCSIInterface\": %q,"+
			"\"FSType\": %q,"+
			"\"ReadOnly\":\"%v\","+
			"\"Portals\":\"%v\","+
			"\"DiscoveryCHAPAuth\":\"%v\","+
			"\"SessionCHAPAuth\":\"%v\","+
			"\"SecretRef\": \"%v\","+
			"\"InitiatorName\":\"%v\""+
			"}",
			volume.Name, iscsi.TargetPortal, iscsi.IQN, iscsi.Lun, iscsi.ISCSIInterface, iscsi.FSType, iscsi.ReadOnly, iscsi.Portals, iscsi.DiscoveryCHAPAuth, iscsi.SessionCHAPAuth, iscsi.SecretRef, initiator)
	case volume.VolumeSource.Glusterfs != nil:
		glusterfs := volume.VolumeSource.Glusterfs
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"Glusterfs\","+
			"\"EndpointsName\": %q,"+
			"\"Path\": \"%v\","+
			"\"ReadOnly\": \"%v\""+
			"}", volume.Name, glusterfs.EndpointsName, glusterfs.Path, glusterfs.ReadOnly)
	case volume.VolumeSource.PersistentVolumeClaim != nil:
		claim := volume.VolumeSource.PersistentVolumeClaim
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"PersistentVolumeClaim\","+
			"\"ClaimName\": %q,"+
			"\"ReadOnly\": \"%v\""+
			"}", volume.Name, claim.ClaimName, claim.ReadOnly)
	case volume.VolumeSource.RBD != nil:
		rbd := volume.VolumeSource.RBD
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"RBD\","+
			"\"CephMonitors\": \"%v\","+
			"\"RBDImage\": %q,"+
			"\"FSType\": %q,"+
			"\"RBDPool\":\"%v\","+
			"\"RadosUser\":\"%v\","+
			"\"Keyring\":\"%v\","+
			"\"SecretRef\": \"%v\","+
			"\"ReadOnly\":\"%v\""+
			"}",
			volume.Name, rbd.CephMonitors, rbd.RBDImage, rbd.FSType, rbd.RBDPool, rbd.RadosUser, rbd.Keyring, rbd.SecretRef, rbd.ReadOnly)
	case volume.VolumeSource.Quobyte != nil:
		quobyte := volume.VolumeSource.Quobyte
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"Quobyte\","+
			"\"Registry\": %q,"+
			"\"Volume\": %q,"+
			"\"ReadOnly\": \"%v\""+
			"}", volume.Name, quobyte.Registry, quobyte.Volume, quobyte.ReadOnly)
	case volume.VolumeSource.DownwardAPI != nil:
		d := volume.VolumeSource.DownwardAPI
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"DownwardAPI\","+
			"\"Mappings\": \"%v\""+
			"}", volume.Name, d.Items)
	case volume.VolumeSource.AzureDisk != nil:
		d := volume.VolumeSource.AzureDisk
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"AzureDisk\","+
			"\"DiskName\": %q,"+
			"\"DiskURI\": %q,"+
			"\"Kind\": %q,"+
			"\"FSType\": %q,"+
			"\"CachingMode\": %q,"+
			"\"ReadOnly\": \"%v\","+
			"}", volume.Name, d.DiskName, d.DataDiskURI, *d.Kind, *d.FSType, *d.CachingMode, *d.ReadOnly)
	case volume.VolumeSource.VsphereVolume != nil:
		vsphere := volume.VolumeSource.VsphereVolume
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"vSphereVolume\","+
			"\"VolumePath\": \"%v\","+
			"\"FSType\": %q,"+
			"\"StoragePolicyName\": %q"+
			"}", volume.Name, vsphere.VolumePath, vsphere.FSType, vsphere.StoragePolicyName)
	case volume.VolumeSource.Cinder != nil:
		cinder := volume.VolumeSource.Cinder
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"Cinder\","+
			"\"VolumeID\": \"%v\","+
			"\"FSType\": %q,"+
			"\"ReadOnly\":\"%v\","+
			"\"SecretRef\": \"%v\""+
			"}",
			volume.Name, cinder.VolumeID, cinder.FSType, cinder.ReadOnly, cinder.SecretRef)
	case volume.VolumeSource.PhotonPersistentDisk != nil:
		photon := volume.VolumeSource.PhotonPersistentDisk
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"PhotonPersistentDisk\","+
			"\"PdID\": \"%v\","+
			"\"FSType\": %q"+
			"}", volume.Name, photon.PdID, photon.FSType)
	case volume.VolumeSource.PortworxVolume != nil:
		portworx := volume.VolumeSource.PortworxVolume
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"PortworxVolume\","+
			"\"VolumeID\": \"%v\""+
			"}", volume.Name, portworx.VolumeID)
	case volume.VolumeSource.ScaleIO != nil:
		sio := volume.VolumeSource.ScaleIO
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"ScaleIO\","+
			"\"Gateway\": %q,"+
			"\"System\": %q,"+
			"\"Projection Domain\": %q,"+
			"\"Storage Pool\": %q,"+
			"\"Storage Mode\": %q,"+
			"\"VolumeName\":\"%v\","+
			"\"FSType\": %q,"+
			"\"ReadOnly\":\"%v\""+
			"}",
			volume.Name, sio.Gateway, sio.System, sio.ProtectionDomain, sio.StoragePool, sio.StorageMode, sio.VolumeName, sio.FSType, sio.ReadOnly)
	case volume.VolumeSource.CephFS != nil:
		cephfs := volume.VolumeSource.CephFS
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"CephFS\","+
			"\"Monitors\": %q,"+
			"\"Path\": %q,"+
			"\"User\": %q,"+
			"\"SecretFile\": %q,"+
			"\"SecretRef\": \"%v\","+
			"\"ReadOnly\":\"%v\""+
			"}",
			volume.Name, cephfs.Monitors, cephfs.Path, cephfs.User, cephfs.SecretFile, cephfs.SecretRef, cephfs.ReadOnly)
	case volume.VolumeSource.StorageOS != nil:
		storageos := volume.VolumeSource.StorageOS
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"StorageOS\","+
			"\"VolumeName\": %q,"+
			"\"VolumeNamespace\": %q,"+
			"\"FSType\": %q,"+
			"\"ReadOnly\":\"%v\""+
			"}",
			volume.Name, storageos.VolumeName, storageos.VolumeNamespace, storageos.FSType, storageos.ReadOnly)
	case volume.VolumeSource.FC != nil:
		fc := volume.VolumeSource.FC
		lun := ValueNone
		if fc.Lun != nil {
			lun = strconv.Itoa(int(*fc.Lun))
		}
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"FC\","+
			"\"TargetWWNs\": %q,"+
			"\"LUN\": %q,"+
			"\"FSType\": %q,"+
			"\"ReadOnly\":\"%v\""+
			"}",
			volume.Name, strings.Join(fc.TargetWWNs, ", "), lun, fc.FSType, fc.ReadOnly)
	case volume.VolumeSource.AzureFile != nil:
		azureFile := volume.VolumeSource.AzureFile
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"AzureFile\","+
			"\"SecretName\": %q,"+
			"\"ShareName\": %q,"+
			"\"ReadOnly\":\"%v\""+
			"}",
			volume.Name, azureFile.SecretName, azureFile.ShareName, azureFile.ReadOnly)
	case volume.VolumeSource.FlexVolume != nil:
		flex := volume.VolumeSource.FlexVolume
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"FlexVolume\","+
			"\"Driver\": %q,"+
			"\"FSType\": %q,"+
			"\"SecretRef\": \"%v\","+
			"\"ReadOnly\":\"%v\","+
			"\"Options\": \"%v\""+
			"}",
			volume.Name, flex.Driver, flex.FSType, flex.SecretRef, flex.ReadOnly, flex.Options)
	case volume.VolumeSource.Flocker != nil:
		flocker := volume.VolumeSource.Flocker
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"Flocker\","+
			"\"DatasetName\": \"%v\","+
			"\"DatasetUUID\": %q"+
			"}", volume.Name, flocker.DatasetName, flocker.DatasetUUID)
	case volume.VolumeSource.Projected != nil:
		projected := volume.VolumeSource.Projected
		sources := "["
		for i, source := range projected.Sources {
			sources += "{"
			if source.Secret != nil {
				sources += "\"Type\": \"Secret\","
				sources += fmt.Sprintf("\"SecretName\": %q", source.Secret.Name)
			} else if source.DownwardAPI != nil {
				sources += "\"Type\": \"DownwardAPI\","
				sources += "\"DownwardAPI\": \"true\""
			} else if source.ConfigMap != nil {
				sources += "\"Type\": \"ConfigMap\","
				sources += fmt.Sprintf("\"ConfigMapName\": %q", source.ConfigMap.Name)
			} else if source.ServiceAccountToken != nil {
				sources += "\"Type\": \"ServiceAccountToken\","
				sources += fmt.Sprintf("\"TokenExpirationSeconds\": \"%v\"", source.ServiceAccountToken.ExpirationSeconds)
			}
			sources += "}"
			if i < len(projected.Sources)-1 {
				sources += ","
			}
		}
		sources += "]"
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"Projected\","+
			"\"sources\": %v"+
			"}", volume.Name, sources)
	case volume.VolumeSource.CSI != nil:
		csi := volume.VolumeSource.CSI
		var readOnly bool
		var fsType string
		if csi.ReadOnly != nil && *csi.ReadOnly {
			readOnly = true
		}
		if csi.FSType != nil {
			fsType = *csi.FSType
		}
		return fmt.Sprintf("{"+
			"\"Name\": %q,"+
			"\"Type\": \"CSI\","+
			"\"Driver\": %q,"+
			"\"FSType\": \"%v\","+
			"\"ReadOnly\": \"%v\""+
			"}", volume.Name, csi.Driver, fsType, readOnly)
	}

	return fmt.Sprintf("{\"name\": %q,\"Type\":\"<unknown>\"}", volume.Name)
}

func formatBytes(c int64) string {
	b := float64(c)

	switch {
	case c > 1<<40:
		return fmt.Sprintf("%.1f Ti", b/(1<<40))
	case c > 1<<30:
		return fmt.Sprintf("%.1f Gi", b/(1<<30))
	case c > 1<<20:
		return fmt.Sprintf("%.1f Mi", b/(1<<20))
	case c > 1<<10:
		return fmt.Sprintf("%.1f Ki", b/(1<<10))
	default:
		return fmt.Sprintf("%d B", c)
	}
}

type MongoDBNode struct {
	Replicas    int64                          `json:"replicas,omitempty"`
	PodTemplate ofst.PodTemplateSpec           `json:"podTemplate,omitempty"`
	Storage     core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}

func mongoDBResources(obj unstructured.Unstructured) (string, error) {
	totalCPU := int64(0)
	totalMemory := int64(0)
	totalStorage := int64(0)

	// Sharded MongoDB
	shardTopology, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "shardTopology")
	if err != nil {
		return "", err
	}
	if found && shardTopology != nil {
		// Shard nodes resources
		shard, err := getMongoDBNodeInfo(obj, "spec", "shardTopology", "shard")
		if err != nil {
			return "", err
		}
		totalCPU += shard.Replicas * shard.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue()
		totalMemory += shard.Replicas * shard.PodTemplate.Spec.Resources.Limits.Memory().Value()
		totalStorage += shard.Replicas * shard.Storage.Resources.Requests.Storage().Value()

		// ConfigServer nodes resources
		configServer, err := getMongoDBNodeInfo(obj, "spec", "shardTopology", "configServer")
		if err != nil {
			return "", err
		}
		totalCPU += configServer.Replicas * configServer.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue()
		totalMemory += configServer.Replicas * configServer.PodTemplate.Spec.Resources.Limits.Memory().Value()
		totalStorage += configServer.Replicas * configServer.Storage.Resources.Requests.Storage().Value()

		// Mongos node resources
		mongos, err := getMongoDBNodeInfo(obj, "spec", "shardTopology", "mongos")
		if err != nil {
			return "", err
		}
		totalCPU += mongos.Replicas * mongos.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue()
		totalMemory += mongos.Replicas * mongos.PodTemplate.Spec.Resources.Limits.Memory().Value()

		// Exporter resources
		cpu, memory, err := exporterResources(obj)
		if err != nil {
			return "", err
		}
		totalCPU += cpu
		totalMemory += memory

		return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
	}

	// MongoDB ReplicaSet
	replicaSet, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicaSet")
	if err != nil {
		return "", err
	}
	if found && replicaSet != nil {
		// ReplicaSet resources
		rs, err := getMongoDBNodeInfo(obj, "spec")
		if err != nil {
			return "", err
		}
		totalCPU += rs.Replicas * rs.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue()
		totalMemory += rs.Replicas * rs.PodTemplate.Spec.Resources.Limits.Memory().Value()
		totalStorage += rs.Replicas * rs.Storage.Resources.Requests.Storage().Value()

		// Exporter resources
		cpu, memory, err := exporterResources(obj)
		if err != nil {
			return "", err
		}
		totalCPU += cpu
		totalMemory += memory

		return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
	}

	// Standalone MongoDB
	db, err := getMongoDBNodeInfo(obj, "spec")
	if err != nil {
		return "", err
	}
	totalCPU += db.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue()
	totalMemory += db.PodTemplate.Spec.Resources.Limits.Memory().Value()
	totalStorage += db.Storage.Resources.Requests.Storage().Value()
	// Exporter resources
	cpu, memory, err := exporterResources(obj)
	if err != nil {
		return "", err
	}
	totalCPU += cpu
	totalMemory += memory

	return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
}

func getMongoDBNodeInfo(obj unstructured.Unstructured, fields ...string) (*MongoDBNode, error) {
	unstructuredNode, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), fields...)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("unable to find path: %s", strings.Join(fields, "."))
	}

	node := new(MongoDBNode)
	data, err := json.Marshal(unstructuredNode)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func exporterResources(obj unstructured.Unstructured) (int64, int64, error) {
	unstructuredExporter, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "monitor", "prometheus", "exporter")
	if err != nil {
		return 0, 0, nil
	}
	if found && unstructuredExporter != nil {
		exporter := new(mona.PrometheusExporterSpec)
		data, err := json.Marshal(unstructuredExporter)
		if err != nil {
			return 0, 0, err
		}
		err = json.Unmarshal(data, &exporter)
		if err != nil {
			return 0, 0, err
		}
		return exporter.Resources.Limits.Cpu().MilliValue(), exporter.Resources.Limits.Memory().Value(), nil
	}
	return 0, 0, nil
}
