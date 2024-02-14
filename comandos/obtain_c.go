package comandos

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

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
	reportes := "MIA/P1/Reports"
	discos := "MIA/P1/Disks"
	nombreArchivo := "MIA/CarpetaImagenes.txt"
	git1 := "MIA/P1/Reports/.gitignore"
	git2 := "MIA/P1/Disks/.gitignore"
	if _, err := os.Stat(nombre); os.IsNotExist(err) {
		err := os.MkdirAll(nombre, 0777)
		if err != nil {
			color.Red("Error al crear carpeta", err)
			return
		}

		color.Green("\t\t\t\t\t\t\tCarpeta MIA/P1 creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\t\t\tCarpeta MIA/P1 ya existente")
	}

	if _, err := os.Stat(reportes); os.IsNotExist(err) {
		err := os.Mkdir(reportes, 0777)
		if err != nil {
			color.Red("Error al crear carpeta", err)
			return
		}
		color.Green("\t\t\t\t\t\t\tCarpeta MIA/P1/Reports creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\t\t\tCarpeta MIA/P1/Reports ya existente")
	}

	if _, err := os.Stat(discos); os.IsNotExist(err) {
		err := os.Mkdir(discos, 0777)
		if err != nil {
			color.Red("Error al crear carpeta", err)
			return
		}
		color.Green("\t\t\t\t\t\t\tCarpeta MIA/P1/Disks creada correctamente")
	} else {
		color.Yellow("\t\t\t\t\t\t\tCarpeta MIA/P1/Disks ya existente")
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
		color.Green("\t\t\t\t\t\t\tArchivo creado correctamente")
	} else {
		color.Yellow("\t\t\t\t\t\t\tArchivo existente")
	}

	if _, err := os.Stat(git1); os.IsNotExist(err) {
		archivo, err := os.Create(git1)
		if err != nil {
			fmt.Println("Error al crear archivo")
			return
		}
		defer archivo.Close()
		color.Green("\t\t\t\t\t\t\tArchivo creado correctamente")
	} else {
		color.Yellow("\t\t\t\t\t\t\tArchivo existente")
	}

	if _, err := os.Stat(git2); os.IsNotExist(err) {
		archivo, err := os.Create(git2)
		if err != nil {
			fmt.Println("Error al crear archivo")
			return
		}
		defer archivo.Close()
		color.Green("\t\t\t\t\t\t\tArchivo creado correctamente")
	} else {
		color.Yellow("\t\t\t\t\t\t\tArchivo existente")
	}
}

func Execute(x []string) []string {
	for _, y := range x {
		var path string
		if strings.HasPrefix(strings.ToLower(y), "path") {
			path = TienePath(y)
		} else {
			y := strings.Split(y, "=")
			color.Red("[EXECUTE] ( \"" + y[0] + "\" ): Comando no reconocido")
			break
		}
		if path == "nil" {
			return nil
		} else {
			return ExecuteFunc(path)
		}
	}
	return nil
}

func TienePath(x string) string {
	y := strings.Split(x, "=")
	fmt.Print("\t\t\t\t\t\t\tBuscando:")
	color.Yellow(y[1])
	if _, err := os.Stat(y[1]); os.IsNotExist(err) {
		color.Red("Archivo No Encontrado")
		return "nil"
	} else {
		color.Green("Archivo Encontrado")
		return y[1]
	}
}

func ExecuteFunc(x string) []string {
	file, err := os.Open(x)
	if err != nil {
		color.Red("Error al abrir archivo", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineas []string

	for scanner.Scan() {
		linea := strings.TrimSpace(scanner.Text())
		if len(linea) > 0 && !strings.HasPrefix(linea, "#") {
			lineas = append(lineas, linea)
		}
	}

	var exportar []string
	reg := regexp.MustCompile(`(.*?)\s*(?:#.*|$)`)
	for _, y := range lineas {
		match := reg.FindStringSubmatch(y)
		//fmt.Println(y, "asdf")
		if len(match) > 1 {
			exportar = append(exportar, match[1])
			//fmt.Println(match[0], "///", match[1])
		}
	}
	//fmt.Println(exportar)
	if err := scanner.Err(); err != nil {
		color.Red("Error en la lectura del archivo:", err)
		return nil
	}

	return exportar
}
