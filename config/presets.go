package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Preset holds all TUI state for a named preset
type Preset struct {
	Name string `toml:"-"` // Not serialized, used internally

	// Master
	MasterEnabled     bool     `toml:"master_enabled"`
	Gain              float32  `toml:"gain"`
	InputFrozen       bool     `toml:"input_frozen"`
	InputFreezeLength float32  `toml:"input_freeze_length"`
	DryWet            float32  `toml:"dry_wet"`
	BlendMode         int      `toml:"blend_mode"`
	EffectsOrder      []string `toml:"effects_order"`

	// Filter
	FilterEnabled   bool    `toml:"filter_enabled"`
	FilterAmount    float32 `toml:"filter_amount"`
	FilterCutoff    float32 `toml:"filter_cutoff"`
	FilterResonance float32 `toml:"filter_resonance"`

	// Overdrive
	OverdriveEnabled bool    `toml:"overdrive_enabled"`
	OverdriveDrive   float32 `toml:"overdrive_drive"`
	OverdriveTone    float32 `toml:"overdrive_tone"`
	OverdriveBias    float32 `toml:"overdrive_bias"`
	OverdriveMix     float32 `toml:"overdrive_mix"`

	// Bitcrush
	BitcrushEnabled    bool    `toml:"bitcrush_enabled"`
	BitDepth           float32 `toml:"bit_depth"`
	BitcrushSampleRate float32 `toml:"bitcrush_sample_rate"`
	BitcrushDrive      float32 `toml:"bitcrush_drive"`
	BitcrushMix        float32 `toml:"bitcrush_mix"`

	// Granular
	GranularEnabled      bool    `toml:"granular_enabled"`
	GranularDensity      float32 `toml:"granular_density"`
	GranularSize         float32 `toml:"granular_size"`
	GranularPitchScatter float32 `toml:"granular_pitch_scatter"`
	GranularPosScatter   float32 `toml:"granular_pos_scatter"`
	GranularMix          float32 `toml:"granular_mix"`
	GranularFrozen       bool    `toml:"granular_frozen"`
	GrainIntensity       string  `toml:"grain_intensity"`

	// Reverb
	ReverbEnabled   bool    `toml:"reverb_enabled"`
	ReverbDecayTime float32 `toml:"reverb_decay_time"`
	ReverbMix       float32 `toml:"reverb_mix"`

	// Delay
	DelayEnabled   bool    `toml:"delay_enabled"`
	DelayTime      float32 `toml:"delay_time"`
	DelayDecayTime float32 `toml:"delay_decay_time"`
	ModRate        float32 `toml:"mod_rate"`
	ModDepth       float32 `toml:"mod_depth"`
	DelayMix       float32 `toml:"delay_mix"`
}

// Hash returns a hash of the preset for dirty detection
func (p *Preset) Hash() string {
	// Simple serialization for hashing
	data, _ := toml.Marshal(p)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// DefaultPreset returns a preset with all effects disabled and values at minimum
func DefaultPreset() Preset {
	return Preset{
		// Master - disabled, zeroed
		MasterEnabled:     false,
		Gain:              0,
		InputFrozen:       false,
		InputFreezeLength: 0.05,
		DryWet:            0,
		BlendMode:         0, // mirror (first option)
		EffectsOrder: []string{
			"filter", "overdrive", "bitcrush",
			"granular", "reverb", "delay",
		},

		// Filter - disabled
		FilterEnabled:   false,
		FilterAmount:    0,
		FilterCutoff:    200,
		FilterResonance: 0,

		// Overdrive - disabled
		OverdriveEnabled: false,
		OverdriveDrive:   0,
		OverdriveTone:    0,
		OverdriveBias:    0,
		OverdriveMix:     0,

		// Bitcrush - disabled
		BitcrushEnabled:    false,
		BitDepth:           4,
		BitcrushSampleRate: 1000,
		BitcrushDrive:      0,
		BitcrushMix:        0,

		// Granular - disabled
		GranularEnabled:      false,
		GranularDensity:      1,
		GranularSize:         0.01,
		GranularPitchScatter: 0,
		GranularPosScatter:   0,
		GranularMix:          0,
		GranularFrozen:       false,
		GrainIntensity:       "subtle", // first option

		// Reverb - disabled
		ReverbEnabled:   false,
		ReverbDecayTime: 0.5,
		ReverbMix:       0,

		// Delay - disabled
		DelayEnabled:   false,
		DelayTime:      0.1,
		DelayDecayTime: 0.1,
		ModRate:        0.1,
		ModDepth:       0,
		DelayMix:       0,
	}
}

// presetsDir returns the presets directory path
func presetsDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "chroma-control", "presets"), nil
}

// ListPresets returns all available preset names (excludes internal _* files)
func ListPresets() ([]string, error) {
	dir, err := presetsDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var presets []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip internal files starting with _
		if strings.HasSuffix(name, ".toml") && !strings.HasPrefix(name, "_") {
			presets = append(presets, strings.TrimSuffix(name, ".toml"))
		}
	}
	return presets, nil
}

// LoadPreset loads a preset by name
func LoadPreset(name string) (Preset, error) {
	if name == "" {
		return DefaultPreset(), fmt.Errorf("preset name cannot be empty")
	}

	dir, err := presetsDir()
	if err != nil {
		return DefaultPreset(), err
	}

	path := filepath.Join(dir, name+".toml")

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultPreset(), fmt.Errorf("preset '%s' not found", name)
	}

	var preset Preset
	if _, err := toml.DecodeFile(path, &preset); err != nil {
		return DefaultPreset(), err
	}

	preset.Name = name
	return preset, nil
}

// SavePreset saves a preset to disk
func SavePreset(preset Preset, name string) error {
	if name == "" {
		return fmt.Errorf("preset name cannot be empty")
	}

	// Validate name (no path traversal, no special chars)
	if strings.ContainsAny(name, "/\\<>:\"|?*") || strings.HasPrefix(name, "_") {
		return fmt.Errorf("invalid preset name")
	}

	dir, err := presetsDir()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, name+".toml")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(preset)
}

// DeletePreset deletes a preset by name
func DeletePreset(name string) error {
	if name == "" || strings.HasPrefix(name, "_") {
		return fmt.Errorf("invalid preset name")
	}

	dir, err := presetsDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, name+".toml")
	return os.Remove(path)
}

// LoadLastPresetName returns the name of the last used preset, or empty if none
func LoadLastPresetName() string {
	dir, err := presetsDir()
	if err != nil {
		return ""
	}

	path := filepath.Join(dir, "_last.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// SaveLastPresetName saves the name of the last used preset
func SaveLastPresetName(name string) error {
	dir, err := presetsDir()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, "_last.toml")
	return os.WriteFile(path, []byte(name), 0644)
}

// LoadAutosave loads the auto-saved session state
func LoadAutosave() (Preset, error) {
	return LoadPreset("_autosave")
}

// SaveAutosave saves current state as auto-saved session
func SaveAutosave(preset Preset) error {
	return SavePreset(preset, "_autosave")
}

// PresetExists checks if a preset exists
func PresetExists(name string) bool {
	_, err := LoadPreset(name)
	return err == nil
}
