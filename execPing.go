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
                exit_err()
        }

	notify := w.(http.CloseNotifier).CloseNotify()

	cmd := exec.Command("ping", "127.0.0.1")

        stdout, err := cmd.StdoutPipe()
	if err != nil {
                println(err)
		exit_err()
        }

	var connected bool = true
	go func() {
		<-notify
		connected = false
	}()

	if err := cmd.Start(); err != nil {
                fmt.Printf("%s\n", err)
		exit_err()
        }

        bf := bufio.NewReader(stdout)
        for {
		if connected == true {
			
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
					exit_err()
				}

				//The connection may have dropped unknown to us while we are inside this 'if' code block.
			default:  
				handle_dropped_connection (cmd)
			}

		}else {
			handle_dropped_connection (cmd)
			}
			
		}
	exit_success()
	return "exec.Ping", 0
}
	

	func handle_dropped_connection (cmd *exec.Cmd) {
		println("The client closed the connection prematurely. Cleaning up.")
		println("killing process ", cmd.Process.Pid)
		if err := cmd.Process.Kill(); err != nil {
			println("failed to kill process: ", err)
		}
		cmd.Wait()
		exit_err()
	}

	func exit_err() (string, int){
		return "execPing.Ping", 1
	}

	func exit_success()  (string, int){
		return "execPing.Ping", 0
	}
