build:
	go build devicefarm-cli.go

install:
	go get github.com/codegangsta/cli
	go get github.com/aws/aws-sdk-go/service/devicefarm
	go get github.com/olekukonko/tablewriter

gox:
	gox -output "dist/devicefarm-cli_0_0_2_{{.OS}}_{{.Arch}}"
