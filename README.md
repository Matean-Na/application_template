- cmd/
    - your-app/
      - main.go 
- config/
    - config.go
    - .env
- internal/
    - auth/
        - handlers/
            - auth_handler.go
        - repositories/
            - auth_repository.go
        - services/
            - auth_service.go
        - usecases/
            - auth_usecase.go
    - user/
        - handlers/
            - user_handler.go
        - models/
            - user.go
        - repositories/
            - user_repository.go
        - services/
            - user_service.go
        - usecases/
            - user_usecase.go
    - integrations/
        - rabbitmq/
            - rabbitmq.go
        - redis/
            - redis.go
