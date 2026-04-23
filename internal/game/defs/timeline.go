package defs

import "encoding/json"

type Timeline struct {
	Clock       string                `nbt:"clock"        json:"clock"`
	PeriodTicks int32                 `nbt:"period_ticks" json:"period_ticks"`
	TimeMarkers map[string]TimeMarker `nbt:"time_markers" json:"time_markers"`
	Tracks      map[string]Track      `nbt:"tracks"       json:"tracks"`
}

type TimeMarker struct {
	ShowInCommands *bool  `nbt:"show_in_commands" json:"show_in_commands,omitempty"`
	Ticks          *int32 `nbt:"ticks"            json:"ticks,omitempty"`
}

func (tm *TimeMarker) UnmarshalJSON(data []byte) error {
	var ticks int32
	if err := json.Unmarshal(data, &ticks); err == nil {
		tm.Ticks = &ticks
		return nil
	}

	var raw struct {
		ShowInCommands bool  `json:"show_in_commands"`
		Ticks          int32 `json:"ticks"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	tm.ShowInCommands = &raw.ShowInCommands
	tm.Ticks = &raw.Ticks
	return nil
}

type Track struct {
	Keyframes []Keyframe    `nbt:"keyframes" json:"keyframes"`
	Modifier  string        `nbt:"modifier"  json:"modifier,omitempty"`
	Ease      *EaseFunction `nbt:"ease"      json:"ease,omitempty"`
}

type Keyframe struct {
	Ticks int32 `nbt:"ticks" json:"ticks"`
	Value any   `nbt:"value" json:"value"`
}

type EaseFunction struct {
	CubicBezier *[4]float32 `nbt:"cubic_bezier" json:"cubic_bezier,omitempty"`
}

func (ef *EaseFunction) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return nil
	}

	var raw struct {
		CubicBezier [4]float32 `json:"cubic_bezier"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	ef.CubicBezier = &raw.CubicBezier
	return nil
}
