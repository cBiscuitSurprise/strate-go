# Strate.Go!

This project implements the Stratego game. We'll have two interfaces to this game. A CLI and a web interface that can be deployed via AWS and integrated into https://casey.boyer.consulting.


## Game

### Rules

This is a two player board game (red vs blue). Each player starts with 40 peices each on a 10x10 board.

Each player can setup their pieces however they wish on their side of the board.

Red starts by moving one of their movable pieces to an open spot on the board.

### Peices


## Project

This project is comprised of a cli which can be used to play the game on the terminal or serve as a websocket application for playing via the web.

### Setup

The generated protobuf interface isn't tracked as part of this repository. Instead it needs to be generated if you're working on this project.

```
./scripts/build_protobuf.sh
```

### Terminal

```
strate-go play
```

### Websockets

```
strate-go serve
```
