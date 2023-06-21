# Mail-Service
--------------

* Serve as service to send e-mail.
* Docker images used (mailhog/mailhog)
* Used lib for mailer
    - github.com/vanng822/go-premailer/premailer
    - github.com/xhit/go-simple-mail/v2
* Mailhog http access (localhost:8025 / docker-machine ip :8025 => http://192.168.99.100:8025/)
* This setting used for development stages
* In prodcution stages mail-service never expose to outside but used/called by other services ex :authentication-services while there are few times failed to login send new email to notice user.