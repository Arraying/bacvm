package bacvm

import "errors"

var (
	// ErrorBufferEmpty is an error.
	ErrorBufferEmpty = errors.New("attempted to access buffer, but buffer is empty")
	// ErrorFeedQuantity is an error.
	ErrorFeedQuantity = errors.New("unexpected feed quantity (too much or too little)")
	// ErrorFunctionReference is an error.
	ErrorFunctionReference = errors.New("unknown native function reference")
	// ErrorScopeMin is an error.
	ErrorScopeMin = errors.New("cannot finalize scope; already at lowest level")
	// ErrorFeedSize is an error.
	ErrorFeedSize = errors.New("invalid feed size")
	// ErrorFeedType is an error.
	ErrorFeedType = errors.New("unknown feed type")
	// ErrorOperationArgument is an error.
	ErrorOperationArgument = errors.New("unexpected operation argument")
	// ErrorOperationUnknown is an error.
	ErrorOperationUnknown = errors.New("unknown operation")
	// ErrorReadingClose is an error.
	ErrorReadingClose = errors.New("there is no active reading session")
	// ErrorVariableExistance is an error.
	ErrorVariableExistance = errors.New("variable referenced does not exist")
	// ErrorVariableType is an error.
	ErrorVariableType = errors.New("variable type error")
)
