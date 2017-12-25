package persistence_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/dandeliondeathray/nona/persistence"
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
	if testing.Short() {
		os.Exit(m.Run())
	} else {
		runWithEtcdContainer(m)
	}
}

type mockAsyncRecoveryHandler struct {
	expectedOnRoundRecovered  *int64
	unmatchedOnRoundRecovered []int64
	chCallReceived            chan int64
	timeout                   time.Duration
}

func newMockAsyncRecoveryHandler(timeout time.Duration) *mockAsyncRecoveryHandler {
	return &mockAsyncRecoveryHandler{nil, make([]int64, 0), make(chan int64, 5), timeout}
}

func (r *mockAsyncRecoveryHandler) OnRoundRecovered(seed int64) {
	r.chCallReceived <- seed
}

func (r *mockAsyncRecoveryHandler) ExpectOnRoundRecovered(seed int64) {
	r.expectedOnRoundRecovered = &seed
}

func (r *mockAsyncRecoveryHandler) Await() error {
	select {
	case seed := <-r.chCallReceived:
		if r.expectedOnRoundRecovered == nil {
			r.unmatchedOnRoundRecovered = append(r.unmatchedOnRoundRecovered, seed)
		} else if seed == *r.expectedOnRoundRecovered {
			r.expectedOnRoundRecovered = nil
		} else {
			r.unmatchedOnRoundRecovered = append(r.unmatchedOnRoundRecovered, seed)
		}
	case <-time.After(r.timeout):
		break
	}

	if r.expectedOnRoundRecovered != nil {
		return fmt.Errorf("No match for expected seed %d. Unexpected calls: %v", *r.expectedOnRoundRecovered, r.unmatchedOnRoundRecovered)
	} else if len(r.unmatchedOnRoundRecovered) > 0 {
		return fmt.Errorf("Unmatched calls to OnRoundRecovery: %v", r.unmatchedOnRoundRecovered)
	}
	return nil
}

func TestRecoverRound_RoundSet_RoundIsRecovered(t *testing.T) {
	if testing.Short() {
		return
	}

	recoveryHandler := newMockAsyncRecoveryHandler(time.Duration(1) * time.Second)
	seed := int64(42)
	recoveryHandler.ExpectOnRoundRecovered(seed)

	// Arrange the persistence to have a current round with seed 42.
	p := persistence.NewPersistence("konsulatet", testingEndpoints)
	p.StoreNewRound(42)

	p2 := persistence.NewPersistence("konsulatet", testingEndpoints)
	p2.Recover(recoveryHandler)
	err := recoveryHandler.Await()
	if err != nil {
		t.Fatalf("%v", err)
	}
}
