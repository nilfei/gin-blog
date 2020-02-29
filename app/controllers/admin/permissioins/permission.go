package permissioins

import (
	"Kronos/app/controllers/admin"
	"Kronos/app/models"
	"Kronos/library/apgs"
	"Kronos/library/page"
	"encoding/json"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"html/template"
)

type PermissionHandler struct {
	admin.AdminBaseHandler
}

func (p *PermissionHandler) Lists(c *gin.Context) {

	all := p.AllParams(c)

	// 条件封装
	where := p.GetWhere(10)
	if all["title"] != nil {
		where["title like"] = all["title"].(string) + "%"
	}

	build, vals, _ := models.WhereBuild(where)
	// 列表页
	//
	var model = &models.Permissions{}
	//查询总条数
	count := model.GetByCount(build, vals)
	// 分页
	pagination := page.NewPagination(c.Request, count, 10)
	// 查询数据绑定到列表slice
	fields := "id, title, http_path, method,slug"
	lists, _ := model.Lists(fields, build, vals, pagination)

	ginview.HTML(c, 200, "permission/lists", gin.H{
		"pagination": template.HTML(pagination.Pages()),
		"total":      pagination.Total,
		"lists":      lists,
		"req":        all,
	})
}

func (p *PermissionHandler) ShowEdit(c *gin.Context) {
	query := p.AllParams(c)

	var model = models.Permissions{}
	// 编辑
	if query["id"] != nil {
		where := p.GetWhere(10)

		where["id"] = query["id"]
		build, vals, _ := models.WhereBuild(where)
		model, _ = model.Get(build, vals)
	}

	ginview.HTML(c, 200, "permission/edit", gin.H{
		"data": model,
		"req":  query,
	})
}

func (p *PermissionHandler) Apply(c *gin.Context) {

	var model = models.Permissions{}
	err := c.ShouldBind(&model)

	if err != nil {
		c.JSON(200, apgs.NewApiReturn(300, "无法获取到参数", nil))
		return
	}

	if model.ID > 0 {
		v := p.GetWhere(10)
		marshal, _ := json.Marshal(model)
		_ = json.Unmarshal(marshal, &v)
		err = model.Update(int(model.ID), v)
		if err != nil {
			c.JSON(200, apgs.NewApiReturn(4003, "无法更新该数据", err))
			return
		}
		c.JSON(200, apgs.NewApiRedirect(200, "更新成功", "/admin/permission/lists"))
		return

	} else {

		err := model.Create()
		if err != nil {
			c.JSON(200, apgs.NewApiReturn(4003, "无法创建该数据", nil))
			return
		}
		c.JSON(200, apgs.NewApiRedirect(200, "创建成功", "/admin/permission/lists"))
		return
	}
}

func (p *PermissionHandler) Delete(c *gin.Context) {

}