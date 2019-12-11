package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Point struct {
	x, y int
}

var nilPoint = Point{-1, -1}

func (this Point) Add(other Point) Point {
	return Point{
		x: this.x + other.x,
		y: this.y + other.y,
	}
}

func (this Point) Sub(other Point) Point {
	return Point{
		x: this.x - other.x,
		y: this.y - other.y,
	}
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (p Point) Dir() (float64, bool) {
	dir := float64(p.y) / float64(p.x)

	if p.x >= 0 {
		return dir, false
	}
	return dir, true
}

func (this Point) LessAngle(other Point) bool {
	td, tl := this.Dir()
	od, ol := other.Dir()

	if tl != ol {
		return ol
	}
	return td < od
}

func detect(m map[Point]bool, pos Point) []Point {
	s := map[Point]bool{}
	cont := true

	try := func(p Point, v Point) (f Point) {
		f = nilPoint

		p = p.Add(v)
		if _, ok := s[p]; ok {
			return
		}

		a, ok := m[p]
		for ok {
			cont = true
			s[p] = true
			if a && (f.x == -1 || f.y == -1) {
				f = p
			}

			p = p.Add(v)
			a, ok = m[p]
		}

		return
	}

	detects := []Point{}

	for off := 1; cont; off++ {
		cont = false

		for x := -off; x <= off; x++ {
			if p := try(pos, Point{x, off}); p != nilPoint {
				detects = append(detects, p)
			}
			if p := try(pos, Point{x, -off}); p != nilPoint {
				detects = append(detects, p)
			}
		}

		for y := 1 - off; y <= off-1; y++ {
			if p := try(pos, Point{off, y}); p != nilPoint {
				detects = append(detects, p)
			}
			if p := try(pos, Point{-off, y}); p != nilPoint {
				detects = append(detects, p)
			}
		}
	}

	return detects
}

func print(m map[Point]bool) {
	max := Point{}
	for p := range m {
		if p.x > max.x {
			max.x = p.x
		}
		if p.y > max.y {
			max.y = p.y
		}
	}

	for y := 0; y <= max.y; y++ {
		for x := 0; x <= max.x; x++ {
			if m[Point{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func ret(m map[Point]bool) {
	max := 0
	for p := range m {
		if p.y > max {
			max = p.y
		}
	}

	fmt.Printf("\033[%dA", max+1)
}

func main() {
	data := `
		..#..###....#####....###........#
		.##.##...#.#.......#......##....#
		#..#..##.#..###...##....#......##
		..####...#..##...####.#.......#.#
		...#.#.....##...#.####.#.###.#..#
		#..#..##.#.#.####.#.###.#.##.....
		#.##...##.....##.#......#.....##.
		.#..##.##.#..#....#...#...#...##.
		.#..#.....###.#..##.###.##.......
		.##...#..#####.#.#......####.....
		..##.#.#.#.###..#...#.#..##.#....
		.....#....#....##.####....#......
		.#..##.#.........#..#......###..#
		#.##....#.#..#.#....#.###...#....
		.##...##..#.#.#...###..#.#.#..###
		.#..##..##...##...#.#.#...#..#.#.
		.#..#..##.##...###.##.#......#...
		...#.....###.....#....#..#....#..
		.#...###..#......#.##.#...#.####.
		....#.##...##.#...#........#.#...
		..#.##....#..#.......##.##.....#.
		.#.#....###.#.#.#.#.#............
		#....####.##....#..###.##.#.#..#.
		......##....#.#.#...#...#..#.....
		...#.#..####.##.#.........###..##
		.......#....#.##.......#.#.###...
		...#..#.#.........#...###......#.
		.#.##.#.#.#.#........#.#.##..#...
		.......#.##.#...........#..#.#...
		.####....##..#..##.#.##.##..##...
		.#.#..###.#..#...#....#.###.#..#.
		............#...#...#.......#.#..
		.........###.#.....#..##..#.##...
	`
	// data := `
	// 	.#..##.###...#######
	// 	##.############..##.
	// 	.#.######.########.#
	// 	.###.#######.####.#.
	// 	#####.##.#.##.###.##
	// 	..#####..#.#########
	// 	####################
	// 	#.####....###.#.#.##
	// 	##.#################
	// 	#####.##.###..####..
	// 	..######..##.#######
	// 	####.##.####...##..#
	// 	.#####..#.######.###
	// 	##...#.##########...
	// 	#.##########.#######
	// 	.####.#.###.###.#.##
	// 	....##.##.###..#####
	// 	.#.#.###########.###
	// 	#.#.#.#####.####.###
	// 	###.##.####.##.#..##
	// `

	m := map[Point]bool{}

	for y, row := range strings.Split(strings.TrimSpace(data), "\n") {
		for x, c := range strings.Split(strings.TrimSpace(row), "") {
			m[Point{x, y}] = c == "#"
		}
	}

	max := 0
	at := Point{}

	for pos, asteroid := range m {
		if !asteroid {
			continue
		}

		if count := len(detect(m, pos)); count > max {
			max = count
			at = pos
		}
	}

	fmt.Printf("%d at %v\n", max, at)

	all := []Point{}

	delay := time.Millisecond * 20

	print(m)
	time.Sleep(delay)

	for {
		d := detect(m, at)
		if len(d) == 0 {
			break
		}

		sort.Slice(d, func(i, j int) bool {
			return d[i].Sub(at).LessAngle(d[j].Sub(at))
		})

		all = append(all, d...)

		for _, p := range d {
			m[p] = false
			ret(m)
			print(m)
			time.Sleep(delay)
		}
	}

	if len(all) >= 200 {
		fmt.Println("200:", all[199])
	}
}
