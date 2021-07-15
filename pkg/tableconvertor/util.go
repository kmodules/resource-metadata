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
	ValueNone                 = "<none>"
	ResourceKindRedis         = "Redis"
	ResourceKindMySQL         = "MySQL"
	ResourceKindMariaDB       = "MariaDB"
	ResourceKindMongoDB       = "MongoDB"
	ResourceKindPostgres      = "Postgres"
	ResourceKindElasticsearch = "Elasticsearch"
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

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

type DBNode struct {
	Replicas    int64                          `json:"replicas,omitempty"`
	PodTemplate ofst.PodTemplateSpec           `json:"podTemplate,omitempty"`
	Storage     core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}

func mongoDBResources(obj unstructured.Unstructured) (string, error) {
	totalCPU := int64(0)
	totalMemory := int64(0)
	totalStorage := int64(0)

	// Sharded MongoDB
	shardTopology, found, err := unstructured.NestedFieldNoCopy(obj.UnstructuredContent(), "spec", "shardTopology")
	if err != nil {
		return "", err
	}
	if found && shardTopology != nil {
		// Shard nodes resources
		shard, err := getDBNodeInfo(obj, "spec", "shardTopology", "shard")
		if err != nil {
			return "", err
		}
		totalCPU += shard.Replicas * shard.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue()
		totalMemory += shard.Replicas * shard.PodTemplate.Spec.Resources.Limits.Memory().Value()
		totalStorage += shard.Replicas * shard.Storage.Resources.Requests.Storage().Value()

		// ConfigServer nodes resources
		configServer, err := getDBNodeInfo(obj, "spec", "shardTopology", "configServer")
		if err != nil {
			return "", err
		}
		totalCPU += configServer.Replicas * configServer.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue()
		totalMemory += configServer.Replicas * configServer.PodTemplate.Spec.Resources.Limits.Memory().Value()
		totalStorage += configServer.Replicas * configServer.Storage.Resources.Requests.Storage().Value()

		// Mongos node resources
		mongos, err := getDBNodeInfo(obj, "spec", "shardTopology", "mongos")
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
	replicaSet, found, err := unstructured.NestedFieldNoCopy(obj.UnstructuredContent(), "spec", "replicaSet")
	if err != nil {
		return "", err
	}
	if found && replicaSet != nil {
		// ReplicaSet resources
		rs, err := getDBNodeInfo(obj, "spec")
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
	db, err := getDBNodeInfo(obj, "spec")
	if err != nil {
		return "", err
	}
	totalCPU += db.Replicas * max(db.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue(), db.PodTemplate.Spec.Resources.Requests.Cpu().MilliValue())
	totalMemory += db.Replicas * max(db.PodTemplate.Spec.Resources.Limits.Memory().Value(), db.PodTemplate.Spec.Resources.Requests.Memory().Value())
	totalStorage += db.Replicas * db.Storage.Resources.Requests.Storage().Value()
	// Exporter resources
	cpu, memory, err := exporterResources(obj)
	if err != nil {
		return "", err
	}
	totalCPU += cpu
	totalMemory += memory

	return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
}

func getMongoDBReplicas(obj unstructured.Unstructured) (string, error) {
	// Sharded MongoDB cluster
	shardTopology, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "shardTopology")
	if err != nil {
		return "", err
	}
	if found && shardTopology != nil {
		shards, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "shardTopology", "shard", "replicas")
		if err != nil {
			return "", err
		}
		configServers, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "shardTopology", "configServer", "replicas")
		if err != nil {
			return "", err
		}
		mongos, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "shardTopology", "mongos", "replicas")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v, %v, %v", shards, configServers, mongos), nil
	}
	// MongoDB ReplicaSet
	replicaSet, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicaSet")
	if err != nil {
		return "", err
	}
	if found && replicaSet != nil {
		replicas, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicas")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", replicas), nil
	}
	// Standalone MongoDB
	return "1", nil
}

func getDBNodeInfo(obj unstructured.Unstructured, fields ...string) (*DBNode, error) {
	unstructuredNode, found, err := unstructured.NestedFieldNoCopy(obj.UnstructuredContent(), fields...)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("unable to find path: %s", strings.Join(fields, "."))
	}

	node := new(DBNode)
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
	unstructuredExporter, found, err := unstructured.NestedFieldNoCopy(obj.UnstructuredContent(), "spec", "monitor", "prometheus", "exporter")
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
		return max(exporter.Resources.Limits.Cpu().MilliValue(), exporter.Resources.Requests.Cpu().MilliValue()),
			max(exporter.Resources.Limits.Memory().Value(), exporter.Resources.Requests.Memory().Value()), nil
	}
	return 0, 0, nil
}

func postgresResources(obj unstructured.Unstructured) (string, error) {
	totalCPU := int64(0)
	totalMemory := int64(0)
	totalStorage := int64(0)

	pg, err := getDBNodeInfo(obj, "spec")
	if err != nil {
		return "", err
	}
	totalCPU += pg.Replicas * max(pg.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue(), pg.PodTemplate.Spec.Resources.Requests.Cpu().MilliValue())
	totalMemory += pg.Replicas * max(pg.PodTemplate.Spec.Resources.Limits.Memory().Value(), pg.PodTemplate.Spec.Resources.Requests.Memory().Value())
	totalStorage += pg.Replicas * pg.Storage.Resources.Requests.Storage().Value()

	// Exporter resources
	cpu, memory, err := exporterResources(obj)
	if err != nil {
		return "", err
	}
	totalCPU += cpu
	totalMemory += memory

	return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
}

type ElasticsearchNode struct {
	Replicas  int64                          `json:"replicas,omitempty"`
	Storage   core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
	Resources core.ResourceRequirements      `json:"resources,omitempty"`

	totalCPU     int64
	totalMemory  int64
	totalStorage int64
}

func elasticsearchDBResources(obj unstructured.Unstructured) (string, error) {
	totalCPU := int64(0)
	totalMemory := int64(0)
	totalStorage := int64(0)

	// Topology
	topology, found, err := unstructured.NestedFieldNoCopy(obj.UnstructuredContent(), "spec", "topology")
	if err != nil {
		return "", err
	}
	if found && topology != nil {
		topology := topology.(map[string]interface{})
		// Elasticsearch master node
		master, _, err := getElasticsearchNodeInfo(topology, "master")
		if err != nil {
			return "", err
		}
		totalCPU += master.totalCPU
		totalMemory += master.totalMemory
		totalStorage += master.totalStorage

		// Elasticsearch ingest node
		ingest, _, err := getElasticsearchNodeInfo(topology, "ingest")
		if err != nil {
			return "", err
		}
		totalCPU += ingest.Replicas * ingest.Resources.Limits.Cpu().MilliValue()
		totalMemory += ingest.Replicas * ingest.Resources.Limits.Memory().Value()
		totalStorage += ingest.Replicas * ingest.Storage.Resources.Requests.Storage().Value()

		// Elasticsearch data node
		data, found, err := getElasticsearchNodeInfo(topology, "data")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += data.totalCPU
			totalMemory += data.totalMemory
			totalStorage += data.totalStorage
		}

		// Elasticsearch dataContent node
		dataContent, found, err := getElasticsearchNodeInfo(topology, "dataContent")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += dataContent.totalCPU
			totalMemory += dataContent.totalMemory
			totalStorage += dataContent.totalStorage
		}

		// Elasticsearch dataHot node
		dataHot, found, err := getElasticsearchNodeInfo(topology, "dataHot")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += dataHot.totalCPU
			totalMemory += dataHot.totalMemory
			totalStorage += dataHot.totalStorage
		}

		// Elasticsearch dataWarm node
		dataWarm, found, err := getElasticsearchNodeInfo(topology, "dataWarm")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += dataWarm.totalCPU
			totalMemory += dataWarm.totalMemory
			totalStorage += dataWarm.totalStorage
		}

		// Elasticsearch dataCold node
		dataCold, found, err := getElasticsearchNodeInfo(topology, "dataCold")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += dataCold.totalCPU
			totalMemory += dataCold.totalMemory
			totalStorage += dataCold.totalStorage
		}

		// Elasticsearch dataFrozen node
		dataFrozen, found, err := getElasticsearchNodeInfo(topology, "dataFrozen")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += dataFrozen.totalCPU
			totalMemory += dataFrozen.totalMemory
			totalStorage += dataFrozen.totalStorage
		}

		// Elasticsearch ml node
		ml, found, err := getElasticsearchNodeInfo(topology, "ml")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += ml.totalCPU
			totalMemory += ml.totalMemory
			totalStorage += ml.totalStorage
		}

		// Elasticsearch transform node
		transform, found, err := getElasticsearchNodeInfo(topology, "transform")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += transform.totalCPU
			totalMemory += transform.totalMemory
			totalStorage += transform.totalStorage
		}

		// Elasticsearch coordinating node
		coordinating, found, err := getElasticsearchNodeInfo(topology, "coordinating")
		if err != nil {
			return "", err
		}
		if found {
			totalCPU += coordinating.totalCPU
			totalMemory += coordinating.totalMemory
			totalStorage += coordinating.totalStorage
		}

		// Exporter resources
		cpu, memory, err := exporterResources(obj)
		if err != nil {
			return "", err
		}
		totalCPU += cpu
		totalMemory += memory

		return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
	}

	// Combined Elasticsearch
	db, err := getDBNodeInfo(obj, "spec")
	if err != nil {
		return "", err
	}
	totalCPU += db.Replicas * max(db.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue(), db.PodTemplate.Spec.Resources.Requests.Cpu().MilliValue())
	totalMemory += db.Replicas * max(db.PodTemplate.Spec.Resources.Limits.Memory().Value(), db.PodTemplate.Spec.Resources.Requests.Memory().MilliValue())
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

func getElasticsearchNodeInfo(obj map[string]interface{}, fields ...string) (*ElasticsearchNode, bool, error) {
	unstructuredNode, found, err := unstructured.NestedFieldNoCopy(obj, fields...)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	node := new(ElasticsearchNode)
	data, err := json.Marshal(unstructuredNode)
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal(data, &node)
	if err != nil {
		return nil, false, err
	}

	node.totalCPU += node.Replicas * max(node.Resources.Limits.Cpu().MilliValue(), node.Resources.Requests.Cpu().MilliValue())
	node.totalMemory += node.Replicas * max(node.Resources.Limits.Memory().Value(), node.Resources.Requests.Memory().Value())
	node.totalStorage += node.Replicas * node.Storage.Resources.Requests.Storage().Value()

	return node, true, nil
}

func getElasticsearchReplicas(obj unstructured.Unstructured) (string, error) {
	topology, found, err := unstructured.NestedMap(obj.UnstructuredContent(), "spec", "topology")
	if err != nil {
		return "", err
	}
	if found && topology != nil {
		var replicas []string

		master, _, err := unstructured.NestedInt64(topology, "master", "replicas")
		if err != nil {
			return "", err
		}

		replicas = append(replicas, fmt.Sprintf("%q: %q", "m", strconv.FormatInt(master, 10)))

		ingest, _, err := unstructured.NestedInt64(topology, "ingest", "replicas")
		if err != nil {
			return "", err
		}
		replicas = append(replicas, fmt.Sprintf("%q: %q", "i", strconv.FormatInt(ingest, 10)))

		data, found, err := unstructured.NestedMap(topology, "data")
		if err != nil {
			return "", err
		}
		if found && data != nil {
			data, _, err := unstructured.NestedInt64(data, "replicas")
			if err != nil {
				return "", err
			}
			replicas = append(replicas, fmt.Sprintf("%q: %q", "d", strconv.FormatInt(data, 10)))
		}

		dataContent, found, err := unstructured.NestedMap(topology, "dataContent")
		if err != nil {
			return "", err
		}
		if found && dataContent != nil {
			dataContent, _, err := unstructured.NestedInt64(dataContent, "replicas")
			if err != nil {
				return "", err
			}
			replicas = append(replicas, fmt.Sprintf("%q: %q", "s", strconv.FormatInt(dataContent, 10)))
		}

		dataHot, found, err := unstructured.NestedMap(topology, "dataHot")
		if err != nil {
			return "", err
		}
		if found && dataHot != nil {
			dataHot, _, err := unstructured.NestedInt64(dataHot, "replicas")
			if err != nil {
				return "", err
			}
			replicas = append(replicas, fmt.Sprintf("%q: %q", "h", strconv.FormatInt(dataHot, 10)))
		}

		dataWarm, found, err := unstructured.NestedMap(topology, "dataWarm")
		if err != nil {
			return "", err
		}
		if found && dataWarm != nil {
			dataWarm, _, err := unstructured.NestedInt64(dataWarm, "replicas")
			if err != nil {
				return "", err
			}
			replicas = append(replicas, fmt.Sprintf("%q: %q", "w", strconv.FormatInt(dataWarm, 10)))
		}

		dataCold, found, err := unstructured.NestedMap(topology, "dataCold")
		if err != nil {
			return "", err
		}
		if found && dataCold != nil {
			dataCold, _, err := unstructured.NestedInt64(dataCold, "replicas")
			if err != nil {
				return "", err
			}
			replicas = append(replicas, fmt.Sprintf("%q: %q", "c", strconv.FormatInt(dataCold, 10)))
		}

		dataFrozen, found, err := unstructured.NestedMap(topology, "dataFrozen")
		if err != nil {
			return "", err
		}
		if found && dataFrozen != nil {
			dataFrozen, _, err := unstructured.NestedInt64(dataFrozen, "replicas")
			if err != nil {
				return "", err
			}
			replicas = append(replicas, fmt.Sprintf("%q: %q", "f", strconv.FormatInt(dataFrozen, 10)))
		}

		ml, found, err := unstructured.NestedMap(topology, "ml")
		if err != nil {
			return "", err
		}
		if found && ml != nil {
			ml, _, err := unstructured.NestedInt64(ml, "replicas")
			if err != nil {
				return "", err
			}
			replicas = append(replicas, fmt.Sprintf("%q: %q", "lr", strconv.FormatInt(ml, 10)))
		}
		transform, found, err := unstructured.NestedMap(topology, "transform")
		if err != nil {
			return "", err
		}
		if found && transform != nil {
			transform, _, err := unstructured.NestedInt64(transform, "replicas")
			if err != nil {
				return "", err
			}
			replicas = append(replicas, fmt.Sprintf("%q: %q", "rt", strconv.FormatInt(transform, 10)))
		}

		return fmt.Sprintf("{%v}", strings.Join(replicas, ", ")), nil
	}

	// Combined mode
	replicas, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicas")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", replicas), nil
}

func mariaDBResources(obj unstructured.Unstructured) (string, error) {
	totalCPU := int64(0)
	totalMemory := int64(0)
	totalStorage := int64(0)

	maria, err := getDBNodeInfo(obj, "spec")
	if err != nil {
		return "", err
	}
	totalCPU += maria.Replicas * max(maria.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue(), maria.PodTemplate.Spec.Resources.Requests.Cpu().MilliValue())
	totalMemory += maria.Replicas * max(maria.PodTemplate.Spec.Resources.Limits.Memory().Value(), maria.PodTemplate.Spec.Resources.Requests.Memory().Value())
	totalStorage += maria.Replicas * maria.Storage.Resources.Requests.Storage().Value()

	// Exporter resources
	cpu, memory, err := exporterResources(obj)
	if err != nil {
		return "", err
	}
	totalCPU += cpu
	totalMemory += memory

	return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
}

func mySQLResources(obj unstructured.Unstructured) (string, error) {
	totalCPU := int64(0)
	totalMemory := int64(0)
	totalStorage := int64(0)

	sql, err := getDBNodeInfo(obj, "spec")
	if err != nil {
		return "", err
	}
	totalCPU += sql.Replicas * max(sql.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue(), sql.PodTemplate.Spec.Resources.Requests.Cpu().MilliValue())
	totalMemory += sql.Replicas * max(sql.PodTemplate.Spec.Resources.Limits.Memory().Value(), sql.PodTemplate.Spec.Resources.Requests.Memory().Value())
	totalStorage += sql.Replicas * sql.Storage.Resources.Requests.Storage().Value()

	innoRouter, err := getMYSQLNodeInfo(obj)
	if err != nil {
		return "", err
	}
	if innoRouter != nil {
		totalCPU += innoRouter.Replicas * max(innoRouter.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue(), innoRouter.PodTemplate.Spec.Resources.Requests.Cpu().MilliValue())
		totalMemory += innoRouter.Replicas * max(innoRouter.PodTemplate.Spec.Resources.Limits.Memory().Value(), innoRouter.PodTemplate.Spec.Resources.Requests.Memory().Value())
	}

	// Exporter resources
	cpu, memory, err := exporterResources(obj)
	if err != nil {
		return "", err
	}
	totalCPU += cpu
	totalMemory += memory

	return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
}

func getMYSQLNodeInfo(obj unstructured.Unstructured) (*DBNode, error) {
	topology, found, err := unstructured.NestedMap(obj.UnstructuredContent(), "spec", "topology", "innoDBCluster", "router", "podTemplate")
	if err != nil {
		return nil, err
	}
	if found && topology != nil {
		mode, found, err := unstructured.NestedString(topology, "mode")
		if err != nil {
			return nil, err
		}
		// Only InnoDBCluster has dedicated Resource
		if found && mode == "InnoDBCluster" {
			inno, found, err := unstructured.NestedMap(topology, "innoDBCluster")
			if err != nil {
				return nil, err
			}
			if found && inno != nil {
				router, found, err := unstructured.NestedFieldNoCopy(inno, "router")
				if err != nil {
					return nil, err
				}

				if found && router != nil {

				}
				node := new(DBNode)
				data, err := json.Marshal(router)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(data, &node)
				if err != nil {
					return nil, err
				}
				return node, nil
			}
		}
	}
	return nil, nil
}

func getMySQLReplicas(obj unstructured.Unstructured) (string, error) {
	topology, found, err := unstructured.NestedMap(obj.UnstructuredContent(), "spec", "topology")
	if err != nil {
		return "", err
	}

	//MySQLClusterModeGroup  MySQLClusterMode = "GroupReplication"
	//InnoDBClusterModeGroup MySQLClusterMode = "InnoDBCluster"
	if found && topology != nil {
		mode, found, err := unstructured.NestedString(topology, "mode")
		if err != nil {
			return "", err
		}
		// Only InnoDBCluster has dedicated replica
		if found && mode == "InnoDBCluster" {
			inno, found, err := unstructured.NestedMap(topology, "innoDBCluster")
			if err != nil {
				return "", err
			}
			if found && inno != nil {
				replica, found, err := unstructured.NestedFieldCopy(inno, "router", "replica")
				if err != nil {
					return "", err
				}
				if found && replica != nil {
					return fmt.Sprintf("%v", replica), nil
				}
			}
		}
	}

	// Standalone or GroupReplication Mode
	replicas, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicas")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", replicas), nil
}

func redisResources(obj unstructured.Unstructured) (string, error) {
	totalCPU := int64(0)
	totalMemory := int64(0)
	totalStorage := int64(0)

	redis, err := getDBNodeInfo(obj, "spec")
	if err != nil {
		return "", err
	}

	// If mode is Cluster, replica number would be
	// summation of master and replicas from RedisClusterSpec
	mode, found, err := unstructured.NestedString(obj.UnstructuredContent(), "spec", "mode")
	if err != nil {
		return "", err
	}
	if found && mode == "Cluster" {
		cluster, found, err := unstructured.NestedMap(obj.UnstructuredContent(), "spec", "cluster")
		if err != nil {
			return "", err
		}
		if found && cluster != nil {
			master, _, err := unstructured.NestedInt64(cluster, "master")
			if err != nil {
				return "", err
			}
			replicas, _, err := unstructured.NestedInt64(cluster, "replicas")
			if err != nil {
				return "", err
			}

			redis.Replicas = master + replicas
		} else {
			return "", fmt.Errorf("failed to get redis resources. Reason: cluster mode not found")
		}
	}

	totalCPU += redis.Replicas * max(redis.PodTemplate.Spec.Resources.Limits.Cpu().MilliValue(), redis.PodTemplate.Spec.Resources.Requests.Cpu().MilliValue())
	totalMemory += redis.Replicas * max(redis.PodTemplate.Spec.Resources.Limits.Memory().Value(), redis.PodTemplate.Spec.Resources.Requests.Memory().Value())
	totalStorage += redis.Replicas * redis.Storage.Resources.Requests.Storage().Value()

	// Exporter resources
	cpu, memory, err := exporterResources(obj)
	if err != nil {
		return "", err
	}
	totalCPU += cpu
	totalMemory += memory

	return fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", core.ResourceCPU, fmt.Sprintf("%dm", totalCPU), core.ResourceMemory, formatBytes(totalMemory), core.ResourceStorage, formatBytes(totalStorage)), nil
}

func getRedisReplicas(obj unstructured.Unstructured) (string, error) {
	mode, found, err := unstructured.NestedString(obj.UnstructuredContent(), "spec", "mode")
	if err != nil {
		return "", err
	}
	if found && mode == "Cluster" {
		cluster, found, err := unstructured.NestedMap(obj.UnstructuredContent(), "spec", "cluster")
		if err != nil {
			return "", err
		}
		if found && cluster != nil {
			master, _, err := unstructured.NestedInt64(cluster, "master")
			if err != nil {
				return "", err
			}
			replicas, _, err := unstructured.NestedInt64(cluster, "replicas")
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%v, %v", master, replicas), nil
		} else {
			return "", fmt.Errorf("failed to detect replica number. Reason: cluster mode not found")
		}
	}

	// Standalone mode
	replicas, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicas")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", replicas), nil
}
