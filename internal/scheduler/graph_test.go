package scheduler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yohamta/dagu/internal/config"
	"github.com/yohamta/dagu/internal/scheduler"
)

func TestCycleDetection(t *testing.T) {
	step1 := &config.Step{}
	step1.Name = "1"
	step1.Depends = []string{"2"}

	step2 := &config.Step{}
	step2.Name = "2"
	step2.Depends = []string{"1"}

	_, err := scheduler.NewExecutionGraph(step1, step2)

	if err == nil {
		t.Fatal("cycle detection should be detected.")
	}
}

func TestRetryExecution(t *testing.T) {
	nodes := []*scheduler.Node{
		{
			Step: &config.Step{Name: "1", Command: "true"},
			NodeState: scheduler.NodeState{
				Status: scheduler.NodeStatusSuccess,
			},
		},
		{
			Step: &config.Step{Name: "2", Command: "true", Depends: []string{"1"}},
			NodeState: scheduler.NodeState{
				Status: scheduler.NodeStatusError,
			},
		},
		{
			Step: &config.Step{Name: "3", Command: "true", Depends: []string{"2"}},
			NodeState: scheduler.NodeState{
				Status: scheduler.NodeStatusCancel,
			},
		},
		{
			Step: &config.Step{Name: "4", Command: "true", Depends: []string{}},
			NodeState: scheduler.NodeState{
				Status: scheduler.NodeStatusSkipped,
			},
		},
		{
			Step: &config.Step{Name: "5", Command: "true", Depends: []string{"4"}},
			NodeState: scheduler.NodeState{
				Status: scheduler.NodeStatusError,
			},
		},
		{
			Step: &config.Step{Name: "6", Command: "true", Depends: []string{"5"}},
			NodeState: scheduler.NodeState{
				Status: scheduler.NodeStatusSuccess,
			},
		},
		{
			Step: &config.Step{Name: "7", Command: "true", Depends: []string{"6"}},
			NodeState: scheduler.NodeState{
				Status: scheduler.NodeStatusSkipped,
			},
		},
		{
			Step: &config.Step{Name: "8", Command: "true", Depends: []string{}},
			NodeState: scheduler.NodeState{
				Status: scheduler.NodeStatusSkipped,
			},
		},
	}
	_, err := scheduler.RetryExecutionGraph(nodes...)
	require.NoError(t, err)
	assert.Equal(t, scheduler.NodeStatusSuccess, nodes[0].Status)
	assert.Equal(t, scheduler.NodeStatusNone, nodes[1].Status)
	assert.Equal(t, scheduler.NodeStatusNone, nodes[2].Status)
	assert.Equal(t, scheduler.NodeStatusSkipped, nodes[3].Status)
	assert.Equal(t, scheduler.NodeStatusNone, nodes[4].Status)
	assert.Equal(t, scheduler.NodeStatusNone, nodes[5].Status)
	assert.Equal(t, scheduler.NodeStatusNone, nodes[6].Status)
	assert.Equal(t, scheduler.NodeStatusSkipped, nodes[7].Status)
}
