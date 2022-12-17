# este codigo es solo un cliente de prueba para establecer comunicacion grpc con el servidor del proveedor
from __future__ import print_function
import logging

import grpc
import stock_pb2
import stock_pb2_grpc
import db_conn

db_inventario = db_conn.DataBase()
def run():
    estado = True
    while estado:
        print("\nComandos de estado")
        print("0 => Salir")
        print("1 => Listar productos")
        print("2 => Hacer solicitud")
        e = str(input("Ingrese estado: "))
        if e == '0':
            estado = False
        elif e == '1':
            print("Listado de productos")
            productos = db_inventario.listar_productos()
            for p in productos:
                print(p.id_producto,p.nombre,p.cantidad_disponible)
        elif e == '2':            
            print("\nIngresar productos a solicitar: ")
            id = (int)(input("id_producto: "))
            n = str(input("nombre: "))
            c = (int)(input("cantidad a solicitar: "))
            with grpc.insecure_channel('localhost:9000') as channel:
                stub = stock_pb2_grpc.StockStub(channel)
                print("\n=> Solicitud Enviada...")
                response = stub.SendStock(stock_pb2.ProductRequest(id=id,nombre=n,cantidadSolicitada=c))
                print("=> Respuesta Recibida:",response.cantidadEnviada,"unidades de",response.nombre )
if __name__ == '__main__':
    logging.basicConfig()
    run()
