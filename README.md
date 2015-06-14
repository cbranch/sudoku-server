# sudoku-server
HTTP service for generating and solving sudoku puzzles, using libsudoku

## Dependencies
Requires [libsudoku](https://github.com/cbranch/libsudoku) compiled & on your path.

## GET endpoints

`/generate[?difficulty=1-99]` - returns a sudoku puzzle with a particular difficulty. The larger the number, the harder the puzzle.

## POST endpoints

`/solve` - solves a puzzle provided in the request body.
