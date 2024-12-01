package colours

// Terminal Colors
const (
	BannerColor  = "\033[1;34m%s\033[0m\033[1;36m%s\033[0m"
	TextColor    = "\033[1;0m%s\033[1;32m%s\n\033[0m"
	InfoColor    = "\033[1;0m%s\033[1;35m%s\033[0m"
	NoticeColor  = "\033[1;0m%s\033[1;34m%s\n\033[0m"
	WarningColor = "\033[1;33m%s%s\033[0m"
	ErrorColor   = "\033[1;31m%s%s\033[0m"
	DebugColor   = "\033[0;36m%s%s\033[0m"
)
