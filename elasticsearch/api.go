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

package elasticsearch

import (
	"context"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
)

var (
	writeRequestIndex       = "kubedb-system"
	writeRequestID          = "info"
	writeRequestType        = "_doc"
	CustomUser              = "custom-user"
	PrivilegeReadKey        = "read"
	PrivilegeWriteKey       = "write"
	PrivilegeCreateIndexKey = "create-index"
	Any                     = "*"
	All                     = "all"
	Names                   = "names"
	Privileges              = "privileges"
	AllowRestrictedIndices  = "allow_restricted_indices"
	Enabled                 = "enabled"
	Cluster                 = "cluster"
	Indices                 = "indices"
	Applications            = "applications"
	Application             = "application"
	RunAs                   = "run_as"
	TransientMetadata       = "transient_metadata"
	Resources               = "resources"
	Kibana                  = "kibana-.kibana"
)

type WriteRequestIndex struct {
	Index WriteRequestIndexBody `json:"index"`
}

type WriteRequestIndexBody struct {
	ID   string `json:"_id"`
	Type string `json:"_type,omitempty"`
}

type ESClient interface {
	ClusterHealthInfo() (map[string]interface{}, error)
	NodesStats() (map[string]interface{}, error)
	GetIndicesInfo() ([]interface{}, error)
	ClusterStatus() (string, error)
	SyncCredentialFromSecret(secret *core.Secret) error
	GetClusterWriteStatus(ctx context.Context, db *api.Elasticsearch) error
	GetClusterReadStatus(ctx context.Context, db *api.Elasticsearch) error
	GetTotalDiskUsage(ctx context.Context) (string, error)
	EnsureDBUserRole(ctx context.Context) error
}
