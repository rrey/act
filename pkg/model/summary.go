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

// StepByID sets the outputs of the step with the given ID
func (js JobSummary) StepByID(Id string) *StepSummary {
	for _, s := range js.Steps {
		if s.ID == Id {
			return &s
		}
	}
	return nil
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

	err = os.WriteFile("workflow-summary.json", summaryJSON, 0o600)
	if err != nil {
		log.Errorf("Failed to write workflow summary to file: %v", err)
		return err
	}
	return nil
}

func (ws *WorkflowSummaries) FromFile(path string) *WorkflowSummaries {
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file: #%v ", err)
	}
	err = json.Unmarshal(jsonFile, ws)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return ws
}

var Summary = make(WorkflowSummaries)
