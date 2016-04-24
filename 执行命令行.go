package main
 
import (
    "os/exec"
    "fmt"
)
 
func main() {
    // cmd := exec.Command("ls") //查看当前目录下的文件
    
    cmd := exec.Command("dig","haosou.com","+short")
    
    out, err := cmd.CombinedOutput()
    if err != nil {  
        fmt.Println(err)  
    }  
    fmt.Println(string(out))  
}