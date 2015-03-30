package ui

type View interface {
	Enter()
	Exit()
	Update(t, dt float64)
	Draw()
}
