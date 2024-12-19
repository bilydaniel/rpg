package assets

import (
	"encoding/json"
	"os"
)

type Tilemap struct {
	Height         int            `json:"height"`
	Width          int            `json:"width"`
	Layers         []TilemapLayer `json:"layers"`
	TilesetName    string
	TilesetRowSize int
	Tilesets       []Tileset
}

type Tileset struct {
	Firstgid int    `json:"firstgid"`
	Source   string `json:"source"`
}
type TilesetData struct {
	Image   string     `json:"image"`
	Columns int        `json:"columns"`
	Name    string     `json:"name"`
	Tiles   []TileData `json:"tiles"`
}
type TileData struct {
	ID          int    `json:"id"`
	Image       string `json:"image"`
	ImageHeight int    `json:"imageheight"`
	ImageWidth  int    `json:"imagewidth"`
}

func InitTilemap() Tilemap {
	return Tilemap{}
}

func LoadTilemap(id string) (*Tilemap, error) {
	tilemap := Tilemap{}

	mapPath := "assets/maps/" + id + ".tmj"
	jsonmap, err := os.ReadFile(mapPath)
	if err != nil {
		return nil, err
	}

	//TODO use gob instead of json?
	err = json.Unmarshal(jsonmap, &tilemap)
	if err != nil {
		return nil, err
	}
	return &tilemap, nil
}

type TilemapLayer struct {
	ID      int      `json:"id"` //TODO dont know if needed
	Name    string   `json:"name"`
	Opacity float64  `json:"opacity"` //TODO not using right now, could be interesting
	Type    string   `json:"type"`
	Visible bool     `json:"visible"`
	Data    []int    `json:"data"`
	Objects []Object `json:"objects"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
}

type Object struct {
	GID      int     `json:"gid"`
	ID       int     `json:"id"`
	X        int     `json:"x"`
	Y        int     `json:"y"`
	Width    int     `json:"Width"`
	Height   int     `json:"Height"`
	Rotation float64 `json:"rotation"`
	Visible  bool    `json:"visible"`
}

/*
func (t *Tilemap) LoadTestMap(assetname string, assets assets.Assets) error {
	//TODO only load up the map json file and depending on what that file uses, load the rest
	//TODO dont hard code
	t.Assets = assets
	t.TilesetName = assetname
	t.TilesetRowSize = 22

	jsonmap, err := os.ReadFile("assets/maps/test_map.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonmap, t)
	if err != nil {
		return err
	}

	return nil
}
*/

/*
func (t *Tilemap) GetTile(id int, assets assets.Assets) *ebiten.Image {
	tileX := ((id - 1) % t.TilesetRowSize) * 16
	tileY := ((id - 1) / t.TilesetRowSize) * 16

	tileset, ok := t.Assets.Tileset[t.TilesetName]
	if !ok {
		return nil
	}
	return tileset.SubImage(image.Rect(tileX, tileY, tileX+config.TileSize, tileY+config.TileSize)).(*ebiten.Image)
}

*/
/*
func (t *Tilemap) Draw(screen *ebiten.Image, camera config.Camera, assets assets.Assets) {
	opts := ebiten.DrawImageOptions{}
	for _, layer := range t.Layers {
		for idx, id := range layer.Data {
			//TODO what can I cashe??
			x := idx % layer.Width
			y := idx / layer.Width

			x *= config.TileSize
			y *= config.TileSize
			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(-camera.X, -camera.Y)
			opts.GeoM.Scale(camera.Scale, camera.Scale)
			tile := t.GetTile(id, assets)
			screen.DrawImage(tile, &opts)

			opts.GeoM.Reset()
		}
	}
}
*/
