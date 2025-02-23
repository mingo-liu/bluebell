package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/pprof"		// 性能分析pprof
)

func SetUpRouter(mode string) (*gin.Engine){
	// 设置gin为 release 模式，启动web程序时，无输出
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} 

	r := gin.New()	
	// 设置全网站限流
	// r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")		// index.html中需要的东西从 /static 找
	r.GET("/", func(c *gin.Context){
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1 := r.Group("api/v1")
	// 注册路由 
	v1.POST("/signup", controller.SignUpHandler)
	// 登录路由
	v1.POST("/login", controller.LoginHandler)

	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)

	v1.GET("/post/:id", controller.GetPostDetailHandler)
	v1.GET("/posts", controller.GetPostListHandler)

	// 根据时间或者分数获取帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)
	// 在某个社区中，根据时间或者分数获取帖子列表
	// 根据community_id的有无，统一了GetPostListHandler2 和GetPostListHandler3
	v1.GET("/posts3", controller.GetPostListHandler3)
	// 在某个社区中，根据时间或者分数获取帖子列表
	v1.GET("/community_posts", controller.GetCommunityPostListHander)

	//启用JWT认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())		

	{
		v1.POST("/post", controller.CreatePostHandler)

		// 投票功能
		v1.POST("/vote", controller.PostVoteHandler)
	}
	
	// r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	// 	// 只有登录的用户可以访问该路由，JWTAuthMiddleware() 判断请求头中是否有有效的JWT
	// 	c.String(http.StatusOK, "pong")
	// })

	pprof.Register(r)	// 注册pprof

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
