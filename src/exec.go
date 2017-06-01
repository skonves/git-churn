package main

import "os/exec"

func main() {
    cmd := exec.Command("git", "log")
    stdout, err := cmd.Output()

    if err != nil {
		print("fail")
        println(err.Error())
        return
    }

    print(string(stdout))
}
