package main

import (
	"image"
	_ "image/png"
	"math"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/salviati/go-tmx/tmx"
)

//GameMap holds all the information on our level
//Note also that the tileset must be an embedded tileset
type GameMap struct {
	LevelMap  *tmx.Map
	TileSet   *TileSetInfo //The parser seems to only support a single tileset per map.
	Colliders []image.Rectangle
}

//TileSetInfo is the information necessary to pull the tileset for the map
type TileSetInfo struct {
	TileHeight   int
	TileWidth    int
	Columns      int
	TileCount    int
	TileSetImage *ebiten.Image
}

//LoadGameMap loads the map from a tmx file
func LoadGameMap(pathToMap string, mapName string, tilesetName string) *GameMap {

	g := GameMap{}

	mapPath := path.Join(pathToMap, mapName)
	myMap, err := tmx.ReadFile(mapPath)
	if err != nil {
		panic(err)

	}
	tilesetPath := path.Join(pathToMap, tilesetName)
	g.LevelMap = myMap
	g.TileSet = newTilesetInfo(myMap, tilesetPath)
	g.Colliders = getColliders(g)

	return &g
}

func newTilesetInfo(myMap *tmx.Map, tileset string) *TileSetInfo {
	tsi := TileSetInfo{}
	tsi.TileHeight = myMap.TileHeight
	tsi.TileWidth = myMap.TileWidth
	//tsi.Columns = myMap     .Columns
	//tsi.TileCount = myMap.Tilecount

	img, _, err := ebitenutil.NewImageFromFile(tileset)
	if err != nil {
		panic(err)
	}
	tsi.TileSetImage = img
	return &tsi
}

func getColliders(g GameMap) []image.Rectangle {
	ret := make([]image.Rectangle, 0)
	for _, layer := range g.LevelMap.Layers {
		if layer.Name != "Collision" {
			continue
		}
		for row := 0; row < g.LevelMap.Height; row++ {
			for column := 0; column < g.LevelMap.Width; column++ {
				if layer.DecodedTiles[row*g.LevelMap.Width+column].Tileset != nil {
					x, y := column*g.LevelMap.TileWidth, row*g.LevelMap.TileHeight
					x2, y2 := x+g.TileSet.TileWidth, y+g.TileSet.TileHeight
					r := image.Rect(x, y, x2, y2)
					ret = append(ret, r)

				}

			}

		}
	}
	return ret
}

//Draw draws the Game Map one layer at a time from 0 up
func (g *GameMap) Draw(screen *ebiten.Image) {
	for index := range g.LevelMap.Layers {
		drawLayer(g, index, screen)
	}
}

//Draw is called each draw cycle and is where we will blit.
func drawLayer(gm *GameMap, layerNumber int, screen *ebiten.Image) {
	floorLayer := gm.LevelMap.Layers[layerNumber]
	if floorLayer.Name == "Collision" {
		return
	}
	mapHeight, mapWidth := gm.LevelMap.Height, gm.LevelMap.Width

	//Optimize to only draw what's in the viewing window -- That will come later
	for row := 0; row < mapHeight; row++ {
		for column := 0; column < mapWidth; column++ {
			id := floorLayer.DecodedTiles[row*mapWidth+column].ID

			r := getTileRect(gm.TileSet.Columns, gm.LevelMap.TileWidth, gm.LevelMap.TileHeight, int(id))
			//Now work out WHERE to draw  it
			op := &ebiten.DrawImageOptions{}
			drawPosX := column * gm.LevelMap.TileWidth
			drawPosY := row * gm.LevelMap.TileHeight

			op.GeoM.Translate(float64(drawPosX), float64(drawPosY))
			if floorLayer.DecodedTiles[row*mapWidth+column].Tileset != nil {
				screen.DrawImage(gm.TileSet.TileSetImage.SubImage(r).(*ebiten.Image), op)
			}

		}

	}

}

func getTileRect(columns int, tileWidth int, tileHeight int, tileNumber int) image.Rectangle {

	row := (math.Floor(float64(tileNumber) / float64(columns)))
	cell := tileNumber - (columns * int(row))

	x1, y1 := (tileWidth * cell), int(row)*tileHeight
	x2, y2 := x1+tileWidth, y1+tileHeight

	return image.Rect(x1, y1, x2, y2)

}
