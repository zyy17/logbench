package distribution

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

	"github.com/zyy17/logbench/pkg/utils"
)

type Distribution struct {
	Ranges []*Range
}

type Range struct {
	Min         int64
	Max         int64
	Probability float64
}

func NewDistribution(inputs []string) (*Distribution, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("empty size range with probability")
	}

	var ranges []*Range
	for _, input := range inputs {
		parts := strings.Split(input, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid size range with probability: %s", input)
		}

		probability, err := utils.ParsePercentage(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid size range with probability: %s", input)
		}

		r, err := NewRange(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid size range with probability: %s", input)
		}
		r.Probability = probability

		ranges = append(ranges, r)
	}

	return &Distribution{
		Ranges: ranges,
	}, nil
}

func (d *Distribution) RandomNumber() int64 {
	sort.Slice(d.Ranges, func(i, j int) bool {
		return d.Ranges[i].Min < d.Ranges[j].Min &&
			d.Ranges[i].Max < d.Ranges[j].Max
	})

	// Calculate the cumulative probability that the generated size will fall into the range.
	for i, sr := range d.Ranges {
		if i > 0 {
			sr.Probability += d.Ranges[i-1].Probability
		}
	}

	// Generate a random number in [0.0, 1.0).
	num := rand.Float64()

	for _, sr := range d.Ranges {
		if num < sr.Probability {
			return utils.RandomNumber(sr.Min, sr.Max)
		}
	}

	// It's impossible to reach here because the sum of all probabilities is 100%. It's already validated before calling this function.
	return 0
}

func (d *Distribution) Validate() error {
	var sum float64
	for i, sr := range d.Ranges {
		// The min should be less than the max.
		if sr.Min >= sr.Max {
			return fmt.Errorf("invalid size range with probability: %v", d.Ranges[i])
		}

		// The probability should be a positive number in (0, 1).
		if sr.Probability < 0 || sr.Probability > 1 {
			return fmt.Errorf("invalid size range with probability: %v", d.Ranges[i])
		}

		// Check if the range overlaps with any other range.
		for j, other := range d.Ranges {
			if i == j {
				continue
			}

			if sr.Max > other.Min && sr.Min < other.Max {
				return fmt.Errorf("the range [%d, %d) overlaps with the range [%d, %d)", sr.Min, sr.Max, other.Min, other.Max)
			}
		}

		sum += sr.Probability
	}

	// The sum of all probabilities should be 1.
	if math.Abs(sum-1) > 0.000001 {
		return fmt.Errorf("invalid size range with probability: %v", d.Ranges)
	}

	return nil
}

func NewRange(input string) (*Range, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid size range: %s", input)
	}

	min, err := utils.ParseSize(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid size range: %s", input)
	}

	max, err := utils.ParseSize(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid size range: %s", input)
	}

	return &Range{
		Min:         min,
		Max:         max,
		Probability: 1.0,
	}, nil
}
