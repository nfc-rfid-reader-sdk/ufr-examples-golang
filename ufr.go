package main

import (
	"C"
	"fmt"
	"syscall"
	"unsafe"
	"bufio"
	"os"
	"strings"
)

var ufr_lib = syscall.NewLazyDLL("lib/windows/x86_64/uFCoder-x86_64")

func ReaderOpening(){
	
	fmt.Printf("1. Simple reader open\n")
	fmt.Printf("2. Advanced reader open\n")
	var open_choice int
	fmt.Scanln(&open_choice)
	
	if open_choice == 1 {
		
		var uFCoder = ufr_lib.NewProc("ReaderOpen")
		ret, _, _ := uFCoder.Call()
	
		fmt.Printf("ReaderOpen: 0x%02X\n", ret)
		
		if ret == 0 {
		
			fmt.Printf("==================================\n")
			fmt.Printf("uFR NFC Reader successfully opened\n")
			fmt.Printf("==================================\n")
			fmt.Printf("Press ENTER to read card UID\n")
		}
		
	}else if open_choice == 2 {
		
		fmt.Printf("Reader type:\n")
		var reader_type int
		fmt.Scanln(&reader_type)
		
		fmt.Printf("Port name:\n")
		reader := bufio.NewReader(os.Stdin)
		var input_port_name string
		input_port_name, _ = reader.ReadString('\n')
		var port_name = strings.TrimRight(input_port_name, "\r\n")
		
		fmt.Printf("Port interface:\n")
		var input_port_interface string
		input_port_interface, _ = reader.ReadString('\n')
		var port_interface int
		
		if strings.TrimRight(input_port_interface, "\r\n") == "T" {
		
			port_interface = 84
		} else if strings.TrimRight(input_port_interface, "\r\n") == "U" {
		
			port_interface = 85
		}
		
		fmt.Printf("Arg:\n")
		var input_arg string
		input_arg, _ = reader.ReadString('\n')
		var arg = strings.TrimRight(input_arg, "\r\n")
		
		port_name_CS := C.CString(port_name)
		argCS := C.CString(arg)
		
		var uFCoder = ufr_lib.NewProc("ReaderOpenEx")
		ret, _, _ := uFCoder.Call(uintptr(byte(reader_type)), uintptr(unsafe.Pointer(port_name_CS)), uintptr(byte(port_interface)), uintptr(unsafe.Pointer(argCS)))
		
		if ret == 0 {
		
			fmt.Printf("==================================\n")
			fmt.Printf("uFR NFC Reader successfully opened\n")
			fmt.Printf("==================================\n")
			fmt.Printf("Press ENTER to read card UID\n")
		}
	}
}

func ReaderUISignal(light byte, beep byte){

	var uFCoder = ufr_lib.NewProc("ReaderUISignal")
	uFCoder.Call(uintptr(light), uintptr(beep))
}

func GetCardIdEx(sak byte, uid [10]byte, uidSize byte){

	var uFCoder = ufr_lib.NewProc("GetCardIdEx")
	ret, _, _ := uFCoder.Call(uintptr(unsafe.Pointer(&sak)), uintptr(unsafe.Pointer(&uid)), uintptr(unsafe.Pointer(&uidSize)))
	
	if ret == 0 {
		
		fmt.Printf("CARD DETECTED [type = 0x%02X", sak)
		fmt.Printf(", uid[%d] = ", uidSize)
		
		for i := 0; i < int(uidSize); i++ {
			fmt.Printf("%02X", uid[i])
		}
		
		fmt.Printf("]")
	
	} else if ret == 8 {
	
		fmt.Printf("NO CARD DETECTED")
	
	} else {
			
		fmt.Printf("GetCardIdEx(): 0x%02X", ret)
	}
}

func HandleUID(){
	
	reader := bufio.NewReader(os.Stdin)
	var input string
    input, _ = reader.ReadString('\n')
	
	if strings.TrimRight(input, "\r\n") == "" {
		
		var sak byte
		var uid [10]byte
		var uidSize byte
		
	    GetCardIdEx(sak, uid, uidSize)
	}
}

func main() {
	
	ReaderOpening()
	ReaderUISignal(1, 1)
	
	for {
		
		HandleUID()
	}

}