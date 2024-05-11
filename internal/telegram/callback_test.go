package telegram

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallbackData_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		cd       CallbackData
		expected string
	}{
		{
			name: "single param",
			cd: CallbackData{
				Route:  "route1",
				Params: []string{"param1"},
			},
			expected: "route1/param1",
		},
		{
			name: "multiple params",
			cd: CallbackData{
				Route:  "route2",
				Params: []string{"param1", "param2"},
			},
			expected: "route2/param1/param2",
		},
		{
			name: "no params",
			cd: CallbackData{
				Route:  "route3",
				Params: []string{},
			},
			expected: "route3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cd.Marshal()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCallbackData_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected CallbackData
	}{
		{
			name: "single param",
			data: "route1/param1",
			expected: CallbackData{
				Route:  "route1",
				Params: []string{"param1"},
			},
		},
		{
			name: "multiple params",
			data: "route2/param1/param2",
			expected: CallbackData{
				Route:  "route2",
				Params: []string{"param1", "param2"},
			},
		},
		{
			name: "no params",
			data: "route3",
			expected: CallbackData{
				Route:  "route3",
				Params: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cd CallbackData
			cd.Unmarshal(tt.data)
			require.Equal(t, tt.expected, cd)
		})
	}
}
