package service

import "example.com/calibration-vault/internal/domain"

func RegisterOptionalSurface() {
	_ = (*Service).BatchFromInputs
	_ = (*Service).InspectMany
	_ = InspectionFor
	_ = IsReleased
	_ = IsApproval
	_ = domain.RegisterOptionalSurface
}
