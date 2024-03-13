package comandos

import "strconv"

func Validar_Permisos(comando string, id_user_padre int32, id_grupo_padre int32, id_user int32, id_grupo int32, permisos int32, accion int32) bool {
	permiso_user := int32(0)
	permiso_grupo := int32(0)
	permiso_otro := int32(0)
	temp := strconv.Itoa(int(permisos))
	if permisos > int32(99) {
		permiso_usera, _ := strconv.Atoi(string(temp[0]))
		permiso_user = int32(permiso_usera)
		permiso_grupoa, _ := strconv.Atoi(string(temp[1]))
		permiso_grupo = int32(permiso_grupoa)
		permiso_otroa, _ := strconv.Atoi(string(temp[2]))
		permiso_otro = int32(permiso_otroa)
	} else if permisos > 9 {
		permiso_user = int32(0)
		permiso_grupoa, _ := strconv.Atoi(string(temp[0]))
		permiso_grupo = int32(permiso_grupoa)
		permiso_otroa, _ := strconv.Atoi(string(temp[1]))
		permiso_otro = int32(permiso_otroa)
	} else if permisos > 0 {
		permiso_user = int32(0)
		permiso_grupo = int32(0)
		permiso_otroa, _ := strconv.Atoi(string(temp[0]))
		permiso_otro = int32(permiso_otroa)
	} else {
		permiso_user = int32(0)
		permiso_grupo = int32(0)
		permiso_otro = int32(0)
	}

	//validación del usuario root
	acion := int(accion)
	if id_user == int32(1) && id_grupo == int32(1) {
		//como es root
		return true
	}

	//En caso de ser el usuario
	if id_user_padre == id_user {
		//validar permiso user
		if permiso_user == int32(7) {
			return true
		} else if permiso_user == int32(6) {
			if acion == 2 || acion == 3 || acion == 4 || acion == 5 || acion == 6 {
				return true
			}
		} else if permiso_user == int32(5) {
			if acion == 1 || acion == 4 || acion == 5 {
				return true
			}
		} else if permiso_user == int32(4) {
			if acion == 4 {
				return true
			}
		} else if permiso_user == int32(3) {
			if acion == 1 || acion == 2 || acion == 3 {
				return true
			}
		} else if permiso_user == int32(2) {
			if acion == 2 {
				return true
			}
		} else if permiso_user == int32(1) {
			if acion == 1 {
				return true
			}
		} else if permiso_user == int32(0) {
			return false
		}
	} else if id_grupo_padre == id_grupo {
		//permisos de grupo
		if permiso_grupo == int32(7) {
			return true
		} else if permiso_grupo == int32(6) {
			if acion == 2 || acion == 3 || acion == 4 || acion == 5 || acion == 6 {
				return true
			}
		} else if permiso_grupo == int32(5) {
			if acion == 1 || acion == 4 || acion == 5 {
				return true
			}
		} else if permiso_grupo == int32(4) {
			if acion == 4 {
				return true
			}
		} else if permiso_grupo == int32(3) {
			if acion == 1 || acion == 2 || acion == 3 {
				return true
			}
		} else if permiso_grupo == int32(2) {
			if acion == 2 {
				return true
			}
		} else if permiso_grupo == int32(1) {
			if acion == 1 {
				return true
			}
		} else if permiso_grupo == int32(0) {
			return false
		}
	} else if id_grupo_padre != id_grupo { // otro
		if permiso_otro == int32(7) {
			return true
		} else if permiso_otro == int32(6) {
			if acion == 2 || acion == 3 || acion == 4 || acion == 5 || acion == 6 {
				return true
			}
		} else if permiso_otro == int32(5) {
			if acion == 1 || acion == 4 || acion == 5 {
				return true
			}
		} else if permiso_otro == int32(4) {
			if acion == 4 {
				return true
			}
		} else if permiso_otro == int32(3) {
			if acion == 1 || acion == 2 || acion == 3 {
				return true
			}
		} else if permiso_otro == int32(2) {
			if acion == 2 {
				return true
			}
		} else if permiso_otro == int32(1) {
			if acion == 1 {
				return true
			}
		} else if permiso_otro == int32(0) {
			return false
		}
	}
	return false
}
