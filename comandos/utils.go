package comandos

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"regexp"
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

func ExisteArchivo(comando string, archivo string) bool {
	if _, err := os.Stat(archivo); os.IsNotExist(err) {
		color.Red("[" + comando + "]: Archivo No Encontrado")
		return false
	}
	return true
}

func TieneSize(comando string, size string) int32 {
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
	if strings.ToUpper(value[1]) == "BF" || strings.ToUpper(value[1]) == "B" {
		return 'B'
	} else if strings.ToUpper(value[1]) == "FF" || strings.ToUpper(value[1]) == "F" {
		return 'F'
	} else if strings.ToUpper(value[1]) == "WF" || strings.ToUpper(value[1]) == "W" {
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

func TieneEntero(valor string) int32 {
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
	return int32(i)
}

func ObFechaInt() int32 {
	fecha := time.Now()
	timestamp := fecha.Unix()
	//fmt.Println(timestamp)
	return int32(timestamp)
}

func IntFechaToStr(fecha int32) string {
	conversion := int64(fecha)
	formato := "2006/01/02 (15:04:05)"
	fech := time.Unix(conversion, 0)
	fechaFormat := fech.Format(formato)
	//fmt.Println(fechaFormat)
	return fechaFormat
}

func ObDiskSignature() int32 {
	source := rand.NewSource(time.Now().UnixNano())
	numberR := rand.New(source)
	signature := numberR.Intn(1000000) + 1
	//fmt.Println(signature)
	return int32(signature)
}

func Tamano(size int32, unit byte) int32 {
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
	if !strings.HasPrefix(strings.ToLower(_type), "type=") {
		return '0'
	}
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

func TieneDriveLetter(comando string, deletter string) byte {
	if !strings.HasPrefix(strings.ToLower(deletter), "driveletter=") {
		color.Red("[" + comando + "]: No tiene driveletter o tiene un valor no valido")
		return '0'
	}
	value := strings.Split(deletter, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene deletter Valido")
		return '0'
	} else {
		valor := []byte(value[1])
		if len(valor) > 1 || len(valor) < 1 {
			color.Red("[" + comando + "]: No tiene driveletter Valido")
			fmt.Println(string(valor))
			return '0'
		} else {
			return valor[0]
		}
	}
}

func TieneNombre(comando string, valor string) string {
	//fmt.Println("Valor ingresado:", valor)
	if !strings.HasPrefix(strings.ToLower(valor), "name=") {
		color.Red("[" + comando + "]: No tiene name o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene name Valido")
		return ""
	} else {
		return value[1]
	}
}

func NameArchivosByte(value string) [10]byte {
	padText := make([]byte, 10)
	for i := range padText {
		padText[i] = '\x00'
	}
	copy(padText[:], []byte(value))
	return [10]byte(padText)
}

func IDParticionByte(value string) [4]byte {
	padText := make([]byte, 4)
	for i := range padText {
		padText[i] = '\x00'
	}
	copy(padText[:], []byte(value))
	return [4]byte(padText)
}

func DevolverNombreByte(value string) [16]byte {
	padText := make([]byte, 16)
	for i := range padText {
		padText[i] = '\x00'
	}
	copy(padText[:], []byte(value))
	return [16]byte(padText)
}

func DevolverContenidoJournal(value string) [150]byte {
	padText := make([]byte, 150)
	for i := range padText {
		padText[i] = '\x00'
	}
	copy(padText[:], []byte(value))
	return [150]byte(padText)
}

func DevolverContenidoArchivo(value string) [64]byte {
	padText := make([]byte, 64)
	for i := range padText {
		padText[i] = '\x00'
	}
	copy(padText[:], []byte(value))
	return [64]byte(padText)
}

func TieneTypeFDISK(valor string) byte {
	if !strings.HasPrefix(strings.ToLower(valor), "type=") {
		return '0'
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[FDISK]: No tiene Type Especificado")
		return '0'
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

func TieneDelete(valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "delete=") {
		color.Red("[FDISK]: No tiene Delete Especificado")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[FDISK]: No tiene Delete Especificado")
		return ""
	}
	if !(strings.ToUpper(value[1]) == "FULL") {
		color.Red("[FDISK]: No tiene Delete valido")
		return ""
	}
	return "FULL"
}

func TieneAdd(valor string) int32 {
	if !strings.HasPrefix(strings.ToLower(valor), "add=") {
		return 0
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		return 0
	}
	num, err := strconv.Atoi(value[1])
	if err != nil {
		color.Red("[FDISK]: valor Add no aceptado")
		return 0
	}
	return int32(num)
}

func TieneID(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "id=") {
		color.Red("[" + comando + "]: No tiene id o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene id Valido")
		return ""
	}
	return value[1]
}

func TieneTypeMKFS(valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "type=") {
		color.Red("[MKFS]: No tiene Type Especificado")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[MKFS]: No tiene Type Especificado")
		return ""
	}
	if !(strings.ToUpper(value[1]) == "FULL") {
		color.Red("[MKSF]: No tiene Type valido")
		return ""
	}
	return "FULL"
}

func TieneFS(valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "fs=") {
		color.Red("[MKFS]: No tiene Type Especificado")
		return ""
	}

	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[MKFS]: No tiene Type Especificado")
		return ""
	}

	if !(strings.ToUpper(value[1]) == "3FS" || strings.ToUpper(value[1]) == "2FS") {
		color.Red("[MKSF]: No tiene Type valido")
		return ""
	}

	if (strings.ToUpper(value[1])) == "3FS" {
		return "3FS"
	} else if (strings.ToUpper(value[1])) == "2FS" {
		return "2FS"
	}

	return "2FS"
}

func TieneUser(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "user=") {
		color.Red("[" + comando + "]: No tiene user o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene user Valido")
		return ""
	}
	return value[1]
}

func TienePassword(comando string, valor string) string {
	if !strings.HasPrefix(strings.ToLower(valor), "pass=") {
		color.Red("[" + comando + "]: No tiene password o tiene un valor no valido")
		return ""
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene password Valido")
		return ""
	}
	return value[1]

}

func TieneGRP(comando string, valor string) {
	if !strings.HasPrefix(strings.ToLower(valor), "grp=") {
		color.Red("[" + comando + "]: No tiene grp o tiene un valor no valido")
		return
	}

}

func TientPath(comando string, valor string) {
	if !strings.HasPrefix(strings.ToLower(valor), "path=") {
		color.Red("[" + comando + "]: No tiene path o tiene un valor no valido")
		return
	}
}

func TieneCont(comando string, valor string) {
	if !strings.HasPrefix(strings.ToLower(valor), "cont=") {
		color.Red("[" + comando + "]: No tiene cont o tiene un valor no valido")
		return
	}
}

func TieneCat(comando string, valor string) {
	re := regexp.MustCompile(`file\d+=`)
	if !re.MatchString(valor) {
		color.Red("[" + comando + "]: No tiene fileN o tiene un valor no valido")
		return
	}
}
