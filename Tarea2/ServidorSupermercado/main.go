package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"math/rand"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

var urlbase = "http://localhost:5000/"

// LOGIN
var session = false
var admin = false
var cliente = false

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func listar_productos() []product {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlbase+"api/productos", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject []product
	json.Unmarshal(bodyBytes, &responseObject)
	return responseObject
}

func update_producto(id_producto int, cambio int) {
	var Producto product
	Producto.Id_producto = id_producto
	Producto.Cantidad_disponible = cambio
	Json, _ := json.Marshal(Producto)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", urlbase+"api/productos/", bytes.NewBuffer(Json))
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var p product
	json.Unmarshal(respBody, &p) //
	return
}
func delete_producto(id int) {
	id_producto := strconv.Itoa(id)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", urlbase+"api/productos/"+id_producto, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var p product
	json.Unmarshal(respBody, &p) //
	return
}
func options_login() {
	println("\nOpciones:")
	println("1. Iniciar sesión como cliente")
	println("2. Iniciar sesión como administrador")
	println("3. Salir")
	print("Ingrese una opción: ")
	var opcion int
	fmt.Scan(&opcion)
	if opcion == 1 {
		login_norm()
	} else if opcion == 2 {
		login_admin()
	} else if opcion == 3 {
		print("\nHasta luego!")
	} else {
		options_login()
	}
}
func login_norm() {
	var id int
	var pwd string
	print("Ingrese su id: ")
	fmt.Scan(&id)
	print("Ingrese su contraseña: ")
	fmt.Scan(&pwd)
	var usuario user
	usuario.Id_cliente = id
	usuario.Contrasena = pwd
	//
	Json, _ := json.Marshal(usuario)
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlbase+"api/clientes/iniciar_sesion", bytes.NewBuffer(Json))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject respuestaLogin
	json.Unmarshal(bodyBytes, &responseObject)
	session = responseObject.Acceso_valido
	// VALIDAR SESION
	if session {
		session = true
		cliente = true
		admin = false
		println("Inicio de sesión exitoso")
		options_cliente(usuario.Id_cliente)
	} else {
		println("Error, no hay ninguna coincidencia con los datos ingresados")
		options_login()
	}
}
func login_admin() {
	print("Ingrese contraseña de administrador: ")
	var pwd string
	fmt.Scan(&pwd)
	// VALIDAR SESION
	if pwd == "1234" {
		session = true
		admin = true
		cliente = false
		println("Inicio de sesión exitoso")
		options_admin()
	} else {
		println("Error, no hay ninguna coincidencia con los datos ingresados")
		options_login()
	}
}

// OPTINOS
func options_admin() {
	println("\nOpciones:")
	println("1. Ver lista de productos")
	println("2. Crear producto")
	println("3. Eliminar producto")
	println("4. Ver estadísticas")
	println("5. Salir")
	print("Ingrese una opción: ")
	var opcion int
	fmt.Scan(&opcion)
	if opcion == 1 {
		exec_1_Listar_Productos(0)
	} else if opcion == 2 {
		exec_2_Crear_Producto()
	} else if opcion == 3 {
		print("Ingrese id producto a eliminar: ")
		var id int
		fmt.Scan(&id)
		exec_3_Eliminar_Producto(id)
	} else if opcion == 4 {
		exec_4_Ver_Estadisticas()
	} else if opcion == 5 {
		options_login()
	}
}
func options_cliente(id int) {
	println("\nOpciones:")
	println("1. Ver lista de productos")
	println("2. Hacer Compra")
	println("3. Consultar despacho")
	println("4. Salir")

	print("Ingrese una opción: ")
	var opcion int
	fmt.Scan(&opcion)
	if opcion == 1 {
		exec_1_Listar_Productos(id)
	} else if opcion == 2 {
		exec_2_Hacer_Compra(id)
	} else if opcion == 3 {
		print("Ingrese el ID del despacho: ")
		var id_despacho int
		fmt.Scan(&id_despacho)
		exec_3_ConsultarDespacho(id, id_despacho)
	} else if opcion == 4 {
		options_login()
	}
}

// ADMIN
func exec_1_Listar_Productos(id int) {
	/* LISTAR PRODUCTOS */
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlbase+"api/productos", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject []product
	json.Unmarshal(bodyBytes, &responseObject)
	//fmt.Printf("API Response as struct %+v\n", responseObject)
	for _, p := range responseObject {
		println(fmt.Sprint(p.Id_producto) + ";" + p.Nombre + ";" + fmt.Sprint(p.Precio_unitario) + " por unidad;" + fmt.Sprint(p.Cantidad_disponible))
	}
	if admin {
		options_admin()
	} else if cliente {
		options_cliente(id)
	}
}
func exec_2_Crear_Producto() {
	var nombre string
	var disponiblidad int
	var precio_unitario int
	print("Ingrese el nombre: ")
	fmt.Scan(&nombre)
	print("Ingrese la disponiblidad: ")
	fmt.Scan(&disponiblidad)
	print("Ingrese el precio unitario: ")
	fmt.Scan(&precio_unitario)
	var Producto product
	Producto.Id_producto = 0
	Producto.Nombre = nombre
	Producto.Precio_unitario = precio_unitario
	Producto.Cantidad_disponible = disponiblidad
	/* INGRESAR PRODUCTO */
	Json, _ := json.Marshal(Producto)
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlbase+"api/productos", bytes.NewBuffer(Json))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject []product
	json.Unmarshal(bodyBytes, &responseObject)
	//fmt.Printf("API Response as struct %+v\n", responseObject)
	println("Producto ingresado correctamente!")
	options_admin()
}

func sendProductRabbit(Id_compra int, id_despacho int) {

	var newDespacho despacho
	newDespacho.Id_despacho = id_despacho
	newDespacho.Estado = "RECIBIDO"
	newDespacho.Id_compra = Id_compra

	conn, err := amqp.Dial("amqp://admin:password@10.10.11.156:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	Json, _ := json.Marshal(newDespacho)

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        Json,
		})

	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", Json)
}
func exec_3_Eliminar_Producto(id int) {
	/* ELIMINAR PRODUCTO */
	id_producto := strconv.Itoa(id)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", urlbase+"api/productos/"+id_producto, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var p responseID
	json.Unmarshal(respBody, &p)

	//
	if p.ID_producto == "0" {
		println("No se encuentra un producto con ese id")
	} else {
		println("Producto eliminado con éxito")
	}

	options_admin()
}
func exec_4_Ver_Estadisticas() {

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlbase+"api/estadisticas", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var newStatistics statistics
	json.Unmarshal(respBody, &newStatistics)

	fmt.Println("Producto mas vendido:", newStatistics.Producto_mas_vendido)
	fmt.Println("Producto mayor ganancia:", newStatistics.Producto_mayor_ganancia)
	fmt.Println("Producto menor ganancia:", newStatistics.Producto_menor_ganancia)
	fmt.Println("Producto menos vendido:", newStatistics.Producto_menos_vendido)
	options_admin()
}

// CLIENTE
func exec_2_Hacer_Compra(id_cliente int) {
	var iteracion int
	var cantidad int
	var pares []string
	var input string
	var productos []product
	productos = listar_productos()
	var monto int
	monto = 0
	print("Ingrese cantidad de productos a comprar: ")
	fmt.Scan(&cantidad)
	for iteracion < cantidad {
		print("Ingrese producto ", iteracion+1, " par id-cantidad: ")
		fmt.Scan(&input)
		pares = append(pares, input)
		iteracion++
	}
	var newCompra compra
	newCompra.Id_cliente = id_cliente
	var elem_comprados int
	elem_comprados = 0
	for _, par := range pares {
		var newDetalle detalle
		var id string
		var cantidad string

		split_par := strings.Split(par, "-")
		id, cantidad = split_par[0], split_par[1]

		id_int, _ := strconv.Atoi(id)
		cantidad_int, _ := strconv.Atoi(cantidad)
		seEncuentra := false
		for _, p := range productos {
			if id_int == p.Id_producto {
				seEncuentra = true
				newDetalle.Id_producto = id_int
				newDetalle.Cantidad = cantidad_int
				if p.Cantidad_disponible-cantidad_int >= 0 {
					newCompra.Productos = append(newCompra.Productos, newDetalle)
					update_producto(id_int, p.Cantidad_disponible-cantidad_int)
					monto += cantidad_int * p.Precio_unitario
					elem_comprados += cantidad_int
					break
				} else {
					enviarSolicitudProveedor(id_int, p.Nombre, cantidad_int-p.Cantidad_disponible)
					fmt.Println("Se ha solitado productos al Proveedor", id_int)
					//fmt.Println("No hay stock suficiente del producto con id", id_int)
					break
				}

				// else : ignorar producto
				//break
			}
		}
		if seEncuentra == false {
			fmt.Println("No se encontro producto con id", id_int)
		}
	}
	Json, _ := json.Marshal(newCompra)

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlbase+"api/compras", bytes.NewBuffer(Json))
	if err != nil {
		fmt.Print(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var respuesta respuestaCompra
	json.Unmarshal(bodyBytes, &respuesta)

	max := 1000000 * 9
	min := 1000000
	id_despacho := rand.Intn(max-min) + min

	println("Gracias por su compra")
	println("Cantidad de productos comprados: ", elem_comprados)
	println("Monto total de la compra: ", monto)
	println("El ID de despacho es: ", id_despacho)

	// Send to rabbitMQ
	sendProductRabbit(respuesta.Id_compra, id_despacho)
	options_cliente(id_cliente)
}

func exec_3_ConsultarDespacho(id int, id_despacho int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlbase+"api/estado_despacho/"+strconv.Itoa(id_despacho), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var despachoResponse despacho
	json.Unmarshal(respBody, &despachoResponse)

	println("El estado del despacho es:", despachoResponse.Estado)
	options_cliente(id)

}

func main() {
	println("Bienvenido")
	rand.Seed(time.Now().UnixNano())
	options_login()
}
