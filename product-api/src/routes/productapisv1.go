package routes

import (
	"github.com/duongnln96/building-microservices-golang/product-api/middleware"
)

func (r *productRouter) ProductRouterV1() {
	v1 := r.router.Group("/api/v1")

	products := v1.Group("/products", middleware.AuthorizeJWT(r.jwtService))
	{
		products.GET("/", r.controller.AllProducts)
		products.GET("/:id", r.controller.FindProductByID)
		products.POST("/", r.controller.CreateProduct)
		products.PUT("/:id", r.controller.UpdateProduct)
		products.DELETE("/:id", r.controller.DeleteProduct)
	}
}
