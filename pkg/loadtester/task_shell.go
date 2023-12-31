/*
Copyright 2020 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package loadtester

import (
	"context"
	"errors"
	"os/exec"
	"strconv"

	"go.uber.org/zap"
)

const TaskTypeShell = "cmd"

func init() {
	taskFactories.Store(TaskTypeShell, func(metadata map[string]string, canary string, logger *zap.SugaredLogger) (Task, error) {
		cmd, ok := metadata["cmd"]
		if !ok {
			return nil, errors.New("cmd not found in metadata")
		}
		logCmdOutput, _ := strconv.ParseBool(metadata["logCmdOutput"])
		return &CmdTask{TaskBase{canary, logger}, cmd, logCmdOutput}, nil
	})
}

type CmdTask struct {
	TaskBase
	command      string
	logCmdOutput bool
}

func (task *CmdTask) Hash() string {
	return hash(task.canary + task.command)
}

func (task *CmdTask) Run(ctx context.Context) *TaskRunResult {
	cmd := exec.CommandContext(ctx, "sh", "-c", task.command)
	out, err := cmd.CombinedOutput()

	if err != nil {
		task.logger.With("canary", task.canary).Errorf("command failed %s %v %s", task.command, err, out)
	} else {
		if task.logCmdOutput {
			task.logger.With("canary", task.canary).Info(string(out))
		}
		task.logger.With("canary", task.canary).Infof("command finished %s", task.command)
	}
	return &TaskRunResult{err == nil, out}
}

func (task *CmdTask) String() string {
	return task.command
}
