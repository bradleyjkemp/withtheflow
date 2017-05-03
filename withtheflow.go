package withtheflow

type FlowId interface{}

type FlowHandler func(args interface{}, runtime Runtime, subFlowResults []interface{}) interface{}

type Runtime interface {
	// Schedules a new flow which combines the results of the given flow ids.
	AddFlow(funcname string, args interface{}, dependentFlows ...FlowId) FlowId
	// Sets the result of this flow to be the result of the given flow id (once it has finished executing)
	// The return of DeferredResult should be returned by the flow handler
	DeferredResult(FlowId) interface{}
}

type WorkflowRunner interface {
	Run(funcname string, args interface{}) interface{}
}
