To create a foreign key in GORM where two keys reside in the same table and the foreign key can be `NULL`, you'll use a self-referential relationship with nullable pointers. Here's how:

**1. Define Your Struct:**

```go
package main

import (
        "log"

        "gorm.io/driver/sqlite"
        "gorm.io/gorm"
)

type Employee struct {
        ID          uint       `gorm:"primaryKey"`
        Name        string     `gorm:"not null"`
        ReportsToID *uint      // Foreign key to Employee's ID (nullable)
        ReportsTo   *Employee  `gorm:"foreignKey:ReportsToID"` // Self-referential relationship
}

func main() {
        db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
        if err != nil {
                log.Fatalf("Failed to connect database: %v", err)
        }

        db.AutoMigrate(&Employee{})

        // Example: Creating employees with a supervisor relationship
        ceo := Employee{Name: "CEO"}
        db.Create(&ceo)

        manager := Employee{Name: "Manager", ReportsToID: &ceo.ID}
        db.Create(&manager)

        employee := Employee{Name: "Employee", ReportsToID: &manager.ID}
        db.Create(&employee)

        // Example: Retrieving an employee with the supervisor
        var retrievedEmployee Employee
        db.Preload("ReportsTo").First(&retrievedEmployee, employee.ID)
        log.Printf("Retrieved employee: %+v", retrievedEmployee)

        // Example: Retrieving CEO who has no reportsTo.
        var retrievedCeo Employee
        db.Preload("ReportsTo").First(&retrievedCeo, ceo.ID)
        log.Printf("Retrieved ceo: %+v", retrievedCeo)
}
```

**Explanation:**

* **`Employee` Struct:**
    * `ID`: The primary key.
    * `ReportsToID *uint`: This is the foreign key column, a nullable pointer to the `ID` of another `Employee` in the same table. It represents the supervisor or manager.
    * `ReportsTo *Employee \`gorm:"foreignKey:ReportsToID"\``: This defines the self-referential relationship.
        * `foreignKey:ReportsToID`: Specifies that `ReportsToID` is the foreign key.
* **`db.AutoMigrate(&Employee{})`:**
    * GORM automatically creates the table and the foreign key constraint.
* **Pointers:**
    * Using pointers (`*uint` and `*Employee`) allows the `ReportsToID` and `ReportsTo` fields to be `NULL`, indicating that an employee may not have a supervisor.
* **`db.Preload("ReportsTo").First(&retrievedEmployee, employee.ID)`:**
    * This shows how to load the related supervisor data when retrieving an employee.
* **Retrieving the CEO:**
    * The last example shows how to retrieve an employee that has no reportsTo value.

**Key Points:**

* **Self-Referential Relationship:** The `ReportsTo` field is a pointer to the `Employee` struct itself, creating the self-referential relationship.
* **Foreign Key Column:** The `ReportsToID` column is the foreign key.
* **Nullable Pointers:** Using pointers allows for optional relationships.
* **`AutoMigrate`:** GORM handles the creation of the foreign key constraint.
* **`Preload`:** The `Preload` function is used to eagerly load the related supervisor data.

**Important Considerations:**

* **Data Types:** Ensure that the data types of the foreign key columns match the data types of the referenced columns.
* **Indexes:** Create indexes on foreign key columns to improve query performance.
* **Referential Integrity:** GORM will generate the necessary foreign key constraints, and your database will enforce referential integrity.
* **Optional Relationships:** Using pointers allows for optional relationships, where an employee may not have a supervisor.
