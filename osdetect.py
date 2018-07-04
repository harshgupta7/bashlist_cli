import platform

import os
dirname, filename = os.path.split(os.path.abspath(__file__))

if platform.system()=='Darwin':
    os.remove(dirname+'/'+'bashls-bin-linux')
    os.rename((dirname+'/'+'bashls-bin-mac'), 'bashlsbin')
elif platform.system()=='Linux':
	os.remove(dirname+'/'+'bashls-bin-linux')
	os.rename(dirname+'/'+'bashls-bin-mac', 'bashlsbin')
