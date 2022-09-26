package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	var sha string
	flag.StringVar(&sha, "sha", "unknown", "SHA sum for start of comparison")
	flag.Parse()

	cmdString := fmt.Sprintf("git rev-list --merges %s..HEAD", sha)
	log.Default().Print(cmdString)

	out, err := exec.Command("git", "rev-list", "--merges", fmt.Sprintf("%s..HEAD", sha)).Output()

	if err != nil {
		log.Default().Print(out)
		panic(err)
	}

	log.Default().Print(string(out))
}
