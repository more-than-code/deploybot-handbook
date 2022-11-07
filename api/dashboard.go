package api

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot/model"
)

func (a *Api) DashboardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pls, _ := a.repo.GetPipelines(ctx, model.GetPipelinesInput{})

		tmpl := template.Must(template.ParseFiles("asset/tasks.html"))

		tmpl.Execute(ctx.Writer, pls)
	}
}
