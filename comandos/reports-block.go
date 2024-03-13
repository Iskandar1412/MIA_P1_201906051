package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func Report_BLOCK(name string, path string, ruta string, id_disco string) {

	conjunto, route, ec := Obtener_Particion_ID(id_disco)
	if !ec {
		return
	}

	// nombre_falta := strings.Split(route, "/")
	// nombre := nombre_falta[len(nombre_falta)-1]
	// nombre_archivo_sin_extension := strings.Split(nombre, ".")
	// nombre_final := nombre_archivo_sin_extension[0]

	superbloque, esb := ReducirSuperBloqueObtener(route, id_disco, conjunto)
	if !esb {
		return
	}

	data := "digraph G {\n\tnode[shape=plaintext fontsize=12];\n\trankdir=LR;\n"

	conexiones := ""
	anterior := ""

	//
	res, eres := Obtener_Codigo_Bloques("REP", route, superbloque, 0, "   ")
	if !eres {
		return
	}
	///

	//lectura archivo
	file, err := os.OpenFile(route, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[REP]: Error al leer Archivo")
		panic(err)
	}
	defer file.Close()
	//bytes_con_inodos := inicio_inodos + numero
	_, err = file.Seek(int64(superbloque.S_bm_inode_start), 0)
	if err != nil {
		color.Red("[REP]: Error al mover el puntero")
		panic(err)
	}

	for i := 0; i < int(superbloque.S_blocks_count); i++ {
		var valor int8
		if err := binary.Read(file, binary.LittleEndian, &valor); err != nil {
			color.Red("[REP]: Error en lectura de bloque")
			panic(err)
		}
		if valor == 1 {
			actual := "bloque" + fmt.Sprint(i)
			if i != 0 && i != 26952 {
				conexiones += anterior + " -> " + actual + ";\n"
			}
			anterior = "bloque" + fmt.Sprint(i)
		}
	}

	data += res
	data += "\n\n" + conexiones

	///
	//

	data += "\n}\n"

	// conversión
	nombre_sin_extension := strings.Split(name, ".")
	rutaB := path + "/" + nombre_sin_extension[0] + ".dot"
	err = os.WriteFile(rutaB, []byte(data), 0644)
	if err != nil {
		color.Red("[REP]: Error al guardar el reporte del MBR")
		return
	}
	imagepath := path + "/" + name

	cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagepath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error to generate img", err)
		return
	}
	color.Green("Report Generate [Block]")
}
