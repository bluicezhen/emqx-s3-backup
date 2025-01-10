package emqx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

/**
 * @description: NewEMQX
 * @return {*EMQX}
 */
func NewEMQX() *EMQX {
	emqxUrl := os.Getenv("EMQX_URL")
	if emqxUrl == "" {
		panic("EMQX_URL is not set")
	}

	apiName := os.Getenv("EMQX_API_NAME")
	if apiName == "" {
		panic("EMQX_API_NAME is not set")
	}

	apiPass := os.Getenv("EMQX_API_PASS")
	if apiPass == "" {
		panic("EMQX_API_PASS is not set")
	}

	return &EMQX{
		logger:  log.New(os.Stdout, "EMQX: ", log.LstdFlags),
		emqxUrl: emqxUrl,
		apiName: apiName,
		apiPass: apiPass,
	}
}

func (e *EMQX) request(method string, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, e.emqxUrl+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(e.apiName, e.apiPass)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.logger.Printf("Failed to request: %s, %v", path, err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var errorResp EmqxErrorResp
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			e.logger.Printf("Failed to decode error response: %v", err)
			return nil, fmt.Errorf("request %s failed with status code: %d", path, resp.StatusCode)
		}
		e.logger.Printf("Request %s failed with status code: %d, error code: %s, message: %s", path, resp.StatusCode, errorResp.Code, errorResp.Message)
		return nil, fmt.Errorf("request %s failed with status code: %d, error code: %s, message: %s", path, resp.StatusCode, errorResp.Code, errorResp.Message)
	}
	return resp, nil
}

/**
 * @description: Export a data backup file
 * @return {*EmqxExportDataResp, error} response
 */
func (e *EMQX) DataExport() (*EmqxExportDataResp, error) {
	payload := map[string]interface{}{
		"root_keys": []string{
			"connectors",
			"actions",
			"sources",
			"rule_engine",
			"schema_registry",
		},
		"table_sets": []string{
			"banned",
			"builtin_authn",
			"builtin_authn_scram",
			"builtin_authz",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		e.logger.Printf("Failed to marshal payload: %v", err)
		return nil, err
	}

	resp, err := e.request("POST", "/api/v5/data/export", jsonData)
	if err != nil {
		e.logger.Printf("Failed to export data: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	var exportResp EmqxExportDataResp
	if err := json.NewDecoder(resp.Body).Decode(&exportResp); err != nil {
		e.logger.Printf("Failed to decode response: %v", err)
		return nil, err
	}

	return &exportResp, nil
}

/**
 * @description: Download a data backup file
 * @param {string} filename
 * @param {string} node
 * @return {*os.File, error} file
 */
func (e *EMQX) DownloadData(filename string, node string) (*os.File, error) {
	resp, err := e.request("GET", "/api/v5/data/files/"+filename+"?node="+node, nil)
	if err != nil {
		e.logger.Printf("Failed to download data: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp("", "emqx-export-*.json")
	if err != nil {
		e.logger.Printf("Failed to create temp file: %v", err)
		return nil, err
	}

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		e.logger.Printf("Failed to write response to file: %v", err)
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, err
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		e.logger.Printf("Failed to seek to start of file: %v", err)
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, err
	}

	return tempFile, nil
}
