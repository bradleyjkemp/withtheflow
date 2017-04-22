package serial

import (
	"bytes"
	"errors"
	"fmt"
)

func (w *workflow) GenerateDotGraph() (string, error) {
	if len(w.workflowState.Scheduled) > 0 {
		return "", errors.New("Workflow not complete")
	}

	var buffer bytes.Buffer

	buffer.WriteString("digraph workflow {\n")
	buffer.WriteString("  id0 [label=\"start\"]\n")

	for _, flowConfig := range w.workflowState.Completed {
		buffer.WriteString(fmt.Sprintf("  id%d [label=\"%s\"]\n", flowConfig.Id, flowConfig.Name))
	}

	for _, flowConfig := range w.workflowState.Completed {
		for _, id := range flowConfig.DependentIds {
			buffer.WriteString(fmt.Sprintf("  id%d->id%d\n", id, flowConfig.Id))
		}
		buffer.WriteString(fmt.Sprintf("  id%d->id%d\n", flowConfig.ParentId, flowConfig.Id))
	}

	buffer.WriteString("}\n")

	return buffer.String(), nil
}
