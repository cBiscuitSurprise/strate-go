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

Once the build completes, we can link the generated proto files back into the project so intellisense works (optional)

```
ln -s $(readlink -e bazel-bin)/api/proto/strategopb_go_proto_/github.com/cBiscuitSurprise/strate-go/api/go/ api/go
```

### Build

```bash
bazel build //:strate-go
```

### Terminal

```bash
bazel run //:strate-go play
```

### Websockets

```bash
bazel run //:strate-go serve
```
