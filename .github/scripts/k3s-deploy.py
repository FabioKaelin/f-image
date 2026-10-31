"""Image-only native k3s delivery, with readiness and bounded image rollback."""
import json, os, re, subprocess, time
from pathlib import Path
UID='da171507-3820-4505-857e-b1dc888dc731'
def kube(*args):
    p=subprocess.run(['kubectl','--request-timeout=20s',*args],capture_output=True,text=True,timeout=30)
    if p.returncode:raise RuntimeError(p.stderr.strip())
    return json.loads(p.stdout) if p.stdout.strip().startswith('{') else None

def patch(p, value):
    if p['resource']=='fdeployment':
        body={'spec':{'tag':value}}
    else:body={'spec':{'template':{'spec':{'containers':[{'name':p['container'],'image':value}]}}}}
    kube('-n',p['namespace'],'patch',p['resource'],p['deployment'],'--type='+('merge' if p['resource']=='fdeployment' else 'strategic'),'-p',json.dumps(body))

def wait(plans, deadline, rollback=False):
    remaining=list(plans)
    while remaining and time.monotonic()<deadline:
        for p in remaining[:]:
            d=kube('-n',p['namespace'],'get','deployment',p['deployment'],'-o','json')
            spec=d['spec'];status=d.get('status',{});desired=spec.get('replicas',1)
            image=p['old_image'] if rollback else p['image_ref']
            container=next(c for c in spec['template']['spec']['containers'] if c['name']==p['container'])
            if container['image']==image and status.get('observedGeneration',0)>=d['metadata']['generation'] and status.get('updatedReplicas',0)==desired and status.get('readyReplicas',0)==desired and status.get('availableReplicas',0)==desired and status.get('replicas',0)==desired:
                print(p['deployment']+': ready',flush=True);remaining.remove(p)
        if remaining:time.sleep(10)
    if remaining:raise RuntimeError('Rollout did not become ready: '+', '.join(p['deployment'] for p in remaining))

def main():
    identity=kube('get','namespace','kube-system','-o','json')['metadata']['uid']
    if identity!=UID:raise RuntimeError('Wrong cluster; refusing deployment')
    pressure=kube('-n','f-operator-system','get','configmap','pressure-state','--ignore-not-found','-o','json')
    state=json.loads(pressure.get('data',{}).get('state.json','{}')) if pressure else {}
    if state.get('active') or state.get('stopped'):raise RuntimeError('Pressure recovery active; deployment deferred')
    plans=json.loads(Path('.github/scripts/k3s-projects.json').read_text())
    for p in plans:
        digest=json.loads(Path('release',p['id']+'.json').read_text())['digest']
        if not re.fullmatch(r'sha256:[a-f0-9]{64}',digest):raise ValueError('Invalid image digest')
        p['image_ref']='ghcr.io/fabiokaelin/'+p['image']+'@'+digest
        d=kube('-n',p['namespace'],'get','deployment',p['deployment'],'-o','json')
        if d['spec'].get('replicas',1)==0:raise RuntimeError('Workload suspended: '+p['deployment'])
        p['old_image']=next(c['image'] for c in d['spec']['template']['spec']['containers'] if c['name']==p['container'])
        p['value']=p['image_ref']
        p['old_value']=p['old_image']
        if p['resource']=='fdeployment':
            cr=kube('-n',p['namespace'],'get','fdeployment',p['deployment'],'-o','json')
            if cr['spec'].get('suspend'):raise RuntimeError('Workload suspended: '+p['deployment'])
            if cr['spec'].get('image',p['deployment'])!=p['image']:raise RuntimeError('Image mapping mismatch')
            p['old_value']=cr['spec']['tag']
            p['value']='sha-'+os.environ['GITHUB_SHA']+'@'+digest
            p['image_ref']='ghcr.io/fabiokaelin/'+p['image']+':'+p['value']
    Path(os.environ['RUNNER_TEMP'],'k3s-previous-images.json').write_text(json.dumps(plans,indent=2))
    changed=[]
    try:
        for p in plans:
            changed.append(p);patch(p,p['value'])
        wait(plans,time.monotonic()+540)
    except Exception:
        print('Deployment failed; restoring previous image references.',flush=True)
        rollback_errors=[]
        for p in reversed(changed):
            try:patch(p,p['old_value'])
            except Exception as error:rollback_errors.append(str(error))
        if rollback_errors:raise RuntimeError('Image rollback incomplete: '+'; '.join(rollback_errors))
        wait(changed,time.monotonic()+180,rollback=True)
        raise
    print('Verified all exact images and ready replicas; native storage and placement retained.')

if __name__=='__main__':main()
