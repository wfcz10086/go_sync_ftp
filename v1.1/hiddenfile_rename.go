package main

import (
	"os"
	"path"
	"strings"
	"path/filepath"
)
func WalkDir(dirPth, suffix string) (files []string, err error) {
 files = make([]string, 0, 30)
 suffix = strings.ToUpper(suffix)
 err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
   if err != nil {
   return err
  }
  if fi.IsDir() {
        return nil
 }


  if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
   if !fi.IsDir() {
        dir := filepath.Dir(filename)
	new_file_name := strings.SplitAfter(dir, "/")[2]
        fullFilename := filepath.Base(filename)
	var filenameWithSuffix string
    	filenameWithSuffix = path.Base(fullFilename)
    	var fileSuffix string
        fileSuffix = path.Ext(filenameWithSuffix)

    	var filenameOnly string
    	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
     	if filenameOnly  == "" {
	new_file_name_m :=dir+"/"+new_file_name+".log"
    	err := os.Rename(filename, new_file_name_m)
	   if err != nil {
            print("ok")
   	 } else {
        	print(filename+"file rename OK!\r\n")
    	}
        }

 }

   files = append(files, filename)
  }
  return nil
 })
 return files, err
}





func main() {
	WalkDir("/syslog",".log" )


}
