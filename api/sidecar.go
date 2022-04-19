package api

import (
	"net/http"
	"strconv"

	"github.com/shuxnhs/istio-dashboard/domain/kube"
	"github.com/shuxnhs/istio-dashboard/domain/sidecar"
	"github.com/shuxnhs/istio-dashboard/model"

	"github.com/gin-gonic/gin"
)

// Check
// @Description 检查边车配置是否同步
// @Summary  检查边车配置是否同步
// @Tags 	sidecar
// @Param	id			query		int64		true		"id"
// @Param	namespace	query		string		true		"namespace"
// @Param	pod			query		string		true		"pod"
// @Success 200 {object} Result  "ok"
// @Router /sidecar/check [get]
func Check(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	kubeConfig, err := model.KubeConfigDB.GetKubeConfigById(id)
	if err != nil {
		ResponseData(ctx, CodeDbError, nil)
		return
	}

	sidecar.NewSidecar(kube.GetConfigStoreKubeConfig(kubeConfig)).
		Check(ctx.Query("namespace"), ctx.Query("pod"))
	if err != nil {
		Response(ctx, http.StatusOK, CodeKubeConnectError, err.Error(), nil)
		return
	}
	ResponseData(ctx, CodeSuccess, nil)
}

// ListEDS
// @Description 获取边车的EDS(端点配置)
// @Summary  获取边车的EDS(端点配置)
// @Tags 	sidecar
// @Param	id			query		int64		true		"id"
// @Param	namespace	query		string		true		"namespace"
// @Param	pod			query		string		true		"pod"
// @Success 200 {object} Result  "ok"
// @Router /sidecar/eds/list [get]
func ListEDS(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	kubeConfig, err := model.KubeConfigDB.GetKubeConfigById(id)
	if err != nil {
		ResponseData(ctx, CodeDbError, nil)
		return
	}

	eds, err := sidecar.NewSidecar(kube.GetConfigStoreKubeConfig(kubeConfig)).
		GetEDS(ctx.Query("namespace"), ctx.Query("pod"))
	if err != nil {
		Response(ctx, http.StatusOK, CodeKubeConnectError, err.Error(), nil)
		return
	}
	ResponseData(ctx, CodeSuccess, eds)
}

// ListCDS
// @Description 获取边车的CDS(集群配置)
// @Summary  获取边车的CDS(集群配置)
// @Tags 	sidecar
// @Param	id			query		int64		true		"id"
// @Param	namespace	query		string		true		"namespace"
// @Param	pod			query		string		true		"pod"
// @Success 200 {object} Result  "ok"
// @Router /sidecar/cds/list [get]
func ListCDS(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	kubeConfig, err := model.KubeConfigDB.GetKubeConfigById(id)
	if err != nil {
		ResponseData(ctx, CodeDbError, nil)
		return
	}

	cds, err := sidecar.NewSidecar(kube.GetConfigStoreKubeConfig(kubeConfig)).
		GetCDS(ctx.Query("namespace"), ctx.Query("pod"))
	if err != nil {
		Response(ctx, http.StatusOK, CodeKubeConnectError, err.Error(), nil)
		return
	}
	ResponseData(ctx, CodeSuccess, cds)
}
