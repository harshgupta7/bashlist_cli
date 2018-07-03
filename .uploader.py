import sys
import os
import requests
if sys.version_info[0]<3
    import ConfigParser
else:
    import configparser as ConfigParser
from collections import OrderedDict
config = ConfigParser.ConfigParser()
inputDict = OrderedDict()
inputDict['acl'] =  'private'
inputDict['key'] = 'fae25fbc-0981-4db7-8918-2138ccfe3aa5__swcli'
inputDict['x-amz-algorithm']= 'AWS4-HMAC-SHA256'
inputDict['x-amz-credential']='AKIAICDC2ZVHRX25EYNQ/20180703/us-east-1/s3/aws4_request'
inputDict['x-amz-date']='20180703T185128Z'
inputDict['policy']='eyJleHBpcmF0aW9uIjogIjIwMTgtMDctMDNUMTk6MDE6MjhaIiwgImNvbmRpdGlvbnMiOiBbeyJhY2wiOiAicHJpdmF0ZSJ9LCBbImNvbnRlbnQtbGVuZ3RoLXJhbmdlIiwgMCwgMTAwMDBdLCB7ImJ1Y2tldCI6ICJiYXNobGlzdC03OCJ9LCB7ImtleSI6ICJmYWUyNWZiYy0wOTgxLTRkYjctODkxOC0yMTM4Y2NmZTNhYTVfX3N3Y2xpIn0sIHsieC1hbXotYWxnb3JpdGhtIjogIkFXUzQtSE1BQy1TSEEyNTYifSwgeyJ4LWFtei1jcmVkZW50aWFsIjogIkFLSUFJQ0RDMlpWSFJYMjVFWU5RLzIwMTgwNzAzL3VzLWVhc3QtMS9zMy9hd3M0X3JlcXVlc3QifSwgeyJ4LWFtei1kYXRlIjogIjIwMTgwNzAzVDE4NTEyOFoifV19'
inputDict['x-amz-signature']='6477f6ca848609c8a94aca520ff3a6eaa55704c8f0815be131a2c031a7c2dd5f'
files = {'file':open('/Users/harsh/desktop/cli/cli/.encontents.bls','rb')}
resp = requests.post(url='https://bashlist-78.s3.amazonaws.com/',data=inputDict,files=files)
print(resp.status_code)
