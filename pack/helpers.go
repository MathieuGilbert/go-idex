package idex

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// UnmarshalErrorOnType is true when the error was a json.UnmarshalTypeError on type t
func UnmarshalErrorOnType(err error, t string) bool {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		if e.Value == t {
			return true
		}
	case *json.UnmarshalFieldError:
		log.Printf("json.UnmarshalFieldError: %+v\n", e)
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
	type Alias TradeInserted
	aux := &struct {
		V string `json:"v"`
		*Alias
	}{
		Alias: (*Alias)(ti),
	}

	if err := json.Unmarshal(b, aux); err != nil {
		if !UnmarshalErrorOnType(err, "int") && !UnmarshalErrorOnType(err, "string") {
			return err
		}
	}

	// update int value if string value was set
	if i, err := strconv.Atoi(aux.V); err == nil {
		ti.V = i
	}

	return nil
}

// UnmarshalJSON custom for OrderInserted to handle V as string or int
func (ti *OrderInserted) UnmarshalJSON(b []byte) error {
	type Alias OrderInserted
	aux := &struct {
		V string `json:"v"`
		*Alias
	}{
		Alias: (*Alias)(ti),
	}

	if err := json.Unmarshal(b, aux); err != nil {
		if !UnmarshalErrorOnType(err, "int") && !UnmarshalErrorOnType(err, "string") {
			return err
		}
	}

	// update int value if string value was set
	if i, err := strconv.Atoi(aux.V); err == nil {
		ti.V = i
	}

	return nil
}

// UnmarshalJSON custom for PushCancel to handle V as string or int
func (ti *PushCancel) UnmarshalJSON(b []byte) error {
	type Alias PushCancel
	aux := &struct {
		V string `json:"v"`
		*Alias
	}{
		Alias: (*Alias)(ti),
	}

	if err := json.Unmarshal(b, aux); err != nil {
		i := UnmarshalErrorOnType(err, "int")
		fmt.Println(i)
		s := UnmarshalErrorOnType(err, "string")
		fmt.Println(s)
		if !i && !s {
			fmt.Println("huh? why still in here?")
			return err
		}
	}

	// update int value if string value was set
	if i, err := strconv.Atoi(aux.V); err == nil {
		ti.V = i
	}

	return nil
}
