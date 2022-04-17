package api

import (
	"github.com/shuxnhs/istio-dashboard/model"

	"github.com/gin-gonic/gin"
)

type Project struct {
	Cid         string
	Description string
	Status      int64
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

	}
	projectRsp := make([]Project, len(*projects))
	for _, project := range *projects {
		projectRsp = append(projectRsp, Project{
			Cid:         project.Cid,
			Description: project.Description,
			Status:      project.Status,
		})
	}
	ResponseData(ctx, CodeSuccess, projectRsp)
}
