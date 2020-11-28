package execPing

// Execute ping, stop the process if the http connection is lost. Return an integer value of 1 on error.

import (
        "fmt"
        "os/exec"
        "io"
        "bufio"
        "net/http"
)



func Ping(w http.ResponseWriter) (string, int) {

        flusher, ok := w.(http.Flusher)
        if !ok {
                http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return "execPing.Ping:: Streaming unsupported!, http.StatusInternalServerError ", 1
        }

	notify := w.(http.CloseNotifier).CloseNotify()

	cmd := exec.Command("ping", "127.0.0.1")

        stdout, err := cmd.StdoutPipe()
	if err != nil {
                println(err)
		return "execPing.Ping:: error returned by cmd.StdoutPipe() ", 1
        }

	var connected bool = true
	go func() {
		<-notify
		connected = false
	}()

	if err := cmd.Start(); err != nil {
                fmt.Printf("%s\n", err)
		return "execPing.Ping:: error returned by cmd.Start()", 1
        }

        bf := bufio.NewReader(stdout)
        for {
		if connected == true {
			
			switch line, err := bf.ReadString('\n'); err {
			case nil:
				// valid line, echo it.  note that line contains trailing \n.
				fmt.Fprintf(w,"%s<br>",line)
				// inject javascript to scroll the browser to the bottom line
				// NOTE: for nested elements use:
				// window.scrollTo(0,document.querySelector(".scrollingContainer").scrollHeight);
				fmt.Fprintf(w,"<script>window.scrollTo(0,document.body.scrollHeight);</script>")
				flusher.Flush()
			case io.EOF:
				if line > "" {
					// last line of file missing \n, but still valid
					fmt.Fprintf(w,"%s\n<br>",line)
				}
				if err := cmd.Process.Kill(); err != nil {
					println("failed to kill process: ", err)
					return "execPing.Ping:: failed to kill process ", 1
				}

				//The connection may drop inside this 'if' code block before the switch statement.
			default:  
				handle_dropped_connection (cmd)
				return "execPing.Ping:: dropped connection", 1
			}

		}else {
			handle_dropped_connection (cmd)
			return "execPing.Ping:: dropped connection", 1

		}
		
	}
	return "execPing.Ping:: normal exit", 0
}


func handle_dropped_connection (cmd *exec.Cmd) {
	println("The client closed the connection prematurely. Cleaning up.")
	println("killing process ", cmd.Process.Pid)
	if err := cmd.Process.Kill(); err != nil {
		println("failed to kill process: ", err)
	}
	cmd.Wait()
	return
}

