package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotX90(t *testing.T) {
	tt := []struct {
		name string
		in   I3
		want I3
	}{
		{
			"z unit vector becomes negative y unit vector",
			I3{Z: 1},
			I3{Y: -1},
		},
		{
			"y unit vector becomes z unit vector",
			I3{Y: 1},
			I3{Z: 1},
		},
		{
			"rotation does not affect x value",
			I3{X: 12, Y: 2, Z: 4},
			I3{X: 12, Y: -4, Z: 2},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			v := tc.in.RotX90()
			assert.Equal(t, tc.want, v)
		})
	}
}

func TestRotY90(t *testing.T) {
	tt := []struct {
		name string
		in   I3
		want I3
	}{
		{
			"z unit vector becomes x unit vector",
			I3{Z: 1},
			I3{X: 1},
		},
		{
			"x unit vector becomes negative z unit vector",
			I3{X: 1},
			I3{Z: -1},
		},
		{
			"rotation does not affect y value",
			I3{Y: 12, Z: 2, X: 4},
			I3{Y: 12, Z: -4, X: 2},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			v := tc.in.RotY90()
			assert.Equal(t, tc.want, v)
		})
	}
}

func TestI3RotZ90(t *testing.T) {
	tt := []struct {
		name string
		in   I3
		want I3
	}{
		{
			"x unit vector becomes y unit vector",
			I3{X: 1},
			I3{Y: 1},
		},
		{
			"y unit vector becomes negative x unit vector",
			I3{Y: 1},
			I3{X: -1},
		},
		{
			"rotation does not affect z value",
			I3{Z: 12, X: 2, Y: 4},
			I3{Z: 12, X: -4, Y: 2},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			v := tc.in.RotZ90()
			assert.Equal(t, tc.want, v)
		})
	}
}
