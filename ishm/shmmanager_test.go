package ishm

import (
	"fmt"
	"math"
	"testing"
)

func TestInit(t *testing.T) {
	sm := NewShmManager(1024)
	err := sm.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	str := "hahaha"
	sm.WriteBlock("zt", []byte(str))
	data, _ := sm.ReadBlock("zt")
	fmt.Println("Recover: ", string(data))
}

type School struct {
	Name string `json:"name"`
}

type Student struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	School School `json:"school"`
}

type Book struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func TestReadWriteBlock(t *testing.T) {
	sm := NewShmManager(1024)
	err := sm.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	students := []string{}
	N := 10
	for i := 1; i <= N; i++ {
		school := School{
			Name: fmt.Sprintf("School_%f", math.Pow(4, float64(i))),
		}
		studentName := fmt.Sprintf("Student_%f", math.Pow(4, float64(i)))
		students = append(students, studentName)
		student := Student{
			Name:   studentName,
			Age:    21,
			School: school,
		}
		studentData, err := Encode(student)
		if err != nil {
			return
		}
		n, err := sm.WriteBlock(studentName, studentData)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Wrote ", n)
	}
	fmt.Println("Read All Students")
	for _, studentName := range students {
		studentData, err := sm.ReadBlock(studentName)
		if err != nil {
			fmt.Println(err.Error())
		}
		student := &Student{}
		fmt.Println("Read: ", string(studentData))
		err = Decode(studentData, student)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("%s is %d years old and study in %s\n", student.Name, student.Age, student.School.Name)
	}
	sm.Show()
}

func TestDeleteThenWrite(t *testing.T) {
	sm := NewShmManager(2048)
	err := sm.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	students := []string{}
	N := 10
	for i := 1; i <= N; i++ {
		school := School{
			Name: fmt.Sprintf("School_%f", math.Pow(4, float64(i))),
		}
		studentName := fmt.Sprintf("Student_%f", math.Pow(4, float64(i)))
		students = append(students, studentName)
		student := Student{
			Name:   studentName,
			Age:    21,
			School: school,
		}
		studentData, err := Encode(student)
		if err != nil {
			return
		}
		n, err := sm.WriteBlock(studentName, studentData)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Wrote ", n)
	}
	fmt.Println("\nAfter deleting some blocks")
	sm.DeleteBlock(students[0])
	sm.DeleteBlock(students[1])
	sm.DeleteBlock(students[5])
	sm.DeleteBlock(students[6])
	sm.DeleteBlock(students[7])
	sm.Show()
	fmt.Println("\nAfter adding some new blocks")
	newStudentName := "Louis George Maurice Adolphe Roche"
	student := Student{
		Name: newStudentName,
		Age:  21,
		School: School{
			Name: "an emerging school in Shenzhen",
		},
	}
	data1, _ := Encode(student)
	sm.WriteBlock(newStudentName, data1)

	newBookName := "Harry Potter and the Sorcerer's Stone"
	book := Book{
		Name: newBookName,
		Desc: `
		Harry Potter has no idea how famous he is. That's 
		because he's being raised by his miserable aunt and
		uncle who are terrified Harry will learn that he's 
		really a wizard.
		`,
	}
	data2, _ := Encode(book)
	fmt.Println(len(data2))
	sm.WriteBlock(newBookName, data2)
	sm.Show()

	fmt.Println("\nRead these two blocks")
	studentData, err := sm.ReadBlock(newStudentName)
	if err != nil {
		fmt.Println(err.Error())
	}
	newstudent := &Student{}
	err = Decode(studentData, newstudent)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%s is %d years old and study in %s\n", newstudent.Name, newstudent.Age, newstudent.School.Name)

	bookData, err := sm.ReadBlock(newBookName)
	if err != nil {
		fmt.Println(err.Error())
	}
	newbook := &Book{}
	err = Decode(bookData, newbook)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%s is a good book. desc:\n%s\n", newbook.Name, newbook.Desc)
}
