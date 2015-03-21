package nes

const frameCounterRate = CPUFrequency / 240.0
const sampleRate = CPUFrequency / 44100.0 / 2

var lengthTable = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

var dutyTable = [][]byte{
	{0, 1, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0},
	{1, 0, 0, 1, 1, 1, 1, 1},
}

var triangleTable = []byte{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}

var noiseTable = []uint16{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}

var pulseTable [31]float64
var tndTable [203]float64

func init() {
	for i := 0; i < 31; i++ {
		pulseTable[i] = 95.52 / (8128.0/float64(i) + 100)
	}
	for i := 0; i < 203; i++ {
		tndTable[i] = 163.67 / (24329.0/float64(i) + 100)
	}
}

// APU

type APU struct {
	console        *Console
	channel        chan byte
	pulse1         Pulse
	pulse2         Pulse
	triangle       Triangle
	noise          Noise
	cycle          uint64
	framePeriod    byte
	frameValue     byte
	frameIRQ       bool
	enablePulse1   bool
	enablePulse2   bool
	enableTriangle bool
	enableNoise    bool
}

func NewAPU(console *Console) *APU {
	apu := APU{}
	apu.console = console
	apu.enablePulse1 = true
	apu.enablePulse2 = true
	apu.enableTriangle = true
	apu.enableNoise = true
	apu.noise.shiftRegister = 1
	return &apu
}

func (apu *APU) SetChannel(channel chan byte) {
	apu.channel = channel
}

func (apu *APU) Step() {
	cycle1 := apu.cycle
	apu.cycle++
	cycle2 := apu.cycle
	apu.stepTimer()
	f1 := int(float64(cycle1) / frameCounterRate)
	f2 := int(float64(cycle2) / frameCounterRate)
	if f1 != f2 {
		apu.stepFrameCounter()
	}
	s1 := int(float64(cycle1) / sampleRate)
	s2 := int(float64(cycle2) / sampleRate)
	if s1 != s2 {
		apu.sendSample()
	}
}

func (apu *APU) sendSample() {
	select {
	case apu.channel <- apu.output():
	default:
	}
}

func (apu *APU) output() byte {
	var p1, p2, t, n, d byte
	if apu.enablePulse1 {
		p1 = apu.pulse1.output()
	}
	if apu.enablePulse2 {
		p2 = apu.pulse2.output()
	}
	if apu.enableTriangle {
		t = apu.triangle.output()
	}
	if apu.enableNoise {
		n = apu.noise.output()
	}
	pulseOut := pulseTable[p1+p2]
	tndOut := tndTable[3*t+2*n+d]
	return byte((pulseOut + tndOut) * 255)
}

// mode 0:    mode 1:       function
// ---------  -----------  -----------------------------
//  - - - f    - - - - -    IRQ (if bit 6 is clear)
//  - l - l    l - l - -    Length counter and sweep
//  e e e e    e e e e -    Envelope and linear counter
func (apu *APU) stepFrameCounter() {
	switch apu.framePeriod {
	case 4:
		apu.frameValue = (apu.frameValue + 1) % 4
		switch apu.frameValue {
		case 0, 2:
			apu.stepEnvelope()
		case 1:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
		case 3:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
			apu.fireIRQ()
		}
	case 5:
		apu.frameValue = (apu.frameValue + 1) % 5
		switch apu.frameValue {
		case 1, 3:
			apu.stepEnvelope()
		case 0, 2:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
		}
	}
}

func (apu *APU) stepTimer() {
	if apu.cycle%2 == 0 {
		apu.pulse1.stepTimer()
		apu.pulse2.stepTimer()
		apu.noise.stepTimer()
	}
	apu.triangle.stepTimer()
}

func (apu *APU) stepEnvelope() {
	apu.pulse1.stepEnvelope()
	apu.pulse2.stepEnvelope()
	apu.triangle.stepCounter()
	apu.noise.stepEnvelope()
}

func (apu *APU) stepSweep() {
	apu.pulse1.stepSweep()
	apu.pulse2.stepSweep()
}

func (apu *APU) stepLength() {
	apu.pulse1.stepLength()
	apu.pulse2.stepLength()
	apu.triangle.stepLength()
	apu.noise.stepLength()
}

func (apu *APU) fireIRQ() {
	if apu.frameIRQ {
		apu.console.CPU.triggerIRQ()
	}
}

func (apu *APU) readRegister(address uint16) byte {
	switch address {
	case 0x4015:
		return apu.readStatus()
		// default:
		// 	log.Fatalf("unhandled apu register read at address: 0x%04X", address)
	}
	return 0
}

func (apu *APU) writeRegister(address uint16, value byte) {
	switch address {
	case 0x4000:
		apu.pulse1.writeControl(value)
	case 0x4001:
		apu.pulse1.writeSweep(value)
	case 0x4002:
		apu.pulse1.writeTimerLow(value)
	case 0x4003:
		apu.pulse1.writeTimerHigh(value)
	case 0x4004:
		apu.pulse2.writeControl(value)
	case 0x4005:
		apu.pulse2.writeSweep(value)
	case 0x4006:
		apu.pulse2.writeTimerLow(value)
	case 0x4007:
		apu.pulse2.writeTimerHigh(value)
	case 0x4008:
		apu.triangle.writeControl(value)
	case 0x4009:
	case 0x400A:
		apu.triangle.writeTimerLow(value)
	case 0x400B:
		apu.triangle.writeTimerHigh(value)
	case 0x400C:
		apu.noise.writeControl(value)
	case 0x400D:
	case 0x400E:
		apu.noise.writePeriod(value)
	case 0x400F:
		apu.noise.writeLength(value)
	case 0x4015:
		apu.writeControl(value)
	case 0x4017:
		apu.writeFrameCounter(value)
		// default:
		// 	log.Fatalf("unhandled apu register write at address: 0x%04X", address)
	}
}

func (apu *APU) readStatus() byte {
	return 0
}

func (apu *APU) writeControl(value byte) {
	apu.enablePulse1 = value&1 == 1
	apu.enablePulse2 = value&2 == 2
	apu.enableTriangle = value&4 == 4
	apu.enableNoise = value&8 == 8
}

func (apu *APU) writeFrameCounter(value byte) {
	apu.framePeriod = 4 + (value>>7)&1
	apu.frameIRQ = (value>>6)&1 == 0
}

// Pulse

type Pulse struct {
	lengthEnabled   bool
	lengthValue     byte
	timerPeriod     uint16
	timerValue      uint16
	dutyMode        byte
	dutyValue       byte
	sweepReload     bool
	sweepEnabled    bool
	sweepNegate     bool
	sweepShift      byte
	sweepPeriod     byte
	sweepValue      byte
	envelopeEnabled bool
	envelopeLoop    bool
	envelopeStart   bool
	envelopePeriod  byte
	envelopeValue   byte
	envelopeVolume  byte
	constantVolume  byte
}

func (p *Pulse) writeControl(value byte) {
	p.dutyMode = (value >> 6) & 3
	p.lengthEnabled = (value>>5)&1 == 0
	p.envelopeLoop = (value>>5)&1 == 1
	p.envelopeEnabled = (value>>4)&1 == 0
	p.envelopePeriod = value & 15
	p.constantVolume = value & 15
	p.envelopeStart = true
}

func (p *Pulse) writeSweep(value byte) {
	p.sweepEnabled = (value>>7)&1 == 1
	p.sweepPeriod = (value >> 4) & 7
	p.sweepNegate = (value>>3)&1 == 1
	p.sweepShift = value & 7
	p.sweepReload = true
}

func (p *Pulse) writeTimerLow(value byte) {
	p.timerPeriod = (p.timerPeriod & 0xFF00) | uint16(value)
}

func (p *Pulse) writeTimerHigh(value byte) {
	p.lengthValue = lengthTable[value>>3]
	p.timerPeriod = (p.timerPeriod & 0x00FF) | (uint16(value&7) << 8)
	p.timerValue = p.timerPeriod
	p.envelopeStart = true
	p.dutyValue = 0
}

func (p *Pulse) stepTimer() {
	if p.timerValue == 0 {
		p.timerValue = p.timerPeriod
		p.dutyValue = (p.dutyValue + 1) % 8
	} else {
		p.timerValue--
	}
}

func (p *Pulse) stepEnvelope() {
	if p.envelopeStart {
		p.envelopeVolume = 15
		p.envelopeValue = p.envelopePeriod
		p.envelopeStart = false
	} else if p.envelopeValue > 0 {
		p.envelopeValue--
	} else {
		if p.envelopeVolume > 0 {
			p.envelopeVolume--
		} else if p.envelopeLoop {
			p.envelopeVolume = 15
			p.envelopeValue = p.envelopePeriod
		}
	}
}

func (p *Pulse) stepSweep() {
	if p.sweepReload {
		if p.sweepEnabled && p.sweepValue == 0 {
			p.sweep()
		}
		p.sweepValue = p.sweepPeriod
		p.sweepReload = false
	} else if p.sweepValue > 0 {
		p.sweepValue--
	} else {
		if p.sweepEnabled {
			p.sweep()
		}
		p.sweepValue = p.sweepPeriod
	}
}

func (p *Pulse) stepLength() {
	if p.lengthEnabled && p.lengthValue > 0 {
		p.lengthValue--
	}
}

func (p *Pulse) sweep() {
	delta := p.timerPeriod >> p.sweepShift
	if p.sweepNegate {
		delta = -delta
	}
	p.timerPeriod += delta
}

func (p *Pulse) output() byte {
	if p.lengthValue == 0 {
		return 0
	}
	if dutyTable[p.dutyMode][p.dutyValue] == 0 {
		return 0
	}
	if p.timerPeriod < 8 || p.timerPeriod > 0x7FF {
		return 0
	}
	if p.envelopeEnabled {
		return p.envelopeVolume
	} else {
		return p.constantVolume
	}
}

// Triangle

type Triangle struct {
	lengthEnabled bool
	lengthValue   byte
	timerPeriod   uint16
	timerValue    uint16
	dutyValue     byte
	counterPeriod byte
	counterValue  byte
	counterReload bool
}

func (t *Triangle) writeControl(value byte) {
	t.lengthEnabled = (value>>7)&1 == 0
	t.counterPeriod = value & 0x7F
}

func (t *Triangle) writeTimerLow(value byte) {
	t.timerPeriod = (t.timerPeriod & 0xFF00) | uint16(value)
}

func (t *Triangle) writeTimerHigh(value byte) {
	t.lengthValue = lengthTable[value>>3]
	t.timerPeriod = (t.timerPeriod & 0x00FF) | (uint16(value&7) << 8)
	t.timerValue = t.timerPeriod
	t.counterReload = true
}

func (t *Triangle) stepTimer() {
	if t.timerValue == 0 {
		t.timerValue = t.timerPeriod
		if t.lengthValue > 0 && t.counterValue > 0 {
			t.dutyValue = (t.dutyValue + 1) % 32
		}
	} else {
		t.timerValue--
	}
}

func (t *Triangle) stepLength() {
	if t.lengthEnabled && t.lengthValue > 0 {
		t.lengthValue--
	}
}

func (t *Triangle) stepCounter() {
	if t.counterReload {
		t.counterValue = t.counterPeriod
	} else if t.counterValue > 0 {
		t.counterValue--
	}
	if t.lengthEnabled {
		t.counterReload = false
	}
}

func (t *Triangle) output() byte {
	if t.lengthValue == 0 {
		return 0
	}
	if t.counterValue == 0 {
		return 0
	}
	return triangleTable[t.dutyValue]
}

// Noise

type Noise struct {
	mode            bool
	shiftRegister   uint16
	lengthEnabled   bool
	lengthValue     byte
	timerPeriod     uint16
	timerValue      uint16
	envelopeEnabled bool
	envelopeLoop    bool
	envelopeStart   bool
	envelopePeriod  byte
	envelopeValue   byte
	envelopeVolume  byte
	constantVolume  byte
}

func (n *Noise) writeControl(value byte) {
	n.lengthEnabled = (value>>5)&1 == 0
	n.envelopeLoop = (value>>5)&1 == 1
	n.envelopeEnabled = (value>>4)&1 == 0
	n.envelopePeriod = value & 15
	n.constantVolume = value & 15
	n.envelopeStart = true
}

func (n *Noise) writePeriod(value byte) {
	n.mode = value&0x80 == 0x80
	n.timerPeriod = noiseTable[value&0x0F]
}

func (n *Noise) writeLength(value byte) {
	n.lengthValue = lengthTable[value>>3]
	n.envelopeStart = true
}

func (n *Noise) stepTimer() {
	if n.timerValue == 0 {
		n.timerValue = n.timerPeriod
		var shift byte
		if n.mode {
			shift = 6
		} else {
			shift = 1
		}
		b1 := n.shiftRegister & 1
		b2 := (n.shiftRegister >> shift) & 1
		n.shiftRegister >>= 1
		n.shiftRegister |= (b1 ^ b2) << 14
	} else {
		n.timerValue--
	}
}

func (n *Noise) stepEnvelope() {
	if n.envelopeStart {
		n.envelopeVolume = 15
		n.envelopeValue = n.envelopePeriod
		n.envelopeStart = false
	} else if n.envelopeValue > 0 {
		n.envelopeValue--
	} else {
		if n.envelopeVolume > 0 {
			n.envelopeVolume--
		} else if n.envelopeLoop {
			n.envelopeVolume = 15
			n.envelopeValue = n.envelopePeriod
		}
	}
}

func (n *Noise) stepLength() {
	if n.lengthEnabled && n.lengthValue > 0 {
		n.lengthValue--
	}
}

func (n *Noise) output() byte {
	if n.lengthValue == 0 {
		return 0
	}
	if n.shiftRegister&1 == 1 {
		return 0
	}
	if n.envelopeEnabled {
		return n.envelopeVolume
	} else {
		return n.constantVolume
	}
}
