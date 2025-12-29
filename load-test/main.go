package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Config holds the configuration for the load test
type Config struct {
	AuthServiceURL    string
	ProductServiceURL string
	CartServiceURL    string
	OrderServiceURL   string
	Concurrency       int
	Duration          time.Duration
}

// Stats holds the metrics collected during the test
type Stats struct {
	TotalRequests int64
	Success       int64
	Failures      int64
	AuthErrors    int64
	BrowseErrors  int64
	CartErrors    int64
	OrderErrors   int64
}

var stats Stats

func main() {
	users := flag.Int("users", 5, "Number of concurrent users")
	duration := flag.Duration("duration", 30*time.Second, "Test duration")
	replenish := flag.Bool("replenish", false, "Replenish stock for all products to 10000")
	flag.Parse()

	cfg := Config{
		AuthServiceURL:    "http://127.0.0.1:8084",
		ProductServiceURL: "http://127.0.0.1:8081",
		CartServiceURL:    "http://127.0.0.1:8082",
		OrderServiceURL:   "http://127.0.0.1:8083",
		Concurrency:       *users,
		Duration:          *duration,
	}

	if *replenish {
		replenishStock(cfg)
		return
	}

	// check if products exist, if not, seed them
	if err := checkAndSeed(cfg); err != nil {
		fmt.Printf("Failed to seed products: %v\n", err)
		return
	}

	fmt.Printf("Starting Load Test with %d users for %v...\n", cfg.Concurrency, cfg.Duration)
	fmt.Println("------------------------------------------------")

	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	// Timer to signal stop
	go func() {
		time.Sleep(cfg.Duration)
		close(stopCh)
	}()

	// Start workers
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			worker(workerID, cfg, stopCh)
		}(i)
	}

	wg.Wait()

	printStats(cfg)
}

func worker(id int, cfg Config, stopCh <-chan struct{}) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a unique user for this worker
	// Username must be alphanumeric
	username := fmt.Sprintf("user%d%d", id, time.Now().UnixNano())
	email := fmt.Sprintf("%s@test.com", username)
	password := "password123"

	// Register & Login once per worker (or could be per iteration if we want to test auth load more)
	token, err := authenticate(client, cfg, username, email, password)
	if err != nil {
		fmt.Printf("Worker %d failed to authenticate: %v\n", id, err)
		atomic.AddInt64(&stats.AuthErrors, 1)
		atomic.AddInt64(&stats.Failures, 1)
		return
	}

	for {
		select {
		case <-stopCh:
			return
		default:
			// Perform user flow
			err := runUserFlow(client, cfg, token)
			atomic.AddInt64(&stats.TotalRequests, 1) // Counting a "flow" as a request unit for simplicity, or we can count individual http calls
			if err != nil {
				atomic.AddInt64(&stats.Failures, 1)
			} else {
				atomic.AddInt64(&stats.Success, 1)
			}

			// Small sleep to prevent overwhelming localhost instantly without thinking
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		}
	}
}

func authenticate(client *http.Client, cfg Config, username, email, password string) (string, error) {
	// Register
	regBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
		"username": username,
	})

	// Try register
	regResp, regErr := client.Post(cfg.AuthServiceURL+"/auth/register", "application/json", bytes.NewBuffer(regBody))
	if regErr != nil {
		fmt.Printf("Register network error for %s: %v\n", email, regErr)
	} else {
		if regResp.StatusCode != http.StatusCreated && regResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(regResp.Body)
			fmt.Printf("Register failed for %s: Status=%d Body=%s\n", email, regResp.StatusCode, string(body))
		}
		regResp.Body.Close()
	}

	// Login
	loginBody, _ := json.Marshal(map[string]string{
		"login":    email,
		"password": password,
	})

	resp, err := client.Post(cfg.AuthServiceURL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	token, ok := res["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}

func runUserFlow(client *http.Client, cfg Config, token string) error {
	// 1. Get Products
	products, err := getProducts(client, cfg)
	if err != nil {
		atomic.AddInt64(&stats.BrowseErrors, 1)
		return err
	}

	// Filter for products with stock > 0
	var inStock []Product
	for _, p := range products {
		if p.Stock > 0 {
			inStock = append(inStock, p)
		}
	}

	if len(inStock) == 0 {
		return fmt.Errorf("no products with stock found")
	}

	// 2. Add Random Product to Cart
	randomProd := inStock[rand.Intn(len(inStock))]
	err = addToCart(client, cfg, token, randomProd.ID)
	if err != nil {
		atomic.AddInt64(&stats.CartErrors, 1)
		return err
	}

	// 3. Checkout
	err = checkout(client, cfg, token)
	if err != nil {
		atomic.AddInt64(&stats.OrderErrors, 1)
		return err
	}

	return nil
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

func getProducts(client *http.Client, cfg Config) ([]Product, error) {
	resp, err := client.Get(cfg.ProductServiceURL + "/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get products failed: %d", resp.StatusCode)
	}

	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}
	return products, nil
}

func addToCart(client *http.Client, cfg Config, token string, productID int) error {
	body, _ := json.Marshal(map[string]interface{}{
		"product_id": productID,
		"quantity":   1,
	})

	req, err := http.NewRequest("POST", cfg.CartServiceURL+"/cart", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("add to cart failed: %d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

func checkout(client *http.Client, cfg Config, token string) error {
	req, err := http.NewRequest("POST", cfg.OrderServiceURL+"/orders", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("checkout failed: %d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

func printStats(cfg Config) {
	fmt.Println("\n------------------------------------------------")
	fmt.Println("Load Test Results")
	fmt.Println("------------------------------------------------")
	fmt.Printf("Duration:        %v\n", cfg.Duration)
	fmt.Printf("Concurrency:     %d users\n", cfg.Concurrency)
	fmt.Printf("Total Flows:     %d\n", stats.TotalRequests)
	fmt.Printf("Successful Flows:%d\n", stats.Success)
	fmt.Printf("Failed Flows:    %d\n", stats.Failures)
	fmt.Println("--- Error Breakdown ---")
	fmt.Printf("Auth Errors:     %d\n", stats.AuthErrors)
	fmt.Printf("Browse Errors:   %d\n", stats.BrowseErrors)
	fmt.Printf("Cart Errors:     %d\n", stats.CartErrors)
	fmt.Printf("Order Errors:    %d\n", stats.OrderErrors)

	rps := float64(stats.TotalRequests) / cfg.Duration.Seconds()
	fmt.Printf("\nRequests (Flows) per Second: %.2f\n", rps)
	fmt.Println("------------------------------------------------")
}

func replenishStock(cfg Config) {
	fmt.Println("Replenishing stock for all products...")
	client := &http.Client{Timeout: 30 * time.Second}
	products, err := getProducts(client, cfg)
	if err != nil {
		fmt.Printf("Failed to get products: %v\n", err)
		return
	}

	for _, p := range products {
		p.Stock = 99999999
		body, _ := json.Marshal(p)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/products/%d", cfg.ProductServiceURL, p.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to update product %d: %v\n", p.ID, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to update product %d: Status %d\n", p.ID, resp.StatusCode)
		}
		resp.Body.Close()
	}
	fmt.Println("Stock replenishment execution complete.")
}

func checkAndSeed(cfg Config) error {
	client := &http.Client{Timeout: 30 * time.Second}
	products, err := getProducts(client, cfg)
	if err != nil {
		return err
	}

	if len(products) > 0 {
		return nil
	}

	fmt.Println("No products found. Seeding 20 test products...")
	for i := 1; i <= 20; i++ {
		p := map[string]interface{}{
			"name":  fmt.Sprintf("Product %d", i),
			"price": rand.Intn(100) + 10,
			"stock": 10000,
		}
		// Fix: API expects an array of products
		body, _ := json.Marshal([]map[string]interface{}{p})
		resp, err := client.Post(cfg.ProductServiceURL+"/products", "application/json", bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("failed to create product %d: %v", i, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to create product %d: status=%d body=%s", i, resp.StatusCode, string(b))
		}
	}
	fmt.Println("Seeding complete.")
	return nil
}
