package netrpc

const (
	HealthCheckService    = "HealthCheck"
	HealthCheckPingNetRPC = "HealthCheck.Ping"
)

type HealthCheck struct{}

func (hc *HealthCheck) Ping(req EmptyRequest, reply *EmptyReply) error {
	return nil
}
