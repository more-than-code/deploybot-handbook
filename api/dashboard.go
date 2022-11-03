package api

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func (a *Api) DashboardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pls, _ := a.repo.GetPipelines(ctx)

		tmpl := template.Must(template.ParseFiles("asset/tasks.html"))

		tmpl.Execute(ctx.Writer, pls)
	}
}
