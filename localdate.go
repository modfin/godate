package localdate

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"time"
)

type LocalDate struct {
	Days  int32
	Valid bool
}

const (
	daysInfinity    = math.MaxInt32
	daysNegInfinity = math.MinInt32
)

func NewLocalDate(year int, month time.Month, day int) LocalDate {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	epochDays := int32(t.Unix() / 86400)
	return LocalDate{Days: epochDays, Valid: true}
}

func InfinityDate() LocalDate {
	return LocalDate{Days: daysInfinity, Valid: true}
}

func NegInfinityDate() LocalDate {
	return LocalDate{Days: daysNegInfinity, Valid: true}
}

func (d LocalDate) Time() time.Time {
	if d.IsInfinity() {
		return time.Time{}
	}
	return time.Unix(int64(d.Days)*86400, 0).UTC()
}

func (d LocalDate) IsInfinity() bool {
	return d.Days == daysInfinity
}

func (d LocalDate) IsNegInfinity() bool {
	return d.Days == daysNegInfinity
}

func (d LocalDate) InfinityModifier() int32 {
	switch d.Days {
	case daysInfinity:
		return 1
	case daysNegInfinity:
		return -1
	default:
		return 0

	}
}

func (d LocalDate) MarshalJSON() ([]byte, error) {
	if d.IsInfinity() {
		return []byte(`"infinity"`), nil
	}
	if d.IsNegInfinity() {
		return []byte(`"-infinity"`), nil
	}
	return json.Marshal(d.Time().Format("2006-01-02"))
}

func (d *LocalDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "infinity":
		d.Days = daysInfinity
		return nil
	case "-infinity":
		d.Days = daysNegInfinity
		return nil
	default:
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
		*d = NewLocalDate(t.Year(), t.Month(), t.Day())
		return nil
	}
}

// SQL scanning
func (d *LocalDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*d = NewLocalDate(v.Year(), v.Month(), v.Day())
		return nil
	case string:
		switch v {
		case "infinity":
			d.Days = daysInfinity
			return nil
		case "-infinity":
			d.Days = daysNegInfinity
			return nil
		default:
			t, err := time.Parse("2006-01-02", v)
			if err != nil {
				return err
			}
			*d = NewLocalDate(t.Year(), t.Month(), t.Day())
			return nil
		}
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported Scan, storing %T into LocalDate", value)
	}
}

// SQL value
func (d LocalDate) Value() (driver.Value, error) {
	if d.IsInfinity() {
		return "infinity", nil
	}
	if d.IsNegInfinity() {
		return "-infinity", nil
	}
	return d.Time(), nil
}

// pgtype conversion
func (d LocalDate) PgDate() pgtype.Date {
	if !d.Valid {
		return pgtype.Date{}
	}
	if modifier := d.InfinityModifier(); modifier != 0 {
		return pgtype.Date{
			Valid:            true,
			InfinityModifier: pgtype.InfinityModifier(modifier),
		}
	}
	return pgtype.Date{
		Valid:            true,
		Time:             d.Time(),
		InfinityModifier: pgtype.Finite,
	}
}

func Today() LocalDate {
	now := time.Now().UTC()
	return NewLocalDate(now.Year(), now.Month(), now.Day())
}

func At(at string) (LocalDate, error) {
	t, err := time.Parse("2006-01-02", at)
	if err != nil {
		return LocalDate{}, err
	}
	return NewLocalDate(t.Year(), t.Month(), t.Day()), nil
}

func ToLocalDate(t time.Time) LocalDate {
	return NewLocalDate(t.Year(), t.Month(), t.Day())
}

func IsEqual(a, b LocalDate) bool {
	return a == b
}
func IsAfter(a, b LocalDate) bool {
	return a.Days > b.Days
}
func IsBefore(a, b LocalDate) bool {
	return a.Days < b.Days
}
func AddDays(a LocalDate, n int) LocalDate {
	if a.IsInfinity() || a.IsNegInfinity() {
		return a
	}

	return LocalDate{Days: a.Days + int32(n), Valid: true}
}

// AddDate wraps/replicate the behavior of time.Time and will handle leap years in the same way
func (t LocalDate) AddDate(years int, months int, days int) LocalDate {
	if t.IsInfinity() || t.IsNegInfinity() {
		return t
	}
	year, month, day := t.Time().Date()
	return NewLocalDate(year+years, month+time.Month(months), day+days)
}

func IsBetween(needle, from, to LocalDate) bool {
	return needle.Days >= from.Days && needle.Days <= to.Days
}
