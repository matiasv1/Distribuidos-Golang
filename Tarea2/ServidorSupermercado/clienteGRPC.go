package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "10.10.11.174:9000", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

type server struct {
	UnimplementedStockServer
}

func enviarSolicitudProveedor(id int, nombre string, cantidad int) {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar con el Servidor del Proveedor: %v", err)
	}
	defer conn.Close()
	c := NewStockClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SendStock(ctx, &ProductRequest{Id: int32(id), Nombre: nombre, CantidadSolicitada: int32(cantidad)})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Se han recibido: %d unidades del Proveedor", r.GetCantidadEnviada())

}
