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

```bash
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

### Running `go`

```bash
bazel run @io_bazel_rules_go//go -- args
bazel run @io_bazel_rules_go//go -- get github.com/spf13/viper
```


## Shorthand notation

The historical moves of the game are captured using a contrived shorthand notation:

`<id> <from:R,C> <to:R,C> <result>`

**Examples**
```
Red:001 - 3,0 1     // piece Red:001 from off board to Row 3 Column 0
Red:001 3,0 3,1 2   // piece Red:001 from Row 3 Column 0 to Row 3 Column 1 taking opponent
Red:001 3,0 3,1 3   // piece Red:001 from Row 3 Column 0 to Row 3 Column 1 losing my piece
Red:001 3,0 3,1 4   // piece Red:001 from Row 3 Column 0 to Row 3 Column 1 losing both pieces
Red:001 3,0 3,1 5   // piece Red:001 from Row 3 Column 0 to Row 3 Column 1 taking flag
```

Where `result` is enumerated like:
* 0: No contest
* 1: Attackee captured
* 2: Attacker captured
* 3: Both captured
* 4: Flag captured
