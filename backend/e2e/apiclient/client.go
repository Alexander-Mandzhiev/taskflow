// Package apiclient — HTTP-клиент для e2e/smoke-проверок API.
//
// Закрытие тела ответа: методы, возвращающие *http.Response (Register, Login, CreateTeam, CreateTask, GetTeam, GetTask и т.д.),
// требуют от вызывающего кода закрыть resp.Body. Методы, возвращающие распарсенные данные (ListTeams, ListTasks и т.п.),
// закрывают тело ответа внутри себя.
package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	apiPath        = "/api/v1"
)

// Client — HTTP-клиент для e2e/smoke-проверок API. После Login() использует cookie для авторизации.
type Client struct {
	baseURL string
	client  *http.Client
}

// New создаёт клиент. baseURL — без слэша в конце (например http://localhost:4000).
func New(baseURL string) (*Client, error) {
	return NewWithRequestDelay(baseURL, 0)
}

// delayTransport выполняет паузу перед каждым запросом (для снижения нагрузки на rate limit в тестах).
type delayTransport struct {
	roundTripper http.RoundTripper
	delay        time.Duration
}

func (t *delayTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.delay > 0 {
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(t.delay):
		}
	}
	return t.roundTripper.RoundTrip(req)
}

// NewWithRequestDelay создаёт клиент с паузой перед каждым запросом.
// delay > 0 используют в e2e-тестах, чтобы не упираться в rate limit (регистрация, логин и т.д.).
func NewWithRequestDelay(baseURL string, delay time.Duration) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar: %w", err)
	}
	base := strings.TrimSuffix(baseURL, "/")
	transport := http.DefaultTransport
	if delay > 0 {
		transport = &delayTransport{roundTripper: http.DefaultTransport, delay: delay}
	}
	return &Client{
		baseURL: base,
		client: &http.Client{
			Timeout:   defaultTimeout,
			Jar:       jar,
			Transport: transport,
		},
	}, nil
}

// NewWithClient создаёт клиент с заданным http.Client (например с другим timeout или jar).
func NewWithClient(baseURL string, c *http.Client) *Client {
	base := strings.TrimSuffix(baseURL, "/")
	if c == nil {
		c = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{baseURL: base, client: c}
}

func (c *Client) url(path string) string {
	if strings.HasPrefix(path, "/") {
		return c.baseURL + path
	}
	return c.baseURL + "/" + path
}

func (c *Client) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var reqBody *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}
	var req *http.Request
	var err error
	if reqBody != nil {
		req, err = http.NewRequestWithContext(ctx, method, c.url(path), reqBody)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, c.url(path), nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req) //nolint:gosec // G704: e2e client, URL from config/test
}

func (c *Client) doGet(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	u := c.url(path)
	if len(query) > 0 {
		u = u + "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return c.client.Do(req) //nolint:gosec // G704: e2e client, URL from config/test
}

// Register — POST /api/v1/register.
func (c *Client) Register(ctx context.Context, req RegisterRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, apiPath+"/register", req)
}

// Login — POST /api/v1/login. При успехе токены сохраняются в cookie jar.
func (c *Client) Login(ctx context.Context, req LoginRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, apiPath+"/login", req)
}

// Logout — POST /api/v1/logout.
func (c *Client) Logout(ctx context.Context) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, apiPath+"/logout", nil)
}

// CreateTeam — POST /api/v1/teams. Возвращает ответ; ID команды в теле.
func (c *Client) CreateTeam(ctx context.Context, req CreateTeamRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, apiPath+"/teams", req)
}

// ListTeams — GET /api/v1/teams.
func (c *Client) ListTeams(ctx context.Context) ([]TeamWithRole, error) {
	resp, err := c.doGet(ctx, apiPath+"/teams", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list teams: %s", resp.Status)
	}
	var out []TeamWithRole
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTeam — GET /api/v1/teams/{id}.
func (c *Client) GetTeam(ctx context.Context, id string) (*http.Response, error) {
	return c.doGet(ctx, apiPath+"/teams/"+url.PathEscape(id), nil)
}

// Invite — POST /api/v1/teams/{id}/invite.
func (c *Client) Invite(ctx context.Context, teamID string, req InviteRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, apiPath+"/teams/"+url.PathEscape(teamID)+"/invite", req)
}

// CreateTask — POST /api/v1/tasks.
func (c *Client) CreateTask(ctx context.Context, req CreateTaskRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, apiPath+"/tasks", req)
}

// ListTasks — GET /api/v1/tasks?team_id=...&status=...&assignee_id=...&limit=...&offset=...
func (c *Client) ListTasks(ctx context.Context, opts ListTasksOpts) (*TaskListResponse, error) {
	q := url.Values{}
	q.Set("team_id", opts.TeamID)
	if opts.Status != "" {
		q.Set("status", opts.Status)
	}
	if opts.AssigneeID != "" {
		q.Set("assignee_id", opts.AssigneeID)
	}
	if opts.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", opts.Offset))
	}
	resp, err := c.doGet(ctx, apiPath+"/tasks", q)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list tasks: %s", resp.Status)
	}
	var out TaskListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTask — GET /api/v1/tasks/{id}.
func (c *Client) GetTask(ctx context.Context, id string) (*http.Response, error) {
	return c.doGet(ctx, apiPath+"/tasks/"+url.PathEscape(id), nil)
}

// UpdateTask — PUT /api/v1/tasks/{id}.
func (c *Client) UpdateTask(ctx context.Context, taskID string, req UpdateTaskRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, apiPath+"/tasks/"+url.PathEscape(taskID), req)
}

// GetTaskHistory — GET /api/v1/tasks/{id}/history.
func (c *Client) GetTaskHistory(ctx context.Context, taskID string) (*http.Response, error) {
	return c.doGet(ctx, apiPath+"/tasks/"+url.PathEscape(taskID)+"/history", nil)
}

// ListComments — GET /api/v1/tasks/{id}/comments.
func (c *Client) ListComments(ctx context.Context, taskID string) (*http.Response, error) {
	return c.doGet(ctx, apiPath+"/tasks/"+url.PathEscape(taskID)+"/comments", nil)
}

// CreateComment — POST /api/v1/tasks/{id}/comments.
func (c *Client) CreateComment(ctx context.Context, taskID string, req CreateCommentRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, apiPath+"/tasks/"+url.PathEscape(taskID)+"/comments", req)
}

// ReportTeamStats — GET /api/v1/reports/team-stats.
func (c *Client) ReportTeamStats(ctx context.Context) (*http.Response, error) {
	return c.doGet(ctx, apiPath+"/reports/team-stats", nil)
}

// ReportTopCreators — GET /api/v1/reports/top-creators.
func (c *Client) ReportTopCreators(ctx context.Context) (*http.Response, error) {
	return c.doGet(ctx, apiPath+"/reports/top-creators", nil)
}

// ReportInvalidAssignees — GET /api/v1/reports/invalid-assignees.
func (c *Client) ReportInvalidAssignees(ctx context.Context) (*http.Response, error) {
	return c.doGet(ctx, apiPath+"/reports/invalid-assignees", nil)
}
