# Pi3d

Creates an 3d model of the digits of pi where the height of each block corresponds to each digit.

The grid can be resized up to 10k blocks with the default pi digits file but it's a bit slow.

![Example Output](snapshot00.png)

### How to run it

Run it with `cat inputs/pi_10000.txt | go run main.go` (or build it if you want).

There is also an option to convert a file of words into a grid of blocks where each letter has a height of it's position the alpabet using the
flag e.g. `echo "foo bar" | go run main.go -input-type=text`.

There are some constants in `main.go` that can be used to alter the output in various ways.