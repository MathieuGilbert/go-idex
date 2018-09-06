package idex

import (
	"encoding/json"
	"strconv"
)

// UnmarshalErrorOnType is true when the error was a json.UnmarshalTypeError on type t
func UnmarshalErrorOnType(err error, t string) bool {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		if e.Value == t {
			return true
		}
	}
	return false
}

// UnmarshalJSON custom for Volume
func (v *Volume) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &v.Markets); err != nil {
		// expecting error on totalETH field having a string value
		if !UnmarshalErrorOnType(err, "string") {
			return err
		}
	}
	delete(v.Markets, "totalETH")

	type total struct {
		TotalETH string `json:"totalETH"`
	}
	t := total{}

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	v.TotalETH = t.TotalETH

	return nil
}

// UnmarshalJSON custom for TradeInserted to handle V as string or int
func (ti *TradeInserted) UnmarshalJSON(b []byte) error {
	// prevent infinite loop
	type Alias TradeInserted
	a := Alias{}

	if err := json.Unmarshal(b, &a); err != nil {
		if UnmarshalErrorOnType(err, "string") {
			// auxiliary struct to wrap problem field
			aux := &struct {
				V string `json:"v"`
				*Alias
			}{
				Alias: (*Alias)(ti),
			}

			if err = json.Unmarshal(b, aux); err != nil {
				return err
			}

			// turn it back into an int
			i, errr := strconv.Atoi(aux.V)
			if errr != nil {
				return errr
			}
			a.V = i
		} else {
			return err
		}
	}

	// reassign alias type
	*ti = TradeInserted(a)

	return nil
}

// UnmarshalJSON custom for OrderInserted to handle V as string or int
func (oi *OrderInserted) UnmarshalJSON(b []byte) error {
	// prevent infinite loop
	type Alias OrderInserted
	a := Alias{}

	if err := json.Unmarshal(b, &a); err != nil {
		if UnmarshalErrorOnType(err, "string") {
			// auxiliary struct to wrap problem field
			aux := &struct {
				V string `json:"v"`
				*Alias
			}{
				Alias: (*Alias)(oi),
			}

			if err = json.Unmarshal(b, aux); err != nil {
				return err
			}

			// turn it back into an int
			i, errr := strconv.Atoi(aux.V)
			if errr != nil {
				return errr
			}
			a.V = i
		} else {
			return err
		}
	}

	// reassign alias type
	*oi = OrderInserted(a)

	return nil
}

// UnmarshalJSON custom for PushCancel to handle V as string or int
func (pc *PushCancel) UnmarshalJSON(b []byte) error {
	// prevent infinite loop
	type Alias PushCancel
	a := Alias{}

	if err := json.Unmarshal(b, &a); err != nil {
		if UnmarshalErrorOnType(err, "string") {
			// auxiliary struct to wrap problem field
			aux := &struct {
				V string `json:"v"`
				*Alias
			}{
				Alias: (*Alias)(pc),
			}

			if err = json.Unmarshal(b, aux); err != nil {
				return err
			}

			// turn it back into an int
			i, errr := strconv.Atoi(aux.V)
			if errr != nil {
				return errr
			}
			a.V = i
		} else {
			return err
		}
	}

	// reassign alias type
	*pc = PushCancel(a)

	return nil
}
