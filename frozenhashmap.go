package frozenhashmap

// #cgo pkg-config: kyotocabinet frozenhashmap
// #include <stdlib.h>
// #include <kclangc.h>
// #include <string.h>
// #include <frozenhashmap/cfhm.h>
import "C"
import "unsafe"
import "fmt"

type FrozenHashMapBuilder struct {
	v *C.struct_CFrozenHashMapBuilder;
}

type FrozenHashMap struct {
	v *C.struct_CFrozenHashMap;
}

func FrozenHashMapBuilderOpen(inmemory bool) (builder *FrozenHashMapBuilder, err error) {
	p, _ := C.CFrozenHashMapBuilderAllocate(C._Bool(inmemory))
	if p == nil {err = fmt.Errorf("Cannot allocate FrozenHashMap Builder"); return}
	ok, _ := C.CFrozenHashMapBuilderOpen(p)
	if ok == false {err = fmt.Errorf("Cannot open FrozenHashMap Builder"); return}
	builder = &FrozenHashMapBuilder{p}
	return
}

func (builder *FrozenHashMapBuilder) PutString(key string, value string) (err error) {
	keyCstr := C.CString(key)
	defer C.free(unsafe.Pointer(keyCstr))
	valueCstr := C.CString(value)
	defer C.free(unsafe.Pointer(valueCstr))

	ok, _ := C.CFrozenHashMapBuilderPutString(builder.v, keyCstr, valueCstr)
	if (!ok) {
		err = fmt.Errorf("Some error")
	}
	return
}

func (builder *FrozenHashMapBuilder) Build(path string) (err error) {
	pathCstr := C.CString(path)
	defer C.free(unsafe.Pointer(pathCstr))

	ok, _ := C.CFrozenHashMapBuilderBuild(builder.v, pathCstr)
	if (!ok) {err = fmt.Errorf("Cannot build index")}
	return
}

func (builder *FrozenHashMapBuilder) Free() {
	C.CFrozenHashMapBuilderFree(builder.v)
}

func FrozenHashMapOpen(path string) (hashmap *FrozenHashMap, err error) {
	p, _ := C.CFrozenHashMapAllocate()
	if p == nil {err = fmt.Errorf("Cannot allocate FrozenHashMap"); return}

	pathCstr := C.CString(path)
	defer C.free(unsafe.Pointer(pathCstr))
	ok, _ := C.CFrozenHashMapOpen(p, pathCstr)
	if ok != true {err = fmt.Errorf("Cannot open database"); return}
	hashmap = &FrozenHashMap{p}
	return
}

func (hashmap *FrozenHashMap) GetString(key string) (value string, err error) {
	keyCstr := C.CString(key)
	keyClen := C.strlen(keyCstr)
	var valueLen C.size_t
	valueCstr := C.CFrozenHashMapGet(hashmap.v, keyCstr, keyClen, &valueLen)
	if (valueCstr == nil) {err = fmt.Errorf("Cannot get string"); return}
	value = C.GoStringN(valueCstr, C.int(valueLen))
	return
}

func (hashmap *FrozenHashMap) Free() {
	C.CFrozenHashMapFree(hashmap.v)
}
