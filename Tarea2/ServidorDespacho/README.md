## Servicio Rabbitmq

Para levantar el servicio se tuvo que crear un usuario de credenciales user: admin y password: password.  
Se puede apreciar en el código "amqp://admin:password@0.0.0.0:5672/" para levantar el servicio. 

| Archivo | Descripción  | 
|---|---|
|consumer.go|Tiene implementado el servicio de Rabbitmq el cual verificara la llegada de mensajes y cada 60 segundos cambiara los estados|
|models.go| Tiene los modelos de datos utilizados|
