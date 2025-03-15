package colours

// Terminal Colors
const (
	BannerColor  = "\033[1;34m%s\033[0m\033[1;36m%s\033[0m"
	TextColor    = "\033[1;0m%s\033[1;32m%s\n\033[0m"
	InfoColor    = "\033[1;37m[\033[1;32mINFO\033[1;37m]\033[0m %s\n"
	NoticeColor  = "\033[1;37m[\033[1;34mNOTICE\033[1;37m]\033[0m %s\n"
	WarningColor = "\033[1;37m[\033[1;34mNOTICE\033[1;37m]\033[0m %s\n"
	ErrorColor   = "\033[1;37m[\033[1;31mERROR\033[1;37m]\033[0m %s\n"
	DebugColor   = "\033[1;37m[\033[0;36mDEBUG\033[1;37m]\033[0m %s\n"
)
