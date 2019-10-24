package execPing

// Execute ping, stop the process if the http connection is lost.

import (
        "fmt"
        "os/exec"
        "io"
        "bufio"
        "net/http"
)


func Ping(w http.ResponseWriter) {

        flusher, ok := w.(http.Flusher)
        if !ok {
                http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
                return
        }

	notify := w.(http.CloseNotifier).CloseNotify()

        cmd := exec.Command("ping", "127.0.0.1")

        stdout, err := cmd.StdoutPipe()
	if err != nil {
                println(err)
		return
        }

	connected := 1
	go func() {
		<-notify
		connected = 0
	}()

	if err := cmd.Start(); err != nil {
                fmt.Printf("%s\n", err)
		return
        }

        bf := bufio.NewReader(stdout)
        for {
		if connected == 1 {
			
			switch line, err := bf.ReadString('\n'); err {
			case nil:
				// valid line, echo it.  note that line contains trailing \n.
				fmt.Fprintf(w,"%s<br>",line)
				flusher.Flush()
			case io.EOF:
				if line > "" {
					// last line of file missing \n, but still valid
					fmt.Fprintf(w,"%s\n<br>",line)
				}
				if err := cmd.Process.Kill(); err != nil {
					println("failed to kill process: ", err)
					fmt.Fprintf(w,"case io.EOF: %d\n<br>", cmd.Process.Pid)
					flusher.Flush()
				}
				return
			default:
				println(err)
			}
		} else {

			println("The client closed the connection prematurely. Cleaning up.")
			println("killing process ", cmd.Process.Pid)
			if err := cmd.Process.Kill(); err != nil {
				println("failed to kill process: ", err)
			}
			cmd.Wait()
			return
		}
		
	}
}

