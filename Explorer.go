package main

/*
	#include <stdlib.h>
	typedef void (*dir_callback)(const char* , unsigned long long , char);
	void static CallDirCallBack(dir_callback c_back , const char* path , unsigned long long size , char is_dir) {
		c_back(path , size , is_dir);
	}
*/
import "C"
import "unsafe"

//export GetCurrDirFiles
func GetCurrDirFiles(driveID C.int, cb C.dir_callback) {

	drive_desc, ok := drives[driveID]
	if !ok {
		return
	}
	files := drive_desc.file_explorer.GetItemList()
	for _, file := range files {
		var is_dir C.char
		if file.IsDir {
			is_dir = 1
		} else {
			is_dir = 0
		}
		name := C.CString(file.Name)
		length := C.ulonglong(len(file.Name))
		defer C.free(unsafe.Pointer(name))
		C.CallDirCallBack(cb, name, length, is_dir)
	}
}
