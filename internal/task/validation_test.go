package task

import "testing"

func TestValidateTask(t *testing.T) {
	tests := []struct {
		name     string
		input    CreateTaskRequest
		expected bool
	}{
		{
			name: "valid task",
			input: CreateTaskRequest{
				Title:    "Learn Go",
				Status:   "todo",
				Priority: "high",
			},
			expected: true,
		},
		{
			name: "empty title",
			input: CreateTaskRequest{
				Title:    "",
				Status:   "todo",
				Priority: "high",
			},
			expected: false,
		},
		{
			name: "whitespace title",
			input: CreateTaskRequest{
				Title:    "   ",
				Status:   "todo",
				Priority: "high",
			},
			expected: false,
		},
		{
			name: "invalid status",
			input: CreateTaskRequest{
				Title:    "Learn Go",
				Status:   "unknown",
				Priority: "high",
			},
			expected: false,
		},
		{
			name: "invalid priority",
			input: CreateTaskRequest{
				Title:    "Learn Go",
				Status:   "todo",
				Priority: "urgent",
			},
			expected: false,
		},
		{
			name: "valid in progress",
			input: CreateTaskRequest{
				Title:    "Build API",
				Status:   "in_progress",
				Priority: "medium",
			},
			expected: true,
		},
		{
			name: "valid completed task",
			input: CreateTaskRequest{
				Title:    "Finish project",
				Status:   "done",
				Priority: "low",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := validateTask(tt.input)

			if actual != tt.expected {
				t.Errorf(
					"Expected %v, got %v",
					tt.expected,
					actual,
				)
			}
		})
	}
}
