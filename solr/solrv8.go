package solr

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/go-resty/resty/v2"
	"k8s.io/klog/v2"
)

type SLClientV8 struct {
	Client *resty.Client
	Config *Config
}

func (sc *SLClientV8) GetClusterStatus() (*Response, error) {
	sc.Config.log.V(5).Info("GETTING CLUSTER STATUS")
	req := sc.Client.R().SetDoNotParseResponse(true)
	req.SetQueryParam("action", "CLUSTERSTATUS")
	res, err := req.Get("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request")
		return nil, err
	}

	clusterResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return clusterResponse, nil
}

func (sc *SLClientV8) ListCollection() (*Response, error) {
	sc.Config.log.V(5).Info(fmt.Sprintf("SEARCHING COLLECTION: %s", writeCollectionName))
	req := sc.Client.R().SetDoNotParseResponse(true)
	req.SetQueryParam("action", "LIST")
	res, err := req.Get("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request while getting colection list")
		return nil, err
	}
	response := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return response, nil
}

func (sc *SLClientV8) CreateCollection() (*Response, error) {
	sc.Config.log.V(5).Info(fmt.Sprintf("CREATING COLLECTION: %s", writeCollectionName))
	req := sc.Client.R().SetDoNotParseResponse(true)
	req.SetHeader("Content-Type", "application/json")
	params := map[string]string{
		Action:            ActionCreate,
		Name:              writeCollectionName,
		NumShards:         "1",
		ReplicationFactor: "1",
	}

	req.SetQueryParams(params)
	res, err := req.Post("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to create a collection")
		return nil, err
	}

	collectionResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return collectionResponse, nil
}

func (sc *SLClientV8) WriteCollection() (*Response, error) {
	sc.Config.log.V(5).Info(fmt.Sprintf("WRITING COLLECTION: %s", writeCollectionName))
	req := sc.Client.R().SetDoNotParseResponse(true)
	req.SetHeader("Content-Type", "application/json")
	data1 := &Data{
		CommitWithin: 5000,
		Overwrite:    true,
		Doc: &Doc{
			Id: 1,
			DB: "elasticsearch",
		},
	}
	add := ADD{
		Add: data1,
	}
	req.SetBody(add)
	res, err := req.Post(fmt.Sprintf("/solr/%s/update", writeCollectionName))
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to add document in collect")
		return nil, err
	}

	writeResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return writeResponse, nil
}

func (sc *SLClientV8) ReadCollection() (*Response, error) {
	sc.Config.log.V(5).Info(fmt.Sprintf("READING COLLECTION: %s", writeCollectionName))
	req := sc.Client.R().SetDoNotParseResponse(true)
	req.SetQueryParam("q", "*:*")
	res, err := req.Get(fmt.Sprintf("/solr/%s/select", writeCollectionName))
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to read a collection")
		return nil, err
	}

	writeResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return writeResponse, nil
}

func (sc *SLClientV8) BackupCollection(ctx context.Context, collection string, backupName string, location string, repository string) (*Response, error) {
	sc.Config.log.V(5).Info(fmt.Sprintf("BACKUP COLLECTION: %s", collection))
	req := sc.Client.R().SetDoNotParseResponse(true).SetContext(ctx)
	req.SetHeader("Content-Type", "application/json")
	params := map[string]string{
		Action:     ActionBackup,
		Name:       backupName,
		Collection: collection,
		Location:   location,
		Repository: repository,
		Async:      fmt.Sprintf("%s-backup", collection),
	}

	req.SetQueryParams(params)

	res, err := req.Post("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to backup a collection")
		return nil, err
	}

	backupResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return backupResponse, nil
}

func (sc *SLClientV8) RestoreCollection(ctx context.Context, collection string, backupName string, location string, repository string, backupId int) (*Response, error) {
	sc.Config.log.V(5).Info(fmt.Sprintf("RESTORE COLLECTION: %s", collection))
	req := sc.Client.R().SetDoNotParseResponse(true).SetContext(ctx)
	req.SetHeader("Content-Type", "application/json")
	params := map[string]string{
		Action:     ActionRestore,
		Name:       backupName,
		Location:   location,
		Collection: collection,
		Repository: repository,
		BackupId:   strconv.Itoa(backupId),
		Async:      fmt.Sprintf("%s-restore", collection),
	}

	req.SetQueryParams(params)

	res, err := req.Post("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to restore a collection")
		return nil, err
	}

	backupResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return backupResponse, nil
}

func (sc *SLClientV8) FlushStatus(asyncId string) (*Response, error) {
	sc.Config.log.V(5).Info("Flush Status")
	req := sc.Client.R().SetDoNotParseResponse(true)
	req.SetHeader("Content-Type", "application/json")

	params := map[string]string{
		Action:    DeleteStatus,
		RequestId: asyncId,
	}
	req.SetQueryParams(params)
	res, err := req.Get("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to flush status")
		return nil, err
	}

	backupResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return backupResponse, nil
}

func (sc *SLClientV8) RequestStatus(asyncId string) (*Response, error) {
	sc.Config.log.V(5).Info("Request Status")
	req := sc.Client.R().SetDoNotParseResponse(true)
	req.SetHeader("Content-Type", "application/json")
	params := map[string]string{
		Action:    RequestStatus,
		RequestId: asyncId,
	}
	req.SetQueryParams(params)
	res, err := req.Get("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to request status")
		return nil, err
	}
	backupResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return backupResponse, nil
}

func (sc *SLClientV8) DeleteBackup(ctx context.Context, backupName string, collection string, location string, repository string, backupId int, snap string) (*Response, error) {
	sc.Config.log.V(5).Info(fmt.Sprintf("DELETE BACKUP ID %d of BACKUP %s", backupId, backupName))
	req := sc.Client.R().SetDoNotParseResponse(true).SetContext(ctx)
	req.SetHeader("Content-Type", "application/json")
	async := fmt.Sprintf("%s-delete", collection)
	if snap != "" {
		async = fmt.Sprintf("%s-%s", async, snap)
	}
	params := map[string]string{
		Action:     ActionDeleteBackup,
		Name:       backupName,
		Location:   location,
		Repository: repository,
		BackupId:   strconv.Itoa(backupId),
		Async:      async,
	}
	req.SetQueryParams(params)

	res, err := req.Delete("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to restore a collection")
		return nil, err
	}

	backupResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return backupResponse, nil
}

func (sc *SLClientV8) PurgeBackup(ctx context.Context, backupName string, collection string, location string, repository string, snap string) (*Response, error) {
	sc.Config.log.V(5).Info(fmt.Sprintf("PURGE BACKUP ID %s", backupName))
	req := sc.Client.R().SetDoNotParseResponse(true).SetContext(ctx)
	req.SetHeader("Content-Type", "application/json")
	async := fmt.Sprintf("%s-purge", collection)
	if snap != "" {
		async = fmt.Sprintf("%s-%s", async, snap)
	}
	params := map[string]string{
		Action:      ActionDeleteBackup,
		Name:        backupName,
		Location:    location,
		Repository:  repository,
		PurgeUnused: "true",
		Async:       async,
	}
	req.SetQueryParams(params)
	res, err := req.Put("/solr/admin/collections")
	if err != nil {
		sc.Config.log.Error(err, "Failed to send http request to restore a collection")
		return nil, err
	}

	backupResponse := &Response{
		Code:   res.StatusCode(),
		header: res.Header(),
		body:   res.RawBody(),
	}
	return backupResponse, nil
}

func (sc *SLClientV8) GetConfig() *Config {
	return sc.Config
}

func (sc *SLClientV8) GetClient() *resty.Client {
	return sc.Client
}

func (sc *SLClientV8) GetLog() logr.Logger {
	return sc.Config.log
}

func (sc *SLClientV8) DecodeBackupResponse(data map[string]interface{}, collection string) ([]byte, error) {
	sc.Config.log.V(5).Info("Decode Backup Data")
	backupResponse, ok := data["response"].([]interface{})
	if !ok {
		err := fmt.Errorf("didn't find status for collection %s\n", collection)
		return nil, err
	}
	mp := make(map[string]interface{})
	for i := 0; i < len(backupResponse); i += 2 {
		a := backupResponse[i].(string)
		b := backupResponse[i+1]
		mp[a] = b
	}
	b, err := json.Marshal(mp)
	if err != nil {
		klog.Error(fmt.Sprintf("Could not format response for collection %s into json", collection))
		return nil, err
	}
	return b, nil

}
