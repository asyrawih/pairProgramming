package helper

import "fmt"

func PanicIfNeed(err error) {
  if err != nil {
    fmt.Println(err)
  }
}
