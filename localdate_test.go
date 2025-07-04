package localdate

import (
	"testing"
	"time"
)

func TestToday(t *testing.T) {
	today := Today()
	now := time.Now().UTC()
	expected := NewLocalDate(now.Year(), now.Month(), now.Day())

	if !IsEqual(today, expected) {
		t.Errorf("Today() = %v, want %v", today, expected)
	}
}

func TestAt(t *testing.T) {
	tests := []struct {
		name    string
		dateStr string
		want    LocalDate
		wantErr bool
	}{
		{
			name:    "valid date",
			dateStr: "2023-05-15",
			want:    NewLocalDate(2023, time.May, 15),
			wantErr: false,
		},
		{
			name:    "invalid format",
			dateStr: "15/05/2023",
			want:    LocalDate{},
			wantErr: true,
		},
		{
			name:    "invalid date",
			dateStr: "2023-13-45",
			want:    LocalDate{},
			wantErr: true,
		},
		{
			name:    "empty string",
			dateStr: "",
			want:    LocalDate{},
			wantErr: true,
		},
		{
			name:    "leap year date",
			dateStr: "2024-02-29",
			want:    NewLocalDate(2024, time.February, 29),
			wantErr: false,
		},
		{
			name:    "non-leap year February 29",
			dateStr: "2023-02-29",
			want:    LocalDate{},
			wantErr: true,
		},
		{
			name:    "boundary case - min date",
			dateStr: "0001-01-01",
			want:    NewLocalDate(1, time.January, 1),
			wantErr: false,
		},
		{
			name:    "boundary case - max date",
			dateStr: "9999-12-31",
			want:    NewLocalDate(9999, time.December, 31),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := At(tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("At() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !IsEqual(got, tt.want) {
				t.Errorf("At() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEqual(t *testing.T) {
	tests := []struct {
		name string
		a    LocalDate
		b    LocalDate
		want bool
	}{
		{
			name: "equal dates",
			a:    NewLocalDate(2023, time.May, 15),
			b:    NewLocalDate(2023, time.May, 15),
			want: true,
		},
		{
			name: "different dates",
			a:    NewLocalDate(2023, time.May, 15),
			b:    NewLocalDate(2023, time.May, 16),
			want: false,
		},
		{
			name: "different months",
			a:    NewLocalDate(2023, time.May, 15),
			b:    NewLocalDate(2023, time.June, 15),
			want: false,
		},
		{
			name: "different years",
			a:    NewLocalDate(2023, time.May, 15),
			b:    NewLocalDate(2024, time.May, 15),
			want: false,
		},
		{
			name: "both infinity",
			a:    InfinityDate(),
			b:    InfinityDate(),
			want: true,
		},
		{
			name: "both negative infinity",
			a:    NegInfinityDate(),
			b:    NegInfinityDate(),
			want: true,
		},
		{
			name: "one infinity, one not",
			a:    InfinityDate(),
			b:    NewLocalDate(2023, time.May, 15),
			want: false,
		},
		{
			name: "one negative infinity, one not",
			a:    NegInfinityDate(),
			b:    NewLocalDate(2023, time.May, 15),
			want: false,
		},
		{
			name: "infinity and negative infinity",
			a:    InfinityDate(),
			b:    NegInfinityDate(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("IsEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsBetween(t *testing.T) {
	tests := []struct {
		name   string
		needle LocalDate
		from   LocalDate
		to     LocalDate
		want   bool
	}{
		{
			name:   "needle is between",
			needle: NewLocalDate(2023, time.May, 15),
			from:   NewLocalDate(2023, time.May, 10),
			to:     NewLocalDate(2023, time.May, 20),
			want:   true,
		},
		{
			name:   "needle equals from",
			needle: NewLocalDate(2023, time.May, 10),
			from:   NewLocalDate(2023, time.May, 10),
			to:     NewLocalDate(2023, time.May, 20),
			want:   true,
		},
		{
			name:   "needle equals to",
			needle: NewLocalDate(2023, time.May, 20),
			from:   NewLocalDate(2023, time.May, 10),
			to:     NewLocalDate(2023, time.May, 20),
			want:   true,
		},
		{
			name:   "needle before range",
			needle: NewLocalDate(2023, time.May, 5),
			from:   NewLocalDate(2023, time.May, 10),
			to:     NewLocalDate(2023, time.May, 20),
			want:   false,
		},
		{
			name:   "needle after range",
			needle: NewLocalDate(2023, time.May, 25),
			from:   NewLocalDate(2023, time.May, 10),
			to:     NewLocalDate(2023, time.May, 20),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBetween(tt.needle, tt.from, tt.to); got != tt.want {
				t.Errorf("IsBetween(%v, %v, %v) = %v, want %v", tt.needle, tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestIsAfter(t *testing.T) {
	tests := []struct {
		name string
		a    LocalDate
		b    LocalDate
		want bool
	}{
		{
			name: "a after b",
			a:    NewLocalDate(2023, time.May, 20),
			b:    NewLocalDate(2023, time.May, 10),
			want: true,
		},
		{
			name: "a equals b",
			a:    NewLocalDate(2023, time.May, 15),
			b:    NewLocalDate(2023, time.May, 15),
			want: false,
		},
		{
			name: "a before b",
			a:    NewLocalDate(2023, time.May, 10),
			b:    NewLocalDate(2023, time.May, 20),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAfter(tt.a, tt.b); got != tt.want {
				t.Errorf("IsAfter(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsBefore(t *testing.T) {
	tests := []struct {
		name string
		a    LocalDate
		b    LocalDate
		want bool
	}{
		{
			name: "a before b",
			a:    NewLocalDate(2023, time.May, 10),
			b:    NewLocalDate(2023, time.May, 20),
			want: true,
		},
		{
			name: "a equals b",
			a:    NewLocalDate(2023, time.May, 15),
			b:    NewLocalDate(2023, time.May, 15),
			want: false,
		},
		{
			name: "a after b",
			a:    NewLocalDate(2023, time.May, 20),
			b:    NewLocalDate(2023, time.May, 10),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBefore(tt.a, tt.b); got != tt.want {
				t.Errorf("IsBefore(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAddDays(t *testing.T) {
	tests := []struct {
		name string
		a    LocalDate
		days int
		want LocalDate
	}{
		{
			name: "add positive days",
			a:    NewLocalDate(2023, time.May, 15),
			days: 5,
			want: NewLocalDate(2023, time.May, 20),
		},
		{
			name: "add negative days",
			a:    NewLocalDate(2023, time.May, 15),
			days: -5,
			want: NewLocalDate(2023, time.May, 10),
		},
		{
			name: "add zero days",
			a:    NewLocalDate(2023, time.May, 15),
			days: 0,
			want: NewLocalDate(2023, time.May, 15),
		},
		{
			name: "add days to infinity",
			a:    InfinityDate(),
			days: 10,
			want: InfinityDate(),
		},
		{
			name: "add days to negative infinity",
			a:    NegInfinityDate(),
			days: 10,
			want: NegInfinityDate(),
		},
		{
			name: "cross month boundary forward",
			a:    NewLocalDate(2023, time.May, 30),
			days: 5,
			want: NewLocalDate(2023, time.June, 4),
		},
		{
			name: "cross month boundary backward",
			a:    NewLocalDate(2023, time.June, 2),
			days: -5,
			want: NewLocalDate(2023, time.May, 28),
		},
		{
			name: "cross year boundary forward",
			a:    NewLocalDate(2023, time.December, 29),
			days: 5,
			want: NewLocalDate(2024, time.January, 3),
		},
		{
			name: "cross year boundary backward",
			a:    NewLocalDate(2024, time.January, 3),
			days: -5,
			want: NewLocalDate(2023, time.December, 29),
		},
		{
			name: "leap year February 28 to 29",
			a:    NewLocalDate(2024, time.February, 28),
			days: 1,
			want: NewLocalDate(2024, time.February, 29),
		},
		{
			name: "leap year February 29 to March 1",
			a:    NewLocalDate(2024, time.February, 29),
			days: 1,
			want: NewLocalDate(2024, time.March, 1),
		},
		{
			name: "large number of days",
			a:    NewLocalDate(2023, time.May, 15),
			days: 365,
			want: NewLocalDate(2024, time.May, 14),
		},
		{
			name: "large negative number of days",
			a:    NewLocalDate(2023, time.May, 15),
			days: -365,
			want: NewLocalDate(2022, time.May, 15),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddDays(tt.a, tt.days)
			if !IsEqual(got, tt.want) {
				t.Errorf("AddDays(%v, %v) = %v, want %v", tt.a, tt.days, got, tt.want)
			}
		})
	}
}

func TestAddDate(t *testing.T) {
	tests := []struct {
		name   string
		date   LocalDate
		years  int
		months int
		days   int
		want   LocalDate
	}{
		{
			name:   "add positive years, months, days",
			date:   NewLocalDate(2023, time.May, 15),
			years:  1,
			months: 2,
			days:   10,
			want:   NewLocalDate(2024, time.July, 25),
		},
		{
			name:   "add negative years, months, days",
			date:   NewLocalDate(2023, time.May, 15),
			years:  -1,
			months: -2,
			days:   -10,
			want:   NewLocalDate(2022, time.March, 5),
		},
		{
			name:   "add zero values",
			date:   NewLocalDate(2023, time.May, 15),
			years:  0,
			months: 0,
			days:   0,
			want:   NewLocalDate(2023, time.May, 15),
		},
		{
			name:   "add years only",
			date:   NewLocalDate(2020, time.February, 29), // leap year
			years:  1,
			months: 0,
			days:   0,
			want:   NewLocalDate(2021, time.March, 1), // Go handles leap year adjustment
		},
		{
			name:   "add months causing year overflow",
			date:   NewLocalDate(2023, time.November, 15),
			years:  0,
			months: 3,
			days:   0,
			want:   NewLocalDate(2024, time.February, 15),
		},
		{
			name:   "add days causing month overflow",
			date:   NewLocalDate(2023, time.January, 25),
			years:  0,
			months: 0,
			days:   10,
			want:   NewLocalDate(2023, time.February, 4),
		},
		{
			name:   "add large positive values",
			date:   NewLocalDate(2000, time.January, 1),
			years:  100,
			months: 12,
			days:   365,
			want:   NewLocalDate(2102, time.January, 1), // 12 months = 1 year, 365 days ≈ 1 year
		},
		{
			name:   "add large negative values",
			date:   NewLocalDate(2023, time.December, 31),
			years:  -10,
			months: -24,
			days:   -100,
			want:   NewLocalDate(2011, time.September, 22), // -24 months = -2 years
		},
		{
			name:   "infinity date remains infinity",
			date:   InfinityDate(),
			years:  1,
			months: 1,
			days:   1,
			want:   InfinityDate(),
		},
		{
			name:   "negative infinity date remains negative infinity",
			date:   NegInfinityDate(),
			years:  1,
			months: 1,
			days:   1,
			want:   NegInfinityDate(),
		},
		{
			name:   "leap year edge case - Feb 29 to non-leap year",
			date:   NewLocalDate(2020, time.February, 29),
			years:  1,
			months: 0,
			days:   0,
			want:   NewLocalDate(2021, time.March, 1), // Go adjusts Feb 29 -> Mar 1 in non-leap year
		},
		{
			name:   "month boundary - January 31 + 1 month",
			date:   NewLocalDate(2023, time.January, 31),
			years:  0,
			months: 1,
			days:   0,
			want:   NewLocalDate(2023, time.March, 3), // Go adjusts Jan 31 + 1 month -> Mar 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.date.AddDate(tt.years, tt.months, tt.days)
			if !IsEqual(got, tt.want) {
				t.Errorf("AddDate(%d, %d, %d) = %v, want %v", tt.years, tt.months, tt.days, got, tt.want)
			}
		})
	}
}

func TestToLocalDate(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want LocalDate
	}{
		{
			name: "UTC time",
			time: time.Date(2023, time.May, 15, 14, 30, 45, 0, time.UTC),
			want: NewLocalDate(2023, time.May, 15),
		},
		{
			name: "time with timezone - EST",
			time: time.Date(2023, time.May, 15, 14, 30, 45, 0, time.FixedZone("EST", -5*3600)),
			want: NewLocalDate(2023, time.May, 15),
		},
		{
			name: "time with timezone - JST",
			time: time.Date(2023, time.May, 15, 14, 30, 45, 0, time.FixedZone("JST", 9*3600)),
			want: NewLocalDate(2023, time.May, 15),
		},
		{
			name: "midnight UTC",
			time: time.Date(2023, time.May, 15, 0, 0, 0, 0, time.UTC),
			want: NewLocalDate(2023, time.May, 15),
		},
		{
			name: "end of day UTC",
			time: time.Date(2023, time.May, 15, 23, 59, 59, 999999999, time.UTC),
			want: NewLocalDate(2023, time.May, 15),
		},
		{
			name: "leap year date",
			time: time.Date(2020, time.February, 29, 12, 0, 0, 0, time.UTC),
			want: NewLocalDate(2020, time.February, 29),
		},
		{
			name: "year 1 date",
			time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: NewLocalDate(1, time.January, 1),
		},
		{
			name: "year 9999 date",
			time: time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
			want: NewLocalDate(9999, time.December, 31),
		},
		{
			name: "zero time",
			time: time.Time{},
			want: NewLocalDate(1, time.January, 1),
		},
		{
			name: "unix epoch",
			time: time.Unix(0, 0).UTC(),
			want: NewLocalDate(1970, time.January, 1),
		},
		{
			name: "time with nanoseconds",
			time: time.Date(2023, time.July, 4, 15, 30, 45, 123456789, time.UTC),
			want: NewLocalDate(2023, time.July, 4),
		},
		{
			name: "time in different timezone that crosses date boundary",
			time: time.Date(2023, time.May, 15, 1, 0, 0, 0, time.FixedZone("HST", -10*3600)), // Hawaii time
			want: NewLocalDate(2023, time.May, 15),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLocalDate(tt.time)
			if !IsEqual(got, tt.want) {
				t.Errorf("ToLocalDate(%v) = %v, want %v", tt.time, got, tt.want)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("IsAfter with infinity dates", func(t *testing.T) {
		// Infinity is after any regular date
		if !IsAfter(InfinityDate(), NewLocalDate(9999, time.December, 31)) {
			t.Errorf("Expected infinity to be after max date")
		}

		// Any regular date is after negative infinity
		if !IsAfter(NewLocalDate(1, time.January, 1), NegInfinityDate()) {
			t.Errorf("Expected min date to be after negative infinity")
		}

		// Infinity is after negative infinity
		if !IsAfter(InfinityDate(), NegInfinityDate()) {
			t.Errorf("Expected infinity to be after negative infinity")
		}

		// Infinity is not after infinity
		if IsAfter(InfinityDate(), InfinityDate()) {
			t.Errorf("Expected infinity not to be after infinity")
		}

		// Negative infinity is not after negative infinity
		if IsAfter(NegInfinityDate(), NegInfinityDate()) {
			t.Errorf("Expected negative infinity not to be after negative infinity")
		}
	})

	t.Run("IsBefore with infinity dates", func(t *testing.T) {
		// Any regular date is before infinity
		if !IsBefore(NewLocalDate(9999, time.December, 31), InfinityDate()) {
			t.Errorf("Expected max date to be before infinity")
		}

		// Negative infinity is before any regular date
		if !IsBefore(NegInfinityDate(), NewLocalDate(1, time.January, 1)) {
			t.Errorf("Expected negative infinity to be before min date")
		}

		// Negative infinity is before infinity
		if !IsBefore(NegInfinityDate(), InfinityDate()) {
			t.Errorf("Expected negative infinity to be before infinity")
		}

		// Infinity is not before infinity
		if IsBefore(InfinityDate(), InfinityDate()) {
			t.Errorf("Expected infinity not to be before infinity")
		}

		// Negative infinity is not before negative infinity
		if IsBefore(NegInfinityDate(), NegInfinityDate()) {
			t.Errorf("Expected negative infinity not to be before negative infinity")
		}
	})

	t.Run("IsBetween with infinity dates", func(t *testing.T) {
		// Any date is between negative infinity and infinity
		if !IsBetween(NewLocalDate(2023, time.May, 15), NegInfinityDate(), InfinityDate()) {
			t.Errorf("Expected date to be between negative infinity and infinity")
		}

		// Infinity is between infinity and infinity
		if !IsBetween(InfinityDate(), InfinityDate(), InfinityDate()) {
			t.Errorf("Expected infinity to be between infinity and infinity")
		}

		// Negative infinity is between negative infinity and negative infinity
		if !IsBetween(NegInfinityDate(), NegInfinityDate(), NegInfinityDate()) {
			t.Errorf("Expected negative infinity to be between negative infinity and negative infinity")
		}

		// No date is between infinity and negative infinity (invalid range)
		if IsBetween(NewLocalDate(2023, time.May, 15), InfinityDate(), NegInfinityDate()) {
			t.Errorf("Expected date not to be between infinity and negative infinity")
		}
	})
}
