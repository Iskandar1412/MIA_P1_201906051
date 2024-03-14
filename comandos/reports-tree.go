package comandos

import (
	"MIA_P1_201906051/structures"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func Report_TREE(name string, path string, ruta string, id_disco string) {
	comando := "REP"
	conjunto, route, ebc := Obtener_Particion_ID(id_disco)
	if !ebc {
		return
	}

	superbloque, esupe := ReducirSuperBloqueObtener(route, id_disco, conjunto)
	if !esupe {
		return
	}

	inodo_raiz, eir := Obtener_Inodo(comando, route, superbloque.S_inode_start, 0)
	if !eir && (inodo_raiz != structures.Inode{}) {
		return
	}

	data := "digraph G{\nrankdir=LR;\n"
	codigotree, ecodigo := Obtener_Codigo_Tree(comando, route, superbloque, 0, id_disco)
	if !ecodigo {
		return
	}
	data += codigotree
	data += "\n}\n"

	//conversion
	nombre_sin_extension := strings.Split(name, ".")
	rutaB := path + "/" + nombre_sin_extension[0] + ".dot"
	err := os.WriteFile(rutaB, []byte(data), 0644)
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
	color.Green("Report Generate [Tree]")
}
