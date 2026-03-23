package controller

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置API路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// API任务相关路由
	apiGroup := r.Group("/api")
	{
		// API任务增删改查
		apiGroup.POST("/jobs", CreateApiJob)
		apiGroup.GET("/jobs", GetApiJobs)
		apiGroup.GET("/jobs/:id", GetApiJob)
		apiGroup.PUT("/jobs/:id", UpdateApiJob)
		apiGroup.DELETE("/jobs/:id", DeleteApiJob)

		// 执行记录相关路由
		apiGroup.GET("/run-records", GetApiRunRecords)
		apiGroup.GET("/run-records/:id", GetApiRunRecord)
		apiGroup.GET("/run-records/api/:api_code", GetApiRunRecordsByApiCode)
	}

	return r
}
