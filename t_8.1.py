import requests
import time

from gmssl import sm2, func,sm3

url = "https://t.jdc.taep.org.cn:17701/api/v1/token/getAccessToken"

playload = {
    "username": "42536518181156892813",
    "password": "LN1c0eZsilRvgl2Mt5bZJIzeqtYkqN",
    "requestTimestamp": "",
    "encrypt_method": "1"
}
#密码通过SM3加密传输，加密方法：
#SM3(username+requestTime+ allocPassword), 生成32的16进制字符串；
playload["requestTimestamp"] = str(int(time.time() ))
xStr=f'{playload["username"]}{playload["requestTimestamp"]}{playload["password"]}'
data_in_bytes = xStr.encode('utf-8')
playload["password"] = sm3.sm3_hash(func.bytes_to_list(data_in_bytes))

print(playload)
print("")
r = requests.post(url, json=playload)
print(r.text)