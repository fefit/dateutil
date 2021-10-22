# dateutil

[![Build Status](https://travis-ci.com/fefit/dateutil.svg?branch=master)](https://travis-ci.com/github/fefit/dateutil)
[![codecov](https://codecov.io/gh/fefit/dateutil/branch/master/graph/badge.svg)](https://codecov.io/gh/fefit/dateutil)

An implementation of PHP methods 'strtotime'(not include the relative formats now), 'date_format' in golang.

## Usage

```go
import (
  "fmt"
  du "github.com/fefit/dateutil"
)
func main(){
   if date, err := du.DateTime("Sep-05-2021 06:07:06pm"); err == nil{
     formatted, _ := du.DateFormat(date, "Y-m-d H:i:s")
     fmt.Printf("%s", formatted) // Output: 2021-09-05 18:07:06
   }
}
```

## License

[MIT License](./LICENSE).
