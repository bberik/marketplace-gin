package routes

import (
	"github.com/bberik/ecom-gin-react/handlers"
	"github.com/bberik/ecom-gin-react/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunServer(address string) {

	userHandler := handlers.NewUserHandler()
	productHandler := handlers.NewProductHandler()
	orderHandler := handlers.NewOrderHandler()
	cartHandler := handlers.NewCartHandler()
	shopHandler := handlers.NewShopHandler()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to Our Mini Ecommerce")
	})

	apiRoutes := router.Group("/api")

	userRoutes := apiRoutes.Group("/user")
	{
		userRoutes.POST("/signup", userHandler.SignUp)
		userRoutes.POST("/signin", userHandler.SignInUser)
	}

	userProtectedRoutes := apiRoutes.Group("/users", middleware.AuthorizeJWT())
	{
		userProtectedRoutes.GET("/user", userHandler.GetUser)
		userProtectedRoutes.PUT("/user", userHandler.UpdateUser)
		userProtectedRoutes.DELETE("/user", userHandler.DeleteUser)
		userProtectedRoutes.PUT("/user/openshop", userHandler.OpenShop)
		userProtectedRoutes.PUT("/user/addaddress", userHandler.AddAddress)
		userProtectedRoutes.PUT("/user/updateaddress", userHandler.UpdateAddress)
	}

	productRoutes := apiRoutes.Group("/products")
	{
		productRoutes.GET("/", productHandler.GetAllProducts)
		productRoutes.GET("/product", productHandler.GetProduct)
		productRoutes.GET("/category", productHandler.GetProductsByCategory)
		productRoutes.GET("/search", productHandler.SearchProduct)
	}

	shopRoutes := apiRoutes.Group("/shop", middleware.AuthorizeJWT())
	{
		shopRoutes.POST("/", shopHandler.CreateProduct)
		shopRoutes.PUT("/product", shopHandler.UpdateProduct)
		shopRoutes.DELETE("/product", shopHandler.DeleteProduct)
		shopRoutes.PUT("/order", shopHandler.UpdateOrderStatus)
	}

	orders := apiRoutes.Group("/order", middleware.AuthorizeJWT())
	{
		orders.POST("/direct", orderHandler.OrderDirect)
		orders.POST("/fromcart", orderHandler.OrderFromCart)
	}

	cart := apiRoutes.Group("/cart", middleware.AuthorizeJWT())
	{
		cart.PUT("/add", cartHandler.AddToCart)
		cart.PUT("/update", cartHandler.UpdateCart)
	}
	router.Run(address)
}
