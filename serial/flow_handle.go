package serial

import (
	"github.com/bradleyjkemp/withtheflow"
	"github.com/bradleyjkemp/withtheflow/serial/proto/flow"
)

type flowHandle struct {
	*workflow
	currentFlow *flow.FlowConfig
}

func (f *flowHandle) NewFlow(flowCall withtheflow.FlowCall) (withtheflow.FlowId, error) {
	flowConfig := &flow.FlowConfig{
		Name: flowCall.FlowName,
		Data: flowCall.Arguments,
	}
	return f.newFlowWithParent(flowConfig, f.currentFlow.Id)
}

func (f *flowHandle) NewBlockedResult(flowReducer withtheflow.FlowReducerCall) (withtheflow.FlowId, error) {

	flowConfig := &flow.FlowConfig{
		Name: flowReducer.FlowReducerName,
		Data: flowReducer.Arguments,
	}

	for _, id := range flowReducer.DependentFlows {
		flowConfig.DependentIds = append(flowConfig.DependentIds, int64(id))
	}

	blockingId, err := f.newFlowWithParent(flowConfig, f.currentFlow.Id)
	if err != nil {
		return 0, err
	}

	f.currentFlow.DelayedResult = true
	f.currentFlow.DependentIds = []int64{int64(blockingId)}

	return blockingId, err
}
