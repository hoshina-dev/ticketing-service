// Command copium-openapi fetches Copium's Swagger 2 doc and writes OpenAPI 3 JSON
// for oapi-codegen (which requires OpenAPI 3 parameter schemas).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

func main() {
	url := flag.String("url", "http://localhost:8081/swagger/doc.json", "Copium swagger doc URL")
	out := flag.String("o", "internal/clients/copium/openapi.json", "output path for OpenAPI 3 JSON")
	flag.Parse()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(*url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "copium-openapi: GET %s: %v\n", *url, err)
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "copium-openapi: GET %s -> %d: %s\n", *url, resp.StatusCode, body)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "copium-openapi: read body: %v\n", err)
		os.Exit(1)
	}

	var doc2 openapi2.T
	if err := json.Unmarshal(body, &doc2); err != nil {
		fmt.Fprintf(os.Stderr, "copium-openapi: decode swagger 2: %v\n", err)
		os.Exit(1)
	}

	doc3, err := openapi2conv.ToV3(&doc2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "copium-openapi: convert to OpenAPI 3: %v\n", err)
		os.Exit(1)
	}

	encoded, err := json.MarshalIndent(doc3, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "copium-openapi: encode OpenAPI 3: %v\n", err)
		os.Exit(1)
	}
	encoded = append(encoded, '\n')

	if err := os.WriteFile(*out, encoded, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "copium-openapi: write %s: %v\n", *out, err)
		os.Exit(1)
	}

	fmt.Printf("copium-openapi: wrote %s from %s\n", *out, *url)
}
