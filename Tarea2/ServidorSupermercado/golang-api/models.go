package main

import (
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	Id_cliente int    `json:"id"`
	Nombre     string `json:"nombre"`
	Contrasena string `json:"contrasena"`
}

type despacho struct {
	Id_despacho int    `json:"id_despacho"`
	Id_compra   int    `json:"id_compra"`
	Estado      string `json:"estado"`
}

type respuestaLogin struct {
	Acceso_valido bool `json:"acceso_valido"`
}
type respuestaCompra struct {
	Id_compra                    int `json:"id_compra"`
	Cantidad_productos_comprados int `json:"cant_prod_comp"`
	Monto_total                  int `json:"monto_total"`
}
type product struct {
	Id_producto         int    `json:"id_producto"`
	Nombre              string `json:"nombre"`
	Cantidad_disponible int    `json:"cantidad_disponible"`
	Precio_unitario     int    `json:"precio_unitario"`
}

type detalle struct {
	Id_producto int `json:"id_producto"`
	Cantidad    int `json:"cantidad"`
}

type compra struct {
	Id_cliente int       `json:"id_cliente"`
	Productos  []detalle `json:"productos"`
}
