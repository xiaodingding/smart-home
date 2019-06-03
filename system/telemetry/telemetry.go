package telemetry

import (
	"github.com/e154/smart-home/system/telemetry/dashboard"
	"github.com/e154/smart-home/system/telemetry/telemetry_map"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/adaptors"
)

type Telemetry struct {
	dashboard *dashboard.Dashboard
	Map       *telemetry_map.Map
}

func NewTelemetry(dashboard *dashboard.Dashboard,
	stream *stream.StreamService,
	adaptors *adaptors.Adaptors) *Telemetry {

	telemetry := &Telemetry{
		dashboard: dashboard,
		Map:       telemetry_map.NewMap(stream, adaptors),
	}
	telemetry.Run()

	return telemetry
}

func (t *Telemetry) Run() {

	t.dashboard.Run()
	t.Map.Run()
}

func (t *Telemetry) Stop() {

	t.dashboard.Stop()
	t.Map.Stop()
}

func (t *Telemetry) Broadcast(pack string) {

	t.dashboard.Broadcast(pack)
	t.Map.Broadcast(pack)
}

func (t *Telemetry) BroadcastOne(pack string, id int64) {

	t.dashboard.BroadcastOne(pack, id)
	t.Map.BroadcastOne(pack, id)
}
