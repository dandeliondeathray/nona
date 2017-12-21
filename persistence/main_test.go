package persistence_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

func runEtcd() *exec.Cmd {
	log.Println("Starting etcd Docker container.")
	cmd := exec.Command("docker-compose", "-f", "testing/docker-compose.yaml", "up")
	cmd.Start()
	return cmd
}

func stopEtcd() {
	log.Println("Stopping the etcd Docker container.")
	cmd := exec.Command("docker-compose", "-f", "testing/docker-compose.yaml", "down")
	err := cmd.Run()
	if err != nil {
		log.Println("Docker Compose failed when stopping the etcd container:", err)
		out, _ := cmd.CombinedOutput()
		log.Println("Output:", string(out))
	}
}

func runWithEtcdContainer(m *testing.M) {
	cmd := runEtcd()
	time.Sleep(2 * time.Second)
	exitCode := m.Run()
	stopEtcd()

	log.Println("Waiting for etcd container to stop...")
	err := cmd.Wait()
	if err != nil {
		log.Printf("Docker Compose had an error while starting etcd: %v", err)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error when getting Docker Compose output: %v", err)
		} else {
			fmt.Println(string(out))
		}
	}

	os.Exit(exitCode)
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		runWithEtcdContainer(m)
		os.Exit(0)
	} else {
		os.Exit(m.Run())
	}
}

func TestRound_NewRoundSet_SameSeedReturnedInGet(t *testing.T) {
}
