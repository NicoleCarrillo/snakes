package scripts

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var characterDO *ebiten.DrawImageOptions

// Snake : Object which the player controls
type Snake struct {
	game          *Game
	numParts      int
	lastDir       string
	headUp        ebiten.Image
	headDown      ebiten.Image
	headLeft      ebiten.Image
	headRight     ebiten.Image
	bodyH         ebiten.Image
	bodyV         ebiten.Image
	parts         [][]float64
	pointsWaiting int
	points        int
	behavior      chan int
}

// CreateSnake : Generates a snake
func CreateSnake(g *Game) *Snake {
	s := Snake{
		game:          g,
		numParts:      0,
		lastDir:       "right",
		pointsWaiting: 0,
	}
	s.behavior = make(chan int)
	s.parts = append(s.parts, []float64{300, 300})
	headUp, _, _ := ebitenutil.NewImageFromFile("images/headSerpentUp.png", ebiten.FilterDefault)
	headDown, _, _ := ebitenutil.NewImageFromFile("images/headSerpentDown.png", ebiten.FilterDefault)
	headLeft, _, _ := ebitenutil.NewImageFromFile("images/headSerpentLeft.png", ebiten.FilterDefault)
	headRight, _, _ := ebitenutil.NewImageFromFile("images/headSerpentRight.png", ebiten.FilterDefault)
	bodyH, _, _ := ebitenutil.NewImageFromFile("images/bodySerpentH.png", ebiten.FilterDefault)
	bodyV, _, _ := ebitenutil.NewImageFromFile("images/bodySerpentV.png", ebiten.FilterDefault)
	s.headUp = *headUp
	s.headDown = *headDown
	s.headLeft = *headLeft
	s.headRight = *headRight
	s.bodyH = *bodyH
	s.bodyV = *bodyV
	return &s
}

//Behavior Pipe
func (s *Snake) Behavior() error {
	dotTime := <-s.behavior
	for {
		s.Update(dotTime)
		dotTime = <-s.behavior
	}
}

// Update : Logical update of the snake
func (s *Snake) Update(dotTime int) error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) && s.lastDir != "right" { //movs the snake by pressing key
		s.lastDir = "right"
		return nil
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && s.lastDir != "down" {
		s.lastDir = "down"
		return nil
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) && s.lastDir != "up" {
		s.lastDir = "up"
		return nil
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && s.lastDir != "left" {
		s.lastDir = "left"
		return nil
	}

	if dotTime == 1 { //snakes collide with the boudings
		xPos, yPos := s.getHeadPos()
		if xPos < 0 || xPos > 580 || yPos < 0 || yPos > 580 || s.collisionWithHimself() {
			s.game.End()
		}
	}
	return nil
}

// Draw the snake
func (s *Snake) Draw(screen *ebiten.Image, dotTime int) error {
	if s.game.playing {
		s.UpdatePos(dotTime)
	}
	characterDO = &ebiten.DrawImageOptions{}

	xPos, yPos := s.getHeadPos()
	characterDO.GeoM.Translate(xPos, yPos)

	if s.lastDir == "up" {
		screen.DrawImage(&s.headUp, characterDO)
	} else if s.lastDir == "down" {
		screen.DrawImage(&s.headDown, characterDO)
	} else if s.lastDir == "right" {
		screen.DrawImage(&s.headRight, characterDO)
	} else if s.lastDir == "left" {
		screen.DrawImage(&s.headLeft, characterDO)
	}

	for i := 0; i < s.numParts; i++ { //create the snakes parts
		partDO := &ebiten.DrawImageOptions{}
		xPos, yPos := s.getPartPos(i)
		partDO.GeoM.Translate(xPos, yPos)
		if s.lastDir == "up" || s.lastDir == "down" {
			screen.DrawImage(&s.bodyH, partDO)
		} else {
			screen.DrawImage(&s.bodyV, partDO)
		}

	}

	return nil
}

// UpdatePos changes position values for the snake head
func (s *Snake) UpdatePos(dotTime int) {
	if dotTime == 1 {
		if s.pointsWaiting > 0 {
			s.numParts++
			s.pointsWaiting--
		}
		switch s.lastDir { //method for parts to follow the main snake
		case "up":
			s.translateHeadPos(0, -10)
		case "down":
			s.translateHeadPos(0, +10)
		case "right":
			s.translateHeadPos(10, 0)
		case "left":
			s.translateHeadPos(-10, 0)
		}

	}
}

func (s *Snake) addPoint() {
	s.points++
	s.pointsWaiting++
}

func (s *Snake) getHeadPos() (float64, float64) { // get position of the head
	return s.parts[0][0], s.parts[0][1]
}

func (s *Snake) getPartPos(pos int) (float64, float64) { //get position of parts
	return s.parts[pos+1][0], s.parts[pos+1][1]
}

func (s *Snake) translateHeadPos(newXPos, newYPos float64) {
	newX := s.parts[0][0] + newXPos
	newY := s.parts[0][1] + newYPos
	s.updateParts(newX, newY)
}

func (s *Snake) updateParts(newX, newY float64) { //update the parts of the snake
	s.parts = append([][]float64{[]float64{newX, newY}}, s.parts...)
	s.parts = s.parts[:s.numParts+1]
}

func (s *Snake) collisionWithHimself() bool { //checks colition with itself
	posX, posY := s.getHeadPos()
	for i := 1; i < len(s.parts); i++ {
		if posX == s.parts[i][0] && posY == s.parts[i][1] {
			return true
		}
	}
	return false
}
