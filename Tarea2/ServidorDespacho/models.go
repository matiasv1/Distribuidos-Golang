package main

type despacho struct {
	Id_despacho int    `json:"id_despacho"`
	Id_compra   int    `json:"id_compra"`
	Estado      string `json:"estado"`
}
