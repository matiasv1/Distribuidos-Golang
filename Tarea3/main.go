package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Estructura del token
type token struct {
	LN []int
	Q  []int
}

// variable global cantidad de procesos
var cantidad int = 0

// Eliminar la primera linea del archivo
func popLine(f *os.File) ([]byte, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	line, err := buf.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	nw, err := io.Copy(f, buf)
	if err != nil {
		return nil, err
	}
	err = f.Truncate(nw)
	if err != nil {
		return nil, err
	}
	err = f.Sync()
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return line, nil
}

// Función recibe string número del puesto a comprar, retorna el precio
func costo(puesto string) int {
	asiento, _ := strconv.Atoi(puesto)
	var precio int
	if asiento <= 16 {
		precio = 8000
	} else if asiento <= 32 {
		precio = 6000
	} else if asiento <= 48 {
		precio = 4000
	}
	return precio
}

// Función actualiza el vector de ganancias, recibe el vector de ganancias actual y el asiento a comprar
func ganancia(vector_ganancias [3]int, puesto string) [3]int {
	asiento, _ := strconv.Atoi(puesto)
	if asiento <= 16 {
		vector_ganancias[0] += costo(puesto)
	} else if asiento <= 32 {
		vector_ganancias[1] += costo(puesto)
	} else if asiento <= 48 {
		vector_ganancias[2] += costo(puesto)
	}
	return vector_ganancias
}

// Función actualiza el mapa de los asientos, recibe el asiento comprado
func marcar_asiento(puesto string) {
	input, err := ioutil.ReadFile("mapa.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte(puesto), []byte("XX"), -1)

	if err = ioutil.WriteFile("mapa.txt", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Función que determina si se puede comprar o no el asiento, recibe el asiento a comprar
func asiento_libre(puesto string) bool {
	readFile, err := os.Open("mapa.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	//Por linea
	for fileScanner.Scan() {
		testArray := strings.Fields(fileScanner.Text())

		//Si es una fila de asientos
		if len((testArray)) == 12 {
			for i := 0; i < 12; i++ {
				//Si es el asiento a comprar
				if puesto == testArray[i] {
					return true
				}
			}
		}
	}
	readFile.Close()
	return false
}

// Función actualiza el archivo de ganancias, recibe el asiento a comprar
func actualizar_ganancia(asiento string) {
	var ganancias [3]int
	//Se lee el archivo ganancias
	readFile, err := os.Open("ganancias.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	//Obtener arreglo de ganancias
	//Por linea
	for fileScanner.Scan() {
		testArray := strings.Fields(fileScanner.Text())
		for i := 0; i < 3; i++ {
			valor, _ := strconv.Atoi(testArray[i])
			ganancias[i] = valor
		}
	}
	//actualizar vector ganancias
	ganancias = ganancia(ganancias, asiento)
	reemplazo := ""
	for i := 0; i < 3; i++ {
		reemplazo += strconv.Itoa(ganancias[i]) + " "
	}

	//Reemplazar en archivo
	err = ioutil.WriteFile("ganancias.txt", []byte(reemplazo), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

// Función agrega el pasajero procesado al archivo procesados, recibe el nombre, asiento a comprar y el proceso que está procesando la información
func add_procesados(nombre string, asiento string, proceso int) {
	//Se abre el archivo
	f, err := os.OpenFile("procesados.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//Se escribe en el archivo
	text := "P" + strconv.Itoa(proceso) + ": " + nombre + " " + asiento + "\n"
	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func critical_section(has_token *int, in_cs *int, id_process int, RN []int, LN []int, queue []int, channels_tokens []chan token) {
	if *has_token == 1 {
		*in_cs = 1
		//Se abre el archivo pasajeros
		readFile, err := os.OpenFile("pasajeros.txt", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
		}
		fileScanner := bufio.NewScanner(readFile)

		fileScanner.Split(bufio.ScanLines)
		var cond bool
		cond = true
		//Se saca la primera linea del archivo
		// Se detiene si el archivo está vacío
		for cond {
			line, err := popLine(readFile)

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			testArray := strings.Fields(string(line[:]))

			if len(testArray) == 2 {
				asiento := testArray[1]
				nombre := testArray[0]

				//Si se pude asignar el asiento
				if asiento_libre(asiento) {
					marcar_asiento(asiento)                                          //actualiza mapa de asientos
					actualizar_ganancia(asiento)                                     // actualiza archivo ganancias
					add_procesados(nombre, asiento, id_process)                      // actualiza archivo procesados
					println("Proceso", id_process, " registrando:", nombre, asiento) // actualiza archivo pasajeros
					cond = false
				}
			} else {
				cond = false
				println("El archivo pasajeros.txt ya no tiene clientes, se termina el proceso para que no quede corriendo infinitamente...")
				os.Exit(1)
			}

		}
		readFile.Close()
		time.Sleep(1 * time.Second)

		//Termina la sección crítica
		release_cs(id_process, LN, RN, queue, has_token, channels_tokens)
		*in_cs = 0
	}
}

// Sacar el elemento tope de la cola
func dequeue(queue []int) (int, []int) {
	element := queue[0]
	if len(queue) == 1 {
		var tmp = []int{}
		return element, tmp
	}
	return element, queue[1:]
}

// Funcion para ver si el proceso se encuentra en la cola
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Liberar la sección critica
func release_cs(id_process int, LN []int, RN []int, queue []int, has_token *int, channels_tokens []chan token) {
	LN[id_process] = RN[id_process]

	//Agregar procesos a la cola
	for i := 0; i < cantidad; i++ {
		if !contains(queue, i) {
			if RN[i] == (LN[i] + 1) {
				queue = append(queue, i)
				//println("Se ha agregado el proceso ", i, " a la cola")
			}
		}
	}
	//Enviando al proceso tope de la cola
	if len(queue) != 0 {
		*has_token = 0
		element, queue2 := dequeue(queue)
		send_token(element, queue2, LN, channels_tokens)
	}
}

// Enviar mensaje de solicitando el token a los caneles de los otros procesos
func send_request(message_rn int, id_process int, channels_message []chan []int) {
	arrayMessage := []int{message_rn, id_process}
	for i := 0; i < cantidad; i++ {
		if i != id_process {
			//fmt.Println("el proceso, ", id_process, "le envia al canal ", i)
			channels_message[i] <- arrayMessage
		}
	}
}

// Enviar el token al canal correspondiente
func send_token(id_request int, queue []int, LN []int, channels_tokens []chan token) {
	var newToken token
	newToken.LN = LN
	newToken.Q = queue
	//println("Enviando token al ", id_request)
	channels_tokens[id_request] <- newToken
}

// Enviar una solicitud a los otros procesos solicitando la sección critica
func request_cs(id_process int, RN []int, in_cs *int, waiting_for_token *int, has_token *int, channels_message []chan []int) {
	if *has_token == 0 {
		RN[id_process] = RN[id_process] + 1
		*waiting_for_token = 1
		//println("Soy el proceso ", id_procces, " quiero el token ", RN[id_procces])
		send_request(RN[id_process], id_process, channels_message)
	}
}

// Función en la cual los procesos esperaran en sus canales recibir el token o request de otros procesos
func receive_request(id_process int, LN []int, RN []int, in_cs *int, waiting_for_token *int, has_token *int, queue []int, channels_message []chan []int, channels_tokens []chan token) {
	for {
		// Esperando a recibir el token en su canal correspondiente
		go func() {
			for n := range channels_tokens[id_process] {
				*has_token = 1
				*waiting_for_token = 0
				//fmt.Println("He recibido el token ", id_process, " ", n)
				LN = n.LN
				queue = n.Q
				critical_section(has_token, in_cs, id_process, RN, LN, queue, channels_tokens)
			}
		}()
		// Esperando a recibir algun mensaje en el canal de mensajes
		go func() {
			var requester_id int
			var sn_value int
			for index, i := range <-channels_message[id_process] {
				if index == 0 {
					sn_value = i
				} else {
					requester_id = i
				}
			}
			//fmt.Println("EL proceso ", id_process, " ha recibido el mensaje: ", requester_id, ",", sn_value)
			if sn_value > RN[requester_id] {
				RN[requester_id] = sn_value
			}
			if *has_token == 1 && *in_cs == 0 && RN[requester_id] == (LN[requester_id]+1) {
				*has_token = 0
				send_token(requester_id, queue, LN, channels_tokens)
			}
		}()
		// Esperando a recibir el token
	}
}

// Se crean n procesos con sus propias variables
func instance(id int, channels_message []chan []int, channels_tokens []chan token) {
	var has_token int = 0
	var in_cs int = 0
	var waiting_for_token int = 0

	RN := []int{}
	LN := []int{}

	var queue = make([]int, 0)

	for i := 0; i < cantidad; i++ {

		RN = append(RN, 0)
		LN = append(LN, 0)
	}

	if id == 0 {
		has_token = 1
	}

	RN[0] = 1

	// Proceso para recibir peticiones
	go receive_request(id, LN, RN, &in_cs, &waiting_for_token, &has_token, queue, channels_message, channels_tokens)

	// Iterar paralelamente  para pedir el token o ejecutar seccion critica, o bien realizar otra tarea
	for {

		// Se llama este time sleep para que el proceso alcance a enviar el token antes de querer ejecutar la sección critica nuevamente
		// aunque se puede eliminar sin problemas

		time.Sleep(1 * time.Second)
		if has_token == 0 {
			// Pedir el token cada 3 segundos
			time.Sleep(3 * time.Second)
			request_cs(id, RN, &in_cs, &waiting_for_token, &has_token, channels_message)
		} else if in_cs == 0 {
			// Ejecutar la sección critica
			critical_section(&has_token, &in_cs, id, RN, LN, queue, channels_tokens)
		}
		for waiting_for_token == 1 {
			//Estado que realizo el request y esta esperando por el token, en esta parte puede ejecutar otra tarea mientras
			time.Sleep(1 * time.Second)
		}
	}
}

// main
func main() {
	var wg sync.WaitGroup
	nProcesos, _ := strconv.Atoi(os.Args[1])
	cantidad = nProcesos

	var channels_message = make([]chan []int, cantidad)
	var channels_tokens = make([]chan token, cantidad)

	wg.Add(cantidad)
	// crear un canal por cada proceso
	for i := range channels_message {
		channels_message[i] = make(chan []int)
		channels_tokens[i] = make(chan token)
	}
	for i := 0; i < cantidad; i++ {
		// LLamar a un proceso entregando su id, los canales de mensajes y canales de tokens
		go instance(i, channels_message, channels_tokens)
		defer wg.Done()
	}
	wg.Wait()
}
