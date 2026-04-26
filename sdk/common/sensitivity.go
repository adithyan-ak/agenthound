package common

type SensitivityLevel string

const (
	SensitivityCritical SensitivityLevel = "critical"
	SensitivityHigh     SensitivityLevel = "high"
	SensitivityMedium   SensitivityLevel = "medium"
	SensitivityLow      SensitivityLevel = "low"
)
