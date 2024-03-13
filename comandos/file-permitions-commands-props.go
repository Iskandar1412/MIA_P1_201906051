package comandos

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func FilesCommandProps(file string, instructions []string) {
	var _er bool
	var _path string //mkfile remove edit rename mkdir copy move find
	//existe -r 		//mkfile mkdir
	var _r bool
	var _size int32  //mkfile
	var _cont string //mkfile
	/*
		var _file []string  //cat
		var _cont string    //edit
		var _name string    //rename find
		var _destino string //move   copy
	*/
	if strings.ToUpper(file) == "MKFILE" {
		_path, _r, _size, _cont, _er = Values_MKFILE(instructions)
		if !_er || _path == "" || len(_path) == 0 {
			color.Yellow("[MKFILE]: Error to asign vallues")
		} else {
			MKFILE_EXECUTE(_path, _r, _size, _cont)
		}

	} else if strings.ToUpper(file) == "CAT" {

	} else if strings.ToUpper(file) == "REMOVE" {

	} else if strings.ToUpper(file) == "EDIT" {

	} else if strings.ToUpper(file) == "RENAME" {

	} else if strings.ToUpper(file) == "MKDIR" {

	} else if strings.ToUpper(file) == "COPY" {

	} else if strings.ToUpper(file) == "MOVE" {

	} else if strings.ToUpper(file) == "FIND" {

	} else {
		color.Red("[FilesCommandProps]: Internal Error")
	}
}

func PermissionsCommandProps(permission string, instructions []string) {
	/*
		var _path string //chown chmod
		var _user string //chown
		//_r			 //chown chmod
		var _user string //chgrp
		var _grp string  //chgrp
		var _ugo string  //chmod
	*/
	if strings.ToUpper(permission) == "CHOWN" {
		fmt.Println("CHOWN")
	} else if strings.ToUpper(permission) == "CHGRP" {
		fmt.Println("CHGRP")
	} else if strings.ToUpper(permission) == "CHMOD" {
		fmt.Println("CHMOD")
	} else if strings.ToUpper(permission) == "PAUSE" {
		fmt.Println("PAUSE")
	} else {
		color.Red("[PermissionsCommandProps]: Internal Error")
	}
}
