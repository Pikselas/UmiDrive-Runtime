package main

import "C"
import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Pikselas/Octodrive/Octo"
	"github.com/Pikselas/Octodrive/Octo/ToOcto"
)

type DriveDesc struct {
	drive *Octo.OctoDrive
	files map[string]*Octo.OctoFile
}

var drives = make(map[C.int]DriveDesc)

//export GetFile
func GetFile(driveID C.int, path *C.char, l C.int) {
	drive_desc, ok := drives[driveID]
	if !ok {
		return
	}
	path_str := C.GoStringN((*C.char)(path), l)
	file, ok := drive_desc.files[path_str]
	if !ok {
		var err error
		file, err = drive_desc.drive.Load(path_str)
		if err != nil {
			return
		}
		drive_desc.files[path_str] = file
	}
	file.GetSize()
}

//export LoadDriver
func LoadDriver(token *C.char, l1 C.int, email *C.char, l2 C.int) C.int {
	token_str := C.GoStringN((*C.char)(token), l1)
	email_str := C.GoStringN((*C.char)(email), l2)

	fmt.Println(email_str)
	fmt.Println(token_str)

	user, err := ToOcto.NewOctoUser(email_str, token_str)
	if err != nil {
		return C.int(-1)
	}
	drive, octo_err := Octo.NewOctoDrive(user, Octo.DefaultFileRegistry)
	if octo_err != nil {
		return C.int(-1)
	}
	rand_src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(rand_src)
	ID := C.int(random.Int31())
	drives[ID] = DriveDesc{drive, make(map[string]*Octo.OctoFile)}
	return ID
}

//export UnloadDriver
func UnloadDriver(ID C.int) {
	delete(drives, ID)
}

func main() {}
