allow_k8s_contexts(['book-planner'])
load('ext://namespace', 'namespace_yaml')
load('./deploy/tilt/operator.Tiltfile', 'operator')
load('./deploy/tilt/app.Tiltfile', 'app')

k8s_yaml(namespace_yaml('book-planner'))

operator('db')

app('apis', 'auth', 'false')
app('pages', 'home', 'true')