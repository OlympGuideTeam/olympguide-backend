package constants

import "time"

const (
	EmailCodeTtl          = 3 * time.Minute
	MaxVerifyCodeAttempts = 3
)
const MaxSessionConn = 10

const RedisDBNum = 0

const (
	LikeUniverPopularIncr = 10
	LikeUniverPopularDecr = -10
	// LikeOlympPopularIncr LikeProgramUniverPopularIncr = 1
	// LikeProgramUniverPopularDecr = -1
	LikeOlympPopularIncr   = 5
	LikeOlympPopularDecr   = -5
	LikeProgramPopularIncr = 5
	LikeProgramPopularDecr = -5
)

const (
	EmailCodeTopic = "auth.email.code"
	PasswordTopic  = "auth.email.password"
	DiplomaTopic   = "user.diplomas.upload"
)
const (
	ContextUserID   = "user_id"
	ContextUniverID = "university_id"
	ContextEmail    = "email"
)

const (
	EmailTokenTTL = time.Minute * 15
	IDTokenTTL    = time.Minute * 15
)

// Варианты message
const (
	UncompletedRegistration = "Uncompleted registration"
	LoggedIn                = "Logged in"
	LoggedOut               = "Logged out"
	Unauthorized            = "Unauthorized"
	AccountDeleted          = "Account deleted"
	Authorized              = "Authorized"
	SignedUp                = "Signed up"
	EmailConfirmed          = "Email confirmed"
	Disliked                = "Disliked"
	AlreadyDisliked         = "Already disliked"
	Liked                   = "Liked"
	AlreadyLiked            = "Already liked"
)

const PwdGenerateSize = 12
