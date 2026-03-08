// Программа для smoke/e2e-проверки API: регистрация, логин, создание команды и задачи.
// Запуск: API_BASE_URL=http://localhost:4000 go run ./e2e/cmd
// Бэкенд должен быть поднят (локально или в тест-контейнере).
//
//nolint:gocritic // exitAfterDefer: один контекст на весь прогон, при log.Fatalf процесс завершается, defer cancel не критичен
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/e2e/apiclient"
)

const defaultBaseURL = "http://localhost:4000"

func main() {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := apiclient.New(baseURL)
	if err != nil {
		log.Fatalf("apiclient.New: %v", err)
	}

	// Генерируем данные через gofakeit (каждый прогон — новые данные).
	reg := apiclient.FakeRegister()
	log.Printf("Register: email=%s", reg.Email)

	resp, err := client.Register(ctx, reg)
	if err != nil {
		log.Fatalf("Register: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatalf("Register: %s", resp.Status)
	}
	log.Println("Register OK")

	resp, err = client.Login(ctx, apiclient.FakeLogin(reg.Email, reg.Password))
	if err != nil {
		log.Fatalf("Login: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Login: %s", resp.Status)
	}
	log.Println("Login OK (cookies set)")

	teamReq := apiclient.FakeCreateTeam()
	log.Printf("CreateTeam: name=%s", teamReq.Name)
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
	log.Printf("CreateTeam OK: id=%s", team.ID)

	teams, err := client.ListTeams(ctx)
	if err != nil {
		log.Fatalf("ListTeams: %v", err)
	}
	log.Printf("ListTeams OK: count=%d", len(teams))

	taskReq := apiclient.FakeCreateTask(team.ID)
	log.Printf("CreateTask: title=%s", taskReq.Title)
	resp, err = client.CreateTask(ctx, taskReq)
	if err != nil {
		log.Fatalf("CreateTask: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		_ = resp.Body.Close()
		log.Fatalf("CreateTask: %s", resp.Status)
	}
	var task apiclient.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		_ = resp.Body.Close()
		log.Fatalf("CreateTask decode: %v", err)
	}
	_ = resp.Body.Close()
	log.Printf("CreateTask OK: id=%s", task.ID)

	list, err := client.ListTasks(ctx, apiclient.ListTasksOpts{TeamID: team.ID})
	if err != nil {
		log.Fatalf("ListTasks: %v", err)
	}
	log.Printf("ListTasks OK: total=%d", list.Total)

	// Проверка: созданная задача действительно есть в списке (по ID).
	var found bool
	for _, t := range list.Items {
		if t.ID == task.ID {
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("ListTasks: created task id=%s not found in list", task.ID)
	}
	log.Printf("ListTasks: created task id=%s found in list", task.ID)

	log.Println("Smoke OK: register → login → team → task → list succeeded.")
}
