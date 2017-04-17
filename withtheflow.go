package withtheflow

type FlowId int64

type FlowCall struct {
	FlowName  string
	Arguments []byte
}

type FlowReducerCall struct {
	FlowReducerName string
	Arguments       []byte
	DependentFlows  []FlowId
}

type FlowResult struct {
	FlowName string
	Result   []byte
}

// Takes a serialised proto and returns a serialised proto
type FlowHandler func([]byte, FlowHandle) ([]byte, error)

// Takes a serialised proto and returns a serialised proto
type DependentFlowHandler func([]FlowResult, []byte, FlowHandle) ([]byte, error)

// Takes a slice of serialised flow results,
type FlowReducer func([]FlowResult, []byte) ([]byte, error)

type Workflow interface {
	NewFlow(FlowCall) (FlowId, error)
	Run() error
	GenerateDotGraph() (string, error)
}

type FlowHandle interface {
	// Schedules the given FlowCall for execution
	NewFlow(FlowCall) (FlowId, error)

	// Schedules a FlowCall which will be called with an array of the results of the given dependent flows
	// and the return value of the current flow will be set to that of the given flow reducer
	NewBlockedResult(FlowReducerCall) (FlowId, error)
}
