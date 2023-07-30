package main

import "C"
import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

//export LoadFile
func LoadFile(driveID C.int, path *C.char, l C.int) C.int {
	drive_desc, ok := drives[driveID]
	if !ok {
		return C.int(-1)
	}
	path_str := C.GoStringN((*C.char)(path), l)
	file, ok := drive_desc.files[path_str]
	if !ok {
		var err error
		file, err = drive_desc.drive.Load(path_str)
		if err != nil {
			return C.int(-1)
		}
		drive_desc.files[path_str] = file
	}
	reader, err := file.GetReader()
	if err != nil {
		return C.int(-1)
	}
	ID := RandomID()
	drive_desc.file_readers[ID] = reader
	fmt.Println("load file", ID)
	return ID
}

var file, _ = os.OpenFile("D:/Virgin's First Love.mp4", os.O_RDONLY, 0666)

//export FLD
func FLD(driveID, fileID C.int, buf unsafe.Pointer, size C.ulonglong) C.longlong {
	buffs := (*[1 << 20]byte)(buf)[:size]
	n, err := file.Read(buffs)
	if err != nil {
		fmt.Println("FLD ERROR", err)
		return C.longlong(-1)
	}
	//fmt.Println("FLD READ", bufs)
	return C.longlong(n)
}

//export FST
func FST(driveID, fileID C.int) {
	fmt.Println("FST WRITING")
	drive_desc, ok := drives[driveID]
	if !ok {
		fmt.Println("FST NO DRIVE")
	}
	reader, ok := drive_desc.file_readers[fileID]
	if !ok {
		fmt.Println("FST NO FILE")
	}
	file, err := os.Create("pattern.mp4")
	if err != nil {
		fmt.Println("FST ERROR CREATION OF FILE")
	}
	defer file.Close()
	n, err := io.Copy(file, reader)
	if err != nil {
		fmt.Println("FST ERROR COPYING")
	}
	fmt.Println("FST WRITTEN", n)
}

//export CloseFile
func CloseFile(driveID, fileID C.int) {
	drive_desc, ok := drives[driveID]
	if !ok {
		return
	}
	reader, ok := drive_desc.file_readers[fileID]
	if !ok {
		return
	}
	reader.Close()
	delete(drive_desc.file_readers, fileID)
}
