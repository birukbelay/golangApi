package global

const (
	Validation          = "VALIDATION_ERROR"
	ItemInitialization  = "ITEM_INITIALIZATION_ERROR"
	CategoriesInitialization  = "GENRES_INITIALIZATION_ERROR"
	Success             = "SUCCESS"
	InvalidData         = "INVALID_DATA"
	//User Status Codes
	EmailOrPassword     = "EMAIL_OR_PASSWORD_WRONG"
	InvalidEmailOrPhone = "INVALID_EMAIL_OR_PHONE"
	PhoneExists         = "PHONE_EXISTS"
	EmailExists         = "EMAIL_EXISTS"
	Password            = "PASSWORD_NOT_STORED"
	Role                = "ROLE_NOT_ASSIGNED"
	Generic             = "SOME_THING_WRONG"

	ParseFile        = "CANT_PARSE_FORM"
	InvalidFile      = "INVALID_FILE"
	InvalidFileType  = "INVALID_FILE_TYPE"
	FileTooBig       = "FILE_TOO_BIG"
	CantReadFileType = "CANT_READ_FILE_TYPE"
	CantReadFile     = "CANT_READ_FILE"
	CantWriteFile    = "CANT_WRITE_FILE"

	ImageUploadError ="IMAGE_UPLOAD_ERROR"
	JsonUnmarshal    =	"JSON_UNMARSHAL_ERROR"

	UserCreated = "USER_CREATED"
	UserNotCreated = "USER_NOT_CREATED"

	//===============- status codes -=============

	StatusOK      = "200"
	StatusCreated = "201"
	//httpMovedPermanently  = "301" // RFC 7231, 6.4.2
	//httpFound             = 302 // RFC 7231, 6.4.3
	//httpSeeOther          = 303
	StatusUnauthorized = "401"
	StatusBadRequest   = "400"
	//httpUnauthorized                 = 401 // RFC 7235, 3.1
	//
	//httpForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound = "404" // RFC 7231, 6.5.4
	//httpMethodNotAllowed             = 405 // RFC 7231, 6.5.5

	//httpRequestTimeout               = 408 // R
	StatusInternalServerError = "500"
	//httpNotImplemented                = 501 // RFC 7231, 6.6.2
	//httpBadGateway                    = 502 // RFC 7231, 6.6.3
	//httpServiceUnavailable            = 503 //

)

var errorText = map[string]string{
	Success:             "Success",
	Validation:          "validation Error",
	ItemInitialization:  "Item ItemInitialization Error",
	CategoriesInitialization :"Category Initialization",

	EmailOrPassword:     "Your email address or password is wrong",
	InvalidEmailOrPhone: " please enter a valid email address or phone number",
	InvalidData:         "This data Is Invalid",
	Generic:             "Some thing may be wrong please try again",
	PhoneExists:         "Phone Already Exists",
	EmailExists:         "Email Already Exists",
	Password:            "Password Could not be stored",
	Role:                "could not assign role to the user",
	ParseFile:           "Could not parse multipart form",
	InvalidFile:         "Invalid File",
	InvalidFileType:     "File Type Not Allowed",
	FileTooBig:          "File is bigger than max allowed size",
	CantReadFileType:    "Cant Read The File Type",
	CantReadFile:        "Cant Read The File",
	CantWriteFile:       "Cant Write The File",

	ImageUploadError :   "Image Upload Error",
	//------------------------------ Entities -----------------
	UserCreated : " User Created",
	UserNotCreated: "User Not Created",
	JsonUnmarshal :"Json Unmarshal Error",

	//............................... HTTP COPY ............................

	StatusOK:                  "OK",
	StatusCreated:             "http Created",
	StatusBadRequest:          "Bad Request",
	StatusNotFound:            "Not Found",
	StatusInternalServerError: "Internal Server Error",
	StatusUnauthorized:        "Unauthorized",
}

func ErrorText(ErrorCode string) string {
	return errorText[ErrorCode]
}
