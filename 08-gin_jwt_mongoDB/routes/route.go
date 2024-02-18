package routes

import (
	"github.com/gin-gonic/gin"
	"go_gin_jwt/controllers"
	"go_gin_jwt/middleware"
	"net/http"
	"os"
)

func Setup() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	v1 := r.Group("/users")
	{
		/*
			这里调用 handler 不需要加(),因为加()会立即调用函数，并将其返回值作为参数传递给 v1.POST 函数
			如果加了括号的函数根本没有返回值,或者返回的不是一个可接受的参数就会报错

			不加括号是直接将函数作为参数传递给 v1.POST 这样可以将函数引用传递给 v1.POST，以便在需要时调用该函数
		*/
		v1.POST("/signup", controllers.Signup)
		v1.POST("/login", controllers.Login)

		v1.Use(middleware.Authenticate)

		v1.GET("users", controllers.GetUsers)
		v1.GET("users/:user_id", controllers.GetUser)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
