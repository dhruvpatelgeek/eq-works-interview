clear
printf "BUILD STARTED @"
date
printf "[PROCRESS STARTED]*********************\n"
export GO111MODULE="off"
go run store/store.go 3000
printf "\n[PROCRESS EXIT]-----------------------\n"