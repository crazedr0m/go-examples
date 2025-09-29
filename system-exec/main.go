package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	ctx := context.Background()

	log.Println("ctx: ", ctx)

	cmd := exec.Command("php", "-i")
//	cmd := exec.Command("which", "php")
//	cmd := exec.Command("ls", "-lah")
//	cmd := exec.Command("pwd")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
	}
	fmt.Printf("Output: %s\n", out.String())	
}
