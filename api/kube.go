package api

import (
	"net/http"
	"strconv"

	"github.com/shuxnhs/istio-dashboard/domain/kube"
	"github.com/shuxnhs/istio-dashboard/model"

	"github.com/gin-gonic/gin"
)

// ListNamespace
// @Description 获取所有命名空间
// @Summary  获取所有命名空间
// @Tags 	kube
// @Param	id		query		int64		true		"ID"
// @Success 200 {object} Result  "ok"
// @Router /kube/namespace/list [get]
func ListNamespace(ctx *gin.Context) {
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

	nsRsp := make([]string, 0)
	ns, err := kube.NewNamespace(kube.NewKubernetesClientSet(kubeConfig)).ListNamespaceByLabel("")
	if err != nil {
		Response(ctx, http.StatusOK, CodeKubeConnectError, err.Error(), nil)
		return
	}
	for _, item := range ns.Items {
		nsRsp = append(nsRsp, item.Name)
	}

	ResponseData(ctx, CodeSuccess, nsRsp)
}
