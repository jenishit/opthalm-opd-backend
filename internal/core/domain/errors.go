package domain

import "errors"

var (
	// ErrInternal is an error for when an internal service fails to process the request
	ErrInternal = errors.New("internal error")
	// ErrDataNotFound is an error for when requested data is not found
	ErrDataNotFound         = errors.New("data not found")
	ErrWrongCurrentPassword = errors.New("current password is wrong")
	// ErrNoUpdatedData is an error for when no data is provided to update
	ErrNoUpdatedData = errors.New("no data to update")
	// ErrConflictingData is an error for when data conflicts with existing data
	ErrConflictingData       = errors.New("data conflicts with existing data in unique column")
	ErrCourseExists          = errors.New("course already exists")
	ErrInvalidUUID           = errors.New("invalid uuid")
	ErrPasswordFormat        = errors.New("password should not contain whitespaces")
	ErrSamePassword          = errors.New("old password cannot be new password")
	ErrExisitingEmail        = errors.New("email already exists")
	ErrBadRequest            = errors.New("bad request")
	ErrEmailPhoneNotVerified = errors.New("email or phone should be verified")
	ErrDayRequest            = errors.New("day should be valid")
	ErrInvalidName           = errors.New("invalid name format")
	ErrInvalidRequest        = errors.New("invalid request")
	ErrMissingField          = errors.New("At least one of the value needs to be selected")

	// ErrInsufficientPayment is an error for when total paid is less than total price
	ErrInsufficientPayment = errors.New("total paid is less than total price")
	// ErrTokenDuration is an error for when the token duration format is invalid
	ErrTokenDuration = errors.New("invalid token duration format")
	// ErrTokenCreation is an error for when the token creation fails
	ErrTokenCreation = errors.New("error creating token")
