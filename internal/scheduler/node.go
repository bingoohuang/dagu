package scheduler

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/yohamta/dagu/internal/config"
	"github.com/yohamta/dagu/internal/utils"
)

type NodeStatus int

const (
	NodeStatusNone NodeStatus = iota
	NodeStatusRunning
	NodeStatusError
	NodeStatusCancel
	NodeStatusSuccess
	NodeStatusSkipped
)

func (s NodeStatus) String() string {
	switch s {
	case NodeStatusRunning:
		return "running"
	case NodeStatusError:
		return "failed"
	case NodeStatusCancel:
		return "canceled"
	case NodeStatusSuccess:
		return "finished"
	case NodeStatusSkipped:
		return "skipped"
	case NodeStatusNone:
		fallthrough
	default:
		return "not started"
	}
}

type Node struct {
	*config.Step
	NodeState
	id         int
	mu         sync.RWMutex
	cmd        *exec.Cmd
	cancelFunc func()
	logFile    *os.File
	logWriter  *bufio.Writer
}

type NodeState struct {
	Status     NodeStatus
	Log        string
	StartedAt  time.Time
	FinishedAt time.Time
	RetryCount int
	DoneCount  int
	Error      error
}

func (n *Node) Execute() error {
	ctx, fn := context.WithCancel(context.Background())
	n.cancelFunc = fn
	cmd := exec.CommandContext(ctx, n.Command, n.Args...)
	n.cmd = cmd
	cmd.Dir = n.Dir
	cmd.Env = append(cmd.Env, n.Variables...)

	if n.logWriter != nil {
		cmd.Stdout = n.logWriter
		cmd.Stderr = n.logWriter
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
	}

	n.Error = cmd.Run()
	return n.Error
}

func (n *Node) clearState() {
	n.NodeState = NodeState{}
}

func (n *Node) ReadStatus() NodeStatus {
	n.mu.RLock()
	defer n.mu.RUnlock()
	ret := n.Status
	return ret
}

func (n *Node) updateStatus(status NodeStatus) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Status = status
}

func (n *Node) signal(sig os.Signal) {
	status := n.ReadStatus()
	if status == NodeStatusRunning {
		n.updateStatus(NodeStatusCancel)
	}
	if n.cmd != nil {
		n.cmd.Process.Signal(sig)
	}
}

func (n *Node) cancel() {
	n.mu.Lock()
	defer n.mu.Unlock()
	status := n.Status
	if status == NodeStatusNone || status == NodeStatusRunning {
		n.Status = NodeStatusCancel
	}
	if n.cancelFunc != nil {
		n.cancelFunc()
	}
}

func (n *Node) setupLog(logDir string) {
	n.StartedAt = time.Now()
	n.Log = filepath.Join(logDir, fmt.Sprintf("%s.%s.log",
		utils.ValidFilename(n.Name, "_"),
		n.StartedAt.Format("20060102.15:04:05"),
	))
}

func (n *Node) openLogFile() error {
	if n.Log == "" {
		return nil
	}
	var err error
	n.logFile, err = utils.OpenOrCreateFile(n.Log)
	if err != nil {
		n.Error = err
		return err
	}
	n.logWriter = bufio.NewWriter(n.logFile)
	return nil
}

func (n *Node) closeLogFile() error {
	var lastErr error = nil
	if n.logWriter != nil {
		lastErr = n.logWriter.Flush()
	}
	if n.logFile != nil {
		if err := n.logFile.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (n *Node) ReadRetryCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.RetryCount
}

func (n *Node) ReadDoneCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.DoneCount
}

func (n *Node) incRetryCount() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.RetryCount++
}

func (n *Node) incDoneCount() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.DoneCount++
}

var nextNodeId = 1

func (n *Node) init() {
	if n.id != 0 {
		return
	}
	n.id = nextNodeId
	nextNodeId++
}
