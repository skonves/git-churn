package main

import "os/exec"

func main() {
    cmd := exec.Command("git", "log", "--pretty=oneline")
    stdout, err := cmd.Output()

    if err != nil {
        println(err.Error())
        return
    }

    print(string(stdout))
}
