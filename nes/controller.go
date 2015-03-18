package nes

const (
	ButtonA = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight
)

type Controller struct {
	buttons [8]bool
	index   byte
	strobe  byte
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) SetButtons(buttons [8]bool) {
	c.buttons = buttons
}

func (c *Controller) Read() byte {
	value := byte(0)
	if c.index < 8 && c.buttons[c.index] {
		value = 1
	}
	c.index++
	if c.strobe&1 == 1 {
		c.index = 0
	}
	return value
}

func (c *Controller) Write(value byte) {
	c.strobe = value
	if c.strobe&1 == 1 {
		c.index = 0
	}
}
