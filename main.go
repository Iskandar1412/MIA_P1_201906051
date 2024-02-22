package main

import (
	"MIA_P1_201906051/comandos"

	"strings"

	"github.com/fatih/color"
)

// MIA_P1_201906051/structures
// MIA_P1_201906051/size
// execute -path=/home/iskandar/Escritorio/prueba.adsj
// execute -path=/home/iskandar/Escritorio/Proyectos/Git/MIA_1S2024_201906051/Proyectos/MIA_P1_201906051/Pruebas/prueba.adsj
func main() {
	color.Blue("PROY1 - 201906051 - Juan Urbina")
	comandos.CrearCarpeta()
	for {
		input := "execute -path=Pruebas/prueba.adsj"
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Ingresar comando EXECUTE (exit para salir): >")
		// input, _ := reader.ReadString('\n')
		// input = strings.TrimSpace(input)
		if strings.ToLower(input) == "exit" {
			color.Cyan("Saliendo del programa")
			break
		} else {
			instrucciones := comandos.ObtenerComandos(input)
			if strings.HasPrefix(strings.ToLower(input), "execute") {
				//fmt.Println(instrucciones)
				ejecutar := comandos.Execute(instrucciones)
				comandos.GlobalCom(ejecutar)
			} else {
				//fmt.Println("comendo erroneo")
				var ejecutar []string
				ejecutar = append(ejecutar, input)
				comandos.GlobalCom(ejecutar)
			}
			//fmt.Println("instruciones", instrucciones)
			//fmt.Println(len(instrucciones))
		}
	}
}
