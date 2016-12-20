package consul

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/laohanlinux/utils/coordinator"
)

const (
	ErrTTLFormat         = `Unexpected response code: 500 (CheckID "%s" does not have associated TTL)`
	ErrSession           = `Unexpected response code: 500 (rpc error: rpc error: invalid session "%s")`
	timestampMaxDelay    = int64(10 * time.Second)
	DefaultIntervalCheck = 3
)

var errCheckID = fmt.Errorf("%s", "the check id is invalid")

func NewAgentServiceRegisterOption(serviceid, name, ip string, port int, tags []string,
	override bool, check *api.AgentServiceCheck) *api.AgentServiceRegistration {
	return &api.AgentServiceRegistration{
		ID:                serviceid,
		Name:              name,
		Address:           ip,
		Port:              port,
		EnableTagOverride: override,
		Check:             check,
		Tags:              tags,
	}
}

func RegisterService(coor *coordinator.CtxGCoordinator, client *api.Client, serverid string,
	registration *api.AgentServiceRegistration, logger log.Logger) {
	cc := consul.NewClient(client)
	registrar := consul.NewRegistrar(cc, registration, logger)
	checkid := fmt.Sprintf("service:%s", serverid)
	registrar.Register()
	go CheckService(coor, registrar, client.Agent(), DefaultIntervalCheck, serverid, checkid, true, logger)
}

func CheckService(coor *coordinator.CtxGCoordinator, reg *consul.Registrar, agent *api.Agent,
	interval int, serverid, checkid string, regAgain bool, logger log.Logger) {
	t := time.NewTicker(time.Second * time.Duration(interval))
	defer func() {
		reg.Deregister()
		t.Stop()
		coor.Sub()
	}()
	for {
		select {
		case <-coor.Done():
			return
		case <-t.C:
			if err := agent.UpdateTTL(checkid, time.Now().String(), "pass"); err != nil {
				if errCheckID == CheckUpdateTTLError(fmt.Sprintf("service:%s",
					serverid), err) {
					if regAgain {
						reg.Register()
					}
				}
				logger.Log("err", err)
			}
		}
	}
}

func CheckUpdateTTLError(serviceID string, err error) error {
	if fmt.Sprintf(ErrTTLFormat, serviceID) == err.Error() {
		return errCheckID
	}
	return err
}
