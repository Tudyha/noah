module noah/client

go 1.25.5

replace noah/pkg => ../pkg

require (
	github.com/sevlyar/go-daemon v0.1.6
	github.com/spf13/cobra v1.10.2
	google.golang.org/protobuf v1.36.10
	noah/pkg v0.0.0-00010101000000-000000000000
)

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
