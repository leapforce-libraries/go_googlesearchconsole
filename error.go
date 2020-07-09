package googlesearchconsole

type GoogleSearchControlError struct {
	Err GoogleSearchControlInnerError `json:"error"`
}
type GoogleSearchControlInnerError struct {
	Code    int                               `json:"code"`
	Message string                            `json:"message"`
	Errors  []GoogleSearchControlInnerstError `json:"errors"`
}

type GoogleSearchControlInnerstError struct {
	Domain  string `json:"domain"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}
