package main

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#include <memory.h>
*/
import "C"
import (
	"unsafe"

	"cosmossdk.io/depinject"
)

//export Supply
func Supply(goRef GoRef, cErr **C.char) C.go_ref {
	toSupply, err := GetGoMem[any](goRef)
	if err != nil {
		*cErr = C.CString(err.Error())
		return 0
	}

	return C.go_ref(SetGoMem(depinject.Supply(toSupply)))
}

//export SupplyMany
func SupplyMany(goRefs *C.go_ref, numGoRefs C.int, cErr **C.char) C.go_ref {
	refs := unsafe.Slice(goRefs, numGoRefs)

	//fmt.Printf(">>> refs: %+v\n", refs)
	//val, err := GetGoMem[any](refs[0])
	//if err != nil {
	//	*cErr = C.CString(err.Error())
	//	return 0
	//}
	//fmt.Printf(">>> ref: %+v\n", val)

	var toSupply []any
	for _, ref := range refs {
		valueToSupply, err := GetGoMem[any](GoRef(ref))
		if err != nil {
			*cErr = C.CString(err.Error())
			//*cErr = C.CString(fmt.Sprintf("%+v", err))
			return C.go_ref(NilGoRef)
		}

		toSupply = append(toSupply, valueToSupply)
	}

	return C.go_ref(SetGoMem(depinject.Supply(toSupply...)))
}

//export Config
func Config(goRefs *C.go_ref, numGoRefs C.int, cErr **C.char) C.go_ref {
	refs := unsafe.Slice(goRefs, numGoRefs)

	var configs []depinject.Config
	for _, ref := range refs {
		cfg, err := GetGoMem[depinject.Config](GoRef(ref))
		if err != nil {
			*cErr = C.CString(err.Error())
			return 0
		}

		configs = append(configs, cfg)
	}

	return C.go_ref(SetGoMem(depinject.Configs(configs...)))
}
