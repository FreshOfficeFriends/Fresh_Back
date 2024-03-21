<div align="center">
    <h1>SSO</h1>
    <h5>
      Один из микросервисов будущего стартапа.
    </h5>
</div>

<details><summary>.env example</summary>
#dev || prod

export ENV=dev

#docker service name or localhost
export DB_HOST=
export DB_PORT=
export DB_NAME=
export DB_USER=
export DB_PASSWORD=
export DB_SSLMODE=

export HASH_SALT=

export GMAIL_PASS=

export DB_DSN=

export JWT_ACCESSTTL=
export JWT_REFRESHTTL=
export JWT_SECRET=

export HTTP_HOST=
export HTTP_PORT=

export redis_host=
