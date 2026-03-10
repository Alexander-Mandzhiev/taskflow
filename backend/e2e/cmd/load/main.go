// Нагрузочный тест API через e2e apiclient.
//
// Один раз регистрирует пользователя, логинится, создаёт команду; затем
// заданное время гоняет параллельные запросы: ListTeams, GetTeam, ListTasks,
// CreateTask, GetTask и т.д. В конце выводит число запросов и ошибок.
//
// Запуск:
//
//	go run ./e2e/cmd/load -base http://localhost:4000
//	go run ./e2e/cmd/load -base http://localhost:4000 -duration 1m -concurrency 10
//
// Увеличение нагрузки: больше воркеров и/или дольше время.
//
//	-concurrency 20   — 20 параллельных воркеров (по умолчанию 5)
//	-concurrency 50   — высокая нагрузка (на Windows следи за сокетами)
//	-duration 2m      — гнать 2 минуты
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/e2e/apiclient"
)

const defaultBaseURL = "http://localhost:4000"

func main() {
	baseURL := flag.String("base", "", "API base URL (default: API_BASE_URL or http://localhost:4000)")
	duration := flag.Duration("duration", 30*time.Second, "how long to run the load")
	concurrency := flag.Int("concurrency", 5, "number of concurrent workers")
	flag.Parse()

	if *baseURL == "" {
		*baseURL = os.Getenv("API_BASE_URL")
	}
	if *baseURL == "" {
		*baseURL = defaultBaseURL
	}

	ctx := context.Background()
	// Пул соединений: ограничиваем и переиспользуем, чтобы на Windows не исчерпать эпиhemeral ports.
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("cookiejar: %v", err)
	}
	// Пул под выбранную concurrency; при большом concurrency лимит сдерживает исчерпание портов на Windows
	maxConns := *concurrency * 2
	if maxConns < 10 {
		maxConns = 10
	}
	transport := &http.Transport{
		MaxIdleConnsPerHost: maxConns,
		MaxConnsPerHost:     maxConns,
	}
	client := apiclient.NewWithClient(*baseURL, &http.Client{
		Timeout:   30 * time.Second,
		Jar:       jar,
		Transport: transport,
	})

	// Один пользователь и одна команда на весь прогон
	reg := apiclient.FakeRegister()
	log.Printf("Register: %s", reg.Email)
	resp, err := client.Register(ctx, reg)
	if err != nil {
		log.Fatalf("Register: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatalf("Register: %s", resp.Status)
	}

	resp, err = client.Login(ctx, apiclient.FakeLogin(reg.Email, reg.Password))
	if err != nil {
		log.Fatalf("Login: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Login: %s", resp.Status)
	}
	log.Println("Login OK (session set)")

	teamReq := apiclient.FakeCreateTeam()
	resp, err = client.CreateTeam(ctx, teamReq)
	if err != nil {
		log.Fatalf("CreateTeam: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		_ = resp.Body.Close()
		log.Fatalf("CreateTeam: %s", resp.Status)
	}
	var team apiclient.Team
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		_ = resp.Body.Close()
		log.Fatalf("CreateTeam decode: %v", err)
	}
	_ = resp.Body.Close()
	log.Printf("Team created: %s", team.ID)

	deadline := time.Now().Add(*duration)
	var total, errors int64
	latencies := newLatencyCollector()
	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(_ int) {
			defer wg.Done()
			runWorker(ctx, client, team.ID, deadline, &total, &errors, latencies)
		}(i)
	}
	wg.Wait()

	totalN := atomic.LoadInt64(&total)
	errN := atomic.LoadInt64(&errors)
	runDuration := *duration
	log.Printf("Done: %d requests, %d errors", totalN, errN)
	if totalN > 0 {
		rps := float64(totalN) / runDuration.Seconds()
		log.Printf("RPS: %.1f", rps)
	}
	latencies.report(runDuration)
}

// latencyCollector собирает длительности успешных запросов для перцентилей.
type latencyCollector struct {
	mu sync.Mutex
	d  []time.Duration
}

func newLatencyCollector() *latencyCollector {
	return &latencyCollector{}
}

func (c *latencyCollector) add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.d = append(c.d, d)
}

func (c *latencyCollector) report(runDuration time.Duration) {
	c.mu.Lock()
	d := c.d
	c.mu.Unlock()
	if len(d) == 0 {
		return
	}
	sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
	p50 := d[len(d)*50/100]
	p95 := d[len(d)*95/100]
	p99 := d[len(d)*99/100]
	rps := float64(len(d)) / runDuration.Seconds()
	log.Printf("Latency (ok requests): count=%d, p50=%v, p95=%v, p99=%v, RPS=%.1f", len(d), p50.Round(time.Millisecond), p95.Round(time.Millisecond), p99.Round(time.Millisecond), rps)
}

func runWorker(ctx context.Context, client *apiclient.Client, teamID string, deadline time.Time, total, errCount *int64, lat *latencyCollector) {
	for time.Now().Before(deadline) {
		// ListTeams
		start := time.Now()
		_, err := client.ListTeams(ctx)
		atomic.AddInt64(total, 1)
		if err != nil {
			atomic.AddInt64(errCount, 1)
			log.Printf("ListTeams: %v", err)
			continue
		}
		lat.add(time.Since(start))

		// GetTeam
		start = time.Now()
		resp, err := client.GetTeam(ctx, teamID)
		atomic.AddInt64(total, 1)
		if err != nil {
			atomic.AddInt64(errCount, 1)
			log.Printf("GetTeam: %v", err)
			continue
		}
		_ = resp.Body.Close()
		lat.add(time.Since(start))

		// ListTasks
		start = time.Now()
		_, err = client.ListTasks(ctx, apiclient.ListTasksOpts{TeamID: teamID, Limit: 20})
		atomic.AddInt64(total, 1)
		if err != nil {
			atomic.AddInt64(errCount, 1)
			log.Printf("ListTasks: %v", err)
			continue
		}
		lat.add(time.Since(start))

		// CreateTask
		start = time.Now()
		taskReq := apiclient.FakeCreateTask(teamID)
		resp, err = client.CreateTask(ctx, taskReq)
		atomic.AddInt64(total, 1)
		if err != nil {
			atomic.AddInt64(errCount, 1)
			log.Printf("CreateTask: %v", err)
			continue
		}
		if resp.StatusCode != http.StatusCreated {
			_ = resp.Body.Close()
			atomic.AddInt64(errCount, 1)
			continue
		}
		var task apiclient.Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			_ = resp.Body.Close()
			atomic.AddInt64(errCount, 1)
			continue
		}
		_ = resp.Body.Close()
		lat.add(time.Since(start))

		// GetTask
		start = time.Now()
		resp, err = client.GetTask(ctx, task.ID)
		atomic.AddInt64(total, 1)
		if err != nil {
			atomic.AddInt64(errCount, 1)
			log.Printf("GetTask: %v", err)
			continue
		}
		_ = resp.Body.Close()
		lat.add(time.Since(start))

		// ReportTeamStats (тяжёлый read)
		start = time.Now()
		resp, err = client.ReportTeamStats(ctx)
		atomic.AddInt64(total, 1)
		if err != nil {
			atomic.AddInt64(errCount, 1)
		} else {
			_ = resp.Body.Close()
			lat.add(time.Since(start))
		}
	}
}
