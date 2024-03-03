package comandos

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func UserComanmandProps(command string, instructions []string) {
	var _user string //login
	var _pass string //login
	var _id string   //login
	/*
	 */
	if strings.ToUpper(command) == "LOGIN" {
		// var err bool
		valor_usuario, err := Values_LOGIN(instructions)
		_user = ToString(valor_usuario.User[:])
		_pass = ToString(valor_usuario.Password[:])
		_id = ToString(valor_usuario.Id_Particion[:])
		if !err {
			color.Red("[LOGIN]: Error to asign values")
		} else {
			LOGIN_EXECUTE(_user, _pass, _id)
		}
		//fmt.Println("Loggeando Usuario")
	} else if strings.ToUpper(command) == "LOGOUT" {
		LOGOUT_EXECUTE()
	} else {
		color.Red("[UserComanmandProps]: Internal Error")
	}

}

func GroupCommandProps(group string, instructions []string) {
	var _name string //mkgrp rmgrp
	var er bool
	/*
		var _user string //mkusr rmusr
		var _pass string //mkusr
		var _grp string  //mkusr
	*/
	if strings.ToUpper(group) == "MKGRP" {
		_name, er = Values_MKGRP(instructions)
		if !er {
			color.Red("[MKGRP]: Error to assing values")
		} else {
			MKGRP_EXECUTE(_name)
		}
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
