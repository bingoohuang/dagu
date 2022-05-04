package models

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yohamta/dagu/internal/config"
	"github.com/yohamta/dagu/internal/scheduler"
	"github.com/yohamta/dagu/internal/utils"
)

type Node struct {
	*config.Step `json:"Step"`
	Log          string               `json:"Log"`
	StartedAt    string               `json:"StartedAt"`
	FinishedAt   string               `json:"FinishedAt"`
	Status       scheduler.NodeStatus `json:"Status"`
	RetryCount   int                  `json:"RetryCount"`
	DoneCount    int                  `json:"DoneCount"`
	Error        string               `json:"Error"`
	StatusText   string               `json:"StatusText"`
}

func (n *Node) ToNode() *scheduler.Node {
	startedAt, _ := utils.ParseTime(n.StartedAt)
	finishedAt, _ := utils.ParseTime(n.FinishedAt)
	var err error = nil
	if n.Error != "" {
		err = fmt.Errorf(n.Error)
	}
	ret := &scheduler.Node{
		Step: n.Step,
		NodeState: scheduler.NodeState{
			Status:     n.Status,
			Log:        n.Log,
			StartedAt:  startedAt,
			FinishedAt: finishedAt,
			RetryCount: n.RetryCount,
			DoneCount:  n.DoneCount,
			Error:      err,
		},
	}
	return ret
}

func FromNode(n *scheduler.Node) *Node {
	node := &Node{
		Step:       n.Step,
		Log:        n.Log,
		StartedAt:  utils.FormatTime(n.StartedAt),
		FinishedAt: utils.FormatTime(n.FinishedAt),
		Status:     n.ReadStatus(),
		StatusText: n.ReadStatus().String(),
		RetryCount: n.ReadRetryCount(),
		DoneCount:  n.ReadDoneCount(),
	}
	if n.Error != nil {
		node.Error = n.Error.Error()
	}
	return node
}

func FromNodes(nodes []*scheduler.Node) []*Node {
	ret := []*Node{}
	for _, n := range nodes {
		ret = append(ret, FromNode(n))
	}
	return ret
}

func FromSteps(steps []*config.Step) []*Node {
	var ret []*Node
	for _, s := range steps {
		ret = append(ret, fromStepWithDefValues(s))
	}
	return ret
}

func StepGraph(steps []*Node, displayStatus bool) string {
	var buf bytes.Buffer
	buf.WriteString("flowchart LR;")
	for _, s := range steps {
		buf.WriteString(fmt.Sprintf("%s(%s)", graphNode(s.Name), s.Name))
		if displayStatus {
			switch s.Status {
			case scheduler.NodeStatusRunning:
				buf.WriteString(":::running")
			case scheduler.NodeStatusError:
				buf.WriteString(":::error")
			case scheduler.NodeStatusCancel:
				buf.WriteString(":::cancel")
			case scheduler.NodeStatusSuccess:
				buf.WriteString(":::done")
			case scheduler.NodeStatusSkipped:
				buf.WriteString(":::skipped")
			default:
				buf.WriteString(":::none")
			}
		} else {
			buf.WriteString(":::none")
		}
		buf.WriteString(";")
		for _, d := range s.Depends {
			buf.WriteString(graphNode(d) + "-->" + graphNode(s.Name) + ";")
		}
	}
	buf.WriteString("classDef none fill:white,stroke:lightblue,stroke-width:2px\n")
	buf.WriteString("classDef running fill:white,stroke:lime,stroke-width:2px\n")
	buf.WriteString("classDef error fill:white,stroke:red,stroke-width:2px\n")
	buf.WriteString("classDef cancel fill:white,stroke:pink,stroke-width:2px\n")
	buf.WriteString("classDef done fill:white,stroke:green,stroke-width:2px\n")
	buf.WriteString("classDef skipped fill:white,stroke:gray,stroke-width:2px\n")
	return buf.String()
}

func graphNode(val string) string {
	return strings.ReplaceAll(val, " ", "_")
}

func fromStepWithDefValues(s *config.Step) *Node {
	if s == nil {
		return nil
	}
	step := &Node{
		Step:       s,
		Log:        "",
		StartedAt:  "-",
		FinishedAt: "-",
		Status:     scheduler.NodeStatusNone,
		StatusText: scheduler.NodeStatusNone.String(),
		RetryCount: 0,
	}
	return step
}
