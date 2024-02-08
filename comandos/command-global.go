package comandos

import (
	"fmt"
	"strings"
)

func GlobalCom(lista []string) {
	for _, comm := range lista {
		// Administracion de Discos
		if strings.HasPrefix(strings.ToLower(comm), "mkdisk") {
			//fmt.Println("MKDISK visto")
			comandos := ObtenerComandos(comm)
			fmt.Println(comandos)
		} else if strings.HasPrefix(strings.ToLower(comm), "fdisk") {
			fmt.Println("FDISK visto")
		} else if strings.HasPrefix(strings.ToLower(comm), "rmdisk") {
			fmt.Println("RMDISK visto")
		} else if strings.HasPrefix(strings.ToLower(comm), "mount") {
			fmt.Println("MOUNT visto")
		} else if strings.HasPrefix(strings.ToLower(comm), "unmount") {
			fmt.Println("UNMOUNT visto")
		} else if strings.HasPrefix(strings.ToLower(comm), "mkfs") {
			fmt.Println("MKFS visto")
		}
	}
}
