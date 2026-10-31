"""Exchange this approved GitHub job identity; never store a cluster password."""
import json, os, pathlib, urllib.parse, urllib.request
endpoint='https://k8s.fabkli.ch/ci/token'
url=os.environ['ACTIONS_ID_TOKEN_REQUEST_URL']+'&audience='+urllib.parse.quote(endpoint,safe='')
req=urllib.request.Request(url,headers={'Authorization':'Bearer '+os.environ['ACTIONS_ID_TOKEN_REQUEST_TOKEN']})
with urllib.request.urlopen(req,timeout=20) as response: identity=json.load(response)['value']
req=urllib.request.Request(endpoint,data=b'',headers={'Authorization':'Bearer '+identity},method='POST')
with urllib.request.urlopen(req,timeout=30) as response: access=json.load(response)
print('::add-mask::'+access['token'])
p=pathlib.Path(os.environ['RUNNER_TEMP'])/'k3s-ci.json'
config={'apiVersion':'v1','kind':'Config','clusters':[{'name':'k3s','cluster':{'server':access['server']}}],
    'contexts':[{'name':'k3s','context':{'cluster':'k3s','user':'ci'}}], 'current-context':'k3s',
    'users':[{'name':'ci','user':{'token':access['token']}}]}
fd=os.open(p,os.O_WRONLY|os.O_CREAT|os.O_TRUNC,0o600)
with os.fdopen(fd,'w') as f: json.dump(config,f)
with open(os.environ['GITHUB_ENV'],'a') as f:f.write('KUBECONFIG='+str(p)+'\n')
print('Scoped deployment credential expires at '+access['expiresAt'])
