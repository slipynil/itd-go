package errors

import "errors"

// InvalidFileExtension возвращается при попытке загрузить файл с неподдерживаемым расширением.
var InvalidFileExtension = errors.New("invalid file extension")

// TooManyFiles возвращается при попытке загрузить больше файлов, чем разрешено API.
var TooManyFiles = errors.New("too many files")

// NoFileExtension возвращается при попытке загрузить файл без расширения.
var NoFileExtension = errors.New("no file extension")
