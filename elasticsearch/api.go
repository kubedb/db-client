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
	core "k8s.io/api/core/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
)

type ESClient interface {
	ClusterHealthInfo() (map[string]interface{}, error)
	NodesStats() (map[string]interface{}, error)
	GetIndicesInfo() ([]interface{}, error)
	ClusterStatus() (string, error)
	SyncCredentialFromSecret(secret *core.Secret) error
	GetClusterWriteStatus(ctx context.Context, db *api.Elasticsearch) error
}
