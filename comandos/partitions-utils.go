package comandos

func Values_FDISK(instructions []string) (int64, byte, []byte, byte, byte, byte, string, string) {
	var _size int64
	var _drivedeletter byte
	var _name []byte
	var _unit byte
	var _type byte
	var _fit byte
	var _delete string
	var _add string

	return _size, _drivedeletter, _name, _unit, _type, _fit, _delete, _add
}

func FDISK_Create(_size int64, _drivedeletter byte, _name []byte, _unit byte, _type byte, _fit byte, _delete string, _add string) {

}
