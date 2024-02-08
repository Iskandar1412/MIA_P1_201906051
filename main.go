package main

import (
	"MIA_P1_201906051/comandos"

	"bufio"
	"fmt"
	"os"
	"strings"
)

// MIA_P1_201906051/structures
// MIA_P1_201906051/size
//execute -path=/home/iskandar/Escritorio/prueba.adsj

func main() {
	fmt.Println("PROY1 - 201906051 - Juan Urbina")
	comandos.CrearCarpeta()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Ingresar comando EXECUTE (exit para salir): >")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "exit" {
			fmt.Println("Saliendo del programa")
			break
		} else {
			instrucciones := comandos.ObtenerComandos(input)
			if strings.HasPrefix(strings.ToLower(input), "execute") {
				//fmt.Println(instrucciones)
				ejecutar := comandos.Execute(instrucciones)
				fmt.Println(ejecutar)
			} else {
				fmt.Println("comendo erroneo")
			}
			//fmt.Println("instruciones", instrucciones)
			//fmt.Println(len(instrucciones))
		}
	}
}
