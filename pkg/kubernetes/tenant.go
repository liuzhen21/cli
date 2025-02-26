package kubernetes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	terrors "github.com/tkeel-io/kit/errors"
	"github.com/tkeel-io/kit/result"
	tenantApi "github.com/tkeel-io/tkeel/api/tenant/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_listTenantsMethodFormat  = "apis/rudder/v1/tenants"
	_createTenantMethodFormat = "apis/rudder/v1/tenants"
	_deleteTenantMethodFormat = "apis/rudder/v1/tenants/%s"
	_infoTenantMethodFormat   = "apis/rudder/v1/tenants/%s"

	_listTenantPluginsMethodFormat  = "apis/rudder/v1/tenants/%s/plugins"
	_createTenantPluginMethodFormat = "apis/rudder/v1/tenants/%s/plugins"
	_deleteTenantPluginMethodFormat = "apis/rudder/v1/tenants/%s/plugins/%s"
	_infoTenantPluginMethodFormat   = "apis/rudder/v1/tenants/%s/plugins/%s"
)

type TenantListOutPut struct {
	ID     string `csv:"ID"`
	Title  string `csv:"TITLE"`
	Remark string `csv:"REMARK"`
}

func TenantCreate(title, remark, adminName, adminPW string) error {
	if len(title) == 0 {
		return errors.New("title param nil")
	}
	tenant := &TenantCreateIn{Title: title, Remark: remark}
	if len(adminName) != 0 && len(adminPW) != 0 {
		admin := TenantAdmin{UserName: adminName, Password: adminPW}
		tenant.Admin = admin
	}

	return CreateTenant(tenant)
}

func CreateTenant(tenant *TenantCreateIn) error {
	token, err := getAdminToken()
	if err != nil {
		return err
	}
	method := _createTenantMethodFormat

	data, err := json.Marshal(tenant) //nolint
	if err != nil {
		return errors.Wrap(err, "marshal plugin request failed")
	}
	resp, err := InvokeByPortForward(_pluginKeel, method, data, http.MethodPost, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "invoke "+method+" error")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "can't unmarshal'")
	}

	if r.Code != terrors.Success.Reason {
		return errors.New("response error: " + r.Msg)
	}

	return nil
}

func TenantList() ([]TenantListOutPut, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting admin token")
	}

	resp, err := InvokeByPortForward(_pluginKeel, _listTenantsMethodFormat, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return nil, errors.New("response error: " + r.Msg)
	}

	listResponse := tenantApi.ListTenantResponse{}
	err = r.Data.UnmarshalTo(&listResponse)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	var list = make([]TenantListOutPut, 0, len(listResponse.Tenants))
	for _, tenant := range listResponse.Tenants {
		list = append(list, TenantListOutPut{tenant.TenantId, tenant.Title, tenant.Remark})
	}
	return list, nil
}

func TenantInfo(tenantId string) ([]TenantListOutPut, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting admin token")
	}
	method := fmt.Sprintf(_infoTenantMethodFormat, tenantId)

	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return nil, errors.New("response error: " + r.Msg)
	}

	tenantResponse := tenantApi.GetTenantResponse{}
	err = r.Data.UnmarshalTo(&tenantResponse)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	var list = make([]TenantListOutPut, 0, 1)
	list = append(list, TenantListOutPut{tenantResponse.TenantId, tenantResponse.Title, tenantResponse.Remark})
	return list, nil
}

func TenantDelete(tenantId string) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "error getting admin token")
	}
	method := fmt.Sprintf(_deleteTenantMethodFormat, tenantId)

	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodDelete, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return errors.New("response error: " + r.Msg)
	}
	return nil
}

// tenant plugin manage
func TenantPluginEnable(tenantId, pluginId string) error {

	return nil
}

func TenantPluginDisable(tenantId, pluginId string) error {

	return nil
}

func TenantPluginList(tenantId string) error {

	return nil
}

type TenantCreateIn struct {
	Title  string      `json:"title"`
	Remark string      `json:"remark"`
	Admin  TenantAdmin `json:"admin"`
}

type TenantAdmin struct {
	UserName string `json:"username"` //nolint
	Password string `json:"password"` //nolint
}
