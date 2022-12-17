package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

var users []user
var products []product

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "admin"
	Contrasenia := "12345678"
	DB := "tarea_1_sd"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+DB)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var context = conexionDB()

func ValidateUser(c *gin.Context) {
	var enterUser user
	var respuesta respuestaLogin
	if err := c.BindJSON(&enterUser); err != nil {
		respuesta.Acceso_valido = false
		c.IndentedJSON(http.StatusOK, respuesta)
		return
	}

	var validUser user

	data, error := context.Query("SELECT * from cliente where id_cliente=? and contrasena=?", enterUser.Id_cliente, enterUser.Contrasena)
	if error != nil {
		panic(error.Error())
	}
	data.Next()
	data.Scan(&validUser.Id_cliente, &validUser.Nombre, &validUser.Contrasena)

	if validUser.Id_cliente != 0 {
		respuesta.Acceso_valido = true
	} else {
		respuesta.Acceso_valido = false
	}

	c.IndentedJSON(http.StatusOK, respuesta)
}

func addProduct(c *gin.Context) {
	var newProduct product
	var id int
	if err := c.BindJSON(&newProduct); err != nil {
		return

	}
	_, error := context.Query("INSERT INTO producto (nombre,cantidad_disponible,precio_unitario) VALUES (?,?,?);", newProduct.Nombre, newProduct.Cantidad_disponible, newProduct.Precio_unitario)

	if error != nil {
		panic(error.Error())
	}

	id_product, _ := context.Query("SELECT id_producto FROM producto ORDER BY id_producto DESC LIMIT 1;")

	id_product.Next()
	id_product.Scan(&id)

	c.IndentedJSON(http.StatusCreated, gin.H{"id_producto": id})
}

func getProducts(c *gin.Context) {

	productos, err := context.Query("SELECT * FROM producto")
	if err != nil {
		panic(err.Error())
	}
	products := []product{}
	for productos.Next() {
		var newProduct product
		err := productos.Scan(&newProduct.Id_producto, &newProduct.Nombre, &newProduct.Cantidad_disponible, &newProduct.Precio_unitario)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, newProduct)
	}
	c.IndentedJSON(http.StatusOK, products)
}

func updateProduct(c *gin.Context) {
	var newProduct product
	err := c.BindJSON(&newProduct)
	if err != nil {
		panic(err.Error())
	}
	_, error := context.Query("UPDATE producto SET cantidad_disponible=? WHERE id_producto= ?", newProduct.Cantidad_disponible, newProduct.Id_producto)

	if error != nil {
		panic(error.Error())
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id_producto": newProduct.Id_producto})
	return
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	var newProduct product
	data, _ := context.Query("SELECT * FROM producto WHERE id_producto=?", id)

	data.Next()
	data.Scan(&newProduct.Id_producto, &newProduct.Nombre, &newProduct.Cantidad_disponible, &newProduct.Precio_unitario)

	if newProduct.Id_producto == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"id_producto": "0"})
	} else {
		_, err := context.Query("DELETE FROM producto WHERE id_producto=?", id)
		if err != nil {
			panic(err.Error())
		}
		c.IndentedJSON(http.StatusOK, gin.H{"id_producto": id})
	}

}

func addCompra(c *gin.Context) {
	var newCompra compra
	var id int

	if err := c.BindJSON(&newCompra); err != nil {
		return

	}
	fmt.Println(newCompra)
	//println(len(newCompra.Productos))
	_, error := context.Query("INSERT INTO compra (id_cliente) VALUES (?);", newCompra.Id_cliente)

	if error != nil {
		panic(error.Error())
	}
	id_compra, _ := context.Query("SELECT id_compra FROM compra ORDER BY id_compra DESC LIMIT 1;")
	id_compra.Next()
	id_compra.Scan(&id)
	for i := 0; i < len(newCompra.Productos); i++ {
		_, error := context.Query("INSERT INTO detalle (id_compra, id_producto, cantidad) VALUES (?,?,?);", id, newCompra.Productos[i].Id_producto, newCompra.Productos[i].Cantidad)
		if error != nil {
			panic(error.Error())
		}
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"id_compra": id})
}

func statistics(c *gin.Context) {
	var vendido_mas int
	var vendido_menos int
	var ganancia_mas int
	var ganancia_menos int

	//Más vendido
	mas_v, _ := context.Query("SELECT id_producto  FROM detalle GROUP BY id_producto ORDER BY SUM(cantidad) DESC LIMIT 1")
	mas_v.Next()
	mas_v.Scan(&vendido_mas)

	//Menos vendido
	menos_v, _ := context.Query("SELECT id_producto  FROM detalle GROUP BY id_producto ORDER BY SUM(cantidad) ASC LIMIT 1")
	menos_v.Next()
	menos_v.Scan(&vendido_menos)

	//Más ganancia
	mas_g, _ := context.Query("SELECT detalle.id_producto FROM detalle LEFT JOIN producto ON detalle.id_producto = producto.id_producto GROUP BY detalle.id_producto ORDER BY SUM(detalle.cantidad* producto.precio_unitario) DESC LIMIT 1")
	mas_g.Next()
	mas_g.Scan(&ganancia_mas)

	//Menos vendido
	menos_g, _ := context.Query("SELECT detalle.id_producto FROM detalle LEFT JOIN producto ON detalle.id_producto = producto.id_producto GROUP BY detalle.id_producto ORDER BY SUM(detalle.cantidad* producto.precio_unitario) ASC LIMIT 1")
	menos_g.Next()
	menos_g.Scan(&ganancia_menos)

	c.IndentedJSON(http.StatusCreated, gin.H{"producto_mas_vendido": vendido_mas, "producto_menos_vendido": vendido_menos, "producto_mayor_ganancia": ganancia_mas, "producto_menor_ganancia": ganancia_menos})

}
