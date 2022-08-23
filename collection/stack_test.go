package collection

import (
	"testing"
)

func TestArrayStack_Cap(t *testing.T) {
	type fields struct {
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"1", fields{1}, 1},
		{"100", fields{100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArrayStack(tt.fields.cap)
			if got := a.Cap(); got != tt.want {
				t.Errorf("Cap() = %v, want %v", got, tt.want)
			}
			if got := a.Len(); got != 0 {
				t.Errorf("Len() = %v, want %v", got, 0)
			}
		})
	}
}

func TestArrayStack_Len(t *testing.T) {
	type fields struct {
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"1", fields{1}, 2},
		{"100", fields{100}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArrayStack(tt.fields.cap)
			a.Push("12")
			a.Push(1)
			a.Push(3.5)

			if got, _ := a.Pop(); got != 3.5 {
				t.Errorf("Pop() = %v, want %v", got, 3.5)
			}
			if got := a.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayStack_Pop(t *testing.T) {
	type fields struct {
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"1", fields{1}, 2},
		{"100", fields{100}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArrayStack(tt.fields.cap)
			a.Push("12")
			a.Push(1)
			a.Push(3.5)

			if got, _ := a.Pop(); got != 3.5 {
				t.Errorf("Pop() = %v, want %v", got, 3.5)
			}
			if got := a.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayStack_Push(t *testing.T) {
	type fields struct {
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"1", fields{1}, 2},
		{"100", fields{100}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArrayStack(tt.fields.cap)
			a.Push("12")
			a.Push(1)
			a.Push(3.5)

			if got, _ := a.Pop(); got != 3.5 {
				t.Errorf("Pop() = %v, want %v", got, 3.5)
			}
			if got := a.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayStack_Top(t *testing.T) {
	type fields struct {
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"1", fields{1}, 2},
		{"100", fields{100}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArrayStack(tt.fields.cap)
			a.Push("12")
			a.Push(1)
			a.Push(3.5)

			if got, _ := a.Top(); got != 3.5 {
				t.Errorf("Top() = %v, want %v", got, 3.5)
			}
			if got, _ := a.Pop(); got != 3.5 {
				t.Errorf("Pop() = %v, want %v", got, 3.5)
			}
			if got, _ := a.Top(); got != 1 {
				t.Errorf("Top() = %v, want %v", got, 1)
			}
			if got := a.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestLinkedStack_Top(t *testing.T) {
	type fields struct {
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"1", fields{1}, 2},
		{"100", fields{100}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewLinkedStack()
			a.Push("12")
			a.Push(1)
			a.Push(3.5)

			if got, _ := a.Top(); got != 3.5 {
				t.Errorf("Top() = %v, want %v", got, 3.5)
			}
			if got, _ := a.Pop(); got != 3.5 {
				t.Errorf("Pop() = %v, want %v", got, 3.5)
			}
			if got, _ := a.Top(); got != 1 {
				t.Errorf("Top() = %v, want %v", got, 1)
			}
			if got := a.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}