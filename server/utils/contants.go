package utils

import socketio "github.com/googollee/go-socket.io"

var SocketServerInstance = socketio.NewServer(nil)

// Status codes
const (
	HTTP_BAD_REQUEST                     int64 = 400
	HTTP_UNAUTHORIZED                    int64 = 401
	HTTP_PAYMENT_REQUIRED                int64 = 402
	HTTP_FORBIDDEN                       int64 = 403
	HTTP_NOT_FOUND                       int64 = 404
	HTTP_METHOD_NOT_ALLOWED              int64 = 405
	HTTP_NOT_ACCEPTABLE                  int64 = 406
	HTTP_PROXY_AUTHENTICATION_REQUIRED   int64 = 407
	HTTP_REQUEST_TIMEOUT                 int64 = 408
	HTTP_CONFLICT                        int64 = 409
	HTTP_GONE                            int64 = 410
	HTTP_LENGTH_REQUIRED                 int64 = 411
	HTTP_PRECONDITION_FAILED             int64 = 412
	HTTP_PAYLOAD_TOO_LARGE               int64 = 413
	HTTP_URI_TOO_LONG                    int64 = 414
	HTTP_UNSUPPORTED_MEDIA_TYPE          int64 = 415
	HTTP_RANGE_NOT_SATISFIABLE           int64 = 416
	HTTP_EXPECTATION_FAILED              int64 = 417
	HTTP_TEAPOT                          int64 = 418
	HTTP_MISDIRECTED_REQUEST             int64 = 421
	HTTP_UNPROCESSABLE_ENTITY            int64 = 422
	HTTP_LOCKED                          int64 = 423
	HTTP_FAILED_DEPENDENCY               int64 = 424
	HTTP_UPGRADE_REQUIRED                int64 = 426
	HTTP_PRECONDITION_REQUIRED           int64 = 428
	HTTP_TOO_MANY_REQUESTS               int64 = 429
	HTTP_REQUEST_HEADER_FIELDS_TOO_LARGE int64 = 431
	HTTP_UNAVAILABLE_FOR_LEGAL_REASONS   int64 = 451
	HTTP_INTERNAL_SERVER_ERROR           int64 = 500
	HTTP_NOT_IMPLEMENTED                 int64 = 501
	HTTP_BAD_GATEWAY                     int64 = 502
	HTTP_SERVICE_UNAVAILABLE             int64 = 503
	HTTP_GATEWAY_TIMEOUT                 int64 = 504
	HTTP_HTTP_VERSION_NOT_SUPPORTED      int64 = 505
	HTTP_VARIANT_ALSO_NEGOTIATES         int64 = 506
	HTTP_INSUFFICIENT_STORAGE            int64 = 507
	HTTP_LOOP_DETECTED                   int64 = 508
	HTTP_NOT_EXTENDED                    int64 = 510
	HTTP_NETWORK_AUTHENTICATION_REQUIRED int64 = 511
	HTTP_OK                              int64 = 200
	HTTP_NO_CONTENT                      int64 = 204
)

const (
	FAILURE       string = "Failure"
	SUCCESS       string = "Success"
	ACCESS_DENIED string = "Access Denied"
	INVALID_TOKEN string = "Token Absent or Invalid token"
	UNAUTHORIZED  string = "Unauthorized"
)
const (
	Authorization string = "Authorization"
)

const (
	UPGRADE_POWER      int64   = 6
	UPGRADE_SHIFT_TIME float64 = 0.1
	UPGRADE_GRIP       float64 = 1.0
)
const (
	PLAYERID  string = "playerId"
	PLAYER_ID string = "player_id"
)

type ARENA_LEVEL int64

const (
	EASY   ARENA_LEVEL = iota + 1 // EnumIndex = 1
	MEDIUM                        // EnumIndex = 2
	HARD                          // EnumIndex = 3
)

// Arena constants
const (
	EASY_PERK   string = "@every 30m"
	MEDIUM_PERK string = "@every 3h"
	HARD_PERK   string = "@every 7h"

	EASY_PERK_MINUTES   int = 30
	MEDIUM_PERK_MINUTES int = 180
	HARD_PERK_MINUTES   int = 420

	EASY_ARENA_SLOT   int64 = 3
	MEDIUM_ARENA_SLOT int64 = 5
	HARD_ARENA_SLOT   int64 = 7
)
const (
	AdminLogin = iota + 1
	PlayerLogin
	GuestLogin
)

const (
	EASY_ARENA_SERIES   int64 = 3
	MEDIUM_ARENA_SERIES int64 = 5
	HARD_ARENA_SERIES   int64 = 7
)

const (
	D = iota + 1
	C
	B
	A
	S
)

const (
	ADD_EMAIL_REWARD int64 = 500
)

const (
	READ = iota + 1
	UNREAD
)

const (
	COINS = iota + 1
	CASH
	REPAIR_PARTS
	REAL_MONEY
)

// car cutomisation mapping
type COLOR int64

const (
	CRED COLOR = iota + 1
	CGREEN
	CPINK
	CYELLOW
	CBLUE
)

type COLOR_TYPE int

const (
	DEFAULT COLOR_TYPE = iota + 1
	FLUORESCENT
	PASTEL
	GUN_METAL
	SATIN
	METAL
	MILITARY
)

type MILITARY_COLOR int

const (
	MCBASIC MILITARY_COLOR = iota + 1
	MCBLACK
	MCDESERT
	MCTRAM
	MCUCP
)

type CALLIPER_COLOR int

const (
	CCBLACK CALLIPER_COLOR = iota + 1
	CCBLUE
	CCGREEN
	CCPINK
	CCRED
	CCYELLOW
)

type INTERIOR_TYPE int

const (
	ITWHITE INTERIOR_TYPE = iota + 1
	ITPINK
	ITGREEN
	ITRED
	ITBLUE
	ITYELLOW
)

const (
	SIGNUP_SUCCESS     string = "Signup Success"
	PASSWORD_NOT_MATCH string = "Password are not same"
)

type GARAGE_TYPE int

const (
	THE_MU GARAGE_TYPE = iota + 1
	REDS_HOTSPOT
	THE_BEARS_HIDEAWAY
	PRINCES_PALACE
	THE_GREAT_SPOT
)

type CLASS_MAX_UPGRADE int64

const (
	D_CLASS CLASS_MAX_UPGRADE = 3
	C_CLASS CLASS_MAX_UPGRADE = 5
	B_CLASS CLASS_MAX_UPGRADE = 7
	A_CLASS CLASS_MAX_UPGRADE = 10
	S_CLASS CLASS_MAX_UPGRADE = 15
)

type CLASS_OR_MULTIPLIER float64

const (
	D_CLASS_OR CLASS_OR_MULTIPLIER = 2.0
	C_CLASS_OR CLASS_OR_MULTIPLIER = 1.7
	B_CLASS_OR CLASS_OR_MULTIPLIER = 1.5
	A_CLASS_OR CLASS_OR_MULTIPLIER = 1.3
	S_CLASS_OR CLASS_OR_MULTIPLIER = 1.2
)
