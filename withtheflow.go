package withtheflow

var flowIdCounter int64

// type FlowId int64

// type FlowConfig struct {
// 	Id             FlowId
// 	HandlerName    string
// 	Arguments      interface{}
// 	DependentFlows []FlowId
// }

type FlowHandler func(args interface{}, runtime Runtime, subFlowResults []interface{}) interface{}

type Runtime interface {
	// Schedules a new flow which combines the results of the given flow ids.
	AddFlow(funcname string, args interface{}, dependentFlows ...int64) int64
	// Sets the result of this flow to be the result of the given flow id (once it has finished executing)
	DeferredResult(int64) interface{}
}

type WorkflowRunner interface {
	Run(funcname string, args interface{}) interface{}
}

// func CreateFlow(funcname string, args interface{}) *FlowConfig {
// 	return &FlowConfig{
// 		Id:          FlowId(atomic.AddInt64(&flowIdCounter, 1)),
// 		HandlerName: funcname,
// 		Arguments:   args,
// 	}
// }

// func CombineFlows(funcname string, args interface{}, dependentFlows ...FlowId) *FlowConfig {
// 	flowConfig := CreateFlow(funcname, args)
// 	flowConfig.DependentFlows = dependentFlows

// 	return flowConfig
// }

// // Takes a serialised proto and returns a serialised proto
// type FlowHandler func([]byte, FlowHandle) ([]byte, error)

// // Takes a serialised proto and returns a serialised proto
// type DependentFlowHandler func([]FlowResult, []byte, FlowHandle) ([]byte, error)

// // Takes a slice of serialised flow results,
// type FlowReducer func([]FlowResult, []byte) ([]byte, error)

// type Workflow interface {
// 	NewFlow(FlowCall) (FlowId, error)
// 	Run() error
// 	GenerateDotGraph() (string, error)
// }

// type FlowHandle interface {
// 	// Schedules the given FlowCall for execution
// 	NewFlow(FlowCall) (FlowId, error)

// 	// Schedules a FlowCall which will be called with an array of the results of the given dependent flows
// 	// and the return value of the current flow will be set to that of the given flow reducer
// 	NewBlockedResult(FlowReducerCall) (FlowId, error)
// }
