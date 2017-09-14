package main
import (
        "log"
        "os"
        "os/exec"
        "time"
	"bytes"
)

func main() {
	lf, err := os.OpenFile("/var/log/sync_ftp.log", os.O_CREATE | os.O_RDWR | os.O_APPEND, 0600)
	  if err != nil {
                os.Exit(1)
        }
	defer lf.Close()
	l := log.New(lf, "", os.O_APPEND)
	 for { 
	        in := bytes.NewBuffer(nil)
    		cmd := exec.Command("sh")
    		cmd.Stdin = in
    		go func() {
        	in.WriteString("go run /home/ftp/hiddenfile_rename.go \n")
      		in.WriteString("go run /home/ftp/sync_ftp.go \n")  
		in.WriteString("exit\n")
    		}()
		err := cmd.Start()

		if err != nil {
                        l.Printf("start sync_ftp fail", time.Now().Format("2006-01-02 15:04:05"), err)
                        time.Sleep(time.Second * 5)
                        continue
                }

		l.Printf("%s sync_ftp start sucess!!", time.Now().Format("2006-01-02 15:04:05"), err)
		err = cmd.Wait()
		l.Printf("%s sync_ftp exit", time.Now().Format("2006-01-02 15:04:05"), err)
                time.Sleep(time.Second * 3600)
	
        }

}
