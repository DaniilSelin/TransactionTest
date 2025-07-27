package validator

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStructWithError валидирует структуру и возвращает читаемую ошибку
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var messages []string
	for _, validationErr := range err.(validator.ValidationErrors) {
		field := validationErr.Field()
		tag := validationErr.Tag()
		param := validationErr.Param()

		message := getValidationMessage(field, tag, param)
		messages = append(messages, message)
	}

	return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
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
	if err := ValidateStruct(req); err != nil {
		return err
	}

	v := reflect.ValueOf(req).Elem()
	from := v.FieldByName("From").String()
	to := v.FieldByName("To").String()

	if from == to {
		return fmt.Errorf("cannot send money to the same address")
	}

	return nil
}

// GetTransactionByInfoRequest проверяет запрос на получения id транзакции по информации о ней
func ValidateGetTransactionByInfoRequest(req interface{}) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}

	// Дополнительная валидация времени
	v := reflect.ValueOf(req).Elem()
	createdAt := v.FieldByName("CreatedAt").String()

	if _, err := time.Parse(time.RFC3339, createdAt); err != nil {
		return fmt.Errorf("created_at must be a valid RFC3339 datetime: %v", err)
	}

	// Проверка на одинаковые адреса
	from := v.FieldByName("From").String()
	to := v.FieldByName("To").String()

	if from == to {
		return fmt.Errorf("from and to must be different")
	}

	return nil
}
