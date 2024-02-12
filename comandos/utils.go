package comandos

import (
	"MIA_P1_201906051/structures"
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func ToByte(str string) []byte {
	result := make([]byte, 1)
	copy(result[:], str)
	return result
}

func ToString(b []byte) string {
	nullIndex := bytes.IndexByte(b, 0)
	if nullIndex == -1 {
		return string(b)
	}
	return string(b[:nullIndex])
}

func TieneSize(comando string, size string) int64 {
	valsize := TieneEntero(size)
	if valsize < 0 {
		fmt.Println("[" + comando + "]: No tiene Size (Obligatorio)")
		return 0
	}
	return valsize
}

func TieneFit(comando string, fit string) byte {
	value := strings.Split(fit, "=")
	if len(value) < 1 {
		return 'F'
	}
	//var val byte = 'F'
	//fmt.Println(val)
	//fmt.Println([]byte(value[1]))
	//fmt.Println(string(val))
	if strings.ToUpper(value[1]) == "BF" {
		return 'B'
	} else if strings.ToUpper(value[1]) == "FF" {
		return 'F'
	} else if strings.ToUpper(value[1]) == "WF" {
		return 'W'
	} else {
		color.Yellow("[" + comando + "]: No tiene Fit Valido")
		return 'F'
	}
}

func TieneUnit(command string, unit string) byte {
	value := strings.Split(unit, "=")
	if len(value) < 1 {
		if command == "fdisk" {
			return 'K'
		} else if command == "mkdisk" {
			return 'M'
		} else {
			color.Red("[" + command + "]: No tiene Unit Valido")
			return 'K'
		}
	}
	if strings.ToUpper(value[1]) == "B" {
		if command == "mkdisk" {
			color.Red("[" + command + "]: No tiene Unit Valido")
			return 'M'
		} else if command == "fdisk" {
			return 'B'
		} else {
			color.Red("[" + command + "]: No tiene Unit Valido")
			return 'K'
		}
	} else if strings.ToUpper(value[1]) == "K" {
		return 'K'
	} else if strings.ToUpper(value[1]) == "M" {
		return 'M'
	} else {
		if command == "mkdisk" {
			color.Red("[" + command + "]: No tiene Unit Valido")
			return 'M'
		} else if command == "fdisk" {
			color.Red("[" + command + "]: No tiene Unit Valido")
			return 'K'
		} else {
			color.Red("[" + command + "]: No tiene Unit Valido")
			return 'K'
		}
	}
}

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

func TieneEntero(valor string) int64 {
	value := strings.Split(valor, "=")
	if len(value) < 1 {
		return 0
	}
	i, err := strconv.Atoi(value[1])
	if err != nil {
		fmt.Println("Error conversion", err)
		return 0
	}
	return int64(i)
}

func ObFechaInt() int64 {
	fecha := time.Now()
	timestamp := fecha.Unix()
	//fmt.Println(timestamp)
	return int64(timestamp)
}

func IntFechaToStr(fecha int64) string {
	formato := "2006-01-02 - 15:04:05"
	fech := time.Unix(fecha, 0)
	fechaFormat := fech.Format(formato)
	//fmt.Println(fechaFormat)
	return fechaFormat
}

func PartitionVacia() structures.Partition {
	var partition structures.Partition
	partition.Part_status = '\x00'
	partition.Part_type = 'P'
	partition.Part_fit = 'F'
	partition.Part_start = -1
	partition.Part_s = -1
	for i := 0; i < len(partition.Part_name); i++ {
		partition.Part_name[i] = '\x00'
	}
	partition.Part_correlative = -1
	for i := 0; i < len(partition.Part_id); i++ {
		partition.Part_id[i] = '\x00'
	}
	return partition
}

func ObDiskSignature() int64 {
	source := rand.NewSource(time.Now().UnixNano())
	numberR := rand.New(source)
	signature := numberR.Intn(1000000) + 1
	//fmt.Println(signature)
	return int64(signature)
}

func Tamano(size int64, unit byte) int64 {
	if unit == 'B' {
		return size
	} else if unit == 'K' {
		return size * 1024
	} else if unit == 'M' {
		return size * 1048576
	} else {
		return 0
	}
}
