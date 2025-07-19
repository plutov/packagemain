This is the `minesweeper` game ported to the
[prototype/draw](https://github.com/gonutz/prototype/) library.

Run the game on Windows, Mac and Linux with:

	go run .

To run the game in your browser using WASM, use the `drawsm` tool:

	go install github.com/gonutz/prototype/cmd/drawsm@latest
	cd packagemain/minesweeper_gonutz
	drawsm run
