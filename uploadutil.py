import sys
import os
import requests
if sys.version_info[0] < 3:
    import ConfigParser
else:
	import configparser as ConfigParser
from collections import OrderedDict
config = ConfigParser.ConfigParser()

path = sys.argv[1]
config.readfp(open('.bashlistuploadconfig.txt'.encode()))



inputDict = OrderedDict()
inputDict['acl'] = config.get('default','acl')
inputDict['key'] = config.get('default','key')
inputDict['x-amz-algorithm'] = config.get('default','x-amz-algorithm')
inputDict['x-amz-credential'] = config.get('default','x-amz-credential')
inputDict['x-amz-date'] = config.get('default','x-amz-date')
inputDict['policy'] = config.get('default','policy')
inputDict['x-amz-signature'] = config.get('default','x-amz-signature')
files = {'file': open(config.get('default','encfile'), 'rb')}
resp = requests.post(url=config.get('default','url'), data=inputDict,files=files)
print(resp.status_code)


