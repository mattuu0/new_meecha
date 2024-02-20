package auth

//ユーザ名とパスワードを検証
func Validate_Name_Password(username string,password string) bool {
	//ユーザネームかパスワードが0文字のとき
	if len(username) == 0 || len(password) == 0 {
		return false
	}

	return true
}