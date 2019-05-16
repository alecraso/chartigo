package chartigo

import (
        "time"
)

// Team represents a team object in Chartio
type Team struct {
        ID        int
        Name      string
        IsOwner   bool
        CreatedAt time.Time
        UpdatedAt time.Time
        Links     Links
        Embedded  []User
}
