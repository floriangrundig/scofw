run:
	go run scofw.go

local-gource:
	tail -F .sco/logs/*.gource.log | gource --highlight-dirs --realtime --log-format custom -
