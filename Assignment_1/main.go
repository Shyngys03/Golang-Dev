package main

import (
    "fmt"
    "employee"
	"manager"
)

func main() {
    emp1 := &manager.Manager{}
    emp1.SetPosition("Manager")
    emp1.SetSalary(60000)
    emp1.SetAddress("123 Main St")

    printEmployeeInfo(emp1)
}

func printEmployeeInfo(emp employee.Employee) {
    fmt.Printf("Position: %s\n", emp.GetPosition())
    fmt.Printf("Salary: $%.2f\n", emp.GetSalary())
    fmt.Printf("Address: %s\n", emp.GetAddress())
}