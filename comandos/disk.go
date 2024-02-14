package comandos

import (
	"strings"

	"github.com/fatih/color"
)

//execute -path=/home/iskandar/Escritorio/prueba.adsj

func DiskCommandProps(command string, instructions []string) {
	var _size int64 //mkdisk  fdisk
	var _fit byte   //mkdisk  fdisk
	var _unit byte  //mkdisk  fdisk
	/*
		var _drivedeletter byte   //rmdisk  fdisk mount
		var _name string          //fdisk   mount
		var _type string          //fdisk   mkfs
		var _delete string        //fdisk
		var _add string           //fdisk
		var _id string            //unmount mkfs
		var _fs string            //mkfs
	*/

	if strings.ToUpper(command) == "MKDISK" {
		_size, _fit, _unit = Values_MKDISK(instructions)
		if _size <= 0 || _fit == '0' || _unit == '0' {
			color.Yellow("[MKDISK]: Error to asign values")
		} else {
			MKDISK_Create(_size, _fit, _unit)
		}
	} else if strings.ToUpper(command) == "FDISK" {

	} else if strings.ToUpper(command) == "RMDISK" {
		_drivedeletter, _error := Values_RMDISK(instructions)
		if _drivedeletter == '0' && _error == false {
			color.Yellow("[RMDISK]: Error to asign values")
		} else {
			RMDISK_EXECUTE(_drivedeletter)
		}
	} else if strings.ToUpper(command) == "MOUNT" {

	} else if strings.ToUpper(command) == "UNMOUNT" {

	} else if strings.ToUpper(command) == "MKFS" {

	} else {
		color.Red("[DiskComandProps]: Internal Error")
	}
}
