package modbus

import (
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/inputs/system"

	mb "github.com/mythay/lark/collector/modbus"
)

type ModbusMetric struct {
	ps     system.PS
	agent  *mb.Agent
	Config string `toml:"config"`
	Log    telegraf.Logger
}

func (_ *ModbusMetric) Description() string {
	return "Read metrics from modbus definition"
}

const sampleConfig = `
  ## yaml config that define the modbus registers
  config = "modbus.yaml"

`

func (_ *ModbusMetric) SampleConfig() string { return sampleConfig }

func (s *ModbusMetric) init() error {
	if s.agent == nil {
		agent, err := mb.NewAgent()
		if err != nil {
			return err
		}
		s.agent = agent
		return s.agent.LoadConfig(s.Config)
	}
	return nil
}

// func (s *ModbusMetric) Start(acc tel.Accumulator) error {
// 	err := s.init()
// 	if err != nil {
// 		return fmt.Errorf("error init modbus: %s", err)
// 	}

// 	return s.agent.Start(acc.AddCounter)
// }

func (s *ModbusMetric) Stop() {
	s.agent.Stop()
}

func (s *ModbusMetric) Gather(acc telegraf.Accumulator) error {

	err := s.init()
	if err != nil {
		return fmt.Errorf("error init modbus: %s", err)
	}

	// s.agent.Gather(func(measurement string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {
	// 	acc.AddGauge(measurement, fields, tags, t...)
	// })

	s.agent.Gather(acc.AddGauge)

	return nil
}

func init() {
	ps := system.NewSystemPS()
	inputs.Add("modbus", func() telegraf.Input {
		return &ModbusMetric{ps: ps}
	})
}
