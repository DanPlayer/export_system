package rtnerr

type RtnError interface {
	Error() string
	Rtn() int
}

type RtnCode int

const (
	UnknownError         RtnCode = -1
	NewRtnCodeExistError RtnCode = -2
	CustomizedError      RtnCode = -3
)

var rtnCodeMap *map[RtnCode]string

func (r RtnCode) Error() string {
	s := *rtnCodeMap
	result, ok := s[r]
	if !ok {
		return "err map no enough"
	}
	return result
}

func (r RtnCode) Rtn() int {
	return int(r)
}

func NewErr(errMap map[RtnCode]string) {
	errMap[UnknownError] = "unknown error"
	errMap[NewRtnCodeExistError] = "code exist, set code exclude -1 and -2"
	rtnCodeMap = &errMap
}

func New(err error) RtnError {
	s := *rtnCodeMap
	s[CustomizedError] = err.Error()
	return CustomizedError
}
