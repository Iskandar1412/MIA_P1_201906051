package comandos

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func UserComanmandProps(user string, instructions []string) {
	/*
		var _user string //login
		var _pass string //login
		var _id string   //login
	*/
	if strings.ToUpper(user) == "LOGIN" {
		fmt.Println("Loggeando Usuario")
	} else if strings.ToUpper(user) == "LOGOUT" {
		fmt.Println("Saliendo del sistema")
	} else {
		color.Red("[UserComanmandProps]: Internal Error")
	}
}

func GroupCommandProps(group string, instructions []string) {
	/*
		var _name string //mkgrp rmgrp
		var _user string //mkusr rmusr
		var _pass string //mkusr
		var _grp string  //mkusr
	*/
	if strings.ToUpper(group) == "MKGRP" {
		fmt.Println("Creando Grupo")
	} else if strings.ToUpper(group) == "RMGRP" {
		fmt.Println("Eliminando Grupo")
	} else if strings.ToUpper(group) == "MKUSR" {
		fmt.Println("Crear usuario en particion")
	} else if strings.ToUpper(group) == "RMUSR" {
		fmt.Println("Eliminando usuario en la parcicion")
	} else {
		color.Red("[GroupCommandProps]: Internal Error")
	}
}
