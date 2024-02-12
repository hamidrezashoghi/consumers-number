package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	// Kafka consumers
	servers := []string{"192.168.1.2:9092", "192.168.1.3:9092", "192.168.1.4:9092"}

	consumersNumber := make(map[string]int, 3)
	groupName := "test-group"
	timeout := "10000"
	kafkaPath := "/usr/local/kafka/bin/"
	consumerServersFile := "/var/lib/node_exporter/textfile_collector/consumer_servers.prom"

	// Get number of consumers per kafka
	for _, server := range servers {
		cmdConsumer := exec.Command(kafkaPath+"kafka-consumer-groups.sh", "--dry-run",
			"--bootstrap-server", server, "--timeout", timeout, "--group="+groupName,
			"--members", "--describe")

		cmdWc := exec.Command("wc", "-l")

		// Create a pipe between two commands
		r, w := io.Pipe()
		cmdConsumer.Stdout = w
		cmdWc.Stdin = r

		go func() {
			defer w.Close()
			if err := cmdConsumer.Run(); err != nil {
				log.Fatalf("kafka-consumer-groups.sh: couldn't get consumers on %s, %v\n", server, err)
			}
		}()

		// Number of consumers
		out, err := cmdWc.Output()
		if err != nil {
			log.Fatalf("wc command: Couldn't get number of consumers on %s, %v\n", server, err)
		}

		cnt, _ := strconv.Atoi(strings.TrimSpace(string(out)))
		cnt -= 2
		consumersNumber[server] = cnt
	}

	// Change server IP to region
	consNumbers := make(map[string]int, 3)
	for server, c := range consumersNumber {
		switch {
		case server == "192.168.1.2:9092":
			consNumbers["local1"] = c
		case server == "192.168.1.3:9092":
			consNumbers["local2"] = c
		case server == "192.168.1.4:9092":
			consNumbers["local3"] = c
		}
	}

	file, err := os.OpenFile("consumer_servers.prom", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Couldn't open consumer_servers.prom file.")
	}

	defer file.Close()

	_ = os.Remove(consumerServersFile)
	_, _ = file.WriteString("# HELP consumer_servers Metric\n")
	_, _ = file.WriteString("# TYPE consumer_servers gauge\n")
	for region, consumers := range consNumbers {
		_, _ = file.WriteString("consumer_servers{region=\"" + region + "\"} " + strconv.Itoa(consumers) + "\n")
	}

	// Change ownership of consumer_servers.prom
	// execute below command to get uid and gid of node_exporter
	// # id node_exporter
	// uid=997(node_exporter) gid=998(node_exporter) groups=998(node_exporter)
	_ = os.Chown("consumer_servers.prom", 997, 998)

	// /var/lib/node_exporter/textfile_collector/
	if err := os.Rename("consumer_servers.prom", consumerServersFile); err != nil {
		log.Fatalln("Couldn't move consumer_servers.prom to /var/lib/node_exporter/textfile_collector/ path.", err)
	}
}
