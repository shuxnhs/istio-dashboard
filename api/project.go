package api

import (
	"github.com/shuxnhs/istio-dashboard/model"

	"github.com/gin-gonic/gin"
)

type Project struct {
	Id          int64  `json:"id"`
	Cid         string `json:"cid"`
	Description string `json:"description"`
	Status      int64  `json:"status"`
}

// ListProjects
// @Description 获取所有的网格
// @Summary  获取所有网格
// @Tags 	project
// @Success 200 {object} Result  "ok"
// @Router /project/list [get]
func ListProjects(ctx *gin.Context) {
	projects, err := model.KubeConfigDB.ListKubeConfig()
	if err != nil {
		ResponseData(ctx, CodeDbError, nil)
		return
	}
	projectRsp := make([]Project, 0)
	for _, project := range *projects {
		projectRsp = append(projectRsp, Project{
			Id:          project.Id,
			Cid:         project.Cid,
			Description: project.Description,
			Status:      project.Status,
		})
	}
	ResponseData(ctx, CodeSuccess, projectRsp)
}
