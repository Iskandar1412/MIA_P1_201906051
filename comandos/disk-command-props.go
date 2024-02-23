package comandos

import (
	"strings"

	"github.com/fatih/color"
)

//execute -path=/home/iskandar/Escritorio/prueba.adsj

func DiskCommandProps(command string, instructions []string) {
	var _size int32       //mkdisk  fdisk
	var _fit byte         //mkdisk  fdisk
	var _unit byte        //mkdisk  fdisk
	var _driveletter byte //rmdisk  fdisk mount
	var _name [16]byte    //fdisk   mount
	var _type byte        //fdisk
	var _delete string    //fdisk
	var _add int32        //fdisk
	/*
		var _type_mkfs   //mkfs
		var _id string            //unmount mkfs
		var _fs string            //mkfs
	*/

	if strings.ToUpper(command) == "MKDISK" {
		_size, _fit, _unit = Values_MKDISK(instructions)
		if _size <= 0 || _fit == '0' || _unit == '0' {
			color.Yellow("[MKDISK]: Error to asign values for '" + string(_name[:]) + "'")
		} else {
			MKDISK_Create(_size, _fit, _unit)
		}
	} else if strings.ToUpper(command) == "FDISK" {
		_size, _driveletter, _name, _unit, _type, _fit, _delete, _add = Values_FDISK(instructions)
		if _size <= 0 || ToString(_name[:]) == "" || _driveletter == '0' {
			if ToString(_name[:]) == "" {
				color.Yellow("[FDISK]: Error to asign values for (unamed) disk")
			} else {
				if _delete == "FULL" {
					FDISK_Create(_size, _driveletter, _name[:], _unit, _type, _fit, _delete, _add)
				} else if _add != 0 {
					FDISK_Create(_size, _driveletter, _name[:], _unit, _type, _fit, _delete, _add)
				} else {
					color.Yellow("[FDISK]: Error to asign values for disk '" + string(_name[:]) + "'")
				}
			}
		} else {
			FDISK_Create(_size, _driveletter, _name[:], _unit, _type, _fit, _delete, _add)
		}
		//fmt.Println(_size, _driveletter, _name, _unit, _type, _fit, _delete, _add)
	} else if strings.ToUpper(command) == "RMDISK" {
		_driveletter, _error := Values_RMDISK(instructions)
		if _driveletter == '0' && !_error {
			color.Yellow("[RMDISK]: Error to asign values")
		} else {
			RMDISK_EXECUTE(_driveletter)
		}
	} else if strings.ToUpper(command) == "MOUNT" {
		color.Green("Mount")
	} else if strings.ToUpper(command) == "UNMOUNT" {

	} else if strings.ToUpper(command) == "MKFS" {

	} else {
		color.Red("[DiskComandProps]: Internal Error")
	}
}
