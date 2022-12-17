import mysql
import mysql.connector
class Producto:
    def __init__(self, id_producto, nombre, cantidad_disponible):
        self.id_producto = id_producto
        self.nombre = nombre
        self.cantidad_disponible = cantidad_disponible

class DataBase:
    def __init__(self):
        self.connection = mysql.connector.connect(
            host="localhost", #ip
            user="admin",
            passwd="12345678",
            db="db_inventario",
            port=3306
        )
        self.cursor = self.connection.cursor()
        print("Conexión a la BD db_inventario fue establecida exitosamente.")
    # CRUD en Python
    # CREATE
    def create_producto(self, id_producto, nombre, cantidad_disponible):
        query_sql = "INSERT INTO inventario (id_producto,nombre,cantidad_disponible) VALUES ({},'{}',{});".format(
            id_producto,
            nombre,
            cantidad_disponible
        )
        try:
            self.cursor.execute(query_sql)
            self.connection.commit() # para mantener los cambios
            return True
        except Exception as e:
            raise
    # READ
    def producto(self, id_producto):
        query_sql = "SELECT * FROM inventario WHERE id_producto = {};".format(id_producto)
        try:
            self.cursor.execute(query_sql)
            producto = self.cursor.fetchone()
            if producto == None:
                producto = Producto(0,"",0)
            else:
                producto = Producto(producto[0],producto[1],producto[2])
            return producto
        except Exception as e:
            raise
    def listar_productos(self):
        query_sql = "SELECT * FROM inventario"
        try:
            self.cursor.execute(query_sql)
            productos = self.cursor.fetchall()
            lista_productos = []
            if (self.cursor.rowcount > 0):
                for p in productos:
                    producto = Producto(p[0],p[1],p[2])
                    lista_productos.append(producto)
            return lista_productos
        except Exception as e:
            raise
    # UPDATE
    def update_producto(self, id_producto, nombre, cantidad_disponible):
        query_sql = "UPDATE inventario SET nombre='{}', cantidad_disponible={} WHERE id_producto={}".format(            
            nombre,
            cantidad_disponible,
            id_producto
        )
        try:
            self.cursor.execute(query_sql)
            self.connection.commit() # para mantener los cambios
            return True
        except Exception as e:
            raise
    # DELETE
    def delete_producto(self, id_producto):
        query_sql = "DELETE FROM inventario WHERE id_producto={}".format(
            id_producto
        )
        try:
            self.cursor.execute(query_sql)
            self.connection.commit() # para mantener los cambios
            return True
        except Exception as e:
            raise