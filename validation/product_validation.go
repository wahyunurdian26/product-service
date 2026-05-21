package validation

import (
	"github.com/wahyunurdian26/product-service/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateProductValidation(req *model.CreateProductRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&req.Price, validation.Required, validation.Min(float64(1))),
		validation.Field(&req.Type, validation.Required, validation.In(model.TypeSayuran, model.TypeProtein, model.TypeBuah, model.TypeSnack)),
	)
}
