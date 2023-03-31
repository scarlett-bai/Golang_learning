package id

// AccountID defines account id object
type AccountID string

func (a AccountID) String() string {
	return string(a)
}

// TripID defines account id object
type TripID string

func (a TripID) String() string {
	return string(a)
}

// IdentityID defines identity id object
type IdentityID string

func (i IdentityID) String() string {
	return string(i)
}

// carID defines car id object
type CarID string

func (i CarID) String() string {
	return string(i)
}
