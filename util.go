package rfe

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"
	"syscall"

	"github.com/cskr/pubsub"
	"github.com/phayes/freeport"
)

// Network
func getFreeHostPort() uint16 {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	return uint16(port)
}

// Process
func killProcessGroup(cmd *exec.Cmd) {
	// More infos: https://man7.org/linux/man-pages/man3/kill.3p.html#DESCRIPTION
	thisProcessAndHisChildren := -cmd.Process.Pid
	syscall.Kill(thisProcessAndHisChildren, syscall.SIGKILL)
}

func makeProcessKillable(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// Stream
func getBothStdoutStderrCombined(cmd *exec.Cmd) io.ReadCloser {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stderr = cmd.Stdout

	return stdout
}

func getStreamReadlinesIterator(stream io.ReadCloser) (<-chan string, error) {
	scanner := bufio.NewScanner(stream)
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	chnl := make(chan string)
	go func() {
		for scanner.Scan() {
			if line := scanner.Text(); stringIsNotEmpty(line) {
				chnl <- line
			}
		}
		close(chnl)
	}()

	return chnl, nil
}

func logPubSubTopic(ps *pubsub.PubSub, topicName string) {
	channel := ps.Sub(topicName)

	for {
		if msg, ok := <-channel; ok {
			log.Printf("%s", msg)
		} else {
			break
		}
	}
}

// Other
func stringIsNotEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
