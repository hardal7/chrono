package apierror

import "errors"

var ErrAlreadyExists = errors.New("resource already exists")
var ErrUnauthorized = errors.New("not authorized to access resource")
