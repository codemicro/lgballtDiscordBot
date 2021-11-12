package pluralkit

import "strconv"

type ErrorCode uint

const (
	ErrorCodeUndefined ErrorCode = iota

	ErrorCodeSystemNotFound ErrorCode = 20001
	ErrorCodeMemberNotFound ErrorCode = 20002
	ErrorCodeMemberRefNotFound ErrorCode = 20003
	ErrorCodeGroupNotFound ErrorCode = 20004
	ErrorCodeGroupRefNotFound ErrorCode = 20005
	ErrorCodeMessageNotFound ErrorCode = 20006
	ErrorCodeSwitchNotFound ErrorCode = 20007
	ErrorCodeSwitchBad ErrorCode = 20008
	ErrorCodeNoSystemGuildSettings ErrorCode = 20009
	ErrorCodeNoMemberGuildSettings ErrorCode = 20010

	ErrorCodeUnauthorizedMemberList ErrorCode = 30001
	ErrorCodeUnauthorizedGroupList ErrorCode = 30002
	ErrorCodeUnauthorizedGroupMemberList ErrorCode = 30003
	ErrorCodeUnauthorizedCurrentFronters ErrorCode = 30004
	ErrorCodeUnauthorizedFrontHistory ErrorCode = 30005
	ErrorCodeMemberNotInSystem ErrorCode = 30006
	ErrorCodeGroupNotInSystem ErrorCode = 30007
	ErrorCodeMemberRefNotInSystem ErrorCode = 30008
	ErrorCodeGroupRefNotInSystem ErrorCode = 30009

	ErrorCodeMissingAutoproxyMember ErrorCode = 40001
	ErrorCodeDuplicateMembersInList ErrorCode = 40002
	ErrorCodeMemberListIdenticalToFronters ErrorCode = 40003
	ErrorCodeSwitchExists ErrorCode = 40004
	ErrorCodeInvalidSwitchID ErrorCode = 40005
)

// DoesHTTPStatusMatchErrorCode returns 1 if the HTTP status and error code match, -1 if they do not match, and 0 if the
// error code was unrecognised.
func DoesHTTPStatusMatchErrorCode(status int, ec ErrorCode) int {

	var firstDigit int
	{
		str := strconv.Itoa(int(ec))
		firstDigit = int(str[0] - 48)
	}

	gi := func(x bool) int {
		if x {
            return 1
        }
        return -1
    }

	var r int

	switch firstDigit {
	case 0:
		r = gi(status == 400 || status == 401 || status == 500)
	case 2:
		r = gi(status == 404)
	case 3:
		r = gi(status == 403)
	case 4:
		r = gi(status == 400)
	}

	return r
}

func DoesErrMatchCode(err error, ec ErrorCode) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ec
	}
	return false
}