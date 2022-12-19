package dtg

import (
	"strings"
	"testing"
	"time"
)

// UTC-12: Y (e.g., Fiji)
// UTC-11: X (American Samoa)
// UTC-10: W (Honolulu, HI)
// UTC-9: V (Juneau, AK)
// UTC-8: U (PST, Los Angeles, CA)
// UTC-7: T (MST, Denver, CO)
// UTC-6: S (CST, Dallas, TX)
// UTC-5: R (EST, New York, NY)
// UTC-4: Q (Halifax, Nova Scotia)
// UTC-3: P (Buenos Aires, Argentina)
// UTC-2: O (Godthab, Greenland)
// UTC-1: N (Azores)
// UTC+-0: Z (Zulu time)
// UTC+1: A (France)
// UTC+2: B (Athens, Greece)
// UTC+3: C (Arab Standard Time, Iraq, Bahrain, Kuwait, Saudi Arabia, Yemen, Qatar)
// UTC+4: D (Used for Moscow, Russia, and Afghanistan, however, Afghanistan is technically +4:30 from UTC)
// UTC+5: E (Pakistan, Kazakhstan, Tajikistan, Uzbekistan, and Turkmenistan)
// UTC+6: F (Bangladesh)
// UTC+7: G (Thailand)
// UTC+8: H (Beijing, China)
// UTC+9: I (Tokyo, Japan)
// UTC+10: K (Brisbane, Australia)
// UTC+11: L (Sydney, Australia)
// UTC+12: M (Wellington, New Zealand)

func Test_DTG_String(t *testing.T) {
	var err error
	dtg := DTG{}

	testTable := []struct {
		input       string
		expectedDTG string
	}{
		{`152359+0000Dec19`, `152359ZDEC19`},
		{`150102+0100Jan01`, `150102AJAN01`},
		{`181920+0200Feb84`, `181920BFEB84`},
		{`181920+0300Feb26`, `181920CFEB26`},
		{`181920+0400Feb26`, `181920DFEB26`},
		{`181920+0500Feb26`, `181920EFEB26`},
		{`181920+0600Feb26`, `181920FFEB26`},
		{`181920+0700Feb26`, `181920GFEB26`},
		{`181920+0800Feb26`, `181920HFEB26`},
		{`181920+0900Feb26`, `181920IFEB26`},
		{`181920+1000Feb26`, `181920KFEB26`},
		{`181920+1100Feb26`, `181920LFEB26`},
		{`181920+1200Feb26`, `181920MFEB26`},
		{`181920-0100Feb26`, `181920NFEB26`},
		{`181920-0200Feb26`, `181920OFEB26`},
		{`181920-0300Feb26`, `181920PFEB26`},
		{`181920-0400Feb26`, `181920QFEB26`},
		{`181920-0500Feb26`, `181920RFEB26`},
		{`181920-0600Feb26`, `181920SFEB26`},
		{`181920-0700Feb26`, `181920TFEB26`},
		{`181920-0800Feb26`, `181920UFEB26`},
		{`181920-0900Feb26`, `181920VFEB26`},
		{`181920-1000Feb26`, `181920WFEB26`},
		{`181920-1100Feb26`, `181920XFEB26`},
		{`181920-1200Feb26`, `181920YFEB26`},
	}

	for _, v := range testTable {
		dtg.Time, err = time.Parse(expandedDtgLayout, v.input)
		if err != nil {
			t.Fatal(err)
		}
		if dtg.String() != v.expectedDTG {
			t.Errorf("Expected \"%s\", but got \"%s\"", v.expectedDTG, dtg.String())
		}
	}
}

func TestParse(t *testing.T) {
	month := strings.ToUpper(time.Now().Format(monthLayout))
	year := time.Now().Format(yearLayout)
	_, offsetHere := time.Now().Zone()
	loc := time.FixedZone(time.Now().Format(numericTimeZoneLayout), offsetHere)

	// Test explicit time zone letter DTGs
	testTable := []struct {
		input               string
		expectedExpandedDTG string
		expectedDTG         string
	}{
		{`241500Z`, `241500+0000` + month + year, `241500Z` + month + year},
		{`010000N`, `010000-0100` + month + year, `010000N` + month + year},
		{`271337V`, `271337-0900` + month + year, `271337V` + month + year},
		{`271337Y`, `271337-1200` + month + year, `271337Y` + month + year},
		{`271337ZJAN29`, `271337+0000JAN29`, `271337ZJAN29`},
		{`271337BDEC10`, `271337+0200DEC10`, `271337BDEC10`},
	}
	for _, v := range testTable {
		dtg, err := Parse(v.input)
		if err != nil {
			t.Fatal(err)
		}
		if dtg.String() != v.expectedDTG {
			t.Errorf("Expected \"%s\", but got \"%s\"", v.expectedDTG, dtg.String())
		}
		ts := strings.ToUpper(dtg.Time.Format(expandedDtgLayout))
		if ts != v.expectedExpandedDTG {
			t.Errorf("Expected \"%s\" for time-zone-letter-expanded formatted time.Time, but got \"%s\"", v.expectedExpandedDTG, ts)
		}
	}
	// Test local time zone DTGs
	testTable2 := []struct {
		input               string
		expectedExpandedDTG string
	}{
		{`142339`, `142339` + loc.String() + month + year},
		{`142339J`, `142339` + loc.String() + month + year},
	}
	for _, v := range testTable2 {
		dtg, err := Parse(v.input)
		if err != nil {
			t.Fatal(err)
		}
		ts := strings.ToUpper(dtg.Time.Format(expandedDtgLayout))
		if ts != v.expectedExpandedDTG {
			t.Errorf("Expected \"%s\" for time-zone-letter-expanded formatted time.Time, but got \"%s\"", v.expectedExpandedDTG, ts)
		}
	}
	// Test invalid DTGs
	invalidDTGs := []string{"0102", "441200", "441200ZDEC29", "442662", "442663AJAN11", "001022", "121212ÖFEB02", "121212AFXB01"}
	for _, invalidDTG := range invalidDTGs {
		_, err := Parse(invalidDTG)
		if err == nil {
			t.Errorf("Expected to fail on invalid DTG \"%s\", but succeeded", invalidDTG)
		}
	}
}

func TestGetNumericTimeZone(t *testing.T) {
	_, offsetHere := time.Now().Zone()

	testTable := []struct {
		letter string
		offset int
	}{
		{`AB`, 0},
		{`Ö`, 0},
		{`Z`, 0},
		{`J`, offsetHere},
		{`N`, -(60 * 60)},
		{`O`, -(2 * 60 * 60)},
		{`W`, -(10 * 60 * 60)},
		{`Y`, -(12 * 60 * 60)},
		{`K`, 10 * 3600},
		{`L`, 11 * 3600},
		{`M`, 12 * 3600},
		{`D`, 4 * 3600},
	}

	dayHourMinuteMonthYears := [][]string{
		{},
		{`15`},
		{`15`, `20`},
		{`10`, `19`, `11`},
		{`09`, `02`, `01`, `OCT`},
		{`03`, `23`, `31`, `SEP`, `06`},
		{``},
		{`15`, ``},
		{``, ``},
		{`15`, ``, `20`},
		{``, `13`, `21`},
		{``, ``, `21`},
		{``, ``, ``},
		{``, ``, ``, ``},
		{``, ``, ``, `MAR`},
		{``, ``, ``, ``, ``, ``},
	}

	for _, v := range testTable {
		for _, dHMmY := range dayHourMinuteMonthYears {
			loc, err := GetNumericTimeZone(v.letter, dHMmY...)
			if err != nil {
				if v.letter == "Ö" || v.letter == "AB" {
					continue
				}
				if len(dHMmY) > 0 {
					empty := false
					for _, mandatory := range dHMmY {
						if len(mandatory) < 2 {
							empty = true
						}
					}
					if empty {
						continue
					}
				}
				t.Fatal(err)
			}
			_, offset := time.Now().In(loc).Zone()
			if v.offset != offset {
				if v.letter == "J" && len(dHMmY) > 0 {
					t.Logf("Letter %s with dayHourMinuteMonthYear=%s resulted in offset %d instead of expected %d (due to DST probably)", v.letter, strings.Join(dHMmY, ", "), offset, v.offset)
				} else {
					t.Errorf("Expected offset %d for letter \"%s\", but got %d (dayHourMinuteMonthYear=%s)", v.offset, v.letter, offset, strings.Join(dHMmY, ", "))
				}
			}
		}
	}
}

func TestValidate(t *testing.T) {
	dtgsOK := []string{`030102`, `131337Z`, `131337m`, `131337bfeb`, `171819udec28`, `171819AAPR12`}
	dtgsFail := []string{`441200J`, `159218`, `102265ADEC12`, ``, `Hello world`, `12024`, `121314ZAP`}
	for _, dtg := range dtgsOK {
		err := Validate(dtg)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, dtg := range dtgsFail {
		err := Validate(dtg)
		if err == nil {
			t.Errorf("Expected to fail validation for \"%s\", but succeeded", dtg)
		}
	}
}
