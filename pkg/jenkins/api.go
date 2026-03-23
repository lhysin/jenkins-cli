package jenkins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) doRequest(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.URL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.User != "" && c.Token != "" {
		req.SetBasicAuth(c.User, c.Token)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) ListJobs() ([]Job, error) {
	data, err := c.doRequest("/api/json?tree=jobs[name,url,color,description]")
	if err != nil {
		return nil, err
	}

	var response JobListResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Jobs, nil
}

func (c *Client) ListBuilds(jobName string) ([]Build, error) {
	endpoint := fmt.Sprintf("/job/%s/api/json?tree=builds[number,url,result,duration,timestamp,displayName]", jobName)
	data, err := c.doRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var response BuildListResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Builds, nil
}

func (c *Client) ListNodes() ([]Node, error) {
	data, err := c.doRequest("/computer/api/json?tree=computer[displayName,offline]")
	if err != nil {
		return nil, err
	}

	var response NodeListResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Computer, nil
}

func (c *Client) GetJobInfo(jobName string) (*Job, error) {
	endpoint := fmt.Sprintf("/job/%s/api/json", jobName)
	data, err := c.doRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var job Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) GetBuildInfo(jobName string, buildNumber int) (*Build, error) {
	endpoint := fmt.Sprintf("/job/%s/%d/api/json", jobName, buildNumber)
	data, err := c.doRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var build Build
	if err := json.Unmarshal(data, &build); err != nil {
		return nil, err
	}

	return &build, nil
}

func (c *Client) doPostRequest(endpoint string, body []byte) error {
	urlStr := fmt.Sprintf("%s%s", c.URL, endpoint)
	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.User, c.Token)
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (c *Client) doDeleteRequest(endpoint string) error {
	urlStr := fmt.Sprintf("%s%s", c.URL, endpoint)
	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.User, c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (c *Client) CreateJob(name string, config string) error {
	endpoint := fmt.Sprintf("/createItem?name=%s", url.QueryEscape(name))
	return c.doPostRequest(endpoint, []byte(config))
}

func (c *Client) DeleteJob(name string) error {
	endpoint := fmt.Sprintf("/job/%s/doDelete", url.QueryEscape(name))
	return c.doDeleteRequest(endpoint)
}

func (c *Client) UpdateJob(name string, config string) error {
	endpoint := fmt.Sprintf("/job/%s/config.xml", url.QueryEscape(name))
	return c.doPostRequest(endpoint, []byte(config))
}

func (c *Client) TriggerBuild(name string) error {
	endpoint := fmt.Sprintf("/job/%s/build", url.QueryEscape(name))
	return c.doPostRequest(endpoint, nil)
}

func (c *Client) GetBuildLogs(jobName string, buildNumber int) (string, error) {
	endpoint := fmt.Sprintf("/job/%s/%d/logText/progressiveText", url.QueryEscape(jobName), buildNumber)
	data, err := c.doRequest(endpoint)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
