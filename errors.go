package e

type Errors struct {
	Code    int
	Message string
}

var (
	ERR_NOERROR = Errors{1, ""}
	ERR_PARAM   = Errors{10100, "参数错误"}
)

// func (err *Errors) GetErr(e Errors) (error Errors) {
// 	error = e
// 	return
// }

// var (
// 	ERR_NOERROR = Errors{1, ""}
// 	ERR_PARAM   = Errors{10100, "参数错误"}
// )

// func GetErr() Errors {
// 	var (
// 		ERR_NOERROR = Errors{1, ""}
// 		ERR_PARAM   = Errors{10100, "参数错误"}
// 	)
// 	return
// }

// type Error struct {
// 	Code    int
// 	Message string
// }

// func GetError(code int) Error {
// 	message := map[int]string{ERR_NOERROR: "错误1", ERR_PARAM: "错误2", 1002: "错误3"}
// 	r := Error{code, message[code]}
// 	return r
// }

// func main() {
//     fmt.Println("Area of r1 is: ", getMessage(ERR_A))
//     fmt.Println("getMessage of r2 is: ", getMessage(ERR_B))
// }

// type ApiErr struct {
// 	ERR_NOERROR Errors
// 	ERR_PARAM   Errors
// }

// var ERR_NOERROR = Errors{1, ""}
// var ERR_PARAM = Errors{10100, "参数错误"}

// var ApiErr = [...]Errors{
// 	ERR_NOERROR: {ERR_NOERROR, ""},
// 	ERR_PARAM: {ERR_PARAM, "参数错误"}
// }

//var ApiErr
// var (
// 	ERR_NOERROR
// )

// func (e Error)  {

// }

// var errors map[int]string

// // var ApiError map[int]string
// var ApiError = map[int]string

// ApiError[ERR_NOERROR] = ""
// ApiError[ERR_PARAM] = "参数错误"
