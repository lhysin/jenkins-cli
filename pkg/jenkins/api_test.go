package jenkins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListJobs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/json", r.URL.Path)
		assert.Equal(t, "Basic dXNlcjp0b2tlbg==", r.Header.Get("Authorization"))

		response := JobListResponse{
			Jobs: []Job{
				{Name: "job1", URL: "http://localhost/job1", Color: "blue"},
				{Name: "job2", URL: "http://localhost/job2", Color: "red"},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	jobs, err := client.ListJobs()

	assert.NoError(t, err)
	assert.Len(t, jobs, 2)
	assert.Equal(t, "job1", jobs[0].Name)
	assert.Equal(t, "blue", jobs[0].Color)
}

func TestListJobs_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	_, err := client.ListJobs()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}

func TestListBuilds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/job/%s/api/json", "my-job")
		assert.Equal(t, expectedPath, r.URL.Path)

		response := BuildListResponse{
			Builds: []Build{
				{Number: 1, URL: "http://localhost/job/my-job/1", Result: "SUCCESS", Duration: 60000},
				{Number: 2, URL: "http://localhost/job/my-job/2", Result: "FAILURE", Duration: 30000},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	builds, err := client.ListBuilds("my-job")

	assert.NoError(t, err)
	assert.Len(t, builds, 2)
	assert.Equal(t, 1, builds[0].Number)
	assert.Equal(t, "SUCCESS", builds[0].Result)
}

func TestListNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/computer/api/json", r.URL.Path)

		response := NodeListResponse{
			Computer: []Node{
				{DisplayName: "master", Offline: false},
				{DisplayName: "agent-1", Offline: true},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	nodes, err := client.ListNodes()

	assert.NoError(t, err)
	assert.Len(t, nodes, 2)
	assert.Equal(t, "master", nodes[0].DisplayName)
	assert.False(t, nodes[0].Offline)
	assert.True(t, nodes[1].Offline)
}

func TestGetBuildInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/job/%s/%d/api/json", "my-job", 42)
		assert.Equal(t, expectedPath, r.URL.Path)

		build := Build{
			Number:    42,
			URL:       "http://localhost/job/my-job/42",
			Result:    "SUCCESS",
			Duration:  120000,
			Timestamp: 1700000000,
		}
		json.NewEncoder(w).Encode(build)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	build, err := client.GetBuildInfo("my-job", 42)

	assert.NoError(t, err)
	assert.Equal(t, 42, build.Number)
	assert.Equal(t, "SUCCESS", build.Result)
	assert.Equal(t, int64(120000), build.Duration)
}

func TestCreateJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/createItem")
		assert.Equal(t, "Basic dXNlcjp0b2tlbg==", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	err := client.CreateJob("test-job", "<xml>config</xml>")

	assert.NoError(t, err)
}

func TestDeleteJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/job/test-job/doDelete", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	err := client.DeleteJob("test-job")

	assert.NoError(t, err)
}

func TestUpdateJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/job/test-job/config.xml", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	err := client.UpdateJob("test-job", "<xml>updated</xml>")

	assert.NoError(t, err)
}

func TestTriggerBuild(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/job/test-job/build", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	err := client.TriggerBuild("test-job")

	assert.NoError(t, err)
}

func TestGetBuildLogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/job/%s/%d/logText/progressiveText", "my-job", 1)
		assert.Equal(t, expectedPath, r.URL.Path)
		w.Write([]byte("Build log output"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	logs, err := client.GetBuildLogs("my-job", 1)

	assert.NoError(t, err)
	assert.Contains(t, logs, "Build log output")
}

func TestListJobs_Anonymous(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		assert.Equal(t, "", auth)

		response := JobListResponse{
			Jobs: []Job{
				{Name: "public-job", URL: "http://localhost/public-job", Color: "blue"},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "")
	jobs, err := client.ListJobs()

	assert.NoError(t, err)
	assert.Len(t, jobs, 1)
	assert.Equal(t, "public-job", jobs[0].Name)
}

func TestGetJobInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/job/%s/api/json", "test-job")
		assert.Equal(t, expectedPath, r.URL.Path)

		job := Job{
			Name:        "test-job",
			URL:         "http://localhost/job/test-job",
			Color:       "blue",
			Description: "Test job description",
		}
		json.NewEncoder(w).Encode(job)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	job, err := client.GetJobInfo("test-job")

	assert.NoError(t, err)
	assert.Equal(t, "test-job", job.Name)
	assert.Equal(t, "Test job description", job.Description)
}

func TestTriggerBuild_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "wrong-token")
	err := client.TriggerBuild("test-job")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}

func TestGetBuildLogs_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "token")
	_, err := client.GetBuildLogs("nonexistent", 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}
