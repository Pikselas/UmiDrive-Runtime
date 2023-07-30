package main

import "C"

import (
	"io"
	"math/rand"
	"time"

	"github.com/Pikselas/Octodrive/Octo"
	"github.com/Pikselas/Octodrive/Octo/ToOcto"
)

type DriveDesc struct {
	drive         *Octo.OctoDrive
	files         map[string]*Octo.OctoFile
	file_readers  map[C.int]io.ReadCloser
	file_explorer *Octo.FileNavigator
}

var drives = make(map[C.int]DriveDesc)

func RandomID() C.int {
	rand_src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(rand_src)
	return C.int(random.Int31())
}

//export LoadDrive
func LoadDrive(token *C.char, l1 C.int, email *C.char, l2 C.int) C.int {
	token_str := C.GoStringN((*C.char)(token), l1)
	email_str := C.GoStringN((*C.char)(email), l2)

	user, err := ToOcto.NewOctoUser(email_str, token_str)
	if err != nil {
		return C.int(-1)
	}
	drive, octo_err := Octo.NewOctoDrive(user, Octo.DefaultFileRegistry)
	if octo_err != nil {
		return C.int(-1)
	}
	ID := RandomID()
	file_nav, octo_err := drive.NewFileNavigator()
	if octo_err != nil {
		return C.int(-1)
	}
	drives[ID] = DriveDesc{drive, make(map[string]*Octo.OctoFile), make(map[C.int]io.ReadCloser), file_nav}
	return ID
}

//export UnloadDrive
func UnloadDrive(ID C.int) {
	delete(drives, ID)
}
