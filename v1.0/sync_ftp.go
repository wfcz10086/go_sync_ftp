package main

import (
    "bufio"
    "path/filepath"
    "io"
    "os"
    "strings"
    ftp "github.com/ftp"
)

type Config struct {
    Mymap  map[string]string
    strcet string
}


func ftpUploadFile(ftpserver, ftpuser, pw, localFile, remoteSavePath, saveName string) {
    ftpConfig := new(Config)
    ftpConfig.InitConfig("config.ini")
    ftpfile_path := ftpConfig.Read("ftp", "ftpfile_path")
    ftp, err := ftp.Connect(ftpserver)
    if err != nil {
        print("System, connect err")
	os.Exit(-1)
    }
    err = ftp.Login(ftpuser, pw)
    if err != nil {
        print("System, Login err \r\n")
	os.Exit(-1)
    }
    //ftp.MakeDir(remoteSavePath)
    ftp.ChangeDir(ftpfile_path)
    //dir, err := ftp.CurrentDir()
    ftp.MakeDir(remoteSavePath)
    //print("make,",remoteSavePath,"\r\n")
    ftp.ChangeDir(remoteSavePath)
   // dir, _ = ftp.CurrentDir()
   // print("System %s ", dir,"\r\n")
    file, err := os.Open(localFile)
    if err != nil {
        print("System,Open err \r\n")
    }
    defer file.Close()
    err = ftp.Stor(saveName, file)
    if err != nil {
        print("System, Stor err \r\n")
    }
    ftp.Logout()
    ftp.Quit()
    print("success upload ",localFile,"\r\n")
}

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
  if !fi.IsDir() {
        dir := filepath.Dir(filename)
	syncname := filepath.Base(filename)
        ftpConfig := new(Config)
	//print(syncname+"\r\n")
        ftpConfig.InitConfig("config.ini")
	ftpserverip := ftpConfig.Read("ftp", "ftp_server_ip")
    	ftpPort := ftpConfig.Read("ftp", "ftp_server_port")
    	ftpuser := ftpConfig.Read("ftp", "ftp_server_name")
    	pw := ftpConfig.Read("ftp", "ftp_server_pwd")
    	ftpserver := ftpserverip + ":" + ftpPort
	ftpUploadFile(ftpserver, ftpuser, pw, filename, dir, syncname)
 }


  if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
   files = append(files, filename)
  }
  return nil
 })
 return files, err
}



func (c *Config) InitConfig(path string) {
    c.Mymap = make(map[string]string)
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    r := bufio.NewReader(f)
    for {
        b, _, err := r.ReadLine()
        if err != nil {
            if err == io.EOF {
                break
            }
            panic(err)
        }
        s := strings.TrimSpace(string(b))
        if strings.Index(s, "#") == 0 {
            continue
        }
        n1 := strings.Index(s, "[")
        n2 := strings.LastIndex(s, "]")
        if n1 > -1 && n2 > -1 && n2 > n1+1 {
            c.strcet = strings.TrimSpace(s[n1+1 : n2])
	    continue
        }
        if len(c.strcet) == 0 {
            continue
        }
        index := strings.Index(s, "=")
        if index < 0 {
            continue
        }
        frist := strings.TrimSpace(s[:index])
        if len(frist) == 0 {
            continue
        }
        second := strings.TrimSpace(s[index+1:])

        pos := strings.Index(second, "\t#")
        if pos > -1 {
            second = second[0:pos]
        }
        pos = strings.Index(second, " #")
        if pos > -1 {
            second = second[0:pos]
        }
        pos = strings.Index(second, "\t//")
        if pos > -1 {
            second = second[0:pos]
        }
        pos = strings.Index(second, " //")
        if pos > -1 {
            second = second[0:pos]
        }
        if len(second) == 0 {
            continue
        }
        key := c.strcet + "=" + frist
        c.Mymap[key] = strings.TrimSpace(second)
    }
}

func (c Config) Read(node, key string) string {
    key = node + "=" + key
    v, found := c.Mymap[key]
    if !found {
        return ""
    }
    return v
}


func main(){
    ftpConfig := new(Config)
    ftpConfig.InitConfig("config.ini")
    sync_path := ftpConfig.Read("path", "file_path")
    sync_type := ftpConfig.Read("file", "file_log")
    WalkDir(sync_path, sync_type)
}
