package routes

import (
	"github.com/duongnln96/building-microservices-golang/product-api/middleware"
)

func (r *productRouter) ProductRouterV1() {
	v1 := r.router.Group("/api/v1")

	authen := v1.Group("/auth")
	{
		authen.POST("/login", r.controllers.AuthenCtrl.Login)
		authen.POST("/register", r.controllers.AuthenCtrl.Register)
	}

	products := v1.Group("/products", middleware.AuthorizeJWT(r.jwtService))
	{
		products.GET("/", r.controllers.ProductCtrl.AllProducts)
		products.GET("/:id", r.controllers.ProductCtrl.FindProductByID)
		products.POST("/", r.controllers.ProductCtrl.CreateProduct)
		products.PUT("/:id", r.controllers.ProductCtrl.UpdateProduct)
		products.DELETE("/:id", r.controllers.ProductCtrl.DeleteProduct)
	}
}
