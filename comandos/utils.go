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
	if valsize <= 0 {
		color.Red("[" + comando + "]: No tiene Size o tiene un valor no valido")
		return 0
	}
	return valsize
}

func TieneFit(comando string, fit string) byte {
	if !strings.HasPrefix(strings.ToLower(fit), "fit=") {
		color.Red("[" + comando + "]: No tiene Fit o tiene un valor no valido")
		return '0'
	}
	value := strings.Split(fit, "=")
	if len(value) < 2 {
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
		return '0'
	}
}

func TieneUnit(command string, unit string) byte {
	if !strings.HasPrefix(strings.ToLower(unit), "unit=") {
		color.Red("[" + command + "]: No tiene Unit o tiene un valor no valido")
		return '0'
	}
	value := strings.Split(unit, "=")
	if len(value) < 2 {
		color.Red("[" + command + "]: No tiene Unit")
		return '0'
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
		color.Red("[" + command + "]: No tiene Unit Valido")
		return '0'
	}
}

func TieneEntero(valor string) int64 {
	if !strings.HasPrefix(strings.ToLower(valor), "size=") {
		return 0
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
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

func Type_FDISK(_type string) byte {
	value := strings.Split(_type, "=")
	if len(value) < 2 {
		color.Red("[FDISK]: No tiene Type Especificado")
		return 'P'
	}
	if strings.ToUpper(value[1]) == "P" {
		return 'P'
	} else if strings.ToUpper(value[1]) == "E" {
		return 'E'
	} else if strings.ToUpper(value[1]) == "L" {
		return 'L'
	} else {
		color.Red("[FDISK]: No reconocido Type")
		return '0'
	}
}

func Type_MKFS(_type string) string {
	if strings.ToUpper(_type) == "FULL" {
		return "FULL"
	} else {
		color.Red("[MKFS]: No reconocido comando Type")
		return ""
	}
}

func TieneDriveDeLetter(comando string, deletter string) byte {
	if !strings.HasPrefix(strings.ToLower(deletter), "drivedeletter=") {
		color.Red("[" + comando + "]: No tiene deletter o tiene un valor no valido")
		return '0'
	}
	value := strings.Split(deletter, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene deletter Valido")
		return '0'
	} else {
		valor := []byte(value[1])
		if len(valor) > 1 || len(valor) < 1 {
			color.Red("[" + comando + "]: No tiene drivedeletter Valido")
			fmt.Println(string(valor))
			return '0'
		} else {
			return valor[0]
		}
	}
}
