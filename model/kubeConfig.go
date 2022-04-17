package model

import (
	"errors"

	"gorm.io/gorm"
)

const (
	K8sAuthTypeUNSAFE = iota
	K8sAuthTypeBASIC
	K8sAuthTypeTLS
	K8sAuthTypeTOKEN
	K8sAuthInCLUSTER
)

type KubeConfig struct {
	Id                       int64  `gorm:"primary_key;column:id"`
	Cid                      string `gorm:"column:cid"`
	Description              string `gorm:"column:description"`
	K8sHost                  string `gorm:"column:k8s_host"`
	K8sAuthType              int    `gorm:"column:k8s_auth_type"`
	K8sAuthBasic             string `gorm:"column:k8s_auth_basic"`
	K8sAuthToken             string `gorm:"column:k8s_auth_token"`
	K8sClusterAuthData       string `gorm:"column:k8s_cluster_auth_data"`
	K8sClientCertificateData string `gorm:"column:k8s_client_certificate_data"`
	K8sClientKeyData         string `gorm:"column:k8s_client_key_data"`
	Status                   int64  `gorm:"column:status"`
	KialiPath                string `gorm:"column:kiali_path;default:'/api/v1/namespaces/istio-system/services/kiali:http/proxy/kiali/api'"`
	JaegerPath               string `gorm:"column:jaeger_path;default:'/api/v1/namespaces/istio-system/services/tracing:http-query/proxy/jaeger/api'"`
	CreateTime               int64  `gorm:"column:create_time"`
	UpdateTime               int64  `gorm:"column:update_time"`
}

var KubeConfigNoExistErr = errors.New("kube-config no exist")

func (k *KubeConfig) TableName() string {
	return KubeConfigTableName
}

func (k *KubeConfig) ListKubeConfig() (*[]KubeConfig, error) {
	whereScopes := func(db *gorm.DB) *gorm.DB {
		return db.Where(map[string]interface{}{})
	}
	projects, err := NewDataModel().GetList(NewKubeConfigModel(), whereScopes, []string{"*"})
	if err != nil {
		return nil, err
	}
	return projects.(*[]KubeConfig), err
}

// GetKubeConfigById 根据集群id获取kube_config配置
func (k *KubeConfig) GetKubeConfigById(id int64) (*KubeConfig, error) {
	whereScopes := func(db *gorm.DB) *gorm.DB {
		return db.Where(map[string]interface{}{"id": id})
	}
	data, err := NewDataModel().GetData(NewKubeConfigModel(), whereScopes, []string{"*"})
	if err == nil {
		kubeConfigData, ok := data.(*KubeConfig)
		if !ok || kubeConfigData.Id == 0 {
			return nil, KubeConfigNoExistErr
		} else {
			return kubeConfigData, nil
		}
	} else {
		return nil, err
	}
}

// GetKubeConfigByCid 根据集群id获取kube_config配置
func (k *KubeConfig) GetKubeConfigByCid(cid string) (*KubeConfig, error) {
	whereScopes := func(db *gorm.DB) *gorm.DB {
		return db.Where(map[string]interface{}{"cid": cid, "status": StatusNormal})
	}
	data, err := NewDataModel().GetData(NewKubeConfigModel(), whereScopes, []string{"*"})
	if err == nil {
		kubeConfigData, ok := data.(*KubeConfig)
		if !ok || kubeConfigData.Id == 0 {
			return nil, KubeConfigNoExistErr
		} else {
			return kubeConfigData, nil
		}
	} else {
		return nil, err
	}
}

// CreateKubeConfig 新增kube_config配置
func (k *KubeConfig) CreateKubeConfig(config *KubeConfig) error {
	_, err := NewDataModel().Insert(config)
	if err != nil {
		return err
	}
	return nil
}

// UpdateKubeConfig 更新kube_config配置
func (k *KubeConfig) UpdateKubeConfig(zid int64, config map[string]interface{}) error {
	whereScopes := func(db *gorm.DB) *gorm.DB {
		return db.Where("zid = ? and status != ? ", zid, StatusDeleted)
	}

	_, err := NewDataModel().UpdateAll(NewKubeConfigModel(), whereScopes, config)
	if err != nil {
		return err
	}
	return nil
}

// KubeConfigModel @业务模型
type KubeConfigModel struct {
	KubeConfig
}

func NewKubeConfigModel() *KubeConfigModel {
	return &KubeConfigModel{KubeConfig{}}
}

func (k *KubeConfigModel) GetTableStruct(isSlice bool) interface{} {
	if isSlice {
		return &[]KubeConfig{}
	}
	return &KubeConfig{}
}
