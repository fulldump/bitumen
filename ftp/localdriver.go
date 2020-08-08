package ftp

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type LocalDriver struct {
	BasePath string
	ReadOnly bool
	AuthHook AuthHook
}

func (d *LocalDriver) absolutePath(p string) string {

	apath := path.Join(d.BasePath, p)

	if !strings.HasPrefix(apath, d.BasePath) {
		return d.BasePath
	}

	return apath
}

func (d *LocalDriver) Authenticate(username string, password string) bool {
	return d.AuthHook(username, password, d)
}

func (d *LocalDriver) Bytes(name string) int64 {

	name = d.absolutePath(name)

	info, err := os.Stat(name)
	if err != nil {
		fmt.Println("ERROR", err)
		return -1
	}

	return info.Size()
}

func (d *LocalDriver) ModifiedTime(name string) (time.Time, bool) {

	name = d.absolutePath(name)

	info, err := os.Stat(name)
	if err != nil {
		fmt.Println("ERROR", err)
		return time.Now(), false
	}

	return info.ModTime(), true

}

func (d *LocalDriver) ChangeDir(name string) bool {

	name = d.absolutePath(name)

	i, err := os.Stat(name)
	if err != nil {
		fmt.Println("ERROR", err)
		return false
	}
	return i.IsDir()
}

func (d *LocalDriver) DirContents(name string) ([]os.FileInfo, bool) {

	name = d.absolutePath(name)

	i, err := os.Stat(name)
	if err != nil {
		fmt.Println("ERROR", err)
		return nil, false
	}

	if !i.IsDir() {
		return nil, false
	}

	files, err := ioutil.ReadDir(name)
	if err != nil {
		return nil, false
	}

	// sort.Sort(&FilesSorter{files}) // TODO: sort files?

	return files, true
}

func (d *LocalDriver) DeleteDir(name string) bool {

	if d.ReadOnly {
		return false
	}

	name = d.absolutePath(name)

	fmt.Println("DeleteDir not implemented", name)

	return false // TODO: implement this

	//if f, ok := d.Files[path]; ok && f.File.IsDir() {
	//	haschildren := false
	//	for p, _ := range d.Files {
	//		if strings.HasPrefix(p, path+"/") {
	//			haschildren = true
	//			break
	//		}
	//	}
	//
	//	if haschildren {
	//		return false
	//	}
	//
	//	delete(d.Files, path)
	//
	//	return true
	//} else {
	//	return false
	//}
}

func (d *LocalDriver) DeleteFile(name string) bool {

	if d.ReadOnly {
		return false
	}

	name = d.absolutePath(name)

	fmt.Println("DeleteFile not implemented", name)

	return false // TODO: implement this

	//if f, ok := d.Files[path]; ok && !f.File.IsDir() {
	//	delete(d.Files, path)
	//	return true
	//} else {
	//	return false
	//}
}

func (d *LocalDriver) Rename(from_name string, to_name string) bool {

	if d.ReadOnly {
		return false
	}

	from_name = d.absolutePath(from_name)
	to_name = d.absolutePath(to_name)

	fmt.Println("Rename not implemented", from_name, to_name)

	return false // TODO: implement this

	//if f, from_path_exists := d.Files[from_path]; from_path_exists {
	//	if _, to_path_exists := d.Files[to_path]; !to_path_exists {
	//		if _, to_path_parent_exists := d.Files[filepath.Dir(to_path)]; to_path_parent_exists {
	//			if f.File.IsDir() {
	//				delete(d.Files, from_path)
	//				d.Files[to_path] = &MemoryFile{graval.NewDirItem(filepath.Base(to_path)), nil}
	//				torename := make([]string, 0)
	//				for p, _ := range d.Files {
	//					if strings.HasPrefix(p, from_path+"/") {
	//						torename = append(torename, p)
	//					}
	//				}
	//				for _, p := range torename {
	//					sf := d.Files[p]
	//					delete(d.Files, p)
	//					np := to_path + p[len(from_path):]
	//					d.Files[np] = sf
	//				}
	//			} else {
	//				delete(d.Files, from_path)
	//				d.Files[to_path] = &MemoryFile{graval.NewFileItem(filepath.Base(to_path), f.File.Size(), f.File.ModTime()), f.Content}
	//			}
	//			return true
	//		} else {
	//			return false
	//		}
	//	} else {
	//		return false
	//	}
	//} else {
	//	return false
	//}
}

func (d *LocalDriver) MakeDir(name string) bool {

	if d.ReadOnly {
		return false
	}

	name = d.absolutePath(name)

	fmt.Println("MakeDir not implemented", name)

	return false // TODO: implement this

	//if _, ok := d.Files[path]; ok {
	//	return false
	//} else {
	//	d.Files[path] = &MemoryFile{graval.NewDirItem(filepath.Base(path)), nil}
	//	return true
	//}
}

func (d *LocalDriver) GetFile(name string, position int64) (io.ReadCloser, bool) {

	name = d.absolutePath(name)

	i, err := os.Stat(name)
	if err != nil {
		fmt.Println("ERROR", err)
		return nil, false
	}

	if i.IsDir() {
		return nil, false
	}

	f, err := os.Open(name)
	if err != nil {
		fmt.Println("ERROR", err)
		return nil, false
	}

	p, err := f.Seek(position, 0)
	if err != nil {
		fmt.Println("ERROR", err)
		return nil, false
	}

	if p != position {
		fmt.Println("ERROR", err)
		return nil, false
	}

	return f, true
}

func (d *LocalDriver) PutFile(name string, src io.Reader) bool {

	if d.ReadOnly {
		return false
	}

	name = d.absolutePath(name)

	f, err := os.Create(name)
	if err != nil {
		fmt.Println("ERROR", err)
		return false
	}

	n, err := io.Copy(f, src)
	if err != nil {
		fmt.Println("ERROR", err)
		return false
	}

	fmt.Printf("INFO: %d bytes copied, path: %s", n, name)

	return true
}
