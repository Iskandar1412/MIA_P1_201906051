package comandos

import (
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"os"

	"github.com/fatih/color"
)

func MKDISK_Create(_size int64, _fit byte, _unit byte) {
	directorio := "MIA/P1/Disks/"
	for i := 0; i < 26; i++ {
		letra := string('A'+i) + ".dsk"
		archivo := directorio + letra
		if _, err := os.Stat(archivo); os.IsNotExist(err) {
			CreateFile(archivo, _size, _fit, _unit)
			color.Green("[MKDISK]: Disco '" + letra + "' Creado")
			break
		} else {
			//color.Yellow("[MKDISK]: Disco Existente")
		}
	}
}

func CreateFile(archivo string, _size int64, _fit byte, _unit byte) {
	file, err := os.Create(archivo)
	if err != nil {
		color.Red("Error al crear el archivo")
		return
	}
	defer file.Close()
	//Escribir datos en archivo
	var estructura structures.MBR
	tamanio := Tamano(_size, _unit)
	estructura.Mbr_tamano = tamanio
	estructura.Mbr_fecha_creacion = ObFechaInt()
	estructura.Mbr_disk_signature = ObDiskSignature()
	estructura.Dsk_fit = _unit
	for i := 0; i < len(estructura.Mbr_partitions); i++ {
		estructura.Mbr_partitions[i] = PartitionVacia()
	}
	//Llenar el archivo
	bytes_llenar := make([]byte, int(tamanio))
	if _, err := file.Write(bytes_llenar); err != nil {
		color.Red("Error al escribir bytes en el archivo")
		return
	}

	//mover de posicion
	if _, err := file.Seek(0, 0); err != nil {
		color.Red("Error al mover puntero del archivo")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &estructura); err != nil {
		color.Red("Error al escribir datos del MBR")
		return
	}
}
