# go-limit

go-limit is a package for controlling the use of resources in your application. 
You can limit your resources by specifying a limit number and a time interval. 


        import (
        	"fmt"
        	"github.com/CyanSol/go-limit"
        	"time"
        )
        
        func main() {
            //this is a new limiter that limits the execution of a specific resource up to 100 times in a minute
        	limiter1 := go_limit.NewLimiter(60*time.Second, 100)
            //Limiter.try checks if the limit is reached and respondes with an error if the limiter is not active (killed) 
            // and with a response struct that contains if the limit is reached and if yes the time left for the limiter to reset
            res, err := limiter1.Try()
       
            if err != nil {
                fmt.Println(err)
                return
            }
            if res.IsReached {
                fmt.Printf("Wait %.2f seconds to unlock\n", res.TimeLeftToUnlock.Seconds())
            } else {
                fmt.Println("Limit not reached")
            }
            
            //Limiter.Kill makes the limiter inactive permanently 
        	err = limiter1.Kill()
       
        
        }

