package mathx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestMax(t *testing.T) {
	m := Max(1, 2, -3, 4, 5-119)
	assert.Equal(t, 4, m)

	m2 := Max(3.0, NaN())
	assert.Equal(t, "NaN", fmt.Sprintf("%v", m2))

	m3 := Max(3.0, InfMax())
	assert.Equal(t, "+Inf", fmt.Sprintf("%v", m3))

	m4 := Max(3.0, InfMin())
	assert.Equal(t, 3.0, m4)
}

func TestMin(t *testing.T) {
	m := Min(1, 2, 5, 6, 7)
	assert.Equal(t, 1, m)

	m2 := Min(3.0, math.NaN())
	assert.Equal(t, "NaN", fmt.Sprintf("%v", m2))

	m3 := Min(3.0, InfMin())
	assert.Equal(t, "-Inf", fmt.Sprintf("%v", m3))

	m4 := Min(3.0, InfMax())
	assert.Equal(t, 3.0, m4)
}
