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
	buttons [8]byte
	index   byte
	strobe  byte
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Press(button int) {
	c.buttons[button] = 1
}

func (c *Controller) Release(button int) {
	c.buttons[button] = 0
}

func (c *Controller) SetPressed(button int, pressed bool) {
	if pressed {
		c.Press(button)
	} else {
		c.Release(button)
	}
}

func (c *Controller) Read() byte {
	var value byte
	if c.index < 8 {
		value = c.buttons[c.index]
	} else {
		value = 0
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
