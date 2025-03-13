package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// APIClient is a simple REST client for interacting with the API.
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	ProjectID  string
	Logger     *log.Logger // Add a logger
}

// NewAPIClient creates a new APIClient.
func NewAPIClient() *APIClient {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		panic("BASE_URL environment variable not set")
	}

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		panic("PROJECT_ID environment variable not set")
	}

	// Use Application Default Credentials (ADC) for authentication
	// This will look for credentials in the environment variables,
	// well-known file locations, or the metadata server.
	// See: https://cloud.google.com/docs/authentication/application-default-credentials
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		panic(fmt.Sprintf("Error finding default credentials: %v", err))
	}

	// Create an HTTP client that uses the token source from the credentials
	httpClient := oauth2.NewClient(ctx, credentials.TokenSource)

	// Create a logger
	logger := log.New(os.Stderr, "APIClient: ", log.LstdFlags|log.Lshortfile)

	return &APIClient{BaseURL: baseURL, HTTPClient: httpClient, ProjectID: projectID, Logger: logger}
}

// ListModels lists all models.
func (c *APIClient) ListModels() ([]Model, error) {
	url := c.BaseURL + "/v1alpha1/models"
	result, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}
	var listModelsResponse ListModelsResponse
	if err := convert(result, &listModelsResponse); err != nil {
		return nil, err
	}

	// Convert the modelNames array to an array of Model objects
	models := make([]Model, len(listModelsResponse.ModelNames))
	for i, name := range listModelsResponse.ModelNames {
		models[i] = Model{Name: name}
	}

	return models, nil
}

// ListModelServers lists all model servers for a given model.
func (c *APIClient) ListModelServers(modelName string) ([]ModelServer, error) {
	u, err := url.Parse(c.BaseURL + "/v1alpha1/modelServers")
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	q := u.Query()
	q.Set("model_name", modelName)
	u.RawQuery = q.Encode()
	result, err := c.getRequest(u.String())
	if err != nil {
		return nil, err
	}
	var listModelServersResponse ListModelServersResponse
	if err := convert(result, &listModelServersResponse); err != nil {
		return nil, err
	}
	// Convert the modelServerNames array to an array of ModelServer objects
	modelServers := make([]ModelServer, len(listModelServersResponse.ModelServerNames))
	for i, name := range listModelServersResponse.ModelServerNames {
		modelServers[i] = ModelServer{Name: name}
	}
	return modelServers, nil
}

// ListModelServerVersions lists all model server versions for a given model and model server.
func (c *APIClient) ListModelServerVersions(modelName, modelServerName string) ([]ModelServerVersion, error) {
	u, err := url.Parse(c.BaseURL + "/v1alpha1/modelServers/" + modelServerName + "/versions")
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	q := u.Query()
	q.Set("model_name", modelName)
	u.RawQuery = q.Encode()
	result, err := c.getRequest(u.String())
	if err != nil {
		return nil, err
	}
	var listModelServerVersionsResponse ListModelServerVersionsResponse
	if err := convert(result, &listModelServerVersionsResponse); err != nil {
		return nil, err
	}
	// Convert the modelServerVersions array to an array of ModelServerVersion objects
	modelServerVersions := make([]ModelServerVersion, len(listModelServerVersionsResponse.ModelServerVersions))
	for i, name := range listModelServerVersionsResponse.ModelServerVersions {
		modelServerVersions[i] = ModelServerVersion{Name: name}
	}
	return modelServerVersions, nil
}

// ListAccelerators lists all accelerators for a given model and model server.
func (c *APIClient) ListAccelerators(modelName, modelServerName string) (*ListAcceleratorsResponse, error) {
	u, err := url.Parse(c.BaseURL + "/v1alpha1/accelerators")
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	q := u.Query()
	q.Set("model_name", modelName)
	q.Set("model_server_name", modelServerName)
	u.RawQuery = q.Encode()
	result, err := c.getRequest(u.String())
	if err != nil {
		return nil, err
	}
	var listAcceleratorsResponse ListAcceleratorsResponse
	if err := convert(result, &listAcceleratorsResponse); err != nil {
		return nil, err
	}
	return &listAcceleratorsResponse, nil
}

// CreateManifest creates a new manifest.
func (c *APIClient) CreateManifest(modelName, modelServerName, modelServerVersion, acceleratorType string, targetNtpotMilliseconds int) (*CreateManifestResponse, error) {
	u, err := url.Parse(c.BaseURL + "/v1alpha1/optimizedManifest")
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	q := u.Query()
	q.Set("model_and_model_server_info.model_name", modelName)
	q.Set("model_and_model_server_info.model_server_name", modelServerName)
	q.Set("model_and_model_server_info.model_server_version", modelServerVersion)
	q.Set("accelerator_type", acceleratorType)
	if targetNtpotMilliseconds > 0 {
		q.Set("target_ntpot_milliseconds", fmt.Sprintf("%d", targetNtpotMilliseconds))
	}
	u.RawQuery = q.Encode()

	result, err := c.getRequest(u.String())
	if err != nil {
		return nil, err
	}

	var createManifestResponse CreateManifestResponse
	if err := convert(result, &createManifestResponse); err != nil {
		return nil, err
	}
	return &createManifestResponse, nil
}

// ListModelsAndServers lists all models and servers.
func (c *APIClient) ListModelsAndServers() ([]ModelsAndServers, error) {
	url := c.BaseURL + "/v1alpha1/modelsAndServers"
	result, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}
	var modelsAndServers []ModelsAndServers
	if err := convert(result, &modelsAndServers); err != nil {
		return nil, err
	}
	return modelsAndServers, nil
}

func (c *APIClient) getRequest(urlStr string) (interface{}, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add the X-Goog-User-Project header
	req.Header.Add("X-Goog-User-Project", c.ProjectID)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		// Log the request headers on error
		dump, dumpErr := httputil.DumpRequestOut(req, false)
		if dumpErr != nil {
			c.Logger.Printf("Error dumping request: %v", dumpErr)
		} else {
			c.Logger.Printf("Request Headers on Error: \n%s", dump)
		}
		return nil, err
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}

func convert(source interface{}, target interface{}) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, target)
	if err != nil {
		return err
	}
	return nil
}
