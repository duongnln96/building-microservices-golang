package routes

func (r *productRouter) ProductRouterV1() {
	v1 := r.router.Group("/api/v1")

	products := v1.Group("/products")
	{
		products.GET("/", r.ctrler.AllProducts)
		products.GET("/:id", r.ctrler.FindProductByID)
		products.POST("/", r.ctrler.CreateProduct)
		products.PUT("/:id", r.ctrler.UpdateProduct)
		products.DELETE("/:id", r.ctrler.DeleteProduct)
	}
}
