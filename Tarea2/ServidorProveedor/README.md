| Archivo | Descripción  | 
|---|---|
|server.py| Se encuentra la lógica del servidor proveedor|
|client.py| Es un archivo de prueba para ejecutar un cliente que consuma el servidor con grcp|
|db_conn.py| Se encarga de establecer la conexión con el servidor de mysql y ejecutar las consultas sql |
|stock_pb2.py y stock_pb2_grcp.py | Son los archivos automáticos generados por grpc que contienen el stub y la serialización/deserialización de los datos |

## Definición de IDL stock.proto

* Servicio **Stock** utiliza un único método **SendStock** que recibe un mensaje de tipo *ProductRequest* y envia como respuesta un mensaje de tipo *ProductReply*.  
* **ProductRequest** y **ProductReply** tienen basicamente la misma estructura, solo que el primero tiene la cantidad solicitada y el segundo la cantidad enviada de productos. Por otra parte, ambos tienen el id del producto y el nombre.

