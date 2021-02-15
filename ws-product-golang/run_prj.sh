printf "BUILD STARTED @"
date
printf "[PROCRESS STARTED]*********************\n"
export GO111MODULE="off"
#     args [client port] [database port]
go run src/server/main.go 3001 3000
printf "\n[PROCRESS EXIT]-----------------------\n"