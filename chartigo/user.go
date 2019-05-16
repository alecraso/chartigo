package chartigo

import (
        "time"
)

// User represents a user object in Chartio
type User struct {
        ID          int
        Email       string
        DisplayName string
        CreatedAt   time.Time
        UpdatedAt   time.Time
        Links       Links
        Embedded    []Team
}
