//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Alexander-Mandzhiev/taskflow/backend/e2e/apiclient"
)

// newAPIClient создаёт клиент с паузой перед каждым запросом (снижает вероятность 429 в тестах).
func newAPIClient() (*apiclient.Client, error) {
	return apiclient.NewWithRequestDelay(env.BackendURL, requestDelay)
}

// newAuthenticatedClient регистрирует пользователя, логинится и возвращает клиент с cookie.
func newAuthenticatedClient(ctx context.Context) *apiclient.Client {
	client, err := newAPIClient()
	Expect(err).ToNot(HaveOccurred())
	reg := apiclient.FakeRegister()
	resp, err := client.Register(ctx, reg)
	Expect(err).ToNot(HaveOccurred())
	defer func() { _ = resp.Body.Close() }()
	Expect(resp.StatusCode).To(Equal(http.StatusCreated), "register: %s", resp.Status)
	resp, err = client.Login(ctx, apiclient.FakeLogin(reg.Email, reg.Password))
	Expect(err).ToNot(HaveOccurred())
	_ = resp.Body.Close()
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	return client
}

// newAuthenticatedClientWithTeam возвращает клиент и созданную команду (ID).
func newAuthenticatedClientWithTeam(ctx context.Context) (*apiclient.Client, string) {
	client := newAuthenticatedClient(ctx)
	teamReq := apiclient.FakeCreateTeam()
	resp, err := client.CreateTeam(ctx, teamReq)
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))
	var team apiclient.Team
	Expect(json.NewDecoder(resp.Body).Decode(&team)).ToNot(HaveOccurred())
	_ = resp.Body.Close()
	Expect(team.ID).ToNot(BeEmpty())
	return client, team.ID
}

var _ = Describe("Auth", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(suiteCtx, requestTimeout)
	})

	AfterEach(func() {
		cancel()
	})

	It("должен регистрировать и логинить пользователя", func() {
		client, err := newAPIClient()
		Expect(err).ToNot(HaveOccurred())

		reg := apiclient.FakeRegister()
		resp, err := client.Register(ctx, reg)
		Expect(err).ToNot(HaveOccurred())
		defer func() { _ = resp.Body.Close() }()
		Expect(resp.StatusCode).To(Equal(http.StatusCreated), "register: %s", resp.Status)

		resp, err = client.Login(ctx, apiclient.FakeLogin(reg.Email, reg.Password))
		Expect(err).ToNot(HaveOccurred())
		defer func() { _ = resp.Body.Close() }()
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "login: %s", resp.Status)
	})

	It("должен выполнять logout", func() {
		client := newAuthenticatedClient(ctx)
		resp, err := client.Logout(ctx)
		Expect(err).ToNot(HaveOccurred())
		defer func() { _ = resp.Body.Close() }()
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})

var _ = Describe("Teams", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(suiteCtx, requestTimeout)
	})

	AfterEach(func() {
		cancel()
	})

	It("CRUD: создание, список, получение по ID", func() {
		client := newAuthenticatedClient(ctx)

		createdTeams := make([]apiclient.Team, 0, createTeamsCount)
		for i := 0; i < createTeamsCount; i++ {
			createReq := apiclient.FakeCreateTeam()
			resp, err := client.CreateTeam(ctx, createReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated), "create team %d: %s", i+1, resp.Status)
			var team apiclient.Team
			Expect(json.NewDecoder(resp.Body).Decode(&team)).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			Expect(team.ID).ToNot(BeEmpty())
			Expect(team.Name).To(Equal(createReq.Name))
			createdTeams = append(createdTeams, team)
		}

		teams, err := client.ListTeams(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(teams)).To(BeNumerically(">=", createTeamsCount))
		for _, created := range createdTeams {
			var found bool
			for _, tt := range teams {
				if tt.ID == created.ID {
					found = true
					Expect(tt.Name).To(Equal(created.Name))
					break
				}
			}
			Expect(found).To(BeTrue(), "команда %s должна быть в списке", created.ID)
		}

		for _, created := range createdTeams {
			resp, err := client.GetTeam(ctx, created.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var got apiclient.TeamWithMembersResponse
			Expect(json.NewDecoder(resp.Body).Decode(&got)).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			Expect(got.Team.ID).To(Equal(created.ID))
			Expect(got.Team.Name).To(Equal(created.Name))
		}
	})

	It("должен приглашать пользователя в команду", func() {
		clientB, err := newAPIClient()
		Expect(err).ToNot(HaveOccurred())
		regB := apiclient.FakeRegister()
		resp, err := clientB.Register(ctx, regB)
		Expect(err).ToNot(HaveOccurred())
		_ = resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(http.StatusCreated), "register second user: %s", resp.Status)

		clientOwner, teamID := newAuthenticatedClientWithTeam(ctx)

		inviteReq := apiclient.FakeInvite(regB.Email, "member")
		resp, err = clientOwner.Invite(ctx, teamID, inviteReq)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated), "invite: %s", resp.Status)
		var inviteResp apiclient.InviteResponse
		Expect(json.NewDecoder(resp.Body).Decode(&inviteResp)).ToNot(HaveOccurred())
		_ = resp.Body.Close()
		Expect(inviteResp.Success).To(BeTrue())
		Expect(inviteResp.Invitation.Email).To(Equal(regB.Email))
		Expect(inviteResp.Invitation.TeamID).To(Equal(teamID), "API должен возвращать team_id в ответе приглашения")
		Expect(inviteResp.Invitation.Role).To(Equal("member"))
		Expect(inviteResp.Invitation.ID).ToNot(BeEmpty())
	})
})

var _ = Describe("Tasks", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(suiteCtx, requestTimeout)
	})

	AfterEach(func() {
		cancel()
	})

	It("CRUD: создание, список, получение, обновление", func() {
		client, teamID := newAuthenticatedClientWithTeam(ctx)

		createdTasks := make([]apiclient.Task, 0, createTasksCount)
		for i := 0; i < createTasksCount; i++ {
			createReq := apiclient.FakeCreateTask(teamID)
			resp, err := client.CreateTask(ctx, createReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated), "create task %d: %s", i+1, resp.Status)
			var task apiclient.Task
			Expect(json.NewDecoder(resp.Body).Decode(&task)).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			Expect(task.ID).ToNot(BeEmpty())
			Expect(task.Title).To(Equal(createReq.Title))
			Expect(task.TeamID).To(Equal(teamID))
			createdTasks = append(createdTasks, task)
		}

		list, err := client.ListTasks(ctx, apiclient.ListTasksOpts{TeamID: teamID})
		Expect(err).ToNot(HaveOccurred())
		Expect(list.Total).To(BeNumerically(">=", createTasksCount))
		for _, created := range createdTasks {
			var foundInList bool
			for _, item := range list.Items {
				if item.ID == created.ID {
					foundInList = true
					break
				}
			}
			Expect(foundInList).To(BeTrue(), "задача %s должна быть в списке", created.ID)
		}

		for _, created := range createdTasks {
			resp, err := client.GetTask(ctx, created.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var got apiclient.Task
			Expect(json.NewDecoder(resp.Body).Decode(&got)).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			Expect(got.ID).To(Equal(created.ID))
			Expect(got.Title).To(Equal(created.Title))
		}

		for _, created := range createdTasks {
			updateReq := apiclient.FakeUpdateTask()
			resp, err := client.UpdateTask(ctx, created.ID, updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "update task %s: %s", created.ID, resp.Status)
			_ = resp.Body.Close()

			resp, err = client.GetTask(ctx, created.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var got apiclient.Task
			Expect(json.NewDecoder(resp.Body).Decode(&got)).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			Expect(got.Title).To(Equal(updateReq.Title))
			Expect(got.Status).To(Equal(updateReq.Status))
		}
	})
})

var _ = Describe("Task history", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(suiteCtx, requestTimeout)
	})

	AfterEach(func() {
		cancel()
	})

	It("должен возвращать историю изменений задачи", func() {
		client, teamID := newAuthenticatedClientWithTeam(ctx)

		createReq := apiclient.FakeCreateTask(teamID)
		resp, err := client.CreateTask(ctx, createReq)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		var task apiclient.Task
		Expect(json.NewDecoder(resp.Body).Decode(&task)).ToNot(HaveOccurred())
		_ = resp.Body.Close()
		Expect(task.ID).ToNot(BeEmpty())

		resp, err = client.GetTaskHistory(ctx, task.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		var history apiclient.TaskHistoryResponse
		Expect(json.NewDecoder(resp.Body).Decode(&history)).ToNot(HaveOccurred())
		_ = resp.Body.Close()
		Expect(history.TaskID).To(Equal(task.ID))

		for i := 0; i < updateHistoryCount; i++ {
			updateReq := apiclient.FakeUpdateTask()
			resp, err = client.UpdateTask(ctx, task.ID, updateReq)
			Expect(err).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		}

		resp, err = client.GetTaskHistory(ctx, task.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(json.NewDecoder(resp.Body).Decode(&history)).ToNot(HaveOccurred())
		_ = resp.Body.Close()
		Expect(history.TaskID).To(Equal(task.ID))
		Expect(len(history.Entries)).To(BeNumerically(">=", updateHistoryCount),
			"после %d обновлений в истории должно быть не меньше %d записей", updateHistoryCount, updateHistoryCount)
	})
})

var _ = Describe("Comments", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(suiteCtx, requestTimeout)
	})

	AfterEach(func() {
		cancel()
	})

	It("CRUD: создание и список комментариев", func() {
		client, teamID := newAuthenticatedClientWithTeam(ctx)

		taskReq := apiclient.FakeCreateTask(teamID)
		resp, err := client.CreateTask(ctx, taskReq)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		var task apiclient.Task
		Expect(json.NewDecoder(resp.Body).Decode(&task)).ToNot(HaveOccurred())
		_ = resp.Body.Close()

		resp, err = client.ListComments(ctx, task.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		var list apiclient.CommentListResponse
		Expect(json.NewDecoder(resp.Body).Decode(&list)).ToNot(HaveOccurred())
		_ = resp.Body.Close()
		Expect(list.Items).To(BeEmpty())

		createdComments := make([]apiclient.Comment, 0, createCommentsCount)
		for i := 0; i < createCommentsCount; i++ {
			commentReq := apiclient.FakeCreateComment()
			resp, err = client.CreateComment(ctx, task.ID, commentReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated), "create comment %d: %s", i+1, resp.Status)
			var comment apiclient.Comment
			Expect(json.NewDecoder(resp.Body).Decode(&comment)).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			Expect(comment.ID).ToNot(BeEmpty())
			Expect(comment.TaskID).To(Equal(task.ID))
			Expect(comment.Content).To(Equal(commentReq.Content))
			createdComments = append(createdComments, comment)
		}

		resp, err = client.ListComments(ctx, task.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(json.NewDecoder(resp.Body).Decode(&list)).ToNot(HaveOccurred())
		_ = resp.Body.Close()
		Expect(list.Items).To(HaveLen(createCommentsCount))
		for _, created := range createdComments {
			var found bool
			for _, c := range list.Items {
				if c.ID == created.ID {
					Expect(c.Content).To(Equal(created.Content))
					Expect(c.TaskID).To(Equal(task.ID))
					found = true
					break
				}
			}
			Expect(found).To(BeTrue(), "комментарий %s должен быть в списке", created.ID)
		}
	})
})

var _ = Describe("Reports", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(suiteCtx, requestTimeout)
	})

	AfterEach(func() {
		cancel()
	})

	It("должны возвращать отчёты по командам и задачам", func() {
		client, teamID := newAuthenticatedClientWithTeam(ctx)

		for i := 0; i < createTeamsCount; i++ {
			teamReq := apiclient.FakeCreateTeam()
			resp, err := client.CreateTeam(ctx, teamReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			var team apiclient.Team
			Expect(json.NewDecoder(resp.Body).Decode(&team)).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			for j := 0; j < createTasksCount; j++ {
				taskReq := apiclient.FakeCreateTask(team.ID)
				resp, err = client.CreateTask(ctx, taskReq)
				Expect(err).ToNot(HaveOccurred())
				_ = resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			}
		}
		for j := 0; j < createTasksCount; j++ {
			taskReq := apiclient.FakeCreateTask(teamID)
			resp, err := client.CreateTask(ctx, taskReq)
			Expect(err).ToNot(HaveOccurred())
			_ = resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		}

		for i := 0; i < reportsCallCount; i++ {
			resp, err := client.ReportTeamStats(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "ReportTeamStats call %d: %s", i+1, resp.Status)
			_ = resp.Body.Close()
		}
		for i := 0; i < reportsCallCount; i++ {
			resp, err := client.ReportTopCreators(ctx)
			Expect(err).ToNot(HaveOccurred())
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "ReportTopCreators call %d: %s body=%s", i+1, resp.Status, string(body))
		}
		for i := 0; i < reportsCallCount; i++ {
			resp, err := client.ReportInvalidAssignees(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "ReportInvalidAssignees call %d: %s", i+1, resp.Status)
			_ = resp.Body.Close()
		}
	})
})
