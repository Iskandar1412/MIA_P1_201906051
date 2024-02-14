package comandos

import (
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"os"
	"strings"

	"github.com/fatih/color"
)

// MKDISK

func Values_MKDISK(instructions []string) (int64, byte, byte) {
	var _size int64
	var _fit byte = 'F'
	var _unit byte = 'M'
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "size") {
			var value = TieneSize("MKDISK", valor)
			_size = value
		} else if strings.HasPrefix(strings.ToLower(valor), "fit") {
			var value = TieneFit("MKDISK", valor)
			_fit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "unit") {
			var value = TieneUnit("mkdisk", valor)
			_unit = value
		} else {
			color.Yellow("[MKDISK]: Atributo no reconocido")
		}
	}
	return _size, _fit, _unit
}

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

// RMDISK

func Values_RMDISK(instructions []string) (byte, bool) {
	var _drivedeletter byte = '0'
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "drivedeletter") {
			var value = TieneDriveDeLetter("RMDISK", valor)
			_drivedeletter = value
			break
		} else {
			color.Yellow("[RMDISK]: Atributo no reconocido")
			_drivedeletter = '0'
			break
		}
	}
	if _drivedeletter == '0' {
		return '0', false
	} else {
		return _drivedeletter, true
	}
}

func RMDISK_EXECUTE(_drivedeletter byte) {
	PATH := "MIA/P1/Disks/" + string(_drivedeletter) + ".dsk"
	if _, err := os.Stat(PATH); os.IsNotExist(err) {
		color.Red("[RMDISK]: No existe el disco")
		return
	}
	err := os.Remove(PATH)
	if err != nil {
		color.Red("[RMDISK]: Error al borrar el disco")
		return
	}
	color.Green("[RMDISK]: Disco '" + string(_drivedeletter) + ".dsk' Borrado")
}

// Globals

var Max_logical_partitions int64 = 23
var Mounted_prefix int64 = 17

// FDISK

func Get_current_partition(id string) {

}
