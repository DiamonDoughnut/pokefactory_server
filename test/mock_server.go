package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", logRequest)
	
	fmt.Println("Mock server listening on :8080")
	fmt.Println("This will log all incoming requests from your Minecraft mod")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n=== REQUEST %s ===\n", time.Now().Format("15:04:05"))
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.String())
	fmt.Printf("Headers:\n")
	for name, values := range r.Header {
		fmt.Printf("  %s: %s\n", name, strings.Join(values, ", "))
	}
	
	body, _ := io.ReadAll(r.Body)
	if len(body) > 0 {
		fmt.Printf("Body: %s\n", string(body))
		
		// Try to pretty print JSON
		var jsonData interface{}
		if json.Unmarshal(body, &jsonData) == nil {
			prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
			fmt.Printf("Formatted JSON:\n%s\n", string(prettyJSON))
		}
	}
	
	// Send appropriate response based on endpoint
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if strings.Contains(r.URL.Path, "/server/auth") {
		w.Write([]byte(`{"token":"fake.jwt.token","expires_in":86400}`))
	} else {
		w.Write([]byte(`{"status":"ok","message":"Request logged"}`))
	}
}