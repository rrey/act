package workflow_test

import (
	"testing"

	"github.com/nektos/act/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowFromFile(t *testing.T) {
	var summaries model.WorkflowSummaries
	summaries.FromFile("workflow-summary.json")

	// Check that the "Demo" workflow is present
	assert.Contains(t, summaries, "Demo")

	// Check that the "Demo" workflow has "job1"
	assert.Contains(t, summaries["Demo"].Jobs, "job1")
	workflow := summaries["Demo"]
	// Missing -> job outcome
	// assert.Equal(t, workflow.GetJob("job1").Outcome, "success")
	job := workflow.Jobs["job1"]
	jobSteps := job.Steps
	// Check that the "job1" has 3 steps
	assert.Len(t, jobSteps, 3)
	// Check that all steps are successful
	for _, step := range jobSteps {
		assert.Equal(t, step.Outcome.String(), "success")
	}

	step1 := job.StepByID("step1")
	// Check the outcome of step with id 'step1' in the "job1" job
	assert.Equal(t, step1.StepResult.Outcome.String(), "success")

	step2 := job.StepByID("step2")
	assert.Equal(t, step2.StepResult.Outcome.String(), "success")
	// Check that 'step2' has an output named 'foo'
	assert.Contains(t, step2.Outputs, "foo", step2.Outputs)
	// Check that 'step2' has an output named 'foo' with value 'bar'
	assert.Contains(t, step2.Outputs["foo"], "bar")

	assert.Contains(t, summaries["Demo"].Jobs, "job2")
}
