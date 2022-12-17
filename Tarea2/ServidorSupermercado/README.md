| Archivo | Descripción  | 
|---|---|
|main.go|Es la aplicación principal que simula la interfaz de supermercado|
|clienteGRPC.go| Se encuentra un función que se encarga de solicitar al servidor proveedor los productos necesarios|
|stock_grpc.pb.go y stock.pb.go| Son los archivos automáticos generados por grpc que contienen el stub y la serialización/deserialización de los datos |

## Definición de IDL stock.proto

* Servicio **Stock** utiliza un único método **SendStock** que recibe un mensaje de tipo *ProductRequest* y envia como respuesta un mensaje de tipo *ProductReply*.  
* **ProductRequest** y **ProductReply** tienen basicamente la misma estructura, solo que el primero tiene la cantidad solicitada y el segundo la cantidad enviada de productos. Por otra parte, ambos tienen el id del producto y el nombre.

