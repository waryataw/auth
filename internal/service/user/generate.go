package user

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i Repository -o ./mocks/ -s "_minimock.go"
