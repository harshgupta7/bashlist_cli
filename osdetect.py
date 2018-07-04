import platform
import os

if platform.system=='Darwin':
	os.remove('bashls-bin-linux')
	os.rename('bashls-bin-mac', 'bashlsbin')
elif platform.system=='Linux':
	os.remove('bashls-bin-linux')
	os.rename('bashls-bin-mac', 'bashlsbin')
