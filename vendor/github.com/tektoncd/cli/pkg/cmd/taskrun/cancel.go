// Copyright © 2019 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package taskrun

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	"github.com/tektoncd/cli/pkg/formatted"
	"github.com/tektoncd/cli/pkg/taskrun"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func cancelCommand(p cli.Params) *cobra.Command {
	eg := `Cancel the TaskRun named 'foo' from namespace 'bar':

    tkn taskrun cancel foo -n bar
`

	c := &cobra.Command{
		Use:               "cancel",
		Short:             "Cancel a TaskRun in a namespace",
		Example:           eg,
		ValidArgsFunction: formatted.ParentCompletion,
		Args:              cobra.ExactArgs(1),
		SilenceUsage:      true,
		Annotations: map[string]string{
			"commandType": "main",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := &cli.Stream{
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}

			return cancelTaskRun(p, s, args[0])
		},
	}

	return c
}

func cancelTaskRun(p cli.Params, s *cli.Stream, trName string) error {
	cs, err := p.Clients()
	if err != nil {
		return fmt.Errorf("failed to create tekton client")
	}

	tr, err := taskrun.Get(cs, trName, metav1.GetOptions{}, p.Namespace())
	if err != nil {
		return fmt.Errorf("failed to find TaskRun: %s", trName)
	}

	if len(tr.Status.Conditions) > 0 {
		if tr.Status.Conditions[0].Status != corev1.ConditionUnknown {
			return fmt.Errorf("failed to cancel TaskRun %s: TaskRun has already finished execution", trName)
		}
	}

	if _, err := taskrun.Patch(cs, trName, metav1.PatchOptions{}, p.Namespace()); err != nil {
		return fmt.Errorf("failed to cancel TaskRun %s: %v", trName, err)
	}

	fmt.Fprintf(s.Out, "TaskRun cancelled: %s\n", tr.Name)
	return nil
}
