# Tiny World

A tiny, slow-paced world and colony building game.

Made with [Arche](https://github.com/mlange-42/arche) and [Ebitengine](https://github.com/hajimehoshi/ebiten).
Very early work in progress!

![Tiny World screenshot](https://github.com/mlange-42/tiny-world/assets/44003176/5a495808-0f7d-4669-b8e4-58e9af563ff7)

## Usage

Currently, you need to clone the repository and run the game with [Go](https://go.dev):

```shell
git clone https://github.com/mlange-42/tiny-world.git
cd tiny-world
go run .
```

## Controls

In the toolbar on the right, the top items are buildings that can be bought by the player for resources.
The natural features in the lower part appear randomly and are replenished when placed by the player.

* Middle mouse button / mouse wheel: pan and zoom.
* Left click with selected terrain or buildable: place it.
* Right click with selected buildable: remove it.
* Ctrl+S: saves the game to `save/autosave.json`

Load a saved game by running with the file as argument:

```shell
go run . save/autosave.json
```
