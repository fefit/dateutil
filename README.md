# dateutil

[![Build Status](https://travis-ci.com/fefit/dateutil.svg?branch=master)](https://travis-ci.com/github/fefit/dateutil)
[![codecov](https://codecov.io/gh/fefit/dateutil/branch/master/graph/badge.svg)](https://codecov.io/gh/fefit/dateutil)


An implementation of php methods strtotime, date_format.

## Usage

```go
import (
  "fmt"
  "github.com/fefit/dateutil"
)
func main(){
   if date, err := dateutil.DateTime("Sep-09-2021"); err == nil{
     formatted, _ := DateFormat(date, "Y-m-d")
     fmt.Printf("%s", formatted) // 2021-09-09
   }
}
```


## License

[MIT License](./LICENSE).
