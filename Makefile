build:
	go build -o dst/client github.com/indeedhat/smon/client
	go build -o dst/server github.com/indeedhat/smon/server
	go build -o dst/local github.com/indeedhat/smon/local

clean:
	rm -rf dst
