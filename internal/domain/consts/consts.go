package consts

const (
	GlobalCacheDb0 = 0
	TimeInForceGTC = "GTC"
	TimeInForceIOK = "IOK"
	TimeInForceFOK = "FOK"

	OrderSideBuy           = "BUY"
	OrderSideSell          = "SELL"
	OrderTypeLimit         = "LIMIT"
	OrderTypeMarket        = "MARKET"
	OrderTypeStopLossLimit = "STOP_LOSS_LIMIT"

	ErrServer         = "server_error"
	ErrUnauthorized   = "unauthorized"
	ErrNotFound       = "not_found"
	ErrWrongInputJson = "wrong_json_format_or_params_type"
	ErrConversion     = "conversion_err"
)
