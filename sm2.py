#在官方案例中,未配置mode = cipherMode,导致与别的程序加密解密不统一.  这里采用C1C3C2模式
from gmssl.sm4 import CryptSM4, SM4_ENCRYPT, SM4_DECRYPT
import binascii
import base64
from gmssl import sm2, func,sm3
# GMSSL for Python
# 安装 pip install gmssl
#
inputText = 'Hello World' #  bytes类型

sm3.sm3_hash()

private_key = 'CB772811F1FEF955CE1B4051130870D86CCA6AFEDE806F1E7C225D7359591D2B'
#private_key = 'f1c804b3e1a3aac5b23f6ae7e6d3fa113b0208dab1ed1014666640193bd45ee8'
#public_key = 'qW4BrYsi0fXEStkYlEoHZxCH4SKqMFtcH1OJikKeUht8vg8Mk4mlRfffqVbmdpRy'
public_key = '0475E60AB5B94860DAD0C2D193551A8B7A628A611DF332E23DFCB42F6ECC348653B8A49418E52FF8872B500EEAF8BE8C43B7389D115E91B7432BB1C939E764D31A'

SM2_PRIVATE_KEY_LENGTH=64
private_key = func.random_hex(SM2_PRIVATE_KEY_LENGTH)
print(f"Generated Private Key: {private_key}")

sm2_crypt = sm2.CryptSM2(
    public_key=None,
    private_key=private_key
)


inputTextBytes = inputText.encode()

def ByteToHex( bins ):
    return ''.join( [ "%02X" % x for x in bins ] ).strip()

def HexToByte( hexStr ):
    return bytes.fromhex(hexStr)


#sm2 cipherMode (c1c3c2)
sm2_crypt = sm2.CryptSM2(
    public_key = public_key,
    private_key = private_key
)

enc_data = sm2_crypt.encrypt(inputTextBytes)
dec_data =sm2_crypt.decrypt(enc_data).decode('utf-8')
assert dec_data == dec_data

#print("SM2加密(bytes):",enc_data)
print("SM2加密 (hex):","04" + ByteToHex(enc_data))
print("SM2加密 (hex):","04" + enc_data.hex())
print("SM2解密(utf-8):",dec_data)