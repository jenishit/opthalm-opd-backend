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
	// ErrExpiredToken is an error for when the access token is expired
	ErrExpiredToken = errors.New("access token has expired")
	ErrExpiredOTP   = errors.New("OTP code invalid")
	// ErrInvalidToken is an error for when the access token is invalid
	ErrInvalidToken = errors.New("access token is invalid")
	// ErrInvalidCredentials is an error for when the credentials are invalid
	ErrInvalidCredentials = errors.New("invalid email or password")
	// ErrEmptyAuthorizationHeader is an error for when the authorization header is empty
	ErrEmptyAuthorizationHeader = errors.New("authorization header is not provided")
	ErrPaymentExceedsDueAmount  = errors.New("payment exceeds due amount")
	ErrInvoiceAlreadyPaid       = errors.New("invoice is already paid")

	ErrEmptyToken = errors.New("token is missing in the url ")

	// ErrInvalidAuthorizationHeader is an error for when the authorization header is invalid
	ErrInvalidAuthorizationHeader = errors.New("authorization header format is invalid")
	// ErrInvalidAuthorizationType is an error for when the authorization type is invalid
	ErrInvalidAuthorizationType = errors.New("authorization type is not supported")
	// ErrUnauthorized is an error for when the user is unauthorized
	ErrUnauthorized = errors.New("user is unauthorized to access the resource")
	// ErrForbidden is an error for when the user is forbidden to access the resource
	ErrForbidden         = errors.New("user is forbidden to access the resource")
	ErrForbiddenToDelete = errors.New("you cannot delete the given live sessions schedule")

	ErrNoRows = errors.New("no rows in result set")
)
