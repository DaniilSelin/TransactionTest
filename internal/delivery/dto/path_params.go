package dto

type AddressPath struct {
    Address string `validate:"required,uuid4"`
}

type CountQuery struct {
    Count int `validate:"required,gt=0"`
}

type TransactionIDPath struct {
    ID int64 `validate:"required,gt=0"`
}
