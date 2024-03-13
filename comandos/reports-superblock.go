package comandos

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Report_SUPERBLOCK(name string, path string, ruta string, id_disco string) {

	conjunto, route, ec := Obtener_Particion_ID(id_disco)
	if !ec {
		return
	}

	nombre_falta := strings.Split(route, "/")
	nombre := nombre_falta[len(nombre_falta)-1]
	nombre_archivo_sin_extension := strings.Split(nombre, ".")
	nombre_final := nombre_archivo_sin_extension[0]

	superbloque, esb := ReducirSuperBloqueObtener(route, id_disco, conjunto)
	if !esb {
		return
	}

	data := "digraph G {\n\tnode[shape=plaintext fontsize=12];\n\trankdir=LR;\n"

	data += "\n\ttable [label=<\n"
	data += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"

	//

	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" BGCOLOR=\"#0E653E\"><FONT COLOR=\"WHITE\"><B>SUPERBLOCK REPORT</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>SB_Nombre_HD</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + nombre_final + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_FileSystem_Type</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_filesistem_type)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>S_Inodes_Count</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_inodes_count)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_Blocks_Count</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_blocks_count)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>S_Free_Blocks_Count</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_free_blocks_count)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_Free_Inodes_Count</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_free_inodes_count)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>S_MTIME</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + IntFechaToStr(superbloque.S_mtime) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_UMTIME</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + IntFechaToStr(superbloque.S_umtime) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>S_MNT_Count</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_mnt_count)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_MAGIC</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_magic)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>S_Inode_S</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_inode_s)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_Block_S</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_block_s)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>S_First_Ino</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_first_ino)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_First_Blo</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_first_blo)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>S_BM_Inode_Start</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_bm_inode_start)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_BM_Block_Start</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_bm_block_start)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>S_Inode_Start</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_inode_start)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>S_Block_Start</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_block_start)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>SB_AP_Bitmap_Bloques</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_bm_block_start)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>SB_AP_Bloques</B></FONT></TD>\n"
	data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" BGCOLOR=\"#2AD488\"><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(superbloque.S_block_start)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"

	//

	data += "\t\t</TABLE>\n"
	data += "\t\n>];\n"
	data += "\n}\n"

	// conversión
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
	color.Green("Report Generate [SuperBlock]")
}
