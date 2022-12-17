package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "admin"
	Contrasenia := "12345678"
	DB := "db_despachos"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+DB)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var context = conexionDB()

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Cambiar el estado del despacho de RECIBIDO a EN_TRANSITO y de EN_TRANSITO a Entregado
func changeState(t time.Time) {

	data, error := context.Query("SELECT * from despacho")

	if error != nil {
		panic(error.Error())
	}

	despachos := []despacho{}
	for data.Next() {
		var newDespacho despacho
		err := data.Scan(&newDespacho.Id_despacho, &newDespacho.Estado, &newDespacho.Id_compra)
		if err != nil {
			panic(err.Error())
		}
		despachos = append(despachos, newDespacho)
	}

	//Iterar los despachos obtenidos de la base de datos para cambiar los estados
	for _, item := range despachos {

		if item.Estado == "RECIBIDO" {
			_, error := context.Query("UPDATE despacho SET estado='EN_TRANSITO' WHERE id_despacho=?", item.Id_despacho)
			if error != nil {
				panic(error.Error())
			}
			
		} else if item.Estado == "EN_TRANSITO" {
			_, error := context.Query("UPDATE despacho SET estado='ENTREGADO' WHERE id_despacho=?", item.Id_despacho)
			if error != nil {
				panic(error.Error())
			}
			
		}
	}

}
// Cada 60 segundos ejecutara la funcion f(x) equivalente a ChangeState
func doEvery(d time.Duration, f func(t time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func main() {
	var newDespacho despacho
	// establecer conexión con RabbitMQ
	conn, err := amqp.Dial("amqp://admin:password@0.0.0.0:5672/")
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}
	
	// Esperando el sistema de mensajería a recibir un mensaje
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			json.Unmarshal(d.Body, &newDespacho)
			_, error := context.Query("INSERT INTO despacho (id_despacho,estado, id_compra) VALUES (?,?,?);", newDespacho.Id_despacho, newDespacho.Estado, newDespacho.Id_compra)
			if error != nil {
				panic(error.Error())
			}

		}
	}()

	// Cada 60 segundos ejecutar en paralelo la función de cambiar estados
	go func() {
		doEvery(60*time.Second, changeState)
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
