package sys

import (
	ares "github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
)

// DoProduction system.
type DoProduction struct {
	time   generic.Resource[ares.Tick]
	update generic.Resource[res.UpdateInterval]
	stock  generic.Resource[res.Stock]

	filter        generic.Filter3[comp.Tile, comp.UpdateTick, comp.Production]
	markerBuilder generic.Map2[comp.Tile, comp.ProductionMarker]
	toCreate      []markerEntry
}

// Initialize the system
func (s *DoProduction) Initialize(world *ecs.World) {
	s.time = generic.NewResource[ares.Tick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.stock = generic.NewResource[res.Stock](world)

	s.filter = *generic.NewFilter3[comp.Tile, comp.UpdateTick, comp.Production]()
	s.markerBuilder = generic.NewMap2[comp.Tile, comp.ProductionMarker](world)
}

// Update the system
func (s *DoProduction) Update(world *ecs.World) {
	stock := s.stock.Get()
	tick := s.time.Get().Tick
	update := s.update.Get()
	tickMod := tick % update.Interval

	query := s.filter.Query(world)
	for query.Next() {
		tile, up, pr := query.Get()

		if up.Tick != tickMod {
			continue
		}
		pr.Countdown -= pr.Amount
		if pr.Countdown < 0 {
			pr.Countdown += update.Countdown
			stock.Res[pr.Type]++
			s.toCreate = append(s.toCreate, markerEntry{Tile: *tile, Resource: pr.Type})
		}
	}

	for _, entry := range s.toCreate {
		s.markerBuilder.NewWith(
			&entry.Tile,
			&comp.ProductionMarker{StartTick: tick, Resource: entry.Resource},
		)
	}
	s.toCreate = s.toCreate[:0]
}

// Finalize the system
func (s *DoProduction) Finalize(world *ecs.World) {}

type markerEntry struct {
	Tile     comp.Tile
	Resource resource.Resource
}
