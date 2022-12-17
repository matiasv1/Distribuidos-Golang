import grpc
import logging
import threading

import stock_pb2
import stock_pb2_grpc
import db_conn

from concurrent import futures
from time import sleep
from random import randint
# Conexion DB
db_inventario = db_conn.DataBase()
# Reponer inventario
def scheduled_task():
    sleep(1) # darle 1 seg despues de correr el servidor
    print("Actualizando productos sin stock cada 1 min...")
    contador_iteraciones = 1
    while True:
        print("\nRevisión de Inventario: min #",contador_iteraciones)
        lista_productos = db_inventario.listar_productos()
        for producto in lista_productos:
            if producto.cantidad_disponible == 0:
                cant_aleatoria = randint(1,10) # criterio = num aleatorio entre 1 y 10
                db_inventario.update_producto(producto.id_producto, producto.nombre, cant_aleatoria)
                print("=> Iteración",contador_iteraciones,"inventario del producto: ",producto.nombre," (#",producto.id_producto,"), se han repuesto",cant_aleatoria,"unidades.")
        contador_iteraciones += 1
        sleep(60) # espera 1 min para la siguiente revision de inventario
class Stock(stock_pb2_grpc.StockServicer):
    def SendStock(self, request, context):
        print("\nSolicitud Recibida ->",request.cantidadSolicitada,"unidades de", request.nombre)
        producto = db_inventario.producto(request.id)
        if producto is None or producto.id_producto == 0: # si el producto no se encuentra en la BD, se retorna un objeto producto con id = 0
            print("- No existe este producto")
            db_inventario.create_producto(
                request.id,
                request.nombre,
                5 # se crea el producto con 5 unidades disponibles
            )
            print("- Se ha creado un producto")
        elif request.cantidadSolicitada <= producto.cantidad_disponible:
            new_cant = producto.cantidad_disponible - request.cantidadSolicitada # se resta la cantidad enviada
            db_inventario.update_producto(request.id, request.nombre,new_cant) # actualiza el inventario restando los productos solicitados
            print("- Stock actualizado, unidades disponibles:",new_cant)
        else:
            new_cant = randint(1,10) # no hay stock, repone cantidad con num aleatoria entre 1 y 10
            db_inventario.update_producto(request.id, request.nombre,new_cant) # actualiza el producto en la BD
            print("- Stock disponible:",new_cant)
        print("Respuesta Enviada:",request.cantidadSolicitada,"unidades de",request.nombre)
        return stock_pb2.ProductReply(id=request.id,nombre=request.nombre,cantidadEnviada=request.cantidadSolicitada) # responde enviando los productos solicitados

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    stock_pb2_grpc.add_StockServicer_to_server(Stock(),server)
    server.add_insecure_port('10.10.11.174:9000')
    server.start()
    print("El servidor 2: Sistema del Proveedor, está corriendo...")
    server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig()
    # Reponer productos sin stock cada 1 min
    hilo_reponer = threading.Thread(target=scheduled_task)
    # Correr el hilo
    hilo_reponer.start()
    # Correr el servidor
    serve()
