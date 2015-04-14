package frozenhashmap

import "testing"
import "io/ioutil"
import "os"
import "path"


type testingInfo struct {
	t *testing.T
	tmpDir string
	dbPath string
}

func tearUp(t *testing.T) (info testingInfo) {
	os.Setenv("FROZENHASH_DEBUG", "1")
	os.Setenv("TMPDIR", "./tmp")
	os.Mkdir("./tmp", 0755)
	
	tmpDir, err := ioutil.TempDir("", "frozenhashmap")
	if err != nil {t.Errorf("Cannot create temporally directory: %s", err)}
	dbPath := path.Join(tmpDir, "hashmap.db")
	info = testingInfo{t, tmpDir, dbPath}
	return
}

func (info *testingInfo) tearDown()  {
	os.Remove(info.dbPath)
	os.Remove(info.tmpDir)
}

func TestDB(t *testing.T) {
	info := tearUp(t)

	// setup data
	builder, err := FrozenHashMapBuilderOpen(false)
	if (err != nil) {t.Error("Cannot open hashmap builder")}
	err = builder.PutString("Hi", "OK-Hi!")
	if (err != nil) {t.Error("Cannot put string 1")}
	err = builder.PutString("Push", "OK-Push!")
	if (err != nil) {t.Error("Cannot put string 2")}
	err = builder.PutString("Echo", "OK-Echo!")
	if (err != nil) {t.Error("Cannot put string 3")}
	err = builder.Build(info.dbPath)
	if (err != nil) {t.Error("Cannot build database")}
	builder.Free()

	// read data
	hashmap, err := FrozenHashMapOpen(info.dbPath)
	if (err != nil) {t.Error("Cannot open hashmap")}

	value, err := hashmap.GetString("Hi")
	if (err != nil) {t.Error("Cannot get value 1")}
	if (value != "OK-Hi!") {t.Error("Invalid data 1")}

	value, err = hashmap.GetString("Push")
	if (err != nil) {t.Error("Cannot get value 2")}
	if (value != "OK-Push!") {t.Error("Invalid data 2")}

	value, err = hashmap.GetString("Echo")
	if (err != nil) {t.Error("Cannot get value 3")}
	if (value != "OK-Echo!") {t.Error("Invalid data 3")}

	value, err = hashmap.GetString("Unknown")
	if (err == nil) {t.Error("No value should found for this")}

	value, err = hashmap.GetString("Unknown1")
	if (err == nil) {t.Error("No value should found for this")}

	value, err = hashmap.GetString("Unknown2")
	if (err == nil) {t.Error("No value should found for this")}

	value, err = hashmap.GetString("Unknown3")
	if (err == nil) {t.Error("No value should found for this")}

	
	hashmap.Free()
}
