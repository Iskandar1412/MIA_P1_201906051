package comandos

import (
	"strings"

	"github.com/fatih/color"
)

// var Particiones_Montadas []string

func GlobalCom(lista []string) {
	for _, comm := range lista {
		// Administracion de Discos
		if strings.HasPrefix(strings.ToLower(comm), "mkdisk") {
			comandos := ObtenerComandos(comm)
			DiskCommandProps("MKDISK", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "fdisk") {
			comandos := ObtenerComandos(comm)
			DiskCommandProps("FDISK", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "rmdisk") {
			comandos := ObtenerComandos(comm)
			DiskCommandProps("RMDISK", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "mount") {
			comandos := ObtenerComandos(comm)
			DiskCommandProps("MOUNT", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "unmount") {
			comandos := ObtenerComandos(comm)
			DiskCommandProps("UNMOUNT", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "mkfs") {
			comandos := ObtenerComandos(comm)
			DiskCommandProps("MKFS", comandos)
			// Reportes
		} else if strings.HasPrefix(strings.ToLower(comm), "rep") {
			comandos := ObtenerComandos(comm)
			ReportCommandProps("REP", comandos)
			// Files
		} else if strings.HasPrefix(strings.ToLower(comm), "mkfile") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("MKFILE", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "cat") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("CAT", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "remove") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("REMOVE", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "edit") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("EDIT", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "rename") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("RENAME", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "mkdir") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("MKDIR", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "copy") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("COPY", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "move") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("MOVE", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "find") {
			comandos := ObtenerComandos(comm)
			FilesCommandProps("FIND", comandos)
			// Permisos
		} else if strings.HasPrefix(strings.ToLower(comm), "chown") {
			comandos := ObtenerComandos(comm)
			PermissionsCommandProps("CHOWN", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "chgrp") {
			comandos := ObtenerComandos(comm)
			PermissionsCommandProps("CHGRP", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "chmod") {
			comandos := ObtenerComandos(comm)
			PermissionsCommandProps("CHMOD", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "pause") {
			comandos := ObtenerComandos(comm)
			PermissionsCommandProps("PAUSE", comandos)
			// Usuarios
		} else if strings.HasPrefix(strings.ToLower(comm), "login") {
			comandos := ObtenerComandos(comm)
			UserComanmandProps("LOGIN", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "logout") {
			comandos := ObtenerComandos(comm)
			UserComanmandProps("LOGOUT", comandos)
			// Grupo
		} else if strings.HasPrefix(strings.ToLower(comm), "mkgrp") {
			comandos := ObtenerComandos(comm)
			GroupCommandProps("MKGRP", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "rmgrp") {
			comandos := ObtenerComandos(comm)
			GroupCommandProps("RMGRP", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "mkusr") {
			comandos := ObtenerComandos(comm)
			GroupCommandProps("MKUSR", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "rmusr") {
			comandos := ObtenerComandos(comm)
			GroupCommandProps("RMUSR", comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "mkgrp") {
			comandos := ObtenerComandos(comm)
			GroupCommandProps("MKGRP", comandos)
		} else {
			color.Red("Comando no Reconocido")
		}
	}
}
