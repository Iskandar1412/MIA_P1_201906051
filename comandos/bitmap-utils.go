package comandos

import (
	"os"

	"github.com/fatih/color"
)

// Modificar bitmap
func Modificar_Bitmap(comando string, path string, inicio_inodos int32, numero int32, accion int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}
	defer file.Close()
	bytes_con_inodos := inicio_inodos + numero
	if _, err := file.Seek(int64(bytes_con_inodos), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return
	}

	if accion == int32(0) {
		if _, err := file.Write([]byte{'\x00'}); err != nil {
			color.Red("[" + comando + "]: No se pudo escribir en el archivo")
			return
		}
	}

	if accion == int32(1) {
		if _, err := file.Write([]byte{'\x01'}); err != nil {
			color.Red("[" + comando + "]: No se pudo escribir en el archivo")
			return
		}
	}
	color.Cyan("[" + comando + "]: Bitmap modificado")
}

func Vaciar_Conjunto_Bitmap(comando string, path string, inicio_bitmap int32, numEstruc int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(inicio_bitmap), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return
	}

	for i := 0; i < int(numEstruc); i++ {
		if _, err := file.Write([]byte{'\x00'}); err != nil {
			color.Red("[" + comando + "]: No se pudo escribir en el archivo")
			return
		}
	}
	color.Green("Bitmap Creado")
}
