package render

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Terrain is a system to render the terrain.
type Terrain struct {
	screen  generic.Resource[game.EbitenImage]
	sprites generic.Resource[game.Sprites]
	terrain generic.Resource[game.Terrain]
	landUse generic.Resource[game.LandUse]
	view    generic.Resource[game.View]
}

// InitializeUI the system
func (s *Terrain) InitializeUI(world *ecs.World) {
	s.screen = generic.NewResource[game.EbitenImage](world)
	s.sprites = generic.NewResource[game.Sprites](world)
	s.terrain = generic.NewResource[game.Terrain](world)
	s.landUse = generic.NewResource[game.LandUse](world)
	s.view = generic.NewResource[game.View](world)
}

// UpdateUI the system
func (s *Terrain) UpdateUI(world *ecs.World) {
	terrain := s.terrain.Get()
	landUse := s.landUse.Get()
	sprites := s.sprites.Get()
	view := s.view.Get()

	canvas := s.screen.Get()
	img := canvas.Image

	off := view.Offset()
	bounds := view.Bounds(canvas.Width, canvas.Height)

	img.Clear()

	op := ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver

	halfWidth := view.TileWidth / 2

	drawSprite := func(grid *game.Grid[terr.Terrain], x, y int, t terr.Terrain, point *image.Point, height int) int {
		idx := sprites.GetTerrainIndex(t)
		sp, info := sprites.Get(idx)
		h := sp.Bounds().Dy() - view.TileHeight

		if info.MultiTile {
			neigh := grid.NeighborsMask(x, y, t)
			idx = sprites.GetMultiTileIndex(t, neigh)
			sp, _ = sprites.Get(idx)
		}

		op.GeoM.Reset()
		op.GeoM.Scale(view.Zoom, view.Zoom)
		op.GeoM.Translate(
			float64(point.X-halfWidth)*view.Zoom-float64(off.X),
			float64(point.Y-h-height-info.YOffset)*view.Zoom-float64(off.Y),
		)
		img.DrawImage(sp, &op)

		return height + info.Height
	}

	mx, my := view.ScreenToGlobal(ebiten.CursorPosition())
	cursor := view.GlobalToTile(mx, my)

	for i := 0; i < terrain.Width(); i++ {
		for j := 0; j < terrain.Height(); j++ {
			point := view.TileToGlobal(i, j)
			if !point.In(bounds) {
				continue
			}

			height := 0
			t := terrain.Get(i, j)
			if t != terr.Air {
				height = drawSprite(&terrain.Grid, i, j, t, &point, height)
			}

			lu := landUse.Get(i, j)
			if lu != terr.Air {
				_ = drawSprite(&landUse.Grid, i, j, lu, &point, height)
			}

			if cursor.X == i && cursor.Y == j {
				_ = drawSprite(nil, i, j, terr.Cursor, &point, 0)
			}
		}
	}
}

// PostUpdateUI the system
func (s *Terrain) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *Terrain) FinalizeUI(world *ecs.World) {}