package validator

import (
	"fmt"
	"time"
	"reflect"
	"strings"

	"TransactionTest/internal/delivery/dto"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult содержит результат валидации
type ValidationResult struct {
	IsValid bool            `json:"is_valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// ValidateStruct валидирует структуру и возвращает ошибку
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidateStructWithDetails валидирует структуру и возвращает детальный результат
func ValidateStructWithDetails(s interface{}) ValidationResult {
	err := validate.Struct(s)
	if err == nil {
		return ValidationResult{IsValid: true}
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		message := getValidationMessage(field, tag, param)
		errors = append(errors, ValidationError{
			Field:   strings.ToLower(field),
			Message: message,
		})
	}

	return ValidationResult{
		IsValid: false,
		Errors:  errors,
	}
}

// getValidationMessage возвращает человекочитаемое сообщение об ошибке
func getValidationMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, tag)
	}
}

// ValidateUUID проверяет, является ли строка валидным UUID
func ValidateUUID(uuid string) bool {
	return validate.Var(uuid, "uuid4") == nil
}

// ValidateAmount проверяет сумму транзакции
func ValidateAmount(amount float64) error {
	return validate.Var(amount, "gt=0")
}

// ValidateBalance проверяет баланс кошелька
func ValidateBalance(balance float64) error {
	return validate.Var(balance, "gte=0")
}

// ValidateLimit проверяет лимит для запросов
func ValidateLimit(limit int) error {
	return validate.Var(limit, "gt=0,lte=1000")
}

// ValidateTransactionID проверяет ID транзакции
func ValidateTransactionID(id int64) error {
	return validate.Var(id, "gt=0")
}

// ValidateWalletAddress проверяет адрес кошелька
func ValidateWalletAddress(address string) error {
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	return validate.Var(address, "uuid4")
}

func ValidateSendMoneyRequest(req interface{}) error {
	// Валидируем структуру
	if err := ValidateStruct(req); err != nil {
		return err
	}

	// Получаем значения полей через reflection для дополнительных проверок
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("invalid type: expected struct, got %s", v.Kind())
	}

	fromField := v.FieldByName("From")
	toField := v.FieldByName("To")
	amountField := v.FieldByName("Amount")

	if fromField.IsValid() && toField.IsValid() {
		from := fromField.String()
		to := toField.String()

		if from == to {
			return fmt.Errorf("cannot send money to the same address")
		}
	}

	if amountField.IsValid() && amountField.Kind() == reflect.Float64 {
		amount := amountField.Float()
		if amount <= 0 {
			return fmt.Errorf("amount must be greater than 0")
		}
	} else {
		return fmt.Errorf("invalid amount field")
	}

	return nil
} 

// GetTransactionByInfoRequest проверяет запрос на получения id транзакции по информации о ней
func ValidateGetTransactionByInfoRequest(req *dto.GetTransactionByInfoRequest) error {
    if err := ValidateStruct(req); err != nil {
        return err
    }

    _, err := time.Parse(time.RFC3339, req.CreatedAt)
    if err != nil {
        return fmt.Errorf("created_at must be a valid RFC3339 datetime: %v", err)
    }

    if req.From == req.To {
        return fmt.Errorf("from and to must be different")
    }

    return nil
}