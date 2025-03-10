package main

import (
 "fmt"
 "time"

 "golang.org/x/exp/rand"
)

func GenerujPESEL(birthDate time.Time, gender string) [11]int {
 var cyfryPESEL [11]int

 year := birthDate.Year() % 100
 month := int(birthDate.Month())
 day := birthDate.Day()
 randomSerial := rand.Intn(900) + 100

 if birthDate.Year() >= 2000 && birthDate.Year() <= 2099 {
  month += 20
 } else if birthDate.Year() >= 2100 && birthDate.Year() <= 2199 {
  month += 40
 } else if birthDate.Year() >= 2200 && birthDate.Year() <= 2299 {
  month += 60
 }

 genderDigit := 0
 if gender == "M" {
  odd := []int{1, 3, 5, 7, 9}
  randomIndex := rand.Intn(len(odd))
  pick := odd[randomIndex]
  genderDigit = pick
 } else {
  even := []int{0, 2, 4, 6, 8}
  randomIndex := rand.Intn(len(even))
  pick := even[randomIndex]
  genderDigit = pick
 }
 
 cyfryPESEL[0] = year / 10
 cyfryPESEL[1] = year % 10
 cyfryPESEL[2] = month / 10
 cyfryPESEL[3] = month % 10
 cyfryPESEL[4] = day / 10
 cyfryPESEL[5] = day % 10
 cyfryPESEL[6] = randomSerial / 100
 cyfryPESEL[7] = (randomSerial / 10) % 10
 cyfryPESEL[8] = randomSerial % 10
 cyfryPESEL[9] = genderDigit
 cyfryPESEL[10] = generateChecksum(cyfryPESEL)

 return cyfryPESEL
}

func generateChecksum(pesel [11]int) int {
 wagi := [10]int{1, 3, 7, 9, 1, 3, 7, 9, 1, 3}
 var sum int
 for i := 0; i < 10; i++ {
  sum += pesel[i] * wagi[i]
 }
 return (10 - (sum % 10)) % 10

}

func WeryfikujPESEL(cyfryPESEL [11]int) bool {

 var czyPESEL bool
 czyPESEL = true
 if len(cyfryPESEL) != 11 {
  czyPESEL = false
 }
 if generateChecksum(cyfryPESEL) != cyfryPESEL[10] {
	czyPESEL = false
 }

 return czyPESEL
}

func main() {
 birthDate := time.Date(2005, 2, 4, 0, 0, 0, 0, time.UTC)
 pesel := GenerujPESEL(birthDate, "M")
 fmt.Println(pesel)
 fmt.Println("Poprawność numeru PESEL: ", WeryfikujPESEL(pesel))
}

 