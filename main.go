package zin1

import (
	"github.com/devisettymahidhar315/zin1/api"
	"github.com/gin-gonic/gin"
)

func Hello() *gin.Engine {

	// Initialize the Gin router and routes
	r := gin.Default()

	r.GET("/:key", api.GetCacheValue)
	r.DELETE("/:key", api.DeleteCacheValue)
	r.POST("/:key/:value/:time", api.SetCacheValue)
	r.GET("/redis/print", api.PrintRedisCache)
	r.GET("/inmemory/print", api.PrintInMemoryCache)
	r.DELETE("/all", api.DeleteAll)

	return r
}
