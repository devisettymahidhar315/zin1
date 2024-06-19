package api

import (
	"net/http"
	"strconv"

	"github.com/devisettymahidhar315/zin1/multi_cache"
	"github.com/gin-gonic/gin"
)

const length = 2

var cache = multi_cache.NewMultiCache()

// Endpoint to retrieve a value by key
func GetCacheValue(ctx *gin.Context) {
	k := ctx.Param("key")
	ctx.JSON(http.StatusOK, cache.Get(k))

}

// Endpoint to delete a value by key
func DeleteCacheValue(ctx *gin.Context) {
	//storing the key value
	k := ctx.Param("key")
	//calling the delete method
	cache.Del(k)
}

// Endpoint to set a key-value pair
func SetCacheValue(ctx *gin.Context) {
	//storing the key and value pair
	k := ctx.Param("key")
	v := ctx.Param("value")
	tstr := ctx.Param("time")
	// Convert time string to integer
	t, err := strconv.Atoi(tstr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid time parameter"})
		return
	}
	//calling the set methods and sending the key,value and length
	cache.Set(k, v, length, t)
}

// Endpoint to print the in-memory cache contents
func PrintInMemoryCache(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, cache.Print_in_mem())
}

// Endpoint to print the redis cache contents
func PrintRedisCache(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, cache.Print_redis())
}

// Endpoint to delete entire data
func DeleteAll(ctx *gin.Context) {
	cache.Del_ALL()
}
