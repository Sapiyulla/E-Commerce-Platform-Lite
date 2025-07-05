package dto

type PaymentStatus struct {
	Status  string
	Message string
}

func New(status, msg string) PaymentStatus {
	return PaymentStatus{Status: status, Message: msg}
}
