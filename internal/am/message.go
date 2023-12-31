package am

// Defined in 'assets/locale/xx.json'

var (
	infoMsg  *InfoMessage
	errorMsg *ErrorMessage
)

func init() {
	infoMsg = newInfoMessage()
	errorMsg = newErrorMessage()
}

type InfoMessage struct {
	CreatedMsg string
	UpdatedMsg string
	DeletedMsg string

	// Auth InfoMsg Messages

	SignedUpMsg  string
	ConfirmedMsg string
	SignedInMsg  string
	SignedOutMsg string
}

type ErrorMessage struct {
	// Generic ErrorMsg Messages
	CreateErr string
	IndexErr  string
	GetErr    string
	UpdateErr string
	DeleteErr string

	// Auth ErrorMsg Messages
	SigninErr            string
	SignupErr            string
	CredentialsErr       string
	SignUpErr            string
	SignInErr            string
	ConfirmErr           string
	ConfirmationTokenErr string

	// Validation ErrorMsg Messages
	ProcessErr     string
	InputValuesErr string
	RequiredErr    string
	MinLengthErr   string
	MaxLengthErr   string
	NotAllowedErr  string
	NotEmailErr    string
	ConfMatchErr   string
}

func newInfoMessage() *InfoMessage {
	return &InfoMessage{
		CreatedMsg:   "created-infoMsg-msg",
		UpdatedMsg:   "updated-infoMsg-msg",
		DeletedMsg:   "deleted-infoMsg-msg",
		SignedUpMsg:  "signed-up-info-msg",
		ConfirmedMsg: "confirmed-info-msg",
		SignedInMsg:  "signed-in-info-msg",
		SignedOutMsg: "signed-out-info-msg",
	}
}

func newErrorMessage() *ErrorMessage {
	return &ErrorMessage{
		CreateErr:            "create-err-msg",
		IndexErr:             "index-err-msg",
		GetErr:               "get-err-msg",
		UpdateErr:            "update-err-msg",
		DeleteErr:            "delete-err-msg",
		SigninErr:            "signin-err-msg",
		SignupErr:            "signup-err-msg",
		CredentialsErr:       "credentials_err_msg",
		SignUpErr:            "signup_err_msg",
		SignInErr:            "signin_err_msg",
		ConfirmErr:           "confirm_user_err_msg",
		ConfirmationTokenErr: "confirmation_token_err_msg",
		ProcessErr:           "process-err-msg",
		InputValuesErr:       "input-values-err-msg",
		RequiredErr:          "required-err-msg",
		MinLengthErr:         "min-length-err-msg",
		MaxLengthErr:         "max-length-err-msg",
		NotAllowedErr:        "not-allowed-err-msg",
		NotEmailErr:          "not-email-err-msg",
		ConfMatchErr:         "conf-match-err-msg",
	}
}

func GetInfoMessage() *InfoMessage {
	return infoMsg
}

func GetErrorMessage() *ErrorMessage {
	return errorMsg
}
