package main

import "C"
import (
	"fmt"
	"io"
	"unsafe"

	"github.com/Pikselas/Octodrive/Octo"
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
	enc_dec := Octo.NewAesEncDecFrom(file.GetUserData())
	file.SetEncDec(enc_dec)
	reader, err := file.GetReader()
	if err != nil {
		return C.int(-1)
	}
	ID := RandomID()
	drive_desc.file_readers[ID] = reader
	fmt.Println("load file", ID)
	return ID
}

//export ReadLoadedFile
func ReadLoadedFile(driveID, fileID C.int, buf unsafe.Pointer, size C.ulonglong) C.longlong {
	drive_desc, ok := drives[driveID]
	if !ok {
		return C.longlong(-1)
	}
	reader, ok := drive_desc.file_readers[fileID]
	if !ok {
		return C.longlong(-1)
	}

	// Casting to a 1GB array (it doesn't actually allocate 1GB of memory) and then slicing it.
	// This is a workaround for the fact that Go doesn't allow casting to a slice

	n, err := reader.Read((*[1 << 30]byte)(buf)[:size])
	if err == io.EOF {
		return C.longlong(-2)
	} else if err != nil {
		return C.longlong(-1)
	}

	return C.longlong(n)
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
