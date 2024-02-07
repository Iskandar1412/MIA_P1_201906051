package comandos

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
)

func ObtenerComandos(x string) []string {
	var comandos []string
	atributos := regexp.MustCompile(`(-|>)(\w+)(?:="([^"]+)"|=(-?/?\w+(?:/[\w.-]+)*))?`).FindAllStringSubmatch(x, -1)
	for _, matches := range atributos {
		atributo := matches[2]
		valorConComillas := matches[3]
		valorSinComillas := matches[4]
		if valorConComillas != "" {
			comandos = append(comandos, fmt.Sprintf("%s=%s", atributo, valorConComillas))
		} else if valorSinComillas != "" {
			comandos = append(comandos, fmt.Sprintf("%s=%s", atributo, valorSinComillas))
		} else {
			comandos = append(comandos, atributo)
		}
	}
	return comandos
}

func CrearCarpeta() {
	nombre := "MIA/P1"
	nombreArchivo := "MIA/CarpetaImagenes.txt"
	if _, err := os.Stat(nombre); os.IsNotExist(err) {
		err := os.MkdirAll(nombre, 0777)
		if err != nil {
			color.Red("Error al crear carpeta", err)
			return
		}

		color.Green("\t\t\t\t\t\t\t\t\t\t\t\tCarpeta MIA/P1 creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\t\t\t\t\t\t\t\tCarpeta MIA/P1 ya existente")
	}

	if _, err := os.Stat(nombreArchivo); os.IsNotExist(err) {
		archivo, err := os.Create(nombreArchivo)
		if err != nil {
			fmt.Println("Error al crear archivo")
			return
		}
		defer archivo.Close()

		content := []byte("Proyecto 1 Manejo e Implementación de Archivos A\n\nCarpeta de Imagenes\n\nPara usar colores para imprimirlos (poner en consola): \"go get -u github.com/fatih/color\"\n\n\t\tCreated by Iskandar")
		_, err = archivo.Write(content)
		if err != nil {
			color.Red("Error escribiendo archivo:", err)
			return
		}
		color.Green("\t\t\t\t\t\t\t\t\t\t\t\tArchivo creado correctamente")
	} else {
		color.Yellow("\t\t\t\t\t\t\t\t\t\t\t\tArchivo existente")
	}
}
