package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Chip8 struct {
	memory     [4096]byte
	V          [16]byte
	I          uint16
	PC         uint16
	stack      [16]uint16
	sp         uint8
	delayTimer uint8
	soundTimer uint8
	keys       [16]byte
	display    [64 * 32]byte
	drawFlag   bool
}

func (c *Chip8) LoadCheckerboard() {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if (x+y)%2 == 0 {
				c.display[y*64+x] = 1
			} else {
				c.display[y*64+x] = 0
			}
		}
	}
}

func (c *Chip8) draw(opcode uint16) {
	xReg := (opcode & 0x0F00) >> 8
	yReg := (opcode & 0x00F0) >> 4
	height := opcode & 0x000F

	x := int(c.V[xReg])
	y := int(c.V[yReg])

	c.V[0xF] = 0 // collision flag

	for row := 0; row < int(height); row++ {
		sprite := c.memory[int(c.I)+row]

		for col := 0; col < 8; col++ {
			// check if bit is set (left to right)
			if (sprite & (0x80 >> col)) != 0 {

				px := (x + col) % 64
				py := (y + row) % 32
				index := py*64 + px

				if c.display[index] == 1 {
					c.V[0xF] = 1 // collision
				}

				c.display[index] ^= 1 // XOR pixel
			}
		}
	}
}

func (c *Chip8) Cycle() {
	opcode := uint16(c.memory[c.PC])<<8 | uint16(c.memory[c.PC+1])
	c.PC += 2

	instruction := opcode & 0xF000

	switch instruction {
	case 0x1000:
		// 1NNN: jump to address NNN
		c.PC = opcode & 0x0FFF

	case 0x6000:
		// 6XNN: set VX = NN
		x := (opcode & 0x0F00) >> 8
		c.V[x] = byte(opcode & 0x00FF)

	case 0x7000:
		// 7XNN: VX += NN
		x := (opcode & 0x0F00) >> 8
		c.V[x] += byte(opcode & 0x00FF)

	case 0xD000:
		c.draw(opcode)
	}

}

func (c *Chip8) LoadROM(path string) error {

	rom := []byte{
		0x60, 0x01, // 6001 -> V0 = 1
		0x70, 0x01, // 7001 -> V0 += 1
		0x12, 0x00, // 1200 -> jump to start
	}

	copy(c.memory[0x200:], rom)
	c.PC = 0x200

	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	return err
	// }

	// start := 0x200

	// for i, b := range data {
	// 	if start+i >= len(c.memory) {
	// 		break
	// 	}
	// 	c.memory[start+i] = b
	// }

	// c.PC = 0x200
	return nil
}

type Game struct {
	chip8 Chip8
}

func (g *Game) Update() error {
	// Write your game's logical update.
	for i := 0; i < 10; i++ {
		g.chip8.Cycle()
	}
	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {
	frame := ebiten.NewImage(64, 32)

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if g.chip8.display[y*64+x] == 1 {
				frame.Set(x, y, color.White)
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(10, 10)
	screen.DrawImage(frame, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 320
}

func main() {
	chip8 := Chip8{}
	chip8.LoadROM("test_opcode.ch8")

	game := &Game{
		chip8: chip8,
	}
	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("CHIP-8 Emulator")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
