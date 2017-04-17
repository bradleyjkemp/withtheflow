package serial

import (
	"github.com/bradleyjkemp/withtheflow"
	"github.com/bradleyjkemp/withtheflow/serial/proto/flow"
	"github.com/golang/protobuf/ptypes"
	"sync/atomic"
	"time"
)

type workflow struct {
	flowCounter           int64
	workflowState         *flow.WorkflowState
	flowConfigs           map[int64]*flow.FlowConfig
	flowHandlers          map[string]withtheflow.FlowHandler
	flowReducers          map[string]withtheflow.FlowReducer
	dependentFlowHandlers map[string]withtheflow.DependentFlowHandler
}

func NewWorkflow(flowHandlers map[string]withtheflow.FlowHandler, flowReducers map[string]withtheflow.FlowReducer) withtheflow.Workflow {
	return &workflow{
		workflowState: &flow.WorkflowState{},
		flowHandlers:  flowHandlers,
		flowReducers:  flowReducers,
		flowConfigs:   map[int64]*flow.FlowConfig{},
	}
}

func (w *workflow) NewFlow(flowCall withtheflow.FlowCall) (withtheflow.FlowId, error) {
	timestamp, _ := ptypes.TimestampProto(time.Now())

	flowConfig := &flow.FlowConfig{
		Name:         flowCall.FlowName,
		Data:         flowCall.Arguments,
		CreationTime: timestamp,
	}
	return w.newFlowWithParent(flowConfig, 0)
}

func (w *workflow) newFlowWithParent(newFlow *flow.FlowConfig, parentId int64) (withtheflow.FlowId, error) {
	timestamp, _ := ptypes.TimestampProto(time.Now())
	newFlow.CreationTime = timestamp
	newFlow.Id = atomic.AddInt64(&w.flowCounter, 1)
	newFlow.ParentId = parentId

	w.flowConfigs[newFlow.Id] = newFlow

	w.workflowState.Scheduled = append(w.workflowState.Scheduled, newFlow)

	return (withtheflow.FlowId(newFlow.Id)), nil
}

func (w *workflow) Run() error {
	for {
		if len(w.workflowState.Scheduled) == 0 {
			break
		}

		flowToRun := w.workflowState.Scheduled[0]
		w.workflowState.Scheduled = w.workflowState.Scheduled[1:]

		w.runHandler(flowToRun)
	}

	return nil
}

func (w *workflow) runHandler(flowToRun *flow.FlowConfig) error {
	var result []byte
	var err error
	if len(flowToRun.DependentIds) > 0 {
		if w.isBlocked(flowToRun) {
			w.workflowState.Scheduled = append(w.workflowState.Scheduled, flowToRun)
			return nil
		}

		results := []withtheflow.FlowResult{}
		for _, id := range flowToRun.DependentIds {
			results = append(results, withtheflow.FlowResult{
				FlowName: w.flowConfigs[id].Name,
				Result:   w.flowConfigs[id].Result,
			})
		}

		result, err = w.flowReducers[flowToRun.Name](results, flowToRun.Data)
	} else {
		result, err = w.flowHandlers[flowToRun.Name](flowToRun.Data, &flowHandle{w, flowToRun})
	}

	if err != nil {
		return err
	}

	flowToRun.Result = result
	flowToRun.CompletionTime, _ = ptypes.TimestampProto(time.Now())

	w.workflowState.Completed = append(w.workflowState.Completed, flowToRun)

	return nil
}

func (w *workflow) isBlocked(flow *flow.FlowConfig) bool {
	if w.checkResultDelayed(flow) {
		return true
	} else {
		for _, id := range flow.DependentIds {
			if w.isBlocked(w.flowConfigs[id]) {
				return true
			}
		}
	}

	return false
}

// updates status of a delayed result and returns bool whether it is still delayed
func (w *workflow) checkResultDelayed(flow *flow.FlowConfig) bool {
	if !flow.DelayedResult {
		return false
	}

	// when a result is delayed DependentIds contains the Id of the FlowReducer
	blockedFlow := w.flowConfigs[flow.DependentIds[0]]

	if blockedFlow.CompletionTime == nil {
		// flow reducer still hasn't run
		return true
	} else {
		flow.Result = blockedFlow.Result
		return false
	}
}
