package helper

import (
    "log"

    "github.com/joho/godotenv"

)

func GetEnvi (key string) string {
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading the .env file")
        return nil
    }
    return os.Getenv(key)
}
