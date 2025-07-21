package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) Option {
	return func(p *GamePerson) {
		copy(p.name[:], name)
	}
}

func WithCoordinates(x, y, z int) Option {
	return func(p *GamePerson) {
		p.x = int32(x)
		p.y = int32(y)
		p.z = int32(z)
	}
}

func WithGold(gold int) Option {
	return func(p *GamePerson) {
		p.gold = uint32(gold)
	}
}

func WithMana(mana int) Option {
	return func(p *GamePerson) {
		if mana < 0 {
			mana = 0
		} else if mana > 1000 {
			mana = 1000
		}

		raw := uint32(p.manaHealth[0]) |
			uint32(p.manaHealth[1])<<8 |
			uint32(p.manaHealth[2])<<16

		raw &^= 0x3FF
		raw |= uint32(mana) & 0x3FF

		p.manaHealth[0] = byte(raw)
		p.manaHealth[1] = byte(raw >> 8)
		p.manaHealth[2] = byte(raw >> 16)
	}
}

func WithHealth(health int) Option {
	return func(p *GamePerson) {
		if health < 0 {
			health = 0
		} else if health > 1000 {
			health = 1000
		}

		raw := uint32(p.manaHealth[0]) |
			uint32(p.manaHealth[1])<<8 |
			uint32(p.manaHealth[2])<<16

		raw &^= 0x3FF << 10
		raw |= (uint32(health) & 0x3FF) << 10

		p.manaHealth[0] = byte(raw)
		p.manaHealth[1] = byte(raw >> 8)
		p.manaHealth[2] = byte(raw >> 16)
	}
}

func WithRespect(respect int) Option {
	return func(p *GamePerson) {
		if respect < 0 {
			respect = 0
		} else if respect > 10 {
			respect = 10
		}
		p.attrs &^= 0xF000
		p.attrs |= uint16(respect) << 12
	}
}

func WithStrength(strength int) Option {
	return func(p *GamePerson) {
		if strength < 0 {
			strength = 0
		} else if strength > 10 {
			strength = 10
		}
		p.attrs &^= 0x0F00
		p.attrs |= uint16(strength) << 8
	}
}

func WithExperience(experience int) Option {
	return func(p *GamePerson) {
		if experience < 0 {
			experience = 0
		} else if experience > 10 {
			experience = 10
		}
		p.attrs &^= 0x00F0
		p.attrs |= uint16(experience) << 4
	}
}

func WithLevel(level int) Option {
	return func(p *GamePerson) {
		if level < 0 {
			level = 0
		} else if level > 10 {
			level = 10
		}
		p.attrs &^= 0x000F
		p.attrs |= uint16(level)
	}
}

func WithHouse() Option {
	return func(p *GamePerson) {
		p.params |= 0b0100
	}
}

func WithGun() Option {
	return func(p *GamePerson) {
		p.params |= 0b0010
	}
}

func WithFamily() Option {
	return func(p *GamePerson) {
		p.params |= 0b0001
	}
}

func WithType(personType int) Option {
	return func(p *GamePerson) {
		typeVal := (uint8(personType) & 0x3) << 4
		p.params &^= 0b00110000
		p.params |= typeVal
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z    int32    // 12 байт: координаты X, Y, Z (диапазон [-2_000_000_000…2_000_000_000])
	gold       uint32   // 4 байта: золото (диапазон [0…2_000_000_000])
	manaHealth [3]byte  // 3 байта: мана (биты 0-9, [0…1000]), здоровье (биты 10-19, [0…1000])
	params     uint8    // 1 байт: флаги (биты 0-2: HasHouse, HasWeapon, HasFamily)
	attrs      uint16   // 2 байта: Respect (биты 0-3), Strength (биты 4-7), Experience (биты 8-11), Level (биты 12-15), диапазон [0…10]
	name       [42]byte // 42 байта: имя пользователя (ASCII, до 42 символов)
}

func NewGamePerson(options ...Option) GamePerson {
	var g = GamePerson{}
	for _, option := range options {
		option(&g)
	}
	return g
}

func (p *GamePerson) Name() string {
	n := 0
	for n < len(p.name) && p.name[n] != 0 {
		n++
	}
	return string(p.name[:n])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	v := uint32(p.manaHealth[0]) | (uint32(p.manaHealth[1]) << 8) | (uint32(p.manaHealth[2]) << 16)
	return int(v & 0x3FF)
}

func (p *GamePerson) Health() int {
	v := uint32(p.manaHealth[0]) | (uint32(p.manaHealth[1]) << 8) | (uint32(p.manaHealth[2]) << 16)
	return int((v >> 10) & 0x3FF)
}

func (p *GamePerson) Respect() int {
	return int(p.attrs & 0xF000 >> 12)
}

func (p *GamePerson) Strength() int {
	return int(p.attrs & 0xF00 >> 8)
}

func (p *GamePerson) Experience() int {
	return int(p.attrs & 0xF0 >> 4)
}

func (p *GamePerson) Level() int {
	return int(p.attrs & 0xF)
}

func (p *GamePerson) HasHouse() bool {
	return (p.params & 0b0100) == 0b0100
}

func (p *GamePerson) HasGun() bool {
	return (p.params & 0b0010) == 0b0010
}

func (p *GamePerson) HasFamilty() bool {
	return (p.params & 0b0001) == 0b0001
}

func (p *GamePerson) Type() int {
	return int(p.params & 0xF0 >> 4)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
