package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	// User's service name
	UserServiceName string

	// User's service's port to serve up health information
	UserServicePort string
)

// Initializes and parses flag parameters before continuing through rest of script
func init() {
	flag.StringVar(&UserServiceName, "svc", "", "A string of the service name you'd like to perform a health check on")
	flag.StringVar(&UserServicePort, "p", "8080", "A string of the port number you'd like to serve the health check over")
	flag.Parse()
}

func handler(w http.ResponseWriter, r *http.Request) {
	if UserServiceName == "" {
		os.Exit(1)
		return
	} else {
		cmd := exec.Command("systemctl", "check", UserServiceName)
		out, err := cmd.CombinedOutput()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Printf("systemctl finished with non-zero: %v\n", exitErr)
			} else {
				fmt.Printf("failed to run systemctl: %v", err)
				os.Exit(1)
			}
		}
		log.Printf("Status is: %s\n", string(out))
		isItActive := strings.Contains(string(out), "active")
		if isItActive {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = "Service is running!"
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		} else {
			w.WriteHeader(http.StatusBadGateway)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = "Service isn't running."
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
	}
}

func main() {
	UserServicePort = ":" + UserServicePort
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(UserServicePort, nil))
}
