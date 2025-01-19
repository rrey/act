package model

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type StepSummary struct {
	StepResult
	EnvVars map[string]string `json:"env"`
	ID      string            `json:"id"`
}

type JobSummary struct {
	EnvVars map[string]string `json:"env"`
	Outputs map[string]string `json:"outputs"`
	Steps   []StepSummary     `json:"steps"`
}

func (js JobSummary) SetEnvVars(envVars map[string]string) {
	for k, v := range envVars {
		js.EnvVars[k] = v
	}
}

func (js JobSummary) SetOutputs(outputs map[string]string) {
	for k, v := range outputs {
		js.Outputs[k] = v
	}
}

func NewJobSummary() JobSummary {
	return JobSummary{
		EnvVars: map[string]string{},
		Outputs: map[string]string{},
		Steps:   []StepSummary{},
	}
}

type WorkflowSummary struct {
	WorkflowName string                `json:"name"`
	Jobs         map[string]JobSummary `json:"jobs"`
}

func NewWorkflowSummary(name string) WorkflowSummary {
	return WorkflowSummary{
		WorkflowName: name,
		Jobs:         map[string]JobSummary{},
	}
}

func (ws WorkflowSummary) AddJob(jobID string) {
	ws.Jobs[jobID] = NewJobSummary()
}

type WorkflowSummaries map[string]WorkflowSummary

func (ws WorkflowSummaries) Write() error {
	summaryJSON, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		log.Errorf("Failed to serialize workflow summary: %v", err)
		return err
	}

	err = os.WriteFile("workflow-summary.json", summaryJSON, 0o644)
	if err != nil {
		log.Errorf("Failed to write workflow summary to file: %v", err)
		return err
	}
	return nil
}

var Summary WorkflowSummaries = make(WorkflowSummaries)
