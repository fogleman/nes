package nes

const frameCounterRate = 7457.3875

var lengthTable = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

// APU

type APU struct {
	console             *Console
	pulse1              Pulse
	pulse2              Pulse
	cycle               uint64
	frame               uint64
	frameCounterMode    byte
	frameCounterInhibit byte
}

func NewAPU(console *Console) *APU {
	apu := APU{}
	apu.console = console
	return &apu
}

func (apu *APU) Step() {
	cycle1 := apu.cycle
	apu.cycle++
	cycle2 := apu.cycle
	f1 := int(float64(cycle1) / frameCounterRate)
	f2 := int(float64(cycle2) / frameCounterRate)
	if f1 != f2 {
		apu.stepFrameCounter()
	}
}

// mode 0:    mode 1:       function
// ---------  -----------  -----------------------------
//  - - - f    - - - - -    IRQ (if bit 6 is clear)
//  - l - l    l - l - -    Length counter and sweep
//  e e e e    e e e e -    Envelope and linear counter
func (apu *APU) stepFrameCounter() {
	apu.frame++
	switch apu.frameCounterMode {
	case 0:
		switch apu.frame % 4 {
		case 0, 2:
			apu.stepEnvelope()
		case 1:
			apu.stepEnvelope()
			apu.stepLength()
		case 3:
			apu.stepEnvelope()
			apu.stepLength()
			apu.fireIRQ()
		}
	case 1:
		switch apu.frame % 5 {
		case 1, 3:
			apu.stepEnvelope()
		case 0, 2:
			apu.stepEnvelope()
			apu.stepLength()
		}
	}
}

func (apu *APU) stepEnvelope() {
}

func (apu *APU) stepLength() {
}

func (apu *APU) fireIRQ() {
	if apu.frameCounterInhibit == 0 {
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
		apu.pulse1.writeRegister0(value)
	case 0x4001:
		apu.pulse1.writeRegister1(value)
	case 0x4002:
		apu.pulse1.writeRegister2(value)
	case 0x4003:
		apu.pulse1.writeRegister3(value)
	case 0x4004:
		apu.pulse2.writeRegister0(value)
	case 0x4005:
		apu.pulse2.writeRegister1(value)
	case 0x4006:
		apu.pulse2.writeRegister2(value)
	case 0x4007:
		apu.pulse2.writeRegister3(value)
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
}

func (apu *APU) writeFrameCounter(value byte) {
	apu.frameCounterMode = (value >> 7) & 1
	apu.frameCounterInhibit = (value >> 6) & 1
}

// Pulse

type Pulse struct {
	duty         byte
	halt         byte
	constant     byte
	volume       byte
	sweepEnabled byte
	sweepPeriod  byte
	sweepNegate  byte
	sweepShift   byte
	timer        uint16
	length       byte
}

// $4000 / $4004	DDLC VVVV	Duty (D), envelope loop / length counter halt (L), constant volume (C), volume/envelope (V)
// $4001 / $4005	EPPP NSSS	Sweep unit: enabled (E), period (P), negate (N), shift (S)
// $4002 / $4006	TTTT TTTT	Timer low (T)
// $4003 / $4007	LLLL LTTT	Length counter load (L), timer high (T)

func (p *Pulse) writeRegister0(value byte) {
	p.duty = (value >> 6) & 3
	p.halt = (value >> 5) & 1
	p.constant = (value >> 4) & 1
	p.volume = value & 15
}

func (p *Pulse) writeRegister1(value byte) {
	p.sweepEnabled = (value >> 7) & 1
	p.sweepPeriod = (value >> 4) & 7
	p.sweepNegate = (value >> 3) & 1
	p.sweepShift = value & 7
}

func (p *Pulse) writeRegister2(value byte) {
	p.timer = (p.timer & 0xFF00) | uint16(value)
}

func (p *Pulse) writeRegister3(value byte) {
	p.timer = (p.timer & 0x00FF) | (uint16(value&7) << 8)
	p.length = lengthTable[value>>3]
}
