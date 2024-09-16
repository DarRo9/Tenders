package repository

import "errors"

var (
	FKViolation = "23503"
	UniqueConstraint = "23505"
)

var (
	ErrUserNotExist = errors.New("user does not exist or is invalid")
	ErrRelationNotExist = errors.New("insufficient rights to perform the action")
)

var (
	ErrOrganizationDepencyNotFound = errors.New("it is impossible to create a tender, since there is no organization with this id")
	ErrTenderNotFound = errors.New("tender nettender not found found")
	ErrTenderORVersionNotFound = errors.New("tender or version not found")
	ErrTenderClosed = errors.New("tender closed")
)

var (
	ErrBidDependencyNotFound = errors.New("Can't create an offer because there is no tender or user")
	ErrBidUnique = errors.New("There can be one proposal from an organization for one tender")
	ErrBidTenderNotFound = errors.New("tender or offer not found")
	ErrBidNotFound = errors.New("offer not found")
	ErrBidORVersionNotFound = errors.New("offer or version not found")
	ErrBidReviewsNotFound = errors.New("no tender or reviews found")
)

