package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	router.POST("/api/clientes/iniciar_sesion", ValidateUser)
	router.POST("/api/productos", addProduct)
	router.POST("/api/compras", addCompra)

	router.GET("/api/estadisticas", statistics)
	router.GET("/api/productos", getProducts)

	router.GET("/api/estado_despacho/:id", stateDespacho)

	router.PUT("/api/productos/", updateProduct)

	router.DELETE("/api/productos/:id", deleteProduct)

	router.Run("localhost:5000")
}
