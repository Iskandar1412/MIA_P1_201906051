package comandos

import (
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

func Report_INODE(name string, path string, ruta string, id_disco string) {
	conjunto, path2, eco := Obtener_Particion_ID(id_disco)
	if !eco {
		return
	}
	if conjunto == nil && path2 == "" {
		return
	}

	superbloque := structures.SuperBlock{}

	logica := structures.EBR{}
	if conjunto[2] != nil {
		if temp_log, ok := conjunto[2].(structures.EBR); ok {
			v := reflect.ValueOf(temp_log)
			reflect.ValueOf(&logica).Elem().Set(v)
			var e bool
			superbloque, e = Obtener_Superbloque("REP", path2, ToString(logica.Name[:]))
			if !e {
				return
			}
		}
		conjunto[0] = nil
		conjunto[1] = nil
	}

	particion := structures.Partition{}
	if conjunto[0] != nil {
		if temp, ok := conjunto[0].(structures.Partition); ok {
			v := reflect.ValueOf(temp)
			reflect.ValueOf(&particion).Elem().Set(v)
			var e bool
			superbloque, e = Obtener_Superbloque("REP", path2, ToString(particion.Part_name[:]))
			if !e {
				return
			}
		}
	}

	data := "digraph G {\n\tnode[shape=plaintext fontsize=12];\n\trankdir=LR;\n"
	conexiones := ""
	anterior := ""

	//lectura archivo
	file, err := os.OpenFile(path2, os.O_RDWR, 0666)
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

	for i := 0; i < int(superbloque.S_inodes_count); i++ {
		var valor int8
		if err := binary.Read(file, binary.LittleEndian, &valor); err != nil {
			color.Red("[REP]: Error en lectura de bloque")
			panic(err)
		}
		if valor == 1 { //Inodo libre
			//obtener inodo
			inodo, ei := Obtener_Inodo("REP", path2, superbloque.S_inode_start, int32(i))
			if !ei {
				return
			}
			actual := "inodo" + fmt.Sprint(i)
			data += "inodo" + fmt.Sprint(i) + "[label=<\n"
			data += "<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
			data += "\t<TR>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>INODO " + fmt.Sprint(i) + "</B></FONT></TD>\n"
			data += "\t</TR>\n"
			data += "\t<TR>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>I_UID</B></FONT></TD>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\">" + fmt.Sprint(inodo.I_uid) + "</FONT></TD>\n"
			data += "\t</TR>\n"
			data += "\t<TR>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>I_GID</B></FONT></TD>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\">" + fmt.Sprint(inodo.I_gid) + "</FONT></TD>\n"
			data += "\t</TR>\n"
			data += "\t<TR>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>I_S</B></FONT></TD>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\">" + fmt.Sprint(inodo.I_s) + "</FONT></TD>\n"
			data += "\t</TR>\n"
			data += "\t<TR>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>I_ATIME</B></FONT></TD>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\">" + IntFechaToStr(inodo.I_atime) + "</FONT></TD>\n"
			data += "\t</TR>\n"
			data += "\t<TR>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>I_CTIME</B></FONT></TD>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\">" + IntFechaToStr(inodo.I_ctime) + "</FONT></TD>\n"
			data += "\t</TR>\n"
			for j := range inodo.I_block {
				data += "\t<TR>\n"
				data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>APT" + fmt.Sprint(j) + "</B></FONT></TD>\n"
				data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\">" + fmt.Sprint(inodo.I_block[j]) + "</FONT></TD>\n"
				data += "\t</TR>\n"
			}
			data += "\t<TR>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>I_TYPE</B></FONT></TD>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\">" + fmt.Sprint(inodo.I_type) + "</FONT></TD>\n"
			data += "\t</TR>\n"
			data += "\t<TR>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>I_PERM</B></FONT></TD>\n"
			data += "\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\">" + fmt.Sprint(inodo.I_perm) + "</FONT></TD>\n"
			data += "\t</TR>\n"
			data += "</TABLE>"
			data += ">];\n"
			if i != 0 {
				conexiones += anterior + " -> " + actual + ";\n"
			}
			anterior = "inodo" + fmt.Sprint(i)

		}
	}

	data += "\n\n" + conexiones
	data += "\n}\n"

	//conversión
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
	color.Green("Report Generate [Inode]")
}

func Report_BM_Inode(name string, path string, ruta string, id_disco string) {
	conjunto, route, ec := Obtener_Particion_ID(id_disco)
	if !ec {
		return
	}

	superbloque, esb := ReducirSuperBloqueObtener(route, id_disco, conjunto)
	if !esb {
		return
	}

	file, err := os.OpenFile(route, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir archivo")
		return
	}
	defer file.Close()

	_, err = file.Seek(int64(superbloque.S_bm_inode_start), 0)
	if err != nil {
		fmt.Println("Error al buscar en archivo")
		return
	}

	rep, err := os.Create(path + "/" + name)
	if err != nil {
		fmt.Println("Error al crear archivo de reporte")
		return
	}
	defer rep.Close()

	err = rep.Truncate(0)
	if err != nil {
		return
	}

	contador := 0
	for i := 0; i < int(superbloque.S_inodes_count); i++ {
		var datos [1]byte
		_, err := file.Read(datos[:])
		if err != nil {
			fmt.Println("Error al leer datos de archivo")
			return
		}

		_, err = rep.WriteString(fmt.Sprintf("%v\t", datos[0]))
		if err != nil {
			fmt.Println("Error al escribir datos del reporte")
			return
		}
		contador++

		if contador == 20 {
			_, err = rep.WriteString("\n")
			if err != nil {
				fmt.Println("Error al escribir datos en archivo de reporte")
				return
			}
			contador = 0
		}
	}
	color.Green("Report Generate [bm_inode]")
}

func Report_BM_Block(name string, path string, ruta string, id_disco string) {
	conjunto, route, ec := Obtener_Particion_ID(id_disco)
	if !ec {
		return
	}

	superbloque, esb := ReducirSuperBloqueObtener(route, id_disco, conjunto)
	if !esb {
		return
	}

	file, err := os.OpenFile(route, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir archivo")
		return
	}
	defer file.Close()

	_, err = file.Seek(int64(superbloque.S_bm_block_start), 0)
	if err != nil {
		fmt.Println("Error al buscar en archivo")
		return
	}

	rep, err := os.Create(path + "/" + name)
	if err != nil {
		fmt.Println("Error al crear archivo de reporte")
		return
	}
	defer rep.Close()

	err = rep.Truncate(0)
	if err != nil {
		return
	}

	contador := 0
	for i := 0; i < int(superbloque.S_blocks_count); i++ {
		var datos [1]byte
		_, err := file.Read(datos[:])
		if err != nil {
			fmt.Println("Error al leer datos de archivo")
			return
		}

		_, err = rep.WriteString(fmt.Sprintf("%v\t", datos[0]))
		if err != nil {
			fmt.Println("Error al escribir datos del reporte")
			return
		}
		contador++

		if contador == 20 {
			_, err = rep.WriteString("\n")
			if err != nil {
				fmt.Println("Error al escribir datos en archivo de reporte")
				return
			}
			contador = 0
		}
	}
	color.Green("Report Generate [bm_inode]")
}
