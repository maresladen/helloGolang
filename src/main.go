package main
 
import (
    "os/exec"
    "fmt"
)
 
func run() {
    cmd := exec.Command("/bin/sh", "-c 3", "ping 127.0.0.1")
    bStr, err := cmd.Output()
    if err != nil {
        panic(err.Error())
    }
 
    if err := cmd.Start(); err != nil {
        panic(err.Error())
    }
 
    if err := cmd.Wait(); err != nil {
        panic(err.Error())
    }
    
    fmt.Println(string(bStr))
}
 
func main() {
    run()
}