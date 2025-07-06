package usecases

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taldoflemis/nume/internal/expressions"
)

type doubleIntegralTestCase struct {
	name               string
	expr               expressions.DualVariableExpr
	leftIntervalX      float64
	rightIntervalX     float64
	leftIntervalY      float64
	rightIntervalY     float64
	numberOfPartitions uint64
	expectedArea       float64
	tolerance          float64
	description        string
}

func TestDoubleIntegralCalculateArea(t *testing.T) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Arrange
	t.Parallel()

	tests := []doubleIntegralTestCase{
		{
			name: "Unit Square",
			expr: func(x, y float64) float64 {
				return 1.0 // Constant function over unit square
			},
			leftIntervalX:      0,
			rightIntervalX:     1,
			leftIntervalY:      0,
			rightIntervalY:     1,
			numberOfPartitions: 100,
			expectedArea:       1.0,
			tolerance:          1e-10,
			description:        "Area of unit square with constant function f(x,y) = 1",
		},
		{
			name: "Rectangle 2x3",
			expr: func(x, y float64) float64 {
				return 1.0 // Constant function over rectangle
			},
			leftIntervalX:      0,
			rightIntervalX:     2,
			leftIntervalY:      0,
			rightIntervalY:     3,
			numberOfPartitions: 100,
			expectedArea:       6.0,
			tolerance:          1e-10,
			description:        "Area of 2x3 rectangle with constant function f(x,y) = 1",
		},
		{
			name: "Circle Approximation with radius 1 and center = 0",
			expr: func(x, y float64) float64 {
				radius := 1.0
				distanceSquared := x*x + y*y

				if distanceSquared <= radius*radius {
					return 1.0 // Inside the circle
				}
				return 0.0 // Outside the circle
			},
			leftIntervalX:      -1,
			rightIntervalX:     1,
			leftIntervalY:      -1,
			rightIntervalY:     1,
			numberOfPartitions: 1000,
			expectedArea:       math.Pi * math.Pow(1, 2), // Area of unit circle is π * r² = π * 1²
			tolerance:          0.01,
			description:        "Area of a circle",
		},
		{
			name: "Area of a Ellipse with semi-major axis 3 and semi-minor axis 2 and center 0",
			expr: func(x, y float64) float64 {
				semiMajorAxisA := 3.0
				semiMinorAxisB := 2.0
				centerX := 0.0
				centerY := 0.0
				
				val := math.Pow(
					(x-centerX)/semiMajorAxisA,
					2,
				) + math.Pow(
					(y-centerY)/semiMinorAxisB,
					2,
				)

				if val <= 1.0 {
					return 1.0
				}
				return 0.0
			},
			leftIntervalX:      -3,
			rightIntervalX:     3,
			leftIntervalY:      -2,
			rightIntervalY:     2,
			numberOfPartitions: 1000,
			expectedArea:       math.Pi * 3 * 2, // Area of ellipse is π * a * b where a = 3, b = 2
			tolerance:          0.01,
			description:        "Area of an ellipse with semi-major axis 3 and semi-minor axis 2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewDoubleIntegralUseCase()

			// Act
			result, err := useCase.CalculateArea(
				t.Context(),
				tc.expr,
				tc.leftIntervalX,
				tc.rightIntervalX,
				tc.leftIntervalY,
				tc.rightIntervalY,
				tc.numberOfPartitions,
			)

			// Assert
			assert.NoError(t, err, "Expected no error for test case: %s", tc.name)
			assert.InDelta(t, tc.expectedArea, result, tc.tolerance,
				"Expected area %v but got %v for %s. Description: %s",
				tc.expectedArea, result, tc.name, tc.description)
		})
	}
}

func TestDoubleIntegralCalculateAreaErrorCases(t *testing.T) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Arrange
	t.Parallel()

	useCase := NewDoubleIntegralUseCase()
	constantFunc := func(x, y float64) float64 {
		return 1.0
	}

	errorTests := []struct {
		name               string
		leftIntervalX      float64
		rightIntervalX     float64
		leftIntervalY      float64
		rightIntervalY     float64
		numberOfPartitions uint64
		expectedError      error
		description        string
	}{
		{
			name:               "Zero width X interval",
			leftIntervalX:      1.0,
			rightIntervalX:     1.0,
			leftIntervalY:      0.0,
			rightIntervalY:     1.0,
			numberOfPartitions: 100,
			expectedError:      ErrZeroWidthInterval,
			description:        "Should return error when X interval has zero width",
		},
		{
			name:               "Zero width Y interval",
			leftIntervalX:      0.0,
			rightIntervalX:     1.0,
			leftIntervalY:      1.0,
			rightIntervalY:     1.0,
			numberOfPartitions: 100,
			expectedError:      ErrZeroWidthInterval,
			description:        "Should return error when Y interval has zero width",
		},
		{
			name:               "Both intervals zero width",
			leftIntervalX:      1.0,
			rightIntervalX:     1.0,
			leftIntervalY:      1.0,
			rightIntervalY:     1.0,
			numberOfPartitions: 100,
			expectedError:      ErrZeroWidthInterval,
			description:        "Should return error when both intervals have zero width",
		},
	}

	for _, tc := range errorTests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := useCase.CalculateArea(
				t.Context(),
				constantFunc,
				tc.leftIntervalX,
				tc.rightIntervalX,
				tc.leftIntervalY,
				tc.rightIntervalY,
				tc.numberOfPartitions,
			)

			// Assert
			assert.Error(t, err, "Expected error for test case: %s", tc.name)
			assert.Equal(
				t,
				tc.expectedError,
				err,
				"Expected specific error for test case: %s",
				tc.name,
			)
			assert.Equal(t, 0.0, result, "Expected zero result when error occurs")

			t.Logf("Test case: %s", tc.name)
			t.Logf("Description: %s", tc.description)
			t.Logf("Expected error: %v", tc.expectedError)
			t.Logf("Actual error: %v", err)
		})
	}
}

func TestDoubleIntegralCalculateAreaZeroPartitions(t *testing.T) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Arrange
	useCase := NewDoubleIntegralUseCase()
	constantFunc := func(x, y float64) float64 {
		return 1.0
	}

	// Act
	result, err := useCase.CalculateArea(
		t.Context(),
		constantFunc,
		0.0, 1.0, // X interval
		0.0, 1.0, // Y interval
		0, // Zero partitions - should be handled gracefully
	)

	// Assert
	assert.NoError(t, err, "Expected no error when partitions is zero")
	assert.Equal(t, 1.0, result, "Expected area of 1.0 for unit square with single partition")

	t.Logf("Zero partitions test - Expected: 1.0, Got: %v", result)
}

func TestDoubleIntegralCalculateAreaBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping benchmark test in short mode")
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn, // Reduce log noise for benchmark
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	useCase := NewDoubleIntegralUseCase()

	// Complex function for performance testing
	complexFunc := func(x, y float64) float64 {
		return math.Sin(x*math.Pi) * math.Cos(y*math.Pi) * math.Exp(-(x*x + y*y))
	}

	partitionSizes := []uint64{10, 50, 100, 500, 1000}

	for _, partitions := range partitionSizes {
		t.Run(fmt.Sprintf("Partitions_%d", partitions), func(t *testing.T) {
			// Act
			result, err := useCase.CalculateArea(
				t.Context(),
				complexFunc,
				-1.0, 1.0, // X interval
				-1.0, 1.0, // Y interval
				partitions,
			)

			// Assert
			assert.NoError(t, err, "Expected no error for partitions: %d", partitions)
			assert.NotZero(t, result, "Expected non-zero result for partitions: %d", partitions)

			t.Logf("Partitions: %d, Result: %v", partitions, result)
		})
	}
}
