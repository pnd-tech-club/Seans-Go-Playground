package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	sc "strconv"
	"svi"
	"time"
)

var w, h, score, length int = 0, 0, 0, 1

type dir int

const (
	U dir = iota
	L
	R
	D
)

var tmpx string
var running, eaten bool = true, false
var sinceeat int

type sprite struct {
	r rune
	X int
	Y int
}

var target sprite

var snake = make([]sprite, 100)

/* has to check if player is on targets; repopulate targets */
func checkTarget() {
	if snake[0].X == target.X && snake[0].Y == target.Y {
		for i := len(snake) - 1; i > 0; i-- {
			if i-1 > 0 {
				snake[i] = snake[i-1]
			} else {
				snake[1].X = target.X
				snake[1].Y = target.Y
			}

		}

		//target = append(target[:i], target[i+1:]...)
		tmpx = sc.Itoa(snake[1].X) + sc.Itoa(snake[1].Y)
		length++
		score++
		eaten = true
		target.X, target.Y = svi.Random(0, 80), svi.Random(1, 24)

	} else {
		if eaten != true {
			eaten = false
		}
	}
}

/* moves the non-head ([0]) parts of the snake */
func shiftParts_old(d dir) {
	var num int = length
	if eaten == true {
		//shift only [1]
		num = 1
		eaten = false
	} else {
		if length > 1 {
			num = length
		} else {
			num = length
		}

	}
	for i := 1; i < num; i++ {
		if d == U && (snake[i].Y-1 > 0) {
			if snake[i-1].X < snake[i].X {
				snake[i].X--
			} else if snake[i-1].X > snake[i].X {
				snake[i].X++
			} else if snake[i-1].Y < snake[i].Y {
				snake[i].Y--
			} else if snake[i-1].Y > snake[i].Y {
				snake[i].Y++
			}

		} else if d == D && (snake[i].Y+1 < h) {
			if snake[i-1].X < snake[i].X {
				snake[i].X--
			} else if snake[i-1].X > snake[i].X {
				snake[i].X++
			} else if snake[i-1].Y > snake[i].Y {
				snake[i].Y++
			} else if snake[i-1].Y < snake[i].Y {
				snake[i].Y--
			}

		} else if d == L && (snake[i].X-1 > -1) {
			if snake[i-1].Y < snake[i].Y {
				snake[i].Y--
			} else if snake[i-1].Y > snake[i].Y {
				snake[i].Y++
			} else if snake[i-1].X > snake[i].X {
				snake[i].X++
			} else if snake[i-1].X < snake[i].X {
				snake[i].X--
			}

		} else if d == R && (snake[i].X+1 < w) {
			if snake[i-1].Y < snake[i].Y {
				snake[i].Y--
			} else if snake[i-1].Y > snake[i].Y {
				snake[i].Y++
			} else if snake[i-1].X < snake[i].X {
				snake[i].X--
			} else if snake[i-1].X > snake[i].X {
				snake[i].X++
			}
		}
	}
}

func shiftParts(oldX, oldY int) {
    var num int = length
    if eaten == true {
        //shift only [1]
        num = 1
        eaten = false
    } else {
        if length > 1 {
            num = length
        } else {
            num = length
        }

    }
    var oneX, oneY int = 0, 0
    if num > 1 {
        oneX, oneY = snake[1].X, snake[1].Y
        snake[1].X, snake[1].Y = oldX, oldY
    }

    if num > 2 {

        for i := length; i > 1; i-- {
            if i > 2 {
                snake[i].X = snake[i-1].X
                snake[i].Y = snake[i-1].Y
            } else {
                snake[2].X, snake[2].Y = oneX, oneY
            }

        }
    }
}

/* move the snake head; invokes shiftParts() */
func moveSnake(d dir) {
	oldX, oldY := snake[0].X, snake[0].Y
    if d == U && (snake[0].Y-1 > 0) {
		snake[0].Y--
	} else if d == D && (snake[0].Y+1 < h) {
		snake[0].Y++
	} else if d == L && (snake[0].X-1 > -1) {
		snake[0].X--
	} else if d == R && (snake[0].X+1 < w) {
		snake[0].X++
	}

	for h, v1 := range snake {
		for i, v2 := range snake {
			if h != i && ((h != 0) || (i != 0)) && (h < length) && (i < length) {
				if v1.X == v2.X && v1.Y == v2.Y {
					running = false
					break
				}
			}
		}
	}

	//shiftParts(d)
        shiftParts(oldX, oldY)
}

/* prints to pos x, y */
func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

/* draws the screen/updates */
func draw(w, h int) {
	for {
		defer termbox.Flush()

		termbox.Clear(termbox.ColorBlack, termbox.ColorBlue)
		for p, v := range snake {
			if p > length {
				break
			}
			tbPrint(v.X, v.Y, termbox.ColorRed, termbox.ColorBlue, "■")
			//tbPrint(35, 0, termbox.ColorWhite, termbox.ColorBlue, "Pos Snake: "+sc.Itoa(p))
		}
		tbPrint(1, 0, termbox.ColorWhite, termbox.ColorBlue, "Score: "+sc.Itoa(score))
		//fixes the printing of extras...which shouldn't happen, but w/e
		tbPrint(0, 0, termbox.ColorBlue, termbox.ColorBlue, "■")
		tbPrint(target.X, target.Y, termbox.ColorCyan, termbox.ColorBlue, string(target.r))
		checkTarget()

		termbox.Flush()
		time.Sleep(5 * time.Millisecond)
	}
}

/* A basic snake game implemented in termbox-go and Golang */

func main() {
	defer func() {
		termbox.Close()
		fmt.Print("Your score: ", score, "\n")
	}()

	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
	}
	termbox.SetInputMode(termbox.InputAlt)
	w, h = termbox.Size()
	snake[0].X, snake[0].Y = w/2, h/2
	target.X, target.Y, target.r = 20, 8, ''
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlue)
	termbox.Flush()
	go draw(w, h)
	termbox.Flush()

	for running = true; running == true; {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			key := string(ev.Ch)

			if ev.Key == termbox.KeyCtrlQ {
				running = false
				termbox.Flush()
			} else if ev.Key == termbox.KeyEnter {
				/* pause button */
				for r := true; r == true; {
					switch ev := termbox.PollEvent(); ev.Type {
					case termbox.EventKey:
						if ev.Key == termbox.KeyEnter {
							r = false
							break
						}
					default:
					}
				}
			} else if key == "w" {
                if snake[0].Y-1 > 0 {
                    go moveSnake(U)
                } else {
                    running = false
                }
			} else if key == "a" {
                if snake[0].X-1 > -1 {
                    go moveSnake(L)
                } else {
                    running = false
                }
			} else if key == "d" {
                if snake[0].X+1 < w {
                    go moveSnake(R)
                } else {
                    running = false
                }
			} else if key == "s" {
                if snake[0].Y+1 < h {
                    go moveSnake(D)
                } else {
                    running = false
                }
			}
		default:
		}
	}
}
