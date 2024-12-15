export APP_PORT=8881
export APP_ENV="dev"
export POSTGRE_DSN=postgres://user:password@localhost:5432/user_registration?sslmode=disable
export POSTGRE_MAX_IDLE_CONS=10
export POSTGRE_MAX_OPEN_CONS=10
export POSTGRE_MAX_LIFETIME_CON=0
export MAILER_USER_NAME=rerednaw1969@yahoo.com
export MAILER_USER_PASSWORD=kpmmoktxnrbaprhm
export MAILER_MAIL_BOX=rerednaw1969@yahoo.com
export MAILER_CONNECT_URL=smtp.mail.yahoo.com:587

cd ../cmd
if [ ! -f ./cmd ]
then
    echo "Service not found"
    go build
fi
./cmd
rm cmd

